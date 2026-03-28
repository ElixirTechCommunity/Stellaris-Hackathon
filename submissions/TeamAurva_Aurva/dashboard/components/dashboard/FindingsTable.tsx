"use client";

import { useMemo } from 'react';
import { motion } from 'framer-motion';
import { useStore } from '@/lib/store';
import { PIIFinding } from '@/lib/api';
import { RiskBadge } from '@/components/ui/RiskBadge';
import { PIIBadge } from '@/components/ui/PIIBadge';
import { ArrowUpDown, FileText } from 'lucide-react';

interface FindingsTableProps { findings: PIIFinding[]; }

export function FindingsTable({ findings }: FindingsTableProps) {
  const { setSelectedFinding, riskFilter, piiTypeFilter, sortColumn, sortDirection, setSorting } = useStore();

  const filteredAndSortedFindings = useMemo(() => {
    let filtered = findings;
    if (riskFilter.length > 0) filtered = filtered.filter(f => riskFilter.includes(f.risk_level));
    if (piiTypeFilter.length > 0) filtered = filtered.filter(f => piiTypeFilter.includes(f.pii_type));
    return [...filtered].sort((a, b) => {
      const aVal = a[sortColumn as keyof PIIFinding];
      const bVal = b[sortColumn as keyof PIIFinding];
      const multiplier = sortDirection === 'asc' ? 1 : -1;
      return aVal > bVal ? multiplier : -multiplier;
    });
  }, [findings, riskFilter, piiTypeFilter, sortColumn, sortDirection]);

  const handleSort = (column: string) => {
    setSorting(column, sortColumn === column && sortDirection === 'asc' ? 'desc' : 'asc');
  };

  const cols = [
    { key: 'resource_name', label: 'Resource' },
    { key: 'resource_type', label: 'Type' },
    { key: 'pii_type', label: 'PII Type' },
    { key: 'risk_level', label: 'Risk' },
    { key: 'dpdp_section', label: 'DPDP Section' },
    { key: 'confidence_score', label: 'Confidence' },
  ];

  return (
    <div className="bg-white border border-slate-200 shadow-[6px_6px_0px_0px_rgba(15,23,42,0.05)] overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr className="bg-slate-50 border-b border-slate-200">
              {cols.map(col => (
                <th
                  key={col.key}
                  onClick={() => handleSort(col.key)}
                  className="px-6 py-4 text-left text-[9px] font-black uppercase tracking-[0.15em] text-slate-400 cursor-pointer select-none whitespace-nowrap transition-colors hover:text-indigo-600"
                >
                  <div className="flex items-center gap-2">
                    {col.label}
                    <ArrowUpDown size={10} className="opacity-50" />
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {filteredAndSortedFindings.map((finding, idx) => (
              <motion.tr
                key={finding.resource_id + finding.pii_type + idx}
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: idx * 0.01 }}
                onClick={() => setSelectedFinding(finding)}
                className="border-b border-slate-100 cursor-pointer transition-all hover:bg-slate-50"
              >
                <td className="px-6 py-4">
                  <div className="flex items-center gap-3">
                    <FileText size={14} className="text-slate-300 flex-shrink-0" />
                    <span className="font-mono text-xs font-bold text-slate-900 max-w-[280px] overflow-hidden text-ellipsis whitespace-nowrap uppercase tracking-tight">
                      {finding.resource_name}
                    </span>
                  </div>
                </td>
                <td className="px-6 py-4">
                  <span className="text-[10px] text-slate-500 font-mono font-bold tracking-wider uppercase">
                    {finding.resource_type}
                  </span>
                </td>
                <td className="px-6 py-4">
                  <PIIBadge type={finding.pii_type} size="sm" />
                </td>
                <td className="px-6 py-4">
                  <RiskBadge risk={finding.risk_level} size="sm" />
                </td>
                <td className="px-6 py-4 text-[11px] text-slate-500 font-mono font-bold">
                  {finding.dpdp_section}
                </td>
                <td className="px-6 py-4">
                  <div className="flex items-center gap-4">
                    <div className="w-20 h-2 bg-slate-100 overflow-hidden">
                      <div className={`h-full transition-all duration-500 ${
                        finding.risk_level === 'critical' ? 'bg-rose-500' : 
                        finding.risk_level === 'high' ? 'bg-orange-500' : 
                        finding.risk_level === 'medium' ? 'bg-indigo-500' : 'bg-emerald-500'
                      }`}
                        style={{ width: `${(isNaN(finding.confidence_score) ? 0 : finding.confidence_score) * 100}%` }}
                      />
                    </div>
                    <span className="text-[10px] text-slate-900 font-black font-mono">
                      {Math.round((isNaN(finding.confidence_score) ? 0 : finding.confidence_score) * 100)}%
                    </span>
                  </div>
                </td>
              </motion.tr>
            ))}
          </tbody>
        </table>
      </div>
      {filteredAndSortedFindings.length === 0 && (
        <div className="text-center py-16 bg-slate-50/50">
          <p className="text-[10px] text-slate-400 font-black uppercase tracking-[0.2em]">Zero Records Match Context</p>
        </div>
      )}
    </div>
  );
}
