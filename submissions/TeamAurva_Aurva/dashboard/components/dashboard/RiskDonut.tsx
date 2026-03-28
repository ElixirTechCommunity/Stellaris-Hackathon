"use client";

import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';

interface RiskDonutProps {
  distribution: {
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
}

export function RiskDonut({ distribution = { critical: 0, high: 0, medium: 0, low: 0 } }: RiskDonutProps) {
  const data = [
    { name: 'Critical', value: distribution.critical, color: '#f43f5e' },
    { name: 'High', value: distribution.high, color: '#f97316' },
    { name: 'Medium', value: distribution.medium, color: '#6366f1' },
    { name: 'Low', value: distribution.low, color: '#10b981' },
  ].filter(item => item.value > 0);

  const total = Object.values(distribution).reduce((a, b) => a + b, 0);

  return (
    <div className="bg-white border border-slate-200 rounded-none p-8 shadow-[6px_6px_0px_0px_rgba(15,23,42,0.05)]">
      <h3 className="text-[10px] font-black uppercase tracking-[0.2em] mb-8 text-slate-400">Risk Matrix</h3>
      <div className="relative">
        <ResponsiveContainer width="100%" height={280}>
          <PieChart>
            <Pie
              data={data}
              cx="50%"
              cy="50%"
              innerRadius={70}
              outerRadius={95}
              paddingAngle={4}
              dataKey="value"
            >
              {data.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={entry.color} stroke="none" />
              ))}
            </Pie>
            <Tooltip
              contentStyle={{
                backgroundColor: '#0f172a',
                border: 'none',
                borderRadius: '0px',
                color: '#fff',
                fontSize: '10px',
                fontWeight: '900',
                textTransform: 'uppercase',
                letterSpacing: '0.1em'
              }}
              itemStyle={{ color: '#fff' }}
              formatter={(value) => [`${value} Records`, 'Count']}
            />
          </PieChart>
        </ResponsiveContainer>
        <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
          <div className="text-4xl font-black font-mono text-slate-900 tracking-tighter">{total}</div>
          <div className="text-[8px] font-black text-slate-400 uppercase tracking-widest">Total Findings</div>
        </div>
      </div>
      <div className="space-y-3 mt-4">
        {data.map((item) => (
          <div key={item.name} className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="w-2.5 h-2.5 rounded-none" style={{ backgroundColor: item.color }} />
              <span className="text-[10px] font-bold text-slate-500 uppercase tracking-wider">{item.name}</span>
            </div>
            <span className="text-[10px] font-black font-mono text-slate-900">{item.value}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
