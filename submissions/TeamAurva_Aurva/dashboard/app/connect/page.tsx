"use client";

import { motion } from 'framer-motion';
import { OnboardingForm } from '@/components/connect/OnboardingForm';
import { TrustSignals } from '@/components/connect/TrustSignals';
import { Terminal, Shield, ArrowLeft } from 'lucide-react';
import { useRouter } from 'next/navigation';

export default function ConnectPage() {
  const router = useRouter();

  return (
    <div className="min-h-screen flex bg-slate-50">
      {/* Left Panel - Brand & Context */}
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        className="w-[450px] bg-slate-900 p-12 flex flex-col justify-between text-white"
      >
        <div>
          <button 
            onClick={() => router.push('/dashboard')}
            className="flex items-center gap-2 text-slate-400 hover:text-white transition-colors mb-12 text-xs font-bold uppercase tracking-widest"
          >
            <ArrowLeft size={14} /> Skip to Demo Dashboard
          </button>

          <div className="flex items-center gap-3 mb-8">
            <div className="w-10 h-10 bg-indigo-600 flex items-center justify-center">
              <Terminal size={20} />
            </div>
            <h1 className="text-2xl font-black uppercase tracking-tighter">Console</h1>
          </div>

          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="text-4xl font-black leading-[1.1] mb-6 uppercase tracking-tight"
          >
            Connect <span className="text-orange-500">Infrastructure</span>.
          </motion.h2>

          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="text-slate-400 text-lg leading-relaxed mb-12"
          >
            Aurva requires read-only IAM access to scan your AWS resources for PII. 
            No data ever leaves your VPC.
          </motion.p>

          <TrustSignals />
        </div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.4 }}
          className="flex items-center gap-2 text-slate-500 text-[10px] font-bold uppercase tracking-[0.2em]"
        >
          <Shield size={12} /> Encrypted & SOC2 Compliant
        </motion.div>
      </motion.div>

      {/* Right Panel - Form */}
      <div className="flex-1 p-12 flex items-center justify-center">
        <div className="w-full max-w-lg bg-white border border-slate-200 p-12 shadow-[12px_12px_0px_0px_rgba(15,23,42,0.05)]">
          <div className="mb-10">
            <h3 className="text-xl font-black text-slate-900 uppercase tracking-tight mb-2">Configure Connection</h3>
            <p className="text-slate-500 text-sm font-medium">Enter your AWS cross-account role details to begin.</p>
          </div>
          <OnboardingForm />
        </div>
      </div>
    </div>
  );
}
