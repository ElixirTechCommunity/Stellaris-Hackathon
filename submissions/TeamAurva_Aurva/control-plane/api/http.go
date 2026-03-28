package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/aurva/compliance-engine/cloud-scanner/workers"
	"github.com/aurva/compliance-engine/control-plane/server"
	pb "github.com/aurva/compliance-engine/proto"
	"github.com/aurva/compliance-engine/shared/pii"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HTTPServer struct {
	db               *sql.DB
	complianceServer *server.ComplianceServer
}

func NewHTTPServer(db *sql.DB, cs *server.ComplianceServer) *HTTPServer {
	return &HTTPServer{
		db:               db,
		complianceServer: cs,
	}
}

func (s *HTTPServer) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/scan", s.handleScan)
	mux.HandleFunc("/api/scan/", s.handleScanStatus)
	mux.HandleFunc("/api/findings", s.handleFindings)
	mux.HandleFunc("/api/report/pdf", s.handleReportPDF)
}

// Health check endpoint
func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Check database connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	if err := s.db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "aurva-control-plane",
		"time":    time.Now().UTC(),
	})
}

// POST /api/scan - Trigger new scan with IAM role
func (s *HTTPServer) handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req struct {
		RoleARN     string `json:"role_arn"`
		AccountID   string `json:"account_id"`
		AccountName string `json:"account_name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.RoleARN == "" || req.AccountID == "" {
		http.Error(w, "role_arn and account_id are required", http.StatusBadRequest)
		return
	}
	
	// Validate IAM role assumption
	if err := s.validateRoleAssumption(req.RoleARN); err != nil {
		log.Printf("Warning: Role validation failed (skipping for demo): %v", err)
		// Don't fail the request - continue with scan
	}
	
	// Create or update cloud account
	_, err := s.db.Exec(`
		INSERT INTO cloud_accounts (account_id, role_arn, account_name, status)
		VALUES ($1, $2, $3, 'active')
		ON CONFLICT (account_id) DO UPDATE 
		SET role_arn = EXCLUDED.role_arn, 
		    account_name = EXCLUDED.account_name,
		    updated_at = CURRENT_TIMESTAMP
	`, req.AccountID, req.RoleARN, req.AccountName)
	
	if err != nil {
		http.Error(w, "Failed to register account", http.StatusInternalServerError)
		return
	}
	
	// Create scan job
	scanID := uuid.New().String()
	_, err = s.db.Exec(`
		INSERT INTO scan_jobs (id, account_id, status, started_at)
		VALUES ($1, $2, 'pending', $3)
	`, scanID, req.AccountID, time.Now())
	
	if err != nil {
		http.Error(w, "Failed to create scan job", http.StatusInternalServerError)
		return
	}
	
	// Audit log
	s.db.Exec(`
		INSERT INTO audit_log (event_type, account_id, event_data)
		VALUES ($1, $2, $3)
	`, "scan_started", req.AccountID, fmt.Sprintf(`{"scan_id": "%s", "role_arn": "%s"}`, scanID, req.RoleARN))
	
	// Spawn scanner as goroutine
	go s.runScanner(scanID, req.AccountID, req.RoleARN)
	
	log.Printf("✓ Created scan job: %s for account: %s", scanID, req.AccountID)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"scan_id":    scanID,
		"account_id": req.AccountID,
		"status":     "pending",
		"message":    "Scan initiated successfully",
	})
}

// GET /api/findings - Get aggregated findings with compliance score
func (s *HTTPServer) handleFindings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		accountID = "demo-account" // Default for demo
	}
	
	// Get compliance summary
	var summary struct {
		TotalResources   int
		TotalFindings    int
		CriticalFindings int
		HighFindings     int
		MediumFindings   int
		LowFindings      int
		ComplianceScore  float64
	}
	
	err := s.db.QueryRow(`
		SELECT 
			COALESCE(COUNT(DISTINCT cr.id), 0) as total_resources,
			COALESCE(COUNT(pf.id), 0) as total_findings,
			COALESCE(COUNT(CASE WHEN pf.risk_level = 'critical' THEN 1 END), 0) as critical_findings,
			COALESCE(COUNT(CASE WHEN pf.risk_level = 'high' THEN 1 END), 0) as high_findings,
			COALESCE(COUNT(CASE WHEN pf.risk_level = 'medium' THEN 1 END), 0) as medium_findings,
			COALESCE(COUNT(CASE WHEN pf.risk_level = 'low' THEN 1 END), 0) as low_findings,
			COALESCE(
				ROUND(
					CASE 
						WHEN COUNT(DISTINCT cr.id) = 0 THEN 100
						ELSE (GREATEST(COUNT(DISTINCT cr.id)::NUMERIC - 
							  COUNT(DISTINCT CASE WHEN pf.risk_level IN ('critical', 'high') THEN cr.id END)::NUMERIC, 0) / 
							 COUNT(DISTINCT cr.id)::NUMERIC) * 100
					END, 
					2
				), 100
			) as compliance_score
		FROM cloud_accounts ca
		LEFT JOIN cloud_resources cr ON ca.account_id = cr.account_id
		LEFT JOIN pii_findings pf ON cr.id = pf.resource_id
		WHERE ca.account_id = $1
		GROUP BY ca.account_id
	`, accountID).Scan(
		&summary.TotalResources,
		&summary.TotalFindings,
		&summary.CriticalFindings,
		&summary.HighFindings,
		&summary.MediumFindings,
		&summary.LowFindings,
		&summary.ComplianceScore,
	)
	
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying summary: %v", err)
		http.Error(w, "Failed to retrieve findings", http.StatusInternalServerError)
		return
	}
	
	// Get detailed findings by resource
	rows, err := s.db.Query(`
		SELECT 
			cr.id as resource_id,
			cr.resource_type,
			cr.resource_name,
			cr.region,
			pf.pii_type,
			pf.confidence_score,
			pf.risk_level,
			pf.dpdp_section,
			pf.sample_data,
			pf.detected_at
		FROM cloud_resources cr
		JOIN pii_findings pf ON cr.id = pf.resource_id
		WHERE cr.account_id = $1
		ORDER BY pf.risk_level DESC, pf.detected_at DESC
		LIMIT 1000
	`, accountID)
	
	if err != nil {
		log.Printf("Error querying findings: %v", err)
		http.Error(w, "Failed to retrieve findings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	
	type Finding struct {
		ResourceID      string  `json:"resource_id"`
		ResourceType    string  `json:"resource_type"`
		ResourceName    string  `json:"resource_name"`
		Region          string  `json:"region"`
		PIIType         string  `json:"pii_type"`
		ConfidenceScore float32 `json:"confidence_score"`
		RiskLevel       string  `json:"risk_level"`
		DPDPSection     string  `json:"dpdp_section"`
		SampleData      string  `json:"sample_data"`
		DetectedAt      time.Time `json:"detected_at"`
	}
	
	var findings []Finding
	for rows.Next() {
		var f Finding
		if err := rows.Scan(&f.ResourceID, &f.ResourceType, &f.ResourceName, &f.Region,
			&f.PIIType, &f.ConfidenceScore, &f.RiskLevel, &f.DPDPSection,
			&f.SampleData, &f.DetectedAt); err != nil {
			log.Printf("Error scanning finding: %v", err)
			continue
		}
		findings = append(findings, f)
	}
	
	// Get PII type breakdown
	typeBreakdown := map[string]int{
		"aadhaar":      0,
		"pan":          0,
		"gstin":        0,
		"phone":        0,
		"voter_id":     0,
		"bank_account": 0,
	}
	bRows, err := s.db.Query(`
		SELECT pf.pii_type, COUNT(*)
		FROM pii_findings pf
		JOIN cloud_resources cr ON pf.resource_id = cr.id
		WHERE cr.account_id = $1
		GROUP BY pf.pii_type
	`, accountID)
	if err == nil {
		for bRows.Next() {
			var pType string
			var count int
			if err := bRows.Scan(&pType, &count); err == nil {
				typeBreakdown[pType] = count
			}
		}
		bRows.Close()
	}

	response := map[string]interface{}{
		"account_id":       accountID,
		"total_resources":  summary.TotalResources,
		"total_pii_count":  summary.TotalFindings,
		"high_risk_count":  summary.CriticalFindings + summary.HighFindings,
		"critical_count":   summary.CriticalFindings,
		"compliance_score": summary.ComplianceScore,
		"risk_distribution": map[string]int{
			"critical": summary.CriticalFindings,
			"high":     summary.HighFindings,
			"medium":   summary.MediumFindings,
			"low":      summary.LowFindings,
		},
		"breakdown": typeBreakdown,
		"findings":  findings,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /api/report/pdf - Generate compliance report
func (s *HTTPServer) handleReportPDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		accountID = "demo-account"
	}

	// ── Fetch full summary ───────────────────────────────────────────────────
	var summary struct {
		AccountName      string
		RoleARN          string
		TotalResources   int
		TotalFindings    int
		CriticalFindings int
		HighFindings     int
		MediumFindings   int
		LowFindings      int
		ComplianceScore  float64
	}

	err := s.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(ca.account_name, ''),
			COALESCE(ca.role_arn, ''),
			COALESCE(COUNT(DISTINCT cr.id), 0),
			COALESCE(COUNT(pf.id), 0),
			COALESCE(COUNT(CASE WHEN pf.risk_level = 'critical' THEN 1 END), 0),
			COALESCE(COUNT(CASE WHEN pf.risk_level = 'high'     THEN 1 END), 0),
			COALESCE(COUNT(CASE WHEN pf.risk_level = 'medium'   THEN 1 END), 0),
			COALESCE(COUNT(CASE WHEN pf.risk_level = 'low'      THEN 1 END), 0),
			COALESCE(ROUND(
				CASE 
					WHEN COUNT(DISTINCT cr.id) = 0 THEN 100
					ELSE (GREATEST(COUNT(DISTINCT cr.id)::NUMERIC - 
						  COUNT(DISTINCT CASE WHEN pf.risk_level IN ('critical', 'high') THEN cr.id END)::NUMERIC, 0) / 
						 COUNT(DISTINCT cr.id)::NUMERIC) * 100
				END, 
				2
			), 100)
		FROM cloud_accounts ca
		LEFT JOIN cloud_resources cr ON ca.account_id = cr.account_id
		LEFT JOIN pii_findings pf ON cr.id = pf.resource_id
		WHERE ca.account_id = $1
		GROUP BY ca.account_id, ca.account_name, ca.role_arn
	`, accountID).Scan(
		&summary.AccountName, &summary.RoleARN,
		&summary.TotalResources, &summary.TotalFindings,
		&summary.CriticalFindings, &summary.HighFindings,
		&summary.MediumFindings, &summary.LowFindings,
		&summary.ComplianceScore,
	)
	if err != nil {
		log.Printf("Error generating report: %v", err)
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	// ── Fetch PII type breakdown ─────────────────────────────────────────────
	type PIIBreakdown struct {
		PIIType  string
		Count    int
		RiskLevel string
	}
	var piiBreakdown []PIIBreakdown
	bRows, _ := s.db.QueryContext(ctx, `
		SELECT pii_type, COUNT(*) as cnt, MAX(risk_level)
		FROM pii_findings pf
		JOIN cloud_resources cr ON pf.resource_id = cr.id
		WHERE cr.account_id = $1
		GROUP BY pii_type ORDER BY cnt DESC
	`, accountID)
	if bRows != nil {
		defer bRows.Close()
		for bRows.Next() {
			var b PIIBreakdown
			bRows.Scan(&b.PIIType, &b.Count, &b.RiskLevel)
			piiBreakdown = append(piiBreakdown, b)
		}
	}

	// ── Fetch violations ─────────────────────────────────────────────────────
	type Violation struct {
		ViolationType    string
		DPDPSection      string
		Severity         string
		Description      string
		RemediationSteps string
		Count            int
	}
	var violations []Violation
	vRows, _ := s.db.QueryContext(ctx, `
		SELECT v.violation_type, v.dpdp_section, v.severity, v.description,
		       COALESCE(v.remediation_steps,''), COUNT(*) as cnt
		FROM violations v
		JOIN pii_findings pf ON v.finding_id = pf.id
		JOIN cloud_resources cr ON pf.resource_id = cr.id
		WHERE cr.account_id = $1
		GROUP BY v.violation_type, v.dpdp_section, v.severity, v.description, v.remediation_steps
		ORDER BY CASE v.severity WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 ELSE 4 END
	`, accountID)
	if vRows != nil {
		defer vRows.Close()
		for vRows.Next() {
			var v Violation
			vRows.Scan(&v.ViolationType, &v.DPDPSection, &v.Severity, &v.Description, &v.RemediationSteps, &v.Count)
			violations = append(violations, v)
		}
	}

	// ── Fetch all findings ───────────────────────────────────────────────────
	type Finding struct {
		ResourceName    string
		ResourceType    string
		Region          string
		PIIType         string
		RiskLevel       string
		DPDPSection     string
		ConfidenceScore float64
		SampleData      string
		DetectedAt      time.Time
	}
	var findings []Finding
	fRows, err := s.db.QueryContext(ctx, `
		SELECT cr.resource_name, cr.resource_type, COALESCE(cr.region,''),
		       pf.pii_type, pf.risk_level, pf.dpdp_section,
		       pf.confidence_score, COALESCE(pf.sample_data,''),
		       pf.detected_at
		FROM cloud_resources cr
		JOIN pii_findings pf ON cr.id = pf.resource_id
		WHERE cr.account_id = $1
		ORDER BY
			CASE pf.risk_level WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 ELSE 4 END,
			pf.detected_at DESC
		LIMIT 1000
	`, accountID)
	if err != nil {
		log.Printf("Error querying findings for PDF: %v", err)
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}
	defer fRows.Close()
	for fRows.Next() {
		var f Finding
		if err := fRows.Scan(&f.ResourceName, &f.ResourceType, &f.Region,
			&f.PIIType, &f.RiskLevel, &f.DPDPSection,
			&f.ConfidenceScore, &f.SampleData, &f.DetectedAt); err == nil {
			findings = append(findings, f)
		} else {
			log.Printf("Error scanning finding row for PDF: %v", err)
		}
	}

	// ── Build PDF ────────────────────────────────────────────────────────────
	pdf := gofpdf.New("P", "mm", "A4", "")
	pageNum := 0
	newPage := func() {
		pdf.AddPage()
		pageNum++
	}
	sectionHeader := func(title string) {
		pdf.SetFont("Arial", "B", 13)
		pdf.SetTextColor(15, 23, 42)
		pdf.Cell(0, 7, title)
		pdf.Ln(8)
		pdf.SetDrawColor(99, 102, 241)
		pdf.SetLineWidth(0.7)
		x, y := pdf.GetXY()
		pdf.Line(x, y, x+190, y)
		pdf.SetLineWidth(0.2)
		pdf.SetDrawColor(0, 0, 0)
		pdf.Ln(4)
	}
	addFooter := func() {
		pdf.SetY(-13)
		pdf.SetFont("Arial", "I", 7)
		pdf.SetTextColor(130, 130, 130)
		pdf.CellFormat(130, 5,
			fmt.Sprintf("Aurva DPDP Act 2023 Compliance Report  |  Account: %s  |  %s",
				accountID, time.Now().Format("02 Jan 2006 15:04 IST")),
			"", 0, "L", false, 0, "")
		pdf.CellFormat(60, 5, fmt.Sprintf("Page %d", pageNum), "", 0, "R", false, 0, "")
		pdf.SetTextColor(0, 0, 0)
	}

	scoreColor := getScoreColor(summary.ComplianceScore)
	riskColor := func(level string) (int, int, int) {
		switch level {
		case "critical":
			return 220, 38, 38
		case "high":
			return 234, 88, 12
		case "medium":
			return 161, 98, 7
		default:
			return 21, 128, 61
		}
	}
	truncate := func(s string, n int) string {
		if len(s) > n {
			return s[:n-3] + "..."
		}
		return s
	}

	// ═══════════════════════════════════════════════════════════════════════
	// PAGE 1 — COVER
	// ═══════════════════════════════════════════════════════════════════════
	newPage()

	// Dark header banner
	pdf.SetFillColor(15, 23, 42)
	pdf.Rect(0, 0, 210, 55, "F")

	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(10, 14)
	pdf.Cell(0, 10, "DPDP Act 2023 Compliance Report")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(148, 163, 184)
	pdf.SetXY(10, 27)
	pdf.Cell(0, 6, "India's Digital Personal Data Protection Act — Automated Audit")

	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(10, 36)
	pdf.Cell(0, 6, fmt.Sprintf("Prepared by Aurva Security Engine  |  %s", time.Now().Format("02 January 2006")))

	// Indigo accent bar
	pdf.SetFillColor(99, 102, 241)
	pdf.Rect(0, 55, 210, 2, "F")

	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(10, 65)

	// Account metadata box
	pdf.SetFillColor(248, 250, 252)
	pdf.RoundedRect(10, 63, 190, 36, 3, "1234", "FD")
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(99, 102, 241)
	pdf.SetXY(15, 67)
	pdf.Cell(0, 5, "ACCOUNT DETAILS")
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(15, 74)
	pdf.Cell(40, 5, "Account Name:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 5, summary.AccountName)
	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(15, 81)
	pdf.Cell(40, 5, "Account ID:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 5, accountID)
	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(15, 88)
	pdf.Cell(40, 5, "IAM Role ARN:")
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 5, truncate(summary.RoleARN, 70))
	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(15, 95)
	pdf.Cell(40, 5, "Report Generated:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 5, time.Now().Format("02 Jan 2006, 15:04:05 IST"))

	pdf.SetTextColor(0, 0, 0)

	// ── Compliance Score big display ─────────────────────────────────────────
	pdf.SetXY(10, 108)
	pdf.SetFillColor(scoreColor[0], scoreColor[1], scoreColor[2])
	pdf.RoundedRect(10, 108, 90, 40, 4, "1234", "F")
	pdf.SetFont("Arial", "B", 28)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(14, 113)
	pdf.Cell(82, 12, fmt.Sprintf("%.1f%%", math.Max(0, summary.ComplianceScore)))
	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(14, 127)
	pdf.Cell(82, 6, "COMPLIANCE SCORE")
	pdf.SetFont("Arial", "", 9)
	pdf.SetXY(14, 134)
	scoreLabel := "Excellent"
	if summary.ComplianceScore < 90 {
		scoreLabel = "Needs Improvement"
	}
	if summary.ComplianceScore < 70 {
		scoreLabel = "Critical — Immediate Action Required"
	}
	pdf.Cell(82, 5, scoreLabel)

	// ── 4 stat boxes on the right ────────────────────────────────────────────
	stats := []struct {
		label string
		value string
		r, g, b int
	}{
		{"Resources Scanned", fmt.Sprintf("%d", summary.TotalResources), 30, 41, 59},
		{"PII Findings", fmt.Sprintf("%d", summary.TotalFindings), 30, 41, 59},
		{"Critical", fmt.Sprintf("%d", summary.CriticalFindings), 220, 38, 38},
		{"High Risk", fmt.Sprintf("%d", summary.HighFindings), 234, 88, 12},
	}
	for i, st := range stats {
		bx := float64(108 + (i%2)*47)
		by := float64(108 + (i/2)*20)
		pdf.SetFillColor(248, 250, 252)
		pdf.RoundedRect(bx, by, 43, 17, 2, "1234", "FD")
		pdf.SetFont("Arial", "B", 14)
		pdf.SetTextColor(st.r, st.g, st.b)
		pdf.SetXY(bx+2, by+2)
		pdf.Cell(39, 7, st.value)
		pdf.SetFont("Arial", "", 7)
		pdf.SetTextColor(100, 116, 139)
		pdf.SetXY(bx+2, by+10)
		pdf.Cell(39, 5, st.label)
	}
	pdf.SetTextColor(0, 0, 0)

	// ── Executive Summary ────────────────────────────────────────────────────
	pdf.SetXY(10, 155)
	sectionHeader("Executive Summary")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(51, 65, 85)
	summary1 := fmt.Sprintf(
		"This automated compliance audit was performed on AWS account '%s' using Aurva's "+
			"VPC-native security engine. A total of %d cloud resources were scanned across S3 "+
			"and RDS. The scan identified %d instances of personal data as defined under India's "+
			"Digital Personal Data Protection Act 2023.",
		summary.AccountName, summary.TotalResources, summary.TotalFindings)
	pdf.MultiCell(190, 5, summary1, "", "L", false)
	pdf.Ln(3)

	riskStatement := ""
	if summary.CriticalFindings > 0 {
		riskStatement = fmt.Sprintf(
			"CRITICAL ALERT: %d finding(s) involve Aadhaar numbers — classified as sensitive personal data "+
				"under Section 8(4) of the DPDP Act. Immediate remediation is mandatory. Failure to protect "+
				"Aadhaar data may attract penalties under Section 66 of the Act.", summary.CriticalFindings)
	} else if summary.HighFindings > 0 {
		riskStatement = fmt.Sprintf(
			"%d high-risk finding(s) were detected involving financial identifiers (PAN, Voter ID, Phone). "+
				"These require encryption and access controls under Sections 8(1)-8(3) of the DPDP Act.",
			summary.HighFindings)
	} else {
		riskStatement = "No critical or high-risk PII was detected. Continue to monitor and schedule periodic scans."
	}
	pdf.SetFont("Arial", "B", 10)
	if summary.CriticalFindings > 0 {
		pdf.SetTextColor(220, 38, 38)
	} else if summary.HighFindings > 0 {
		pdf.SetTextColor(234, 88, 12)
	} else {
		pdf.SetTextColor(21, 128, 61)
	}
	pdf.MultiCell(190, 5, riskStatement, "", "L", false)
	pdf.SetTextColor(0, 0, 0)

	// ── PII Type Breakdown table ─────────────────────────────────────────────
	pdf.Ln(5)
	sectionHeader("PII Type Breakdown")

	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(30, 41, 59)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(60, 7, "PII Category", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, "Occurrences", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, "Risk Level", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 7, "DPDP Act Reference", "1", 1, "C", true, 0, "")

	piiDPDP := map[string]string{
		"aadhaar":  "Section 8(4) — Sensitive Personal Data",
		"pan":      "Section 8(3) — Financial Identifier",
		"phone":    "Section 8(1) — Contact Information",
		"voter_id": "Section 8(2) — Government ID",
		"gstin":    "Section 7(1) — Business Identifier",
	}

	pdf.SetFont("Arial", "", 9)
	for i, b := range piiBreakdown {
		if i%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		r, g, bl := riskColor(b.RiskLevel)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(60, 7, strings.ToUpper(b.PIIType), "1", 0, "L", true, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%d", b.Count), "1", 0, "C", true, 0, "")
		pdf.SetTextColor(r, g, bl)
		pdf.CellFormat(40, 7, strings.ToUpper(b.RiskLevel), "1", 0, "C", true, 0, "")
		pdf.SetTextColor(0, 0, 0)
		dpdp := piiDPDP[b.PIIType]
		if dpdp == "" {
			dpdp = "Section 7 — Personal Data"
		}
		pdf.CellFormat(60, 7, dpdp, "1", 1, "L", true, 0, "")
	}

	addFooter()

	// ═══════════════════════════════════════════════════════════════════════
	// PAGE 2 — DPDP ACT OBLIGATIONS & VIOLATIONS
	// ═══════════════════════════════════════════════════════════════════════
	newPage()

	sectionHeader("DPDP Act 2023 — Relevant Obligations")

	obligations := []struct{ section, title, text string }{
		{
			"Section 4",
			"Grounds for Processing Personal Data",
			"Personal data may only be processed for a lawful purpose for which the data principal has given consent, or for certain legitimate uses specified in the Act.",
		},
		{
			"Section 8",
			"Obligations of Data Fiduciaries",
			"Every data fiduciary must ensure accuracy, completeness, and consistency of personal data. Appropriate technical and organisational measures must be implemented to protect personal data from breaches.",
		},
		{
			"Section 8(4)",
			"Sensitive Personal Data — Aadhaar",
			"Aadhaar numbers are classified as sensitive personal data. Data fiduciaries must implement enhanced security controls, explicit consent mechanisms, and strict access limitations.",
		},
		{
			"Section 17",
			"Data Localisation",
			"Certain categories of personal data may be subject to data localisation requirements as notified by the Central Government. Cross-border transfers require compliance with Section 16.",
		},
		{
			"Section 25 / Section 66",
			"Penalties",
			"Non-compliance may attract penalties up to INR 250 crore per instance of violation. The Data Protection Board of India may investigate and adjudicate complaints.",
		},
	}

	pdf.SetFont("Arial", "", 9)
	for _, ob := range obligations {
		pdf.SetFillColor(239, 246, 255)
		pdf.SetFont("Arial", "B", 9)
		pdf.SetTextColor(30, 64, 175)
		pdf.CellFormat(190, 6, fmt.Sprintf("  %s — %s", ob.section, ob.title), "LTR", 1, "L", true, 0, "")
		pdf.SetFont("Arial", "", 9)
		pdf.SetTextColor(30, 41, 59)
		pdf.SetFillColor(255, 255, 255)
		pdf.MultiCell(190, 5, "  "+ob.text, "LBR", "L", true)
		pdf.Ln(2)
	}

	pdf.Ln(4)
	sectionHeader("Detected Violations")

	if len(violations) == 0 {
		pdf.SetFont("Arial", "I", 10)
		pdf.SetTextColor(100, 116, 139)
		pdf.Cell(0, 8, "No violations recorded.")
		pdf.SetTextColor(0, 0, 0)
	} else {
		pdf.SetFont("Arial", "B", 9)
		pdf.SetFillColor(30, 41, 59)
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(50, 7, "Violation Type", "1", 0, "L", true, 0, "")
		pdf.CellFormat(25, 7, "Severity", "1", 0, "C", true, 0, "")
		pdf.CellFormat(40, 7, "DPDP Section", "1", 0, "L", true, 0, "")
		pdf.CellFormat(15, 7, "Count", "1", 0, "C", true, 0, "")
		pdf.CellFormat(60, 7, "Description", "1", 1, "L", true, 0, "")
		pdf.SetTextColor(0, 0, 0)

		for i, v := range violations {
			if i%2 == 0 {
				pdf.SetFillColor(248, 250, 252)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}
			r, g, b := riskColor(v.Severity)
			pdf.SetTextColor(0, 0, 0)
			pdf.CellFormat(50, 7, truncate(v.ViolationType, 30), "1", 0, "L", true, 0, "")
			pdf.SetTextColor(r, g, b)
			pdf.CellFormat(25, 7, strings.ToUpper(v.Severity), "1", 0, "C", true, 0, "")
			pdf.SetTextColor(0, 0, 0)
			pdf.CellFormat(40, 7, truncate(v.DPDPSection, 25), "1", 0, "L", true, 0, "")
			pdf.CellFormat(15, 7, fmt.Sprintf("%d", v.Count), "1", 0, "C", true, 0, "")
			pdf.CellFormat(60, 7, truncate(v.Description, 38), "1", 1, "L", true, 0, "")
		}
	}

	addFooter()

	// ═══════════════════════════════════════════════════════════════════════
	// PAGE 3+ — DETAILED FINDINGS TABLE
	// ═══════════════════════════════════════════════════════════════════════
	newPage()
	sectionHeader("Detailed PII Findings")

	tableHeader := func() {
		pdf.SetFont("Arial", "B", 8)
		pdf.SetFillColor(30, 41, 59)
		pdf.SetTextColor(255, 255, 255)
		pdf.CellFormat(55, 7, "Resource / Location", "1", 0, "L", true, 0, "")
		pdf.CellFormat(16, 7, "Type", "1", 0, "C", true, 0, "")
		pdf.CellFormat(15, 7, "Risk", "1", 0, "C", true, 0, "")
		pdf.CellFormat(55, 7, "DPDP Section", "1", 0, "L", true, 0, "")
		pdf.CellFormat(20, 7, "Confidence", "1", 0, "C", true, 0, "")
		pdf.CellFormat(29, 7, "Masked Value", "1", 1, "L", true, 0, "")
		pdf.SetTextColor(0, 0, 0)
	}
	tableHeader()

	pdf.SetFont("Arial", "", 8)
	rowCount := 0
	for i, f := range findings {
		// New page if needed
		if pdf.GetY() > 270 {
			addFooter()
			newPage()
			sectionHeader(fmt.Sprintf("Detailed PII Findings (continued)"))
			tableHeader()
			pdf.SetFont("Arial", "", 8)
		}

		if i%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		r, g, b := riskColor(f.RiskLevel)
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(55, 7, truncate(f.ResourceName, 34), "1", 0, "L", true, 0, "")
		pdf.CellFormat(16, 7, strings.ToUpper(f.PIIType), "1", 0, "C", true, 0, "")
		pdf.SetTextColor(r, g, b)
		pdf.CellFormat(15, 7, strings.ToUpper(f.RiskLevel), "1", 0, "C", true, 0, "")
		pdf.SetTextColor(0, 0, 0)
		pdf.CellFormat(55, 7, truncate(f.DPDPSection, 34), "1", 0, "L", true, 0, "")
		pdf.CellFormat(20, 7, fmt.Sprintf("%.0f%%", f.ConfidenceScore*100), "1", 0, "C", true, 0, "")
		pdf.CellFormat(29, 7, truncate(f.SampleData, 18), "1", 1, "L", true, 0, "")
		rowCount++
	}

	if rowCount == 0 {
		pdf.SetFont("Arial", "I", 10)
		pdf.SetTextColor(100, 116, 139)
		pdf.Cell(0, 8, "No findings recorded for this account.")
		pdf.SetTextColor(0, 0, 0)
	}

	addFooter()

	// ═══════════════════════════════════════════════════════════════════════
	// FINAL PAGE — REMEDIATION PLAYBOOK
	// ═══════════════════════════════════════════════════════════════════════
	newPage()
	sectionHeader("Remediation Playbook")

	playbook := []struct{ title, steps string }{
		{
			"Aadhaar (Critical — Section 8(4))",
			"1. Encrypt all Aadhaar fields using AES-256 at rest and TLS 1.2+ in transit.\n" +
				"2. Apply field-level encryption before storing in databases.\n" +
				"3. Implement role-based access control — only authorised personnel may access.\n" +
				"4. Obtain explicit written consent from data principals before collection.\n" +
				"5. Enable comprehensive audit logging of every access event.\n" +
				"6. Set a data retention policy — delete after purpose is served.",
		},
		{
			"PAN Card (High — Section 8(3))",
			"1. Mask PAN in all application logs and UI displays (show only last 4 chars).\n" +
				"2. Store only the hash (SHA-256) where full PAN is not needed.\n" +
				"3. Restrict database-level access using column-level permissions.\n" +
				"4. Implement a data retention and purge schedule.\n" +
				"5. Log all access to raw PAN values.",
		},
		{
			"Phone Numbers (High — Section 8(1))",
			"1. Tokenise phone numbers using a secure vault (e.g. HashiCorp Vault).\n" +
				"2. Never log raw phone numbers in application or access logs.\n" +
				"3. Enforce MFA for any system accessing contact data in bulk.\n" +
				"4. Conduct quarterly access reviews.",
		},
		{
			"Voter ID (High — Section 8(2))",
			"1. Encrypt at rest with AES-256 and enforce TLS in transit.\n" +
				"2. Implement strict RBAC — government ID data access must be need-based.\n" +
				"3. Build and maintain a complete audit trail for every access.\n" +
				"4. Verify and document consent records for every data principal.",
		},
		{
			"GSTIN (Medium — Section 7(1))",
			"1. Apply field-level encryption for stored GSTIN values.\n" +
				"2. Limit access to finance and compliance teams only.\n" +
				"3. Log all queries that return GSTIN in results.\n" +
				"4. Ensure data is not replicated to non-compliant environments.",
		},
	}

	for _, p := range playbook {
		if pdf.GetY() > 255 {
			addFooter()
			newPage()
			sectionHeader("Remediation Playbook (continued)")
		}
		pdf.SetFillColor(239, 246, 255)
		pdf.SetFont("Arial", "B", 10)
		pdf.SetTextColor(30, 64, 175)
		pdf.CellFormat(190, 7, "  "+p.title, "LTR", 1, "L", true, 0, "")
		pdf.SetFillColor(255, 255, 255)
		pdf.SetFont("Arial", "", 9)
		pdf.SetTextColor(30, 41, 59)
		pdf.MultiCell(190, 5, p.steps, "LBR", "L", true)
		pdf.Ln(3)
	}

	// ── Next Steps ───────────────────────────────────────────────────────────
	pdf.Ln(4)
	sectionHeader("Next Steps & Recommended Schedule")
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(30, 41, 59)
	nextSteps := "1. Remediate all CRITICAL findings within 72 hours.\n" +
		"2. Remediate HIGH findings within 2 weeks.\n" +
		"3. Schedule a re-scan after remediation to verify findings are resolved.\n" +
		"4. Appoint a Data Protection Officer (DPO) as required under the DPDP Act.\n" +
		"5. Register with the Data Protection Board of India once notified.\n" +
		"6. Implement a Privacy Impact Assessment (PIA) process for new data pipelines.\n" +
		"7. Run Aurva scans on a weekly schedule for continuous compliance monitoring."
	pdf.MultiCell(190, 6, nextSteps, "", "L", false)

	addFooter()

	// ── Audit log ────────────────────────────────────────────────────────────
	s.db.Exec(`
		INSERT INTO audit_log (event_type, account_id, event_data)
		VALUES ($1, $2, $3)
	`, "report_generated", accountID, `{"format":"pdf","version":"2.0"}`)

	// ── Stream PDF ───────────────────────────────────────────────────────────
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=aurva-dpdp-compliance-%s-%s.pdf",
			accountID, time.Now().Format("2006-01-02")))

	if err := pdf.Output(w); err != nil {
		log.Printf("Error outputting PDF: %v", err)
	}
	log.Printf("✓ Generated %d-page compliance report for account: %s (%d findings)", pageNum, accountID, rowCount)
}

func (s *HTTPServer) handleScanStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	scanID := strings.TrimPrefix(r.URL.Path, "/api/scan/")
	if scanID == "" {
		http.Error(w, "scan_id required", http.StatusBadRequest)
		return
	}
	var scan struct {
		ID                string         `json:"id"`
		AccountID         string         `json:"account_id"`
		Status            string         `json:"status"`
		ResourcesDiscovered int          `json:"resources_discovered"`
		PIIFound          int            `json:"pii_found"`
		CurrentWorker     string         `json:"current_worker"`
		ErrorMessage      sql.NullString `json:"-"`
		ErrMsg            string         `json:"error_message"`
	}
	err := s.db.QueryRow(`
		SELECT id, account_id, status, 
		       COALESCE(resources_scanned, 0),
		       COALESCE(findings_count, 0),
		       COALESCE(error_message, '')
		FROM scan_jobs WHERE id = $1
	`, scanID).Scan(&scan.ID, &scan.AccountID, &scan.Status,
		&scan.ResourcesDiscovered, &scan.PIIFound, &scan.ErrMsg)
	if err != nil {
		http.Error(w, "Scan not found", http.StatusNotFound)
		return
	}
	scan.CurrentWorker = "S3 + RDS Scanner"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scan)
}

