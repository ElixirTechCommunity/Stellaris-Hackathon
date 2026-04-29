# Aurva - Quick Start Guide

## Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Docker & Docker Compose (optional)
- AWS Account with appropriate IAM role

## Setup Instructions

### 1. Start PostgreSQL Database

Using Docker Compose:
```bash
cd aurva
docker-compose up -d postgres
```

Or use your own PostgreSQL instance and update connection string in `control-plane/main.go`.

### 2. Initialize Database Schema

```bash
psql -h localhost -U aurva -d aurva -f schema.sql
# Password: aurva_dev_password
```

### 3. Build Go Modules

```bash
go mod tidy
go mod download
```

### 4. Start Control Plane

```bash
cd control-plane
go run main.go
```

You should see:
```
✓ Connected to PostgreSQL
✓ gRPC server listening on :50051
✓ HTTP server listening on :9090
✓ Aurva Control Plane is running
```

### 5. Start Dashboard

```bash
cd dashboard
npm install
npm run dev
```

Dashboard will be available at: http://localhost:3000

### 6. Setup AWS IAM Role

Create an IAM role in your AWS account with this policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListAllMyBuckets",
        "s3:ListBucket",
        "s3:GetObject",
        "rds:DescribeDBInstances",
        "rds:DescribeDBClusters"
      ],
      "Resource": "*"
    }
  ]
}
```

Trust relationship (replace YOUR_ACCOUNT_ID):
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::YOUR_ACCOUNT_ID:root"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

### 7. Run a Scan

1. Open http://localhost:3000
2. Enter your AWS Account ID
3. Enter the IAM Role ARN created above
4. Click "Start Compliance Scan"
5. Wait for redirect to dashboard
6. View findings and download PDF report

### 8. Manual Scanner Execution (Optional)

If you need to run the scanner manually:

```bash
cd cloud-scanner
go run main.go \
  -role-arn arn:aws:iam::123456789012:role/AurvaReadOnly \
  -scan-id $(uuidgen) \
  -control-plane localhost:50051
```

## Testing the PII Classifier

```bash
cd shared/pii
go test -v
```

## API Endpoints

### Control Plane (HTTP :9090)
- `GET /health` - Health check
- `POST /api/scan` - Trigger new scan
- `GET /api/findings?account_id=X` - Get findings
- `GET /api/report/pdf?account_id=X` - Download PDF report

### Control Plane (gRPC :50051)
- `ReportCloudFinding` - Submit PII finding
- `ReportScanProgress` - Update scan progress
- `CompleteScan` - Mark scan complete

## Architecture

```
┌─────────────────────────────────┐
│     Customer Browser            │
│  localhost:3000 (Next.js)       │
└──────────────┬──────────────────┘
               │ HTTP
┌──────────────▼──────────────────┐
│     Control Plane               │
│  gRPC :50051 + HTTP :9090       │
└──────┬───────────────┬──────────┘
       │ gRPC          │ SQL
┌──────▼──────┐  ┌────▼────────┐
│Cloud Scanner│  │  PostgreSQL  │
│AWS SDK      │  │  Port 5432   │
└──────┬──────┘  └─────────────┘
       │ AWS SDK (AssumeRole)
┌──────▼──────────────────────────┐
│     AWS Infrastructure          │
│     S3 Buckets + RDS            │
└─────────────────────────────────┘
```

## Compliance Score Formula

```
compliance_score = ((total_resources - critical_count - high_count) / total_resources) * 100
```

## PII Types Detected

1. **Aadhaar** (12 digits, Verhoeff checksum validated) - CRITICAL
2. **PAN** (10 chars, 4th char validated) - HIGH
3. **GSTIN** (15 chars) - MEDIUM
4. **Phone** (10 digits, +91 prefix) - HIGH
5. **Voter ID** (3 letters + 7 digits) - HIGH

## Troubleshooting

### Database connection failed
- Ensure PostgreSQL is running: `docker-compose ps`
- Check connection string in `control-plane/main.go`
- Verify password: `aurva_dev_password`

### AWS role assumption failed
- Verify IAM role exists and has correct policy
- Check trust relationship allows your account
- Ensure AWS credentials are configured locally

### gRPC connection refused
- Verify control plane is running
- Check firewall allows port 50051
- Scanner and control plane must be network-reachable

### Dashboard not loading findings
- Check control plane HTTP server is running on :9090
- Open browser console for error details
- Verify CORS headers are present

## Production Deployment Notes

- Use environment variables for database credentials
- Store IAM role ARNs in Secrets Manager
- Enable TLS for gRPC communication
- Add authentication to HTTP endpoints
- Set up monitoring and alerting
- Configure backup schedule for PostgreSQL
- Rate limit API endpoints
- Implement job queue for large scans

## Support

For issues or questions, review the main README.md and source code comments.
