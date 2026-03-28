import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        bg: '#07090F',
        surface: '#0E1118',
        surface2: '#141821',
        surface3: '#1A1F2E',
        border: 'rgba(255,255,255,0.07)',
        amber: '#F59E0B',
        sky: '#38BDF8',
        critical: '#F87171',
        high: '#FB923C',
        medium: '#FBBF24',
        low: '#34D399',
        muted: '#4B5563',
        secondary: '#9CA3AF',
        accent: '#6366F1',
      },
      fontFamily: {
        sans: ['DM Sans', 'system-ui', 'sans-serif'],
        mono: ['DM Mono', 'monospace'],
      },
      borderRadius: {
        'xl': '12px',
        '2xl': '16px',
      },
      animation: {
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'fade-in': 'fadeIn 0.4s ease-out',
        'slide-up': 'slideUp 0.4s ease-out',
      },
      keyframes: {
        fadeIn: { '0%': { opacity: '0' }, '100%': { opacity: '1' } },
        slideUp: { '0%': { opacity: '0', transform: 'translateY(12px)' }, '100%': { opacity: '1', transform: 'translateY(0)' } },
      },
      boxShadow: {
        'glow-amber': '0 0 24px rgba(245,158,11,0.08)',
        'glow-critical': '0 0 24px rgba(248,113,113,0.08)',
        'glow-sky': '0 0 24px rgba(56,189,248,0.08)',
        'card': '0 1px 3px rgba(0,0,0,0.4), 0 1px 2px rgba(0,0,0,0.3)',
      }
    },
  },
  plugins: [],
};

export default config;
