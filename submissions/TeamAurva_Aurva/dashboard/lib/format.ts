import { formatDistanceToNow, format } from 'date-fns';

export function formatRelativeTime(date: string | Date): string {
  if (!date) return "just now";
  const d = new Date(date);
  if (isNaN(d.getTime())) return "just now";
  return formatDistanceToNow(d, { addSuffix: true });
}

export function formatTimestamp(date: string | Date): string {
  return format(new Date(date), 'MMM dd, yyyy HH:mm:ss');
}

export function formatNumber(num: number): string {
  return new Intl.NumberFormat('en-US').format(num);
}

export function formatPercentage(num: number): string {
  return `${num.toFixed(1)}%`;
}

export function maskPII(type: string, value: string): string {
  switch (type.toLowerCase()) {
    case 'aadhaar':
      return value.slice(-4).padStart(12, 'X');
    case 'pan':
      return value.slice(0, 2) + 'XXX' + value.slice(-2);
    case 'phone':
      return 'XXXXX' + value.slice(-5);
    case 'gstin':
      return value.slice(0, 2) + 'XXX' + value.slice(-3);
    default:
      return value.slice(0, 3) + 'XXX' + value.slice(-3);
  }
}

export function getDPDPDescription(section: string): string {
  const descriptions: Record<string, string> = {
    'Section 5': 'Notice and consent for processing personal data',
    'Section 6': 'Grounds for processing without consent',
    'Section 7': 'Purpose limitation - data processing boundaries',
    'Section 8': 'Data retention and accuracy obligations',
    'Section 9': 'Data subject rights - access and correction',
    'Section 10': 'Security safeguards for personal data',
    'Section 11': 'Additional obligations for significant data fiduciaries',
    'Section 12': 'Cross-border data transfer provisions',
  };
  return descriptions[section] || 'DPDP Act compliance requirement';
}

export function getComplianceColor(score: number): string {
  if (score < 60) return 'text-rose-600';
  if (score < 80) return 'text-orange-600';
  return 'text-emerald-600';
}

export function getRiskColor(risk: string): string {
  const colors: Record<string, string> = {
    critical: 'bg-rose-50 text-rose-600 border-rose-200',
    high:     'bg-orange-50 text-orange-600 border-orange-200',
    medium:   'bg-indigo-50 text-indigo-600 border-indigo-200',
    low:      'bg-emerald-50 text-emerald-600 border-emerald-200',
  };
  return colors[risk.toLowerCase()] || colors.low;
}
