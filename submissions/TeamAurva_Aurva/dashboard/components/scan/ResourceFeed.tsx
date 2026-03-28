"use client";

import { motion, AnimatePresence } from 'framer-motion';
import { useState, useEffect } from 'react';

interface ResourceEvent {
  id: string;
  timestamp: string;
  resource_type: string;
  resource_name: string;
  status: 'scanning' | 'clean' | 'pii_found';
  pii_type?: string;
}

interface ResourceFeedProps {
  scanId: string;
  resourceCount: number;
}

export function ResourceFeed({ scanId, resourceCount }: ResourceFeedProps) {
  const [events, setEvents] = useState<ResourceEvent[]>([]);

  // Simulate resource discovery for demo
  useEffect(() => {
    const resources = [
      { type: 'S3', name: 's3://prod-user-data/customers.csv', pii: 'AADHAAR' },
      { type: 'S3', name: 's3://prod-logs/app-2024-03.log', pii: null },
      { type: 'S3', name: 's3://backups/db-dump-march.json', pii: 'PAN' },
      { type: 'RDS', name: 'prod-mysql/users_table', pii: 'PHONE' },
      { type: 'RDS', name: 'prod-mysql/orders_table', pii: null },
      { type: 'S3', name: 's3://analytics/user-events.csv', pii: 'GSTIN' },
      { type: 'S3', name: 's3://static-assets/images/', pii: null },
      { type: 'RDS', name: 'staging-postgres/kyc_data', pii: 'AADHAAR' },
    ];

    let idx = 0;
    const interval = setInterval(() => {
      if (idx >= Math.min(resourceCount, resources.length)) {
        clearInterval(interval);
        return;
      }

      const resource = resources[idx % resources.length];
      const event: ResourceEvent = {
        id: `${scanId}-${idx}`,
        timestamp: new Date().toISOString(),
        resource_type: resource.type,
        resource_name: resource.name,
        status: 'scanning',
        pii_type: resource.pii || undefined,
      };

      setEvents((prev) => [event, ...prev]);

      // Update status after 1 second
      setTimeout(() => {
        setEvents((prev) =>
          prev.map((e) =>
            e.id === event.id
              ? { ...e, status: resource.pii ? 'pii_found' : 'clean' }
              : e
          )
        );
      }, 1000);

      idx++;
    }, 1500);

    return () => clearInterval(interval);
  }, [scanId, resourceCount]);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'scanning':
        return 'text-slate-400';
      case 'clean':
        return 'text-pink-400';
      case 'pii_found':
        return 'text-pink-600';
      default:
        return 'text-slate-800';
    }
  };

  const getStatusText = (event: ResourceEvent) => {
    switch (event.status) {
      case 'scanning':
        return 'SCANNING...';
      case 'clean':
        return '✓ CLEAN';
      case 'pii_found':
        return `⚠ ${event.pii_type} FOUND`;
      default:
        return '';
    }
  };

  return (
    <div className="flex-1 bg-white border border-pink-200 rounded-none p-6 font-mono text-sm overflow-y-auto max-h-[calc(100vh-16rem)] shadow-[2px_2px_0px_0px_rgba(244,114,182,0.2)]">
      <AnimatePresence initial={false}>
        {events.map((event) => (
          <motion.div
            key={event.id}
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0 }}
            className={`mb-2 ${getStatusColor(event.status)}`}
          >
            <span className="text-slate-400">
              [{new Date(event.timestamp).toLocaleTimeString()}]
            </span>{' '}
            <span className="text-pink-400">{event.resource_type}</span>{' '}
            <span className="text-slate-800">{event.resource_name}</span>
            {' — '}
            <span className={getStatusColor(event.status)}>
              {getStatusText(event)}
            </span>
          </motion.div>
        ))}
      </AnimatePresence>
      
      {events.length === 0 && (
        <div className="text-slate-400 text-center py-12">
          Initializing scan workers...
        </div>
      )}
    </div>
  );
}
