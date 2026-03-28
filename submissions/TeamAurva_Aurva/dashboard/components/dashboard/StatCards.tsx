"use client";

import { StatCard } from '@/components/ui/StatCard';
import { formatRelativeTime } from '@/lib/format';

interface StatCardsProps {
  complianceScore: number;
  totalPII: number;
  criticalCount: number;
  lastScanAt: string;
  totalResources: number;
}

export function StatCards({ complianceScore, totalPII, criticalCount, lastScanAt, totalResources }: StatCardsProps) {
  const score = Math.max(0, complianceScore);
  const scoreColor = score >= 80 ? 'text-emerald-600' : score >= 50 ? 'text-indigo-600' : 'text-rose-600';
  const scoreAccent = score >= 80 ? '#10b981' : score >= 50 ? '#6366f1' : '#f43f5e';

  return (
    <div className="grid grid-cols-4 gap-6 mb-8">
      <StatCard title="Compliance Score" value={score} suffix="%" color={scoreColor} accent={scoreAccent} animate />
      <StatCard title="PII Records Found" value={totalPII} subtitle={`across ${totalResources} resources`} color="text-orange-600" accent="#f97316" animate />
      <StatCard title="Critical Risk" value={criticalCount} subtitle="require immediate action" color="text-rose-600" accent="#ef4444" animate />
      <StatCard title="Last Scan" value={formatRelativeTime(lastScanAt)} color="text-slate-600" accent="#475569" animate={false} />
    </div>
  );
}