// runScanner spawns the scanner as a goroutine
func (s *HTTPServer) runScanner(scanID, accountID, roleARN string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()
	
	// Update status to running
	s.db.Exec(`UPDATE scan_jobs SET status = 'running' WHERE id = $1`, scanID)
	log.Printf("✓ Scanner started for scan: %s", scanID)
	
	// Load AWS config and assume role
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Failed to load AWS config: %v", err)
		s.db.Exec(`UPDATE scan_jobs SET status = 'failed', error_message = $1 WHERE id = $2`, 
			"Failed to load AWS config", scanID)
		return
	}
	
	// Assume the customer's IAM role
	stsClient := sts.NewFromConfig(cfg)
	provider := stscreds.NewAssumeRoleProvider(stsClient, roleARN, func(o *stscreds.AssumeRoleOptions) {
		o.RoleSessionName = "aurva-scanner"
		o.Duration = 3600 * time.Second // 1 hour max
	})
	
	assumedCfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(provider))
	if err != nil {
		log.Printf("Failed to assume role: %v", err)
		// For demo: continue with default config if role assumption fails
		assumedCfg = cfg
	}
	
	// Initialize PII classifier
	classifier := pii.NewClassifier()
	
	// Create gRPC client connection to self
	conn, err := grpc.Dial("localhost:50055",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Printf("Failed to connect to gRPC: %v", err)
		s.db.Exec(`UPDATE scan_jobs SET status = 'failed', error_message = $1 WHERE id = $2`, 
			"Failed to connect to control plane", scanID)
		return
	}
	defer conn.Close()
	
	client := pb.NewComplianceServiceClient(conn)
	
	// Create scanners
	s3Scanner := workers.NewS3Scanner(assumedCfg, classifier, client, scanID, accountID)
	rdsScanner := workers.NewRDSScanner(assumedCfg, classifier, client, scanID, accountID)
	
	// Run scans in parallel
	var wg sync.WaitGroup
	var scanErrors []error
	var mu sync.Mutex
	
	wg.Add(2)
	
	go func() {
		defer wg.Done()
		if err := s3Scanner.Scan(ctx); err != nil {
			log.Printf("S3 scan error: %v", err)
			mu.Lock()
			scanErrors = append(scanErrors, err)
			mu.Unlock()
		}
	}()
	
	go func() {
		defer wg.Done()
		if err := rdsScanner.Scan(ctx); err != nil {
			log.Printf("RDS scan error: %v", err)
			mu.Lock()
			scanErrors = append(scanErrors, err)
			mu.Unlock()
		}
	}()
	
	wg.Wait()
	
	// Report completion
	success := len(scanErrors) == 0
	var errorMsg string
	if !success {
		errorMsg = "Scan completed with errors"
	}
	
	totalResources := s3Scanner.GetResourceCount() + rdsScanner.GetResourceCount()
	totalFindings := s3Scanner.GetFindingCount() + rdsScanner.GetFindingCount()
	
	_, err = client.CompleteScan(ctx, &pb.CompleteRequest{
		ScanId:         scanID,
		Success:        success,
		ErrorMessage:   errorMsg,
		TotalResources: int32(totalResources),
		TotalFindings:  int32(totalFindings),
	})
	
	if err != nil {
		log.Printf("Failed to report completion: %v", err)
	}
	
	log.Printf("✓ Scan complete: %s (%d resources, %d findings)", scanID, totalResources, totalFindings)
}

// validateRoleAssumption validates that we can assume the provided IAM role
func (s *HTTPServer) validateRoleAssumption(roleARN string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("no AWS credentials available: %w", err)
	}
	
	stsClient := sts.NewFromConfig(cfg)
	creds := stscreds.NewAssumeRoleProvider(stsClient, roleARN, func(o *stscreds.AssumeRoleOptions) { o.Duration = 3600 * time.Second })
	
	// Test the credentials
	_, err = creds.Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to assume role: %w", err)
	}
	
	return nil
}

func getScoreColor(score float64) [3]int {
	if score >= 90 {
		return [3]int{40, 167, 69} // Green
	} else if score >= 70 {
		return [3]int{255, 193, 7} // Yellow
	} else {
		return [3]int{220, 53, 69} // Red
	}
}
