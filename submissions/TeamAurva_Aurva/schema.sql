-- Aurva DPDP Act 2023 Compliance Engine Database Schema

-- 1. Cloud Accounts Table
-- Stores customer AWS account information and IAM role details
CREATE TABLE IF NOT EXISTS cloud_accounts (
    id SERIAL PRIMARY KEY,
    account_id VARCHAR(12) NOT NULL UNIQUE,
    role_arn TEXT NOT NULL,
    account_name VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    last_scan_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cloud_accounts_account_id ON cloud_accounts(account_id);
CREATE INDEX idx_cloud_accounts_status ON cloud_accounts(status);

-- 2. Cloud Resources Table
-- Tracks all AWS resources scanned (S3 buckets, RDS instances)
CREATE TABLE IF NOT EXISTS cloud_resources (
    id TEXT PRIMARY KEY,  -- Format: s3://bucket/path or rds://instance/database/table
    account_id VARCHAR(12) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,  -- 's3_object', 'rds_table'
    resource_name TEXT NOT NULL,
    region VARCHAR(50),
    size_bytes BIGINT,
    last_scanned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES cloud_accounts(account_id) ON DELETE CASCADE
);

CREATE INDEX idx_cloud_resources_account ON cloud_resources(account_id);
CREATE INDEX idx_cloud_resources_type ON cloud_resources(resource_type);
CREATE INDEX idx_cloud_resources_scanned ON cloud_resources(last_scanned_at);

-- 3. PII Findings Table
-- Core table storing all detected PII instances
-- UNIQUE constraint ensures idempotent scanning
CREATE TABLE IF NOT EXISTS pii_findings (
    id SERIAL PRIMARY KEY,
    resource_id TEXT NOT NULL,
    pii_type VARCHAR(50) NOT NULL,  -- 'aadhaar', 'pan', 'gstin', 'phone', 'voter_id'
    confidence_score FLOAT NOT NULL CHECK (confidence_score >= 0 AND confidence_score <= 1),
    sample_data TEXT,  -- Masked sample for verification
    location_info JSONB,  -- {line_number, column_name, offset, etc}
    risk_level VARCHAR(20) NOT NULL,  -- 'critical', 'high', 'medium', 'low'
    dpdp_section TEXT NOT NULL,  -- DPDP Act 2023 section reference
    detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (resource_id) REFERENCES cloud_resources(id) ON DELETE CASCADE,
    UNIQUE(resource_id, pii_type, sample_data)  -- Ensures idempotency while allowing multiple distinct samples
);

CREATE INDEX idx_pii_findings_resource ON pii_findings(resource_id);
CREATE INDEX idx_pii_findings_type ON pii_findings(pii_type);
CREATE INDEX idx_pii_findings_risk ON pii_findings(risk_level);
CREATE INDEX idx_pii_findings_detected ON pii_findings(detected_at);

-- 4. Violations Table
-- Maps findings to specific DPDP Act violations
CREATE TABLE IF NOT EXISTS violations (
    id SERIAL PRIMARY KEY,
    finding_id INTEGER NOT NULL,
    violation_type VARCHAR(100) NOT NULL,
    dpdp_section TEXT NOT NULL,
    severity VARCHAR(20) NOT NULL,  -- 'critical', 'high', 'medium', 'low'
    description TEXT NOT NULL,
    remediation_steps TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (finding_id) REFERENCES pii_findings(id) ON DELETE CASCADE,
    UNIQUE(finding_id, violation_type)
);

CREATE INDEX idx_violations_finding ON violations(finding_id);
CREATE INDEX idx_violations_severity ON violations(severity);

-- 5. Audit Log Table
-- Complete audit trail for compliance purposes
CREATE TABLE IF NOT EXISTS audit_log (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,  -- 'scan_started', 'finding_detected', 'report_generated'
    account_id VARCHAR(12),
    resource_id TEXT,
    user_id TEXT,
    event_data JSONB,
    ip_address INET,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_log_event_type ON audit_log(event_type);
CREATE INDEX idx_audit_log_account ON audit_log(account_id);
CREATE INDEX idx_audit_log_created ON audit_log(created_at);

-- 6. Scan Jobs Table (for tracking scan progress)
CREATE TABLE IF NOT EXISTS scan_jobs (
    id TEXT PRIMARY KEY,  -- UUID
    account_id VARCHAR(12) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',  -- 'pending', 'running', 'completed', 'failed'
    progress INTEGER DEFAULT 0,  -- 0-100
    resources_scanned INTEGER DEFAULT 0,
    findings_count INTEGER DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES cloud_accounts(account_id) ON DELETE CASCADE
);

CREATE INDEX idx_scan_jobs_account ON scan_jobs(account_id);
CREATE INDEX idx_scan_jobs_status ON scan_jobs(status);
CREATE INDEX idx_scan_jobs_created ON scan_jobs(created_at);

-- Functions for automatic timestamp updates
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_cloud_accounts_updated_at BEFORE UPDATE ON cloud_accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Views for compliance reporting
CREATE OR REPLACE VIEW compliance_summary AS
SELECT 
    ca.account_id,
    ca.account_name,
    COUNT(DISTINCT cr.id) as total_resources,
    COUNT(pf.id) as total_findings,
    COUNT(CASE WHEN pf.risk_level = 'critical' THEN 1 END) as critical_findings,
    COUNT(CASE WHEN pf.risk_level = 'high' THEN 1 END) as high_findings,
    COUNT(CASE WHEN pf.risk_level = 'medium' THEN 1 END) as medium_findings,
    COUNT(CASE WHEN pf.risk_level = 'low' THEN 1 END) as low_findings,
    ROUND(
        CAST(((COUNT(DISTINCT cr.id)::FLOAT - 
          COUNT(CASE WHEN pf.risk_level IN ('critical', 'high') THEN 1 END)::FLOAT) / 
         NULLIF(COUNT(DISTINCT cr.id)::FLOAT, 0)) * 100 AS NUMERIC), 
        2
    ) as compliance_score
FROM cloud_accounts ca
LEFT JOIN cloud_resources cr ON ca.account_id = cr.account_id
LEFT JOIN pii_findings pf ON cr.id = pf.resource_id
GROUP BY ca.account_id, ca.account_name;

-- View for detailed findings report
CREATE OR REPLACE VIEW findings_report AS
SELECT 
    cr.id as resource_id,
    cr.resource_type,
    cr.resource_name,
    cr.region,
    pf.pii_type,
    pf.confidence_score,
    pf.risk_level,
    pf.dpdp_section,
    pf.detected_at,
    v.violation_type,
    v.severity,
    v.description as violation_description
FROM cloud_resources cr
JOIN pii_findings pf ON cr.id = pf.resource_id
LEFT JOIN violations v ON pf.id = v.finding_id
ORDER BY pf.risk_level DESC, pf.detected_at DESC;
