package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	pb "github.com/aurva/compliance-engine/proto"
)

type ComplianceServer struct {
	pb.UnimplementedComplianceServiceServer
	db *sql.DB
}

func NewComplianceServer(db *sql.DB) *ComplianceServer {
	return &ComplianceServer{db: db}
}

// ReportCloudFinding handles incoming PII findings from scanners
func (s *ComplianceServer) ReportCloudFinding(ctx context.Context, req *pb.FindingRequest) (*pb.FindingResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return &pb.FindingResponse{Success: false, Message: "database error"}, err
	}
	defer tx.Rollback()
	
	// 1. Upsert cloud resource
	var resourceExists bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM cloud_resources WHERE id = $1)
	`, req.ResourceId).Scan(&resourceExists)
	
	if err != nil {
		return &pb.FindingResponse{Success: false, Message: "resource check failed"}, err
	}
	
	accountID := req.AccountId
	if accountID == "" {
		accountID = extractAccountID(ctx, s.db, req.ScanId)
	}
	
	if !resourceExists {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO cloud_resources (id, account_id, resource_type, resource_name, region, size_bytes, last_scanned_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, req.ResourceId, accountID, req.ResourceType, req.ResourceName, req.Region, req.SizeBytes, time.Now())
		
		if err != nil {
			return &pb.FindingResponse{Success: false, Message: "resource insert failed"}, err
		}
	} else {
		// Update last_scanned_at
		_, err = tx.ExecContext(ctx, `
			UPDATE cloud_resources SET last_scanned_at = $1 WHERE id = $2
		`, time.Now(), req.ResourceId)
		
		if err != nil {
			return &pb.FindingResponse{Success: false, Message: "resource update failed"}, err
		}
	}
	
	// 2. Insert PII finding (UPSERT for idempotency)
	pii := req.PiiDetection
	locationJSON, _ := json.Marshal(map[string]interface{}{
		"line_number":  pii.Location.LineNumber,
		"column_name":  pii.Location.ColumnName,
		"offset":       pii.Location.Offset,
		"field_path":   pii.Location.FieldPath,
	})
	
	var findingID int32
	err = tx.QueryRowContext(ctx, `
		INSERT INTO pii_findings (resource_id, pii_type, confidence_score, sample_data, location_info, risk_level, dpdp_section)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (resource_id, pii_type, sample_data) DO UPDATE 
		SET confidence_score = EXCLUDED.confidence_score,
		    location_info = EXCLUDED.location_info,
		    risk_level = EXCLUDED.risk_level,
		    dpdp_section = EXCLUDED.dpdp_section,
		    detected_at = CURRENT_TIMESTAMP
		RETURNING id
	`, req.ResourceId, pii.PiiType, pii.ConfidenceScore, pii.SampleData, locationJSON, pii.RiskLevel, pii.DpdpSection).Scan(&findingID)
	
	if err != nil {
		return &pb.FindingResponse{Success: false, Message: "finding insert failed"}, err
	}
	
	// 3. Create violation record
	violationType := fmt.Sprintf("Unprotected %s detected", pii.PiiType)
	description := fmt.Sprintf("Indian PII (%s) found without adequate protection mechanisms as required by DPDP Act 2023 %s", 
		pii.PiiType, pii.DpdpSection)
	
	_, err = tx.ExecContext(ctx, `
		INSERT INTO violations (finding_id, violation_type, dpdp_section, severity, description, remediation_steps)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING
	`, findingID, violationType, pii.DpdpSection, pii.RiskLevel, description, getRemediationSteps(pii.PiiType))
	
	if err != nil {
		log.Printf("Warning: Failed to create violation: %v", err)
	}
	
	// 4. Audit log
	_, err = tx.ExecContext(ctx, `
		INSERT INTO audit_log (event_type, account_id, resource_id, event_data)
		VALUES ($1, $2, $3, $4)
	`, "finding_detected", accountID, req.ResourceId, locationJSON)
	
	if err != nil {
		log.Printf("Warning: Failed to create audit log: %v", err)
	}
	
	if err := tx.Commit(); err != nil {
		return &pb.FindingResponse{Success: false, Message: "commit failed"}, err
	}
	
	log.Printf("✓ Recorded finding: %s in %s (risk: %s)", pii.PiiType, req.ResourceId, pii.RiskLevel)
	
	return &pb.FindingResponse{
		Success:   true,
		Message:   "Finding recorded",
		FindingId: findingID,
	}, nil
}

// ReportScanProgress updates scan job progress
func (s *ComplianceServer) ReportScanProgress(ctx context.Context, req *pb.ProgressRequest) (*pb.ProgressResponse, error) {
	_, err := s.db.ExecContext(ctx, `
		UPDATE scan_jobs 
		SET progress = $1, resources_scanned = $2, findings_count = $3
		WHERE id = $4
	`, req.ProgressPercent, req.ResourcesScanned, req.FindingsCount, req.ScanId)
	
	if err != nil {
		return &pb.ProgressResponse{Success: false, Message: "progress update failed"}, err
	}
	
	return &pb.ProgressResponse{Success: true, Message: "Progress updated"}, nil
}

// CompleteScan marks a scan job as complete
func (s *ComplianceServer) CompleteScan(ctx context.Context, req *pb.CompleteRequest) (*pb.CompleteResponse, error) {
	status := "completed"
	if !req.Success {
		status = "failed"
	}
	
	_, err := s.db.ExecContext(ctx, `
		UPDATE scan_jobs 
		SET status = $1, progress = 100, resources_scanned = $2, findings_count = $3, 
		    error_message = $4, completed_at = $5
		WHERE id = $6
	`, status, req.TotalResources, req.TotalFindings, req.ErrorMessage, time.Now(), req.ScanId)
	
	if err != nil {
		return &pb.CompleteResponse{Success: false, Message: "completion update failed"}, err
	}
	
	// Audit log
	s.db.ExecContext(ctx, `
		INSERT INTO audit_log (event_type, event_data)
		VALUES ($1, $2)
	`, "scan_completed", fmt.Sprintf(`{"scan_id": "%s", "resources": %d, "findings": %d}`, 
		req.ScanId, req.TotalResources, req.TotalFindings))
	
	log.Printf("✓ Scan completed: %s (resources: %d, findings: %d)", req.ScanId, req.TotalResources, req.TotalFindings)
	
	return &pb.CompleteResponse{Success: true, Message: "Scan completed"}, nil
}

func extractAccountID(ctx context.Context, db *sql.DB, scanID string) string {
	var accountID string
	err := db.QueryRowContext(ctx, `SELECT account_id FROM scan_jobs WHERE id = $1`, scanID).Scan(&accountID)
	if err != nil {
		log.Printf("Warning: Failed to lookup account_id for scan %s: %v", scanID, err)
		return "demo-account" // Fallback
	}
	return accountID
}

func getRemediationSteps(piiType string) string {
	steps := map[string]string{
		"aadhaar":  "1. Encrypt at rest using AES-256\n2. Implement access controls\n3. Enable audit logging\n4. Obtain explicit user consent",
		"pan":      "1. Mask in logs and UI\n2. Encrypt sensitive fields\n3. Restrict database access\n4. Implement data retention policy",
		"gstin":    "1. Apply field-level encryption\n2. Limit access to authorized users\n3. Log all access attempts",
		"phone":    "1. Implement tokenization\n2. Enable MFA for access\n3. Regular access audits",
		"voter_id": "1. Encrypt at rest\n2. Implement RBAC\n3. Enable comprehensive audit trail\n4. Verify consent records",
	}
	
	if steps[piiType] != "" {
		return steps[piiType]
	}
	return "Implement appropriate security controls per DPDP Act 2023 requirements"
}
