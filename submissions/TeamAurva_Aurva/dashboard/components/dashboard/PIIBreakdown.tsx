"use client";

import { BarChart, Bar, XAxis, YAxis, ResponsiveContainer, Cell, Tooltip } from 'recharts';

interface PIIBreakdownProps {
  breakdown: {
    aadhaar: number;
    pan: number;
    gstin: number;
    phone: number;
    voter_id: number;
    bank_account: number;
  };
}

export function PIIBreakdown({ breakdown }: PIIBreakdownProps) {
  const data = [
    { name: 'AADHAAR', value: breakdown.aadhaar, color: '#6366f1' },
    { name: 'PAN', value: breakdown.pan, color: '#0ea5e9' },
    { name: 'GSTIN', value: breakdown.gstin, color: '#10b981' },
    { name: 'BANK ACCOUNT', value: breakdown.bank_account, color: '#a855f7' },
    { name: 'PHONE', value: breakdown.phone, color: '#f97316' },
    { name: 'VOTER ID', value: breakdown.voter_id, color: '#f43f5e' },
  ].filter(item => item.value > 0);

  return (
    <div className="bg-white border border-slate-200 rounded-none p-8 shadow-[6px_6px_0px_0px_rgba(15,23,42,0.05)]">
      <h3 className="text-[10px] font-black uppercase tracking-[0.2em] mb-8 text-slate-400">Classification Breakdown</h3>
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={data} layout="vertical" margin={{ left: -20, right: 10 }}>
          <XAxis type="number" hide />
          <YAxis 
            type="category" 
            dataKey="name" 
            stroke="#94a3b8" 
            style={{ fontSize: '9px', fontWeight: '900', fontFamily: 'monospace' }} 
            width={70}
            axisLine={false}
            tickLine={false}
          />
          <Tooltip
            contentStyle={{
              backgroundColor: '#0f172a',
              border: 'none',
              borderRadius: '0px',
              color: '#fff',
              fontSize: '10px',
              fontWeight: '900',
              textTransform: 'uppercase',
            }}
            itemStyle={{ color: '#fff' }}
            cursor={{ fill: 'rgba(15,23,42,0.05)' }}
            formatter={(value) => [`${value} instances`, 'Count']}
          />
          <Bar dataKey="value" radius={[0, 0, 0, 0]} barSize={20}>
            {data.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={entry.color} />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
