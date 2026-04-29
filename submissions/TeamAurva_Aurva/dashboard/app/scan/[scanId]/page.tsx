"use client";

import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { useRouter, useParams } from 'next/navigation';
import { useScan } from '@/hooks/useScan';
import { ScanDot } from '@/components/ui/ScanDot';
import { ResourceFeed } from '@/components/scan/ResourceFeed';
import { ScanStats } from '@/components/scan/ScanStats';
import { AlertCircle } from 'lucide-react';

export default function ScanPage() {
  const params = useParams();
  const router = useRouter();
  const scanId = params.scanId as string;
  const [elapsedTime, setElapsedTime] = useState(0);

  const { data: scan, isLoading, isError } = useScan(scanId);

  useEffect(() => {
    const interval = setInterval(() => {
      setElapsedTime((prev) => prev + 1);
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    if (scan?.status === 'completed') {
      setTimeout(() => router.push('/dashboard'), 3000);
    }
  }, [scan?.status, router]);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-slate-400">Loading scan...</div>
      </div>
    );
  }

  if (isError || !scan) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="bg-white border border-pink-200 rounded-none p-6 max-w-md shadow-[2px_2px_0px_0px_rgba(244,114,182,0.2)]">
          <div className="flex items-center gap-3 mb-4">
            <AlertCircle className="w-6 h-6 text-pink-600" />
            <h2 className="text-lg font-semibold text-pink-600">Scan Error</h2>
          </div>
          <p className="text-slate-400">{scan?.error_message || 'Failed to load scan status'}</p>
        </div>
      </div>
    );
  }

  const formatTime = (s: number) => `${Math.floor(s/60).toString().padStart(2,'0')}:${(s%60).toString().padStart(2,'0')}`;

  return (
    <div className="min-h-screen p-8">
      <motion.div initial={{ opacity: 0, y: -20 }} animate={{ opacity: 1, y: 0 }} className="mb-8">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold mb-1 text-pink-600">{scan.account_id}</h1>
            <p className="text-slate-400 text-sm">Scan ID: {scanId}</p>
          </div>
          <div className="flex items-center gap-8">
            <div className="text-right">
              <div className="text-slate-400 text-sm mb-1">Elapsed Time</div>
              <div className="text-2xl font-mono font-bold text-pink-600">{formatTime(elapsedTime)}</div>
            </div>
            {scan.status === 'running' ? <ScanDot /> : scan.status === 'completed' ? (
              <div className="flex items-center gap-2">
                <div className="w-2.5 h-2.5 rounded-none bg-pink-400" />
                <span className="text-pink-400 font-medium text-sm uppercase tracking-wide">Complete</span>
              </div>
            ) : null}
          </div>
        </div>
      </motion.div>

      <div className="flex gap-6">
        <ResourceFeed scanId={scanId} resourceCount={scan.resources_discovered} />
        <ScanStats resourcesDiscovered={scan.resources_discovered} piiFound={scan.pii_found} currentWorker={scan.current_worker} />
      </div>

      {scan.status === 'completed' && (
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} className="mt-8 bg-white border border-pink-200 rounded-none p-6 text-center shadow-[2px_2px_0px_0px_rgba(244,114,182,0.2)]">
          <p className="text-pink-600 text-lg font-medium">Scan complete. Redirecting to dashboard...</p>
        </motion.div>
      )}
    </div>
  );
}
