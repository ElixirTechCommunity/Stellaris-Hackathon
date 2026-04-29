# Aurva - India's First DPDP Act 2023 Compliance Engine

## 👥 Team Name
Team Aurva

## 🧑‍💻 Team Members
| Name | Role | GitHub |
|------|------|--------|
| Aditya Bisht | Frontend / Lead | [@cherry-bisht](https://github.com/cherry-bisht) |
| Abhyuday Pundir | Backend / Security | [@abhyudaypundir](https://github.com/abhyudaypundir) |
| Aabhas Viswas | Cloud Infrastructure | [@aabhasviswas](https://github.com/aabhasviswas) |

## 💡 Problem Statement
The Indian Digital Personal Data Protection (DPDP) Act 2023 mandates strict compliance for handling personal data. Most existing cloud security tools are not optimized for Indian-specific PII (like Aadhaar, PAN, GSTIN) and do not provide audit-ready compliance reports aligned with Indian laws.

**Aurva** solves this by providing an autonomous SaaS bot that:
- Crawls AWS S3 and RDS to detect Indian PII.
- Validates data using Verhoeff and Luhn checksums to minimize false positives.
- Generates DPDP Act-compliant risk scores and audit reports.

## 🛠️ Tech Stack
- **Languages:** Go, TypeScript
- **Frontend:** Next.js 15, Tailwind CSS, Lucide React
- **Backend:** Go (gRPC for scanner communication, HTTP for dashboard)
- **Database:** PostgreSQL (with complex JSONB support for finding aggregation)
- **Cloud:** AWS (SDK v2, STS for cross-account AssumeRole)
- **PII Engine:** Regex + Checksum Validation (Verhoeff for Aadhaar, etc.)

## 🔗 Links
- **Repository:** [https://github.com/cherry-bisht/aurva](https://github.com/cherry-bisht/aurva)
- **Presentation (PPT/PDF):** [Link to Project Documentation](https://github.com/cherry-bisht/aurva/blob/main/README.md)

## 📸 Screenshots
![Dashboard Overview](https://raw.githubusercontent.com/cherry-bisht/aurva/main/README.md)
*(Detailed dashboard view includes Risk Donut charts, PII breakdown, and resource feeds)*

## 🚀 How to Run Locally

### 1. Database Setup
```bash
docker-compose up -d postgres
# Run the schema SQL to initialize tables
psql -h localhost -U postgres -d compliance -f schema.sql
```

### 2. Control Plane (Backend)
```bash
cd control-plane
go run main.go
```
The control plane starts a gRPC server on `:50051` and an HTTP API on `:9090`.

### 3. Cloud Scanner
```bash
cd cloud-scanner
# Requires AWS credentials with STS AssumeRole permissions
go run main.go -role-arn arn:aws:iam::ACCOUNT_ID:role/AurvaReadOnly
```

### 4. Dashboard (Frontend)
```bash
cd dashboard
npm install
npm run dev
```
Open [http://localhost:3000](http://localhost:3000) to view the dashboard.

## 🛡️ Key Features & Implementation Details

- **PII Classification:** Custom-built detection engine for Aadhaar (checksum validated), PAN, GSTIN, Phone numbers, and Voter ID.
- **Cross-Account Scanning:** Securely assumes customer IAM roles via STS, ensuring no long-term credentials are stored.
- **Sampling Logic:** High performance achieved by sampling the first 10KB of S3 files and 100 rows per RDS table.
- **Compliance Score:** Real-time calculation based on resource risk density:
  `((total_resources - critical_resources - high_resources) / total_resources) * 100`
- **Audit-Ready Reports:** Generates structured PDF reports summarizing PII violations and recommended remediation steps.
