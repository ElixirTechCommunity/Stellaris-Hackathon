interface PIIBadgeProps { type: string; size?: 'sm' | 'md'; }

const piiClasses: Record<string, string> = {
  aadhaar: 'bg-indigo-50 text-indigo-600 border-indigo-200',
  pan:     'bg-sky-50 text-sky-600 border-sky-200',
  gstin:   'bg-emerald-50 text-emerald-600 border-emerald-200',
  phone:   'bg-orange-50 text-orange-600 border-orange-200',
};

export function PIIBadge({ type, size = 'md' }: PIIBadgeProps) {
  const normalizedType = type?.toLowerCase() || 'pan';
  const classes = piiClasses[normalizedType] || piiClasses.pan;
  const padding = size === 'sm' ? 'px-1.5 py-0.5' : 'px-2.5 py-1';
  const fontSize = size === 'sm' ? 'text-[9px]' : 'text-[10px]';

  return (
    <span className={`inline-flex items-center border rounded-none font-bold tracking-widest uppercase font-mono ${classes} ${padding} ${fontSize}`}>
      {type}
    </span>
  );
}
