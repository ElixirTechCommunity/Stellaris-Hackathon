package workers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"

	pb "github.com/aurva/compliance-engine/proto"
	"github.com/aurva/compliance-engine/shared/pii"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	_ "github.com/lib/pq"
)

type RDSScanner struct {
	cfg           aws.Config
	classifier    *pii.Classifier
	client        pb.ComplianceServiceClient
	scanID        string
	accountID     string
	resourceCount int32
	findingCount  int32
}

func NewRDSScanner(cfg aws.Config, classifier *pii.Classifier, client pb.ComplianceServiceClient, scanID, accountID string) *RDSScanner {
	return &RDSScanner{
		cfg:        cfg,
		classifier: classifier,
		client:     client,
		scanID:     scanID,
		accountID:  accountID,
	}
}

func (r *RDSScanner) Scan(ctx context.Context) error {
	rdsClient := rds.NewFromConfig(r.cfg)

	// List all RDS instances
	describeResult, err := rdsClient.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
	if err != nil {
		return fmt.Errorf("failed to list RDS instances: %w", err)
	}

	log.Printf("Found %d RDS instances to scan", len(describeResult.DBInstances))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10) // Limit concurrent scans

	for _, instance := range describeResult.DBInstances {
		// Only scan PostgreSQL and MySQL for now
		engine := strings.ToLower(*instance.Engine)
		if engine != "postgres" && engine != "mysql" {
			log.Printf("Skipping unsupported engine: %s", engine)
			continue
		}

		wg.Add(1)
		go func(inst types.DBInstance) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := r.scanInstance(ctx, inst); err != nil {
				log.Printf("Error scanning instance %s: %v", *inst.DBInstanceIdentifier, err)
			}
		}(instance)
	}

	wg.Wait()
	return nil
}

func (r *RDSScanner) scanInstance(ctx context.Context, instance types.DBInstance) error {
	instanceID := *instance.DBInstanceIdentifier
	if instance.Endpoint == nil {
		log.Printf("Skipping stopped instance: %s", instanceID)
		return nil
	}
	endpoint := *instance.Endpoint.Address
	port := *instance.Endpoint.Port
	engine := strings.ToLower(*instance.Engine)
	region := *instance.AvailabilityZone // Simplified - should extract region

	log.Printf("Scanning RDS instance: %s (%s)", instanceID, engine)

	// Note: In production, credentials should come from Secrets Manager or parameter store
	// For demo purposes, we'll skip actual DB connection and simulate findings
	
	// Simulate scanning databases
	r.simulateDatabaseScan(ctx, instanceID, endpoint, int(port), engine, region)
	
	return nil
}

// simulateDatabaseScan simulates scanning RDS tables
// In production, this would connect to the database and sample rows
func (r *RDSScanner) simulateDatabaseScan(ctx context.Context, instanceID, endpoint string, port int, engine, region string) {
	// Simulated databases and tables
	simulatedData := map[string][]string{
		"production": {"users", "orders", "payments"},
		"staging":    {"customers", "transactions"},
	}

	for dbName, tables := range simulatedData {
		for _, tableName := range tables {
			atomic.AddInt32(&r.resourceCount, 1)
			
			// Simulate sampling data that might contain PII
			// In production, would execute: SELECT * FROM table LIMIT 100
			simulatedRows := r.generateSimulatedData(tableName)
			
			// Scan for PII
			for rowIdx, rowData := range simulatedRows {
				detections := r.classifier.Scan(rowData)
				
				if len(detections) > 0 {
					log.Printf("Found %d PII detections in %s.%s.%s", len(detections), instanceID, dbName, tableName)
				}
				
				// Report each detection
				for _, detection := range detections {
					atomic.AddInt32(&r.findingCount, 1)
					
					resourceID := fmt.Sprintf("rds://%s/%s/%s", instanceID, dbName, tableName)
					
					_, err := r.client.ReportCloudFinding(ctx, &pb.FindingRequest{
						ScanId:       r.scanID,
						ResourceId:   resourceID,
						ResourceType: "rds_table",
						ResourceName: fmt.Sprintf("%s.%s.%s", instanceID, dbName, tableName),
						Region:       region,
						SizeBytes:    0, // Unknown for tables
						PiiDetection: &pb.PIIDetection{
							PiiType:         string(detection.Type),
							ConfidenceScore: detection.ConfidenceScore,
							SampleData:      detection.Value,
							RiskLevel:       detection.RiskLevel,
							DpdpSection:     detection.DPDPSection,
							Location: &pb.LocationInfo{
								LineNumber: int32(rowIdx),
								ColumnName: r.guessColumnName(detection.Type),
								FieldPath:  fmt.Sprintf("%s.%s", dbName, tableName),
							},
						},
					})
					
					if err != nil {
						log.Printf("Failed to report finding: %v", err)
					}
				}
			}
		}
	}
}

// generateSimulatedData creates sample data for demonstration
// In production, this would be actual database rows
func (r *RDSScanner) generateSimulatedData(tableName string) []string {
	data := []string{}
	switch tableName {
	case "users", "customers":
		data = append(data, "user_id: 12345, name: Rajesh Kumar, phone: 9876543210, aadhaar: 499118665519")
		data = append(data, "user_id: 12346, name: Priya Sharma, phone: +91 8765432109, pan: BQRPH1234C")
		data = append(data, "user_id: 12347, name: Amit Patel, phone: 9123456700, aadhaar: 234123412340")
	case "orders":
		data = append(data, "order_id: 9001, customer_phone: 9123456789, pan: ABCPE1234F")
		data = append(data, "order_id: 9002, customer_phone: 9234567890, aadhaar: 499118665519")
	case "payments":
		data = append(data, "payment_id: 5001, pan_card: PQRPH5678G, gstin: 29ABCPE1234F1Z5")
		data = append(data, "payment_id: 5002, pan_card: DKLPH9012B, gstin: 27DKLPH9012B1Z3")
	case "transactions":
		data = append(data, "txn_id: 7001, voter_id: ABC1234567, mobile: 9988776655")
		data = append(data, "txn_id: 7002, voter_id: XYZ9876543, pan: CKTPF5678G")
	}
	return data
}

func (r *RDSScanner) guessColumnName(piiType pii.PIIType) string {
	columnNames := map[pii.PIIType]string{
		pii.PIITypeAadhaar:  "aadhaar_number",
		pii.PIITypePAN:      "pan_card",
		pii.PIITypeGSTIN:    "gstin",
		pii.PIITypePhone:    "phone_number",
		pii.PIITypeVoterID:  "voter_id",
		pii.PIITypeAccount:  "account_number",
	}
	
	if col, ok := columnNames[piiType]; ok {
		return col
	}
	return "unknown_column"
}

func (r *RDSScanner) GetResourceCount() int {
	return int(atomic.LoadInt32(&r.resourceCount))
}

func (r *RDSScanner) GetFindingCount() int {
	return int(atomic.LoadInt32(&r.findingCount))
}
