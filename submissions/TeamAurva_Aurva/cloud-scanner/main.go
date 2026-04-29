package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aurva/compliance-engine/cloud-scanner/workers"
	pb "github.com/aurva/compliance-engine/proto"
	"github.com/aurva/compliance-engine/shared/pii"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Parse flags
	roleARN := flag.String("role-arn", "", "AWS IAM role ARN to assume")
	accountID := flag.String("account-id", "", "AWS Account ID being scanned")
	scanID := flag.String("scan-id", "", "Scan job ID")
	controlPlane := flag.String("control-plane", "localhost:50055", "Control plane gRPC address")
	flag.Parse()

	if *roleARN == "" || *scanID == "" || *accountID == "" {
		log.Fatal("Usage: cloud-scanner -role-arn <ARN> -account-id <ID> -scan-id <ID>")
	}

	log.Printf("🚀 Aurva Cloud Scanner starting...")
	log.Printf("   Role ARN: %s", *roleARN)
	log.Printf("   Account ID: %s", *accountID)
	log.Printf("   Scan ID: %s", *scanID)
	log.Printf("   Control Plane: %s", *controlPlane)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Load AWS config and assume role
	log.Println("Assuming IAM role...")
	awsCfg, err := assumeRole(ctx, *roleARN)
	if err != nil {
		log.Fatalf("Failed to assume role: %v", err)
	}
	log.Println("✓ IAM role assumed successfully")

	// Connect to control plane gRPC
	conn, err := grpc.Dial(*controlPlane, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to control plane: %v", err)
	}
	defer conn.Close()
	client := pb.NewComplianceServiceClient(conn)
	log.Println("✓ Connected to control plane")

	// Initialize PII classifier
	classifier := pii.NewClassifier()
	log.Println("✓ PII classifier initialized")

	// Create scanners
	s3Scanner := workers.NewS3Scanner(awsCfg, classifier, client, *scanID, *accountID)
	rdsScanner := workers.NewRDSScanner(awsCfg, classifier, client, *scanID, *accountID)

	// Run scans
	log.Println("Starting AWS infrastructure scan...")
	
	s3Results := make(chan error, 1)
	rdsResults := make(chan error, 1)

	// S3 scan
	go func() {
		s3Results <- s3Scanner.Scan(ctx)
	}()

	// RDS scan
	go func() {
		rdsResults <- rdsScanner.Scan(ctx)
	}()

	// Wait for both scanners
	var scanErrors []error
	for i := 0; i < 2; i++ {
		select {
		case err := <-s3Results:
			if err != nil {
				log.Printf("S3 scan error: %v", err)
				scanErrors = append(scanErrors, err)
			} else {
				log.Println("✓ S3 scan completed")
			}
		case err := <-rdsResults:
			if err != nil {
				log.Printf("RDS scan error: %v", err)
				scanErrors = append(scanErrors, err)
			} else {
				log.Println("✓ RDS scan completed")
			}
		case <-ctx.Done():
			log.Println("Scan interrupted")
			return
		}
	}

	// Report completion
	success := len(scanErrors) == 0
	var errorMsg string
	if !success {
		errorMsg = "Scan completed with errors"
	}

	totalResources := s3Scanner.GetResourceCount() + rdsScanner.GetResourceCount()
	totalFindings := s3Scanner.GetFindingCount() + rdsScanner.GetFindingCount()

	_, err = client.CompleteScan(ctx, &pb.CompleteRequest{
		ScanId:         *scanID,
		Success:        success,
		ErrorMessage:   errorMsg,
		TotalResources: int32(totalResources),
		TotalFindings:  int32(totalFindings),
	})

	if err != nil {
		log.Printf("Failed to report completion: %v", err)
	}

	log.Printf("✓ Scan complete: %d resources, %d findings", totalResources, totalFindings)
}

// assumeRole uses STS to assume the provided IAM role and returns configured AWS config
func assumeRole(ctx context.Context, roleARN string) (aws.Config, error) {
	// Load default config (uses ambient credentials or instance profile)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, err
	}

	// Create STS client
	stsClient := sts.NewFromConfig(cfg)

	// Assume the role
	provider := stscreds.NewAssumeRoleProvider(stsClient, roleARN, func(o *stscreds.AssumeRoleOptions) {
		o.RoleSessionName = "aurva-scanner"
		o.Duration = 3600 * time.Second // 1 hour
	})

	// Create new config with assumed role credentials
	assumedCfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(provider),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return assumedCfg, nil
}
