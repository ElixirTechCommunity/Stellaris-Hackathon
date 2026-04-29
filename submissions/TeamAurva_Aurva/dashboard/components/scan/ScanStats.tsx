"use client";

import { motion } from 'framer-motion';
import { StatCard } from '@/components/ui/StatCard';

interface ScanStatsProps {
  resourcesDiscovered: number;
  piiFound: number;
  currentWorker: string;
}

export function ScanStats({ resourcesDiscovered, piiFound, currentWorker }: ScanStatsProps) {
  return (
    <div className="w-80 space-y-4">
      <motion.div
        initial={{ opacity: 0, x: 20 }}
        animate={{ opacity: 1, x: 0 }}
      >
        <StatCard
          title="Resources Discovered"
          value={resourcesDiscovered}
          color="text-pink-600"
        />
      </motion.div>

      <motion.div
        initial={{ opacity: 0, x: 20 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ delay: 0.1 }}
      >
        <StatCard
          title="PII Findings"
          value={piiFound}
          color="text-pink-500"
        />
      </motion.div>

      <motion.div
        initial={{ opacity: 0, x: 20 }}
        animate={{ opacity: 1, x: 0 }}
        transition={{ delay: 0.2 }}
        className="bg-white border border-pink-200 rounded-none p-6 shadow-[2px_2px_0px_0px_rgba(244,114,182,0.2)]"
      >
        <div className="text-slate-400 text-sm font-medium mb-2">Current Worker</div>
        <div className="text-pink-600 font-mono text-lg">{currentWorker}</div>
      </motion.div>
    </div>
  );
}
