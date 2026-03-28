"use client";

import { motion, AnimatePresence } from 'framer-motion';
import { useStore } from '@/lib/store';
import { RiskBadge } from '@/components/ui/RiskBadge';
import { PIIBadge } from '@/components/ui/PIIBadge';
import { getDPDPDescription, maskPII } from '@/lib/format';
import { X, FileText, Terminal } from 'lucide-react';

export function FindingDrawer() {
  const { selectedFinding, setSelectedFinding, isDrawerOpen } = useStore();

  if (!selectedFinding) return null;

  const remediation = [
    "Encrypt data at rest using AWS KMS",
    "Implement field-level encryption for PII",
    "Enable audit logging for all access",
    "Review and update IAM policies",
    "Consider data minimization strategies",
  ];

  return (
    <AnimatePresence>
      {isDrawerOpen && (
        <>
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={() => setSelectedFinding(null)}
            className="fixed inset-0 bg-slate-900/60 z-40 backdrop-blur-sm"
          />
          
          <motion.div
            initial={{ x: '100%' }}
            animate={{ x: 0 }}
            exit={{ x: '100%' }}
            transition={{ type: 'spring', damping: 25, stiffness: 200 }}
            className="fixed right-0 top-0 h-full w-[520px] bg-white border-l border-slate-200 z-50 overflow-y-auto shadow-[-8px_0px_30px_rgba(15,23,42,0.1)]"
          >
            <div className="p-8">
              <div className="flex items-start justify-between mb-10 pb-6 border-b border-slate-100">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-slate-900 flex items-center justify-center">
                    <Terminal size={20} className="text-white" />
                  </div>
                  <div>
                    <h2 className="text-xl font-black text-slate-900 uppercase tracking-tighter">Record Intelligence</h2>
                    <p className="text-[10px] text-slate-400 font-bold uppercase tracking-widest">Finding ID: {selectedFinding.resource_id.slice(-8)}</p>
                  </div>
                </div>
                <button
                  onClick={() => setSelectedFinding(null)}
                  className="p-2 hover:bg-slate-100 transition-colors"
                >
                  <X className="w-6 h-6 text-slate-400" />
                </button>
              </div>

              <div className="space-y-10">
                <div>
                  <label className="text-slate-400 text-[10px] font-black uppercase tracking-[0.2em] block mb-3">Resource Location</label>
                  <div className="text-slate-900 font-mono text-xs font-bold bg-slate-50 p-5 border border-slate-200 break-all shadow-[4px_4px_0px_0px_rgba(15,23,42,0.05)]">
                    {selectedFinding.resource_name}
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-8">
                  <div>
                    <label className="text-slate-400 text-[10px] font-black uppercase tracking-[0.2em] block mb-3">PII Classification</label>
                    <PIIBadge type={selectedFinding.pii_type} />
                  </div>
                  <div>
                    <label className="text-slate-400 text-[10px] font-black uppercase tracking-[0.2em] block mb-3">Risk Assessment</label>
                    <RiskBadge risk={selectedFinding.risk_level} />
                  </div>
                </div>

                <div>
                  <label className="text-slate-400 text-[10px] font-black uppercase tracking-[0.2em] block mb-2">Detection Confidence</label>
                  <div className="text-indigo-600 font-mono text-4xl font-black tracking-tighter">
                    {(selectedFinding.confidence_score * 100).toFixed(1)}%
                  </div>
                </div>

                <div>
                  <label className="text-slate-400 text-[10px] font-black uppercase tracking-[0.2em] block mb-3">Masked Payload Sample</label>
                  <div className="text-orange-600 font-mono text-sm font-bold bg-orange-50/50 p-5 border border-orange-200 shadow-[4px_4px_0px_0px_rgba(249,115,22,0.1)]">
                    {maskPII(selectedFinding.pii_type, selectedFinding.sample_data)}
                  </div>
                </div>

                <div>
                  <label className="text-slate-400 text-[10px] font-black uppercase tracking-[0.2em] block mb-3">Regulatory Context (DPDP Act)</label>
                  <div className="bg-slate-900 p-6 border border-slate-800 shadow-[4px_4px_0px_0px_rgba(15,23,42,0.3)]">
                    <div className="text-white font-mono font-black text-sm mb-3 uppercase tracking-wider">
                      {selectedFinding.dpdp_section}
                    </div>
                    <div className="text-slate-400 text-xs font-medium leading-relaxed uppercase tracking-wide">
                      {getDPDPDescription(selectedFinding.dpdp_section)}
                    </div>
                  </div>
                </div>

                <div>
                  <label className="text-slate-400 text-[10px] font-black uppercase tracking-[0.2em] block mb-4">Autonomous Remediation Plan</label>
                  <ol className="space-y-4">
                    {remediation.map((step, idx) => (
                      <li key={idx} className="flex gap-4 text-xs text-slate-600 font-bold uppercase tracking-tight">
                        <span className="text-indigo-500 font-mono font-black">{idx + 1}.</span>
                        <span>{step}</span>
                      </li>
                    ))}
                  </ol>
                </div>

                <button className="w-full h-14 bg-indigo-600 hover:bg-indigo-700 text-white font-black uppercase tracking-[0.2em] text-[10px] transition-all flex items-center justify-center gap-3 shadow-[6px_6px_0px_0px_rgba(79,70,229,0.3)] mt-8">
                  <FileText className="w-4 h-4" />
                  Commit to Audit Report
                </button>
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}
