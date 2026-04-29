"use client";

import { Shield, Lock, FileCheck } from 'lucide-react';

export function TrustSignals() {
  const signals = [
    {
      icon: Lock,
      text: "Read-only IAM access",
    },
    {
      icon: Shield,
      text: "Zero data leaves your VPC",
    },
    {
      icon: FileCheck,
      text: "DPDP Act 2023 mapped",
    },
  ];

  return (
    <div className="flex flex-col gap-4 mt-12">
      {signals.map((signal, idx) => (
        <div key={idx} className="flex items-center gap-3 text-slate-400">
          <signal.icon className="w-5 h-5 text-pink-400" />
          <span className="text-sm">{signal.text}</span>
        </div>
      ))}
    </div>
  );
}
