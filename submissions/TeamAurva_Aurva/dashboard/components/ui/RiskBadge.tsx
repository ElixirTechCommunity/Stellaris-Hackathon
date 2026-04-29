interface RiskBadgeProps { risk: string; size?: 'sm' | 'md' | 'lg'; }

const riskClasses: Record<string, string> = {
  critical: 'bg-rose-50 text-rose-600 border-rose-200',
  high:     'bg-orange-50 text-orange-600 border-orange-200',
  medium:   'bg-indigo-50 text-indigo-600 border-indigo-200',
  low:      'bg-emerald-50 text-emerald-600 border-emerald-200',
};

export function RiskBadge({ risk, size = 'md' }: RiskBadgeProps) {
  const normalizedRisk = risk?.toLowerCase() || 'low';
  const classes = riskClasses[normalizedRisk] || riskClasses.low;
  const padding = size === 'sm' ? 'px-1.5 py-0.5' : 'px-2.5 py-1';
  const fontSize = size === 'sm' ? 'text-[9px]' : 'text-[10px]';

  return (
    <span className={`inline-flex items-center border rounded-none font-bold tracking-widest uppercase font-mono ${classes} ${padding} ${fontSize}`}>
      {risk}
    </span>
  );
}
