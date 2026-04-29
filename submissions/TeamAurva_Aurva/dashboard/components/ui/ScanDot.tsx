"use client";

import { motion } from 'framer-motion';

export function ScanDot() {
  return (
    <div className="flex items-center gap-2">
      <motion.div
        className="w-2.5 h-2.5 rounded-none bg-pink-500"
        animate={{
          opacity: [1, 0.3, 1],
          scale: [1, 1.2, 1],
        }}
        transition={{
          duration: 2,
          repeat: Infinity,
          ease: "easeInOut",
        }}
      />
      <span className="text-pink-500 font-medium text-sm uppercase tracking-wide">Scanning</span>
    </div>
  );
}
