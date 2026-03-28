"use client";

import { motion } from 'framer-motion';
import { useRouter } from 'next/navigation';
import { useFindings } from '@/hooks/useScan';
import { usePDFDownload } from '@/hooks/usePDFDownload';
import { StatCards } from '@/components/dashboard/StatCards';
import { FindingsTable } from '@/components/dashboard/FindingsTable';
import { FindingDrawer } from '@/components/dashboard/FindingDrawer';
import { RiskDonut } from '@/components/dashboard/RiskDonut';
import { PIIBreakdown } from '@/components/dashboard/PIIBreakdown';
import { useStore } from '@/lib/store';
import { Download, Terminal, ScanLine, Shield } from 'lucide-react';

const riskColors: Record<string, string> = {
  critical: '#f43f5e', high: '#f97316', medium: '#6366f1', low: '#10b981',
};
const piiColors: Record<string, string> = {
  AADHAAR: '#6366f1', PAN: '#0ea5e9', GSTIN: '#10b981', BANK_ACCOUNT: '#a855f7', PHONE: '#f97316', VOTER_ID: '#f43f5e',
};

export default function DashboardPage() {
  const router = useRouter();
  const { data: findings, isLoading } = useFindings();
  const { download } = usePDFDownload();
  const { riskFilter, setRiskFilter, piiTypeFilter, setPIITypeFilter } = useStore();

  if (isLoading) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center gap-4 bg-white">
        <div className="w-10 h-10 border-2 border-slate-100 border-t-indigo-600 rounded-none animate-spin" />
        <span className="text-slate-400 text-xs font-bold uppercase tracking-widest">Initialising...</span>
      </div>
    );
  }

  if (!findings || !findings.findings || findings.findings.length === 0) {
    return (
      <div className="min-h-screen flex items-center justify-center p-8">
        <div className="text-center max-w-md bg-white border border-slate-200 p-12 shadow-[8px_8px_0px_0px_rgba(15,23,42,0.05)]">
          <div className="w-16 h-16 bg-slate-50 border border-slate-200 flex items-center justify-center mx-auto mb-6">
            <Shield size={32} className="text-indigo-600" />
          </div>
          <h2 className="text-2xl font-black text-slate-900 mb-3 uppercase tracking-tight">System Offline</h2>
          <p className="text-slate-500 text-sm mb-8 leading-relaxed font-medium">
            Connect your AWS infrastructure via the Console to begin autonomous PII discovery.
          </p>
          <button 
            onClick={() => router.push('/connect')} 
            className="h-12 px-8 bg-indigo-600 text-white font-bold text-xs uppercase tracking-widest transition-all hover:bg-indigo-700 flex items-center gap-3 mx-auto shadow-[4px_4px_0px_0px_rgba(79,70,229,0.3)]"
          >
            <Terminal size={16} /> Open Console
          </button>
        </div>
      </div>
    );
  }

  const toggleRiskFilter = (risk: string) => setRiskFilter(riskFilter.includes(risk) ? riskFilter.filter(r => r !== risk) : [...riskFilter, risk]);
  const togglePIIFilter = (type: string) => setPIITypeFilter(piiTypeFilter.includes(type) ? piiTypeFilter.filter(t => t !== type) : [...piiTypeFilter, type]);

  return (
    <div className="min-h-screen p-10 max-w-[1440px] mx-auto bg-transparent">

      {/* Header */}
      <motion.div initial={{ opacity: 0, y: -16 }} animate={{ opacity: 1, y: 0 }} className="flex items-center justify-between mb-10">
        <div className="flex items-center gap-4">
          <div className="w-11 h-11 bg-slate-900 flex items-center justify-center shadow-[4px_4px_0px_0px_rgba(15,23,42,0.2)]">
            <Shield size={22} color="#fff" />
          </div>
          <div>
            <div className="text-xl font-black text-slate-900 uppercase tracking-tighter">Aurva</div>
            <div className="text-[10px] text-slate-400 font-bold tracking-[0.2em] uppercase">DPDP Compliance Control</div>
          </div>
        </div>

        <div className="flex items-center gap-4">
          <button onClick={download} className="h-11 px-6 bg-white border border-slate-200 text-slate-900 text-[10px] font-bold uppercase tracking-widest hover:border-slate-400 transition-all flex items-center gap-2 shadow-[4px_4px_0px_0px_rgba(15,23,42,0.05)]">
            <Download size={15} /> Export Report
          </button>
          <button onClick={() => router.push('/connect')} className="h-11 px-6 bg-orange-600 text-white text-[10px] font-bold uppercase tracking-widest hover:bg-orange-700 transition-all flex items-center gap-2 shadow-[4px_4px_0px_0px_rgba(249,115,22,0.3)]">
            <Terminal size={15} /> Open Console
          </button>
        </div>
      </motion.div>

      {/* Stat Cards */}
      <StatCards
        complianceScore={findings.compliance_score}
        totalPII={findings.total_pii_count}
        criticalCount={findings.critical_count}
        lastScanAt={findings.last_scan_at}
        totalResources={findings.total_resources}
      />

      {/* Filters */}
      <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 0.2 }}
        className="flex items-center gap-8 mb-6 flex-wrap">
        <div className="flex items-center gap-3">
          <span className="text-[10px] text-slate-400 font-bold uppercase tracking-widest">Filter Risk</span>
          <div className="flex gap-2">
            {['critical', 'high', 'medium', 'low'].map(risk => {
              const active = riskFilter.includes(risk);
              const col = riskColors[risk];
              return (
                <button key={risk} onClick={() => toggleRiskFilter(risk)} 
                  className={`px-3 py-1.5 text-[9px] font-bold uppercase tracking-widest transition-all font-mono border ${
                    active ? 'bg-slate-900 border-slate-900 text-white' : 'bg-white border-slate-200 text-slate-400'
                  }`}
                  style={active ? { backgroundColor: col, borderColor: col } : {}}
                >
                  {risk}
                </button>
              );
            })}
          </div>
        </div>
        <div className="w-px h-6 bg-slate-200" />
        <div className="flex items-center gap-3">
          <span className="text-[10px] text-slate-400 font-bold uppercase tracking-widest">PII Type</span>
          <div className="flex gap-2">
            {['AADHAAR', 'PAN', 'GSTIN', 'BANK_ACCOUNT', 'PHONE', 'VOTER_ID'].map(type => {
              const active = piiTypeFilter.includes(type);
              const col = piiColors[type];
              return (
                <button key={type} onClick={() => togglePIIFilter(type)} 
                  className={`px-3 py-1.5 text-[9px] font-bold uppercase tracking-widest transition-all font-mono border ${
                    active ? 'bg-slate-900 border-slate-900 text-white' : 'bg-white border-slate-200 text-slate-400'
                  }`}
                  style={active ? { backgroundColor: col, borderColor: col } : {}}
                >
                  {type.replace('_', ' ')}
                </button>
              );
            })}
          </div>
        </div>
      </motion.div>

      {/* Main grid */}
      <div className="grid grid-cols-[1fr_320px] gap-8 items-start">
        <FindingsTable findings={findings.findings} />
        <div className="flex flex-col gap-6">
          <RiskDonut distribution={findings.risk_distribution || { critical: 0, high: 0, medium: 0, low: 0 }} />
          <PIIBreakdown breakdown={findings.breakdown} />
        </div>
      </div>

      <FindingDrawer />
    </div>
  );
}
