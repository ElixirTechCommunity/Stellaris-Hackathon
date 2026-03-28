"use client";

import { motion } from 'framer-motion';
import { useEffect, useState } from 'react';

interface StatCardProps {
  title: string;
  value: number | string;
  subtitle?: string;
  color?: string;
  prefix?: string;
  suffix?: string;
  animate?: boolean;
  icon?: React.ReactNode;
  accent?: string;
}

export function StatCard({ 
  title, value, subtitle, color = 'text-slate-900', 
  prefix = '', suffix = '', animate = true, accent = 'rgba(71, 85, 105, 0.1)'
}: StatCardProps) {
  const [displayValue, setDisplayValue] = useState(0);
  const numericValue = typeof value === 'number' ? Math.max(0, value) : parseFloat(value as string) || 0;

  useEffect(() => {
    if (!animate || typeof value !== 'number') { setDisplayValue(numericValue); return; }
    let start = 0;
    const duration = 1200;
    const increment = numericValue / (duration / 16);
    const timer = setInterval(() => {
      start += increment;
      if (start >= numericValue) { setDisplayValue(numericValue); clearInterval(timer); }
      else setDisplayValue(Math.floor(start));
    }, 16);
    return () => clearInterval(timer);
  }, [numericValue, animate, value]);

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, ease: [0.23, 1, 0.32, 1] }}
      className="relative overflow-hidden rounded-none p-6 border border-slate-200 bg-white shadow-[4px_4px_0px_0px_rgba(15,23,42,0.05)]"
    >
      <div
        className="absolute top-0 left-0 right-0 h-1"
        style={{ background: accent }}
      />
      <div className="text-[10px] font-bold uppercase tracking-[0.2em] mb-3 text-slate-400">
        {title}
      </div>
      <div className={`text-4xl font-black font-mono tracking-tighter ${color} mb-1`}>
        {prefix}{typeof value === 'number' ? displayValue.toLocaleString() : value}{suffix}
      </div>
      {subtitle && (
        <div className="text-xs mt-2 text-slate-500 font-medium">{subtitle}</div>
      )}
    </motion.div>
  );
}
