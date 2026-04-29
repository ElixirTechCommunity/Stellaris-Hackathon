package workers

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/aurva/compliance-engine/proto"
	"github.com/aurva/compliance-engine/shared/pii"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Scanner struct {
	cfg            aws.Config
	classifier     *pii.Classifier
	client         pb.ComplianceServiceClient
	scanID         string
	accountID      string
	resourceCount  int32
	findingCount   int32
}

func NewS3Scanner(cfg aws.Config, classifier *pii.Classifier, client pb.ComplianceServiceClient, scanID, accountID string) *S3Scanner {
	return &S3Scanner{
		cfg:        cfg,
		classifier: classifier,
		client:     client,
		scanID:     scanID,
		accountID:  accountID,
	}
}

func (s *S3Scanner) Scan(ctx context.Context) error {
	s3Client := s3.NewFromConfig(s.cfg)

	// List all buckets
	listResult, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("failed to list S3 buckets: %w", err)
	}

	log.Printf("Found %d S3 buckets to scan", len(listResult.Buckets))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50)

	for _, bucket := range listResult.Buckets {
		wg.Add(1)
		go func(bucketName string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := s.scanBucket(ctx, s3Client, bucketName); err != nil {
				log.Printf("Error scanning bucket %s: %v", bucketName, err)
			}
		}(*bucket.Name)
	}

	wg.Wait()
	return nil
}

func (s *S3Scanner) scanBucket(ctx context.Context, client *s3.Client, bucketName string) error {
	regionResp, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return err
	}

	region := string(regionResp.LocationConstraint)
	if region == "" {
		region = "us-east-1"
	}

	client = s3.NewFromConfig(s.cfg, func(o *s3.Options) { o.Region = region })

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int32(1000),
	})

	var objectCount int32
	var objWg sync.WaitGroup
	objSem := make(chan struct{}, 20)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			if !s.shouldScanObject(*obj.Key) {
				continue
			}

			// Atomic check to enforce cap properly across goroutines
			if atomic.LoadInt32(&objectCount) >= 500 {
				break
			}

			objWg.Add(1)
			go func(key string, size int64) {
				defer objWg.Done()
				objSem <- struct{}{}
				defer func() { <-objSem }()

				if atomic.AddInt32(&objectCount, 1) > 500 {
					return
				}

				if err := s.scanObject(ctx, client, bucketName, key, region, size); err != nil {
					log.Printf("Error scanning %s/%s: %v", bucketName, key, err)
				}
			}(*obj.Key, *obj.Size)
		}

		if atomic.LoadInt32(&objectCount) >= 500 {
			break
		}
	}

	objWg.Wait()
	return nil
}

func (s *S3Scanner) tryS3Select(ctx context.Context, client *s3.Client, bucket, key string) (string, error) {
	result, err := client.SelectObjectContent(ctx, &s3.SelectObjectContentInput{
		Bucket:         aws.String(bucket),
		Key:            aws.String(key),
		ExpressionType: s3types.ExpressionTypeSql,
		Expression:     aws.String("SELECT * FROM S3Object LIMIT 1000"),
		InputSerialization: &s3types.InputSerialization{
			CSV: &s3types.CSVInput{FileHeaderInfo: s3types.FileHeaderInfoUse},
		},
		OutputSerialization: &s3types.OutputSerialization{
			CSV: &s3types.CSVOutput{},
		},
	})
	if err != nil {
		return "", err
	}
	defer result.GetStream().Close()

	var sb strings.Builder
	for event := range result.GetStream().Events() {
		if v, ok := event.(*s3types.SelectObjectContentEventStreamMemberRecords); ok {
			sb.Write(v.Value.Payload)
		}
	}
	return sb.String(), nil
}

func (s *S3Scanner) scanObject(ctx context.Context, client *s3.Client, bucket, key, region string, size int64) error {
	atomic.AddInt32(&s.resourceCount, 1)
	
	// Ensure per-object timeout doesn't exceed 60s
	objCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var content string
	var err error

	if strings.HasSuffix(strings.ToLower(key), ".csv") {
		content, err = s.tryS3Select(objCtx, client, bucket, key)
		if err != nil {
			content, err = s.downloadObject(objCtx, client, bucket, key, size)
		}
	} else {
		content, err = s.downloadObject(objCtx, client, bucket, key, size)
	}

	if err != nil {
		return err
	}

	detections := s.classifier.Scan(content)
	
	if len(detections) > 0 {
		resourceID := fmt.Sprintf("s3://%s/%s", bucket, key)
		s.reportFindings(objCtx, detections, resourceID, "s3_object", fmt.Sprintf("%s/%s", bucket, key), region, size)
	}

	return nil
}

func (s *S3Scanner) downloadObject(ctx context.Context, client *s3.Client, bucket, key string, size int64) (string, error) {
	maxBytes := int64(524288)
	if size < maxBytes {
		maxBytes = size
	}

	getResp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Range:  aws.String(fmt.Sprintf("bytes=0-%d", maxBytes-1)),
	})
	if err != nil {
		return "", err
	}
	defer getResp.Body.Close()

	buf := make([]byte, 65536)
	var fullContent strings.Builder
	bytesRead := 0
	for int64(bytesRead) < maxBytes {
		n, err := getResp.Body.Read(buf)
		if n > 0 {
			fullContent.Write(buf[:n])
			bytesRead += n
		}
		if err != nil {
			break
		}
	}

	return fullContent.String(), nil
}

func (s *S3Scanner) reportFindings(ctx context.Context, detections []pii.Detection, resourceID, resourceType, resourceName, region string, size int64) {
	for _, detection := range detections {
		atomic.AddInt32(&s.findingCount, 1)
		retries := 3
		for i := 0; i < retries; i++ {
			_, err := s.client.ReportCloudFinding(ctx, &pb.FindingRequest{
				ScanId:       s.scanID,
				AccountId:    s.accountID, // Pass account ID directly to skip DB lookup
				ResourceId:   resourceID,
				ResourceType: resourceType,
				ResourceName: resourceName,
				Region:       region,
				SizeBytes:    size,
				PiiDetection: &pb.PIIDetection{
					PiiType:         string(detection.Type),
					ConfidenceScore: detection.ConfidenceScore,
					SampleData:      detection.Value,
					RiskLevel:       detection.RiskLevel,
					DpdpSection:     detection.DPDPSection,
					Location: &pb.LocationInfo{
						FieldPath: resourceID,
					},
				},
			})
			if err == nil {
				break
			}
			time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
		}
	}
}

func (s *S3Scanner) shouldScanObject(key string) bool {
	key = strings.ToLower(key)
	skipExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".mp4", ".mp3",
		".zip", ".tar", ".pdf", ".parquet", ".orc",
		".avro", ".pb", ".bin", ".exe", ".dll", ".so",
		".wasm", ".class", ".pyc",
	}
	for _, ext := range skipExtensions {
		if strings.HasSuffix(key, ext) {
			return false
		}
	}
	
	scanExtensions := []string{
		".csv", ".json", ".jsonl", ".txt", ".log",
		".xml", ".yaml", ".yml", ".tsv", ".sql", ".ndjson",
	}
	for _, ext := range scanExtensions {
		if strings.HasSuffix(key, ext) {
			return true
		}
	}
	return false
}

func (s *S3Scanner) GetResourceCount() int {
	return int(atomic.LoadInt32(&s.resourceCount))
}

func (s *S3Scanner) GetFindingCount() int {
	return int(atomic.LoadInt32(&s.findingCount))
}
