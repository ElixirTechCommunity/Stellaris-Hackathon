const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:9090";

export interface ScanRequest {
  account_id: string;
  role_arn: string;
  account_nickname?: string;
}

export interface ScanResponse {
  scan_id: string;
  status: string;
  account_id: string;
}

export interface ScanStatus {
  scan_id: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  account_id: string;
  resources_discovered: number;
  pii_found: number;
  current_worker: string;
  started_at: string;
  completed_at?: string;
  error_message?: string;
}

export interface PIIFinding {
  id: string;
  resource_id: string;
  resource_type: string;
  resource_name: string;
  pii_type: string;
  risk_level: 'critical' | 'high' | 'medium' | 'low';
  confidence_score: number;
  dpdp_section: string;
  sample_data: string;
  location_info: any;
  detected_at: string;
}

export interface FindingsResponse {
  account_id: string;
  total_resources: number;
  total_pii_count: number;
  high_risk_count: number;
  critical_count: number;
  compliance_score: number;
  last_scan_at: string;
  findings: PIIFinding[];
  breakdown: {
    aadhaar: number;
    pan: number;
    gstin: number;
    phone: number;
    voter_id: number;
    bank_account: number;
  };
  risk_distribution: {
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
}

export interface ResourceScanEvent {
  timestamp: string;
  resource_type: string;
  resource_name: string;
  status: 'scanning' | 'clean' | 'pii_found';
  pii_type?: string;
}

export const api = {
  async triggerScan(data: ScanRequest): Promise<ScanResponse> {
    const res = await fetch(`${API_BASE_URL}/api/scan`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error('Failed to trigger scan');
    return res.json();
  },

  async getScanStatus(scanId: string): Promise<ScanStatus> {
    const res = await fetch(`${API_BASE_URL}/api/scan/${scanId}`);
    if (!res.ok) throw new Error('Failed to fetch scan status');
    return res.json();
  },

  async getFindings(): Promise<FindingsResponse> {
    const res = await fetch(`${API_BASE_URL}/api/findings?account_id=814023042898`);
    if (!res.ok) throw new Error('Failed to fetch findings');
    return res.json();
  },

  async downloadReport(): Promise<Blob> {
    const res = await fetch(`${API_BASE_URL}/api/report/pdf?account_id=814023042898`);
    if (!res.ok) throw new Error('Failed to download report');
    return res.blob();
  },

  async checkHealth(): Promise<boolean> {
    try {
      const res = await fetch(`${API_BASE_URL}/health`);
      return res.ok;
    } catch {
      return false;
    }
  },
};
