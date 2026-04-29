"use client";

import { motion } from 'framer-motion';
import { useRouter } from 'next/navigation';
import { 
  Shield, 
  Terminal, 
  Zap, 
  Search, 
  Lock, 
  ArrowRight, 
  ChevronDown,
  Globe,
  MessageSquare,
  Database,
  Activity,
  CheckCircle2
} from 'lucide-react';

export default function LandingPage() {
  const router = useRouter();

  return (
    <div className="min-h-screen bg-white">
      {/* Navigation */}
      <nav className="fixed top-0 w-full z-50 bg-white/80 backdrop-blur-md border-b border-slate-100">
        <div className="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
          <div className="flex items-center gap-10">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 bg-slate-900 flex items-center justify-center">
                <Shield size={18} className="text-white" />
              </div>
              <span className="text-xl font-black uppercase tracking-tighter text-slate-900">Aurva</span>
            </div>
            
            <div className="hidden md:flex items-center gap-8 text-[10px] font-bold uppercase tracking-widest text-slate-500">
              <button className="hover:text-indigo-600 transition-colors flex items-center gap-1">Product <ChevronDown size={12} /></button>
              <button className="hover:text-indigo-600 transition-colors">Solutions</button>
              <button className="hover:text-indigo-600 transition-colors">Developers</button>
              <button className="hover:text-indigo-600 transition-colors">Company</button>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <button onClick={() => router.push('/connect')} className="text-[10px] font-bold uppercase tracking-widest text-slate-900 px-6 py-2 hover:bg-slate-50 transition-colors">Login</button>
            <button onClick={() => router.push('/connect')} className="bg-slate-900 text-white text-[10px] font-bold uppercase tracking-widest px-6 py-3 shadow-[4px_4px_0px_0px_rgba(15,23,42,0.2)] hover:translate-x-[1px] hover:translate-y-[1px] hover:shadow-none transition-all">Book a Demo</button>
          </div>
        </div>
      </nav>

      <main className="pt-20">
        {/* Hero Section */}
        <section className="relative py-24 px-6 overflow-hidden">
          <div className="max-w-7xl mx-auto grid lg:grid-cols-2 gap-16 items-center">
            <motion.div 
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.5 }}
            >
              <div className="inline-flex items-center gap-2 px-3 py-1 bg-indigo-50 border border-indigo-100 text-indigo-600 text-[10px] font-bold uppercase tracking-widest mb-8">
                <Zap size={12} /> VPC-Native Security Engine
              </div>
              <h1 className="text-6xl font-black text-slate-900 leading-[1.05] uppercase tracking-tight mb-8">
                Secure Your <span className="text-indigo-600">Cloud Data</span>. Comply with Confidence.
              </h1>
              <p className="text-lg text-slate-500 font-medium leading-relaxed mb-10 max-w-xl">
                Aurva is the VPC-native security engine built for modern data privacy. 
                Get real-time observability, control data flows, and automate compliance—
                without clunky agents slowing down your infrastructure.
              </p>
              
              <div className="flex flex-wrap gap-4 mb-12">
                <button 
                  onClick={() => router.push('/connect')}
                  className="h-14 px-8 bg-indigo-600 text-white font-black text-xs uppercase tracking-[0.2em] shadow-[6px_6px_0px_0px_rgba(79,70,229,0.3)] hover:translate-x-[1px] hover:translate-y-[1px] hover:shadow-none transition-all flex items-center gap-3"
                >
                  Start Free Trial <ArrowRight size={16} />
                </button>
                <button className="h-14 px-8 bg-white border border-slate-200 text-slate-900 font-black text-xs uppercase tracking-[0.2em] shadow-[6px_6px_0px_0px_rgba(15,23,42,0.05)] hover:translate-x-[1px] hover:translate-y-[1px] hover:shadow-none transition-all flex items-center gap-3">
                  View the Docs
                </button>
              </div>

              <div className="bg-slate-900 p-4 inline-flex items-center gap-4 border-l-4 border-indigo-500 shadow-xl">
                <Terminal size={16} className="text-indigo-400" />
                <code className="text-indigo-100 font-mono text-sm tracking-tight">
                  helm install aurva aurva/security-engine
                </code>
              </div>
            </motion.div>

            <motion.div 
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: 0.2 }}
              className="relative"
            >
              <div className="bg-white border border-slate-200 shadow-[20px_20px_0px_0px_rgba(15,23,42,0.03)] p-8">
                <div className="flex items-center justify-between mb-8 border-b border-slate-100 pb-4">
                  <div className="flex items-center gap-3">
                    <div className="w-3 h-3 bg-rose-400" />
                    <div className="w-3 h-3 bg-amber-400" />
                    <div className="w-3 h-3 bg-emerald-400" />
                  </div>
                  <div className="text-[10px] font-black text-slate-300 uppercase tracking-widest">Aurva Control Plane</div>
                </div>
                <div className="space-y-6">
                  <div className="h-4 bg-slate-50 w-3/4" />
                  <div className="h-4 bg-slate-50 w-1/2" />
                  <div className="grid grid-cols-3 gap-4">
                    <div className="h-20 bg-indigo-50 border border-indigo-100" />
                    <div className="h-20 bg-orange-50 border border-orange-100" />
                    <div className="h-20 bg-emerald-50 border border-emerald-100" />
                  </div>
                  <div className="h-32 bg-slate-900 shadow-inner" />
                </div>
              </div>
            </motion.div>
          </div>
        </section>

        {/* Social Proof */}
        <section className="py-20 border-y border-slate-100 bg-slate-50/50">
          <div className="max-w-7xl mx-auto px-6">
            <p className="text-center text-[10px] font-black uppercase tracking-[0.3em] text-slate-400 mb-12">
              Trusted by forward-thinking security and engineering teams at
            </p>
            <div className="flex flex-wrap justify-center items-center gap-16 opacity-50 grayscale contrast-125">
              <span className="text-2xl font-black text-slate-900 tracking-tighter">PHOENIX.AI</span>
              <span className="text-2xl font-black text-slate-900 tracking-tighter">DATAFLOW</span>
              <span className="text-2xl font-black text-slate-900 tracking-tighter">CYBERUNIT</span>
              <span className="text-2xl font-black text-slate-900 tracking-tighter">NEURALINK</span>
            </div>
          </div>
        </section>

        {/* Why Aurva */}
        <section className="py-32 px-6">
          <div className="max-w-3xl mx-auto text-center">
            <h2 className="text-4xl font-black text-slate-900 uppercase tracking-tight mb-8">
              Why <span className="text-orange-600">Aurva</span>?
            </h2>
            <p className="text-xl text-slate-500 font-medium leading-relaxed">
              Data privacy isn't just a perimeter problem anymore. With strict regulations 
              like the DPDP Act and complex cloud architectures, traditional security tools 
              leave blind spots. 
              <br /><br />
              Aurva sits natively inside your VPC. We give you deep visibility into exactly 
              who—or what—is accessing your sensitive data, catching breaches and compliance 
              violations at runtime. <span className="text-slate-900 font-bold underline decoration-indigo-500 decoration-4">Zero friction for your DevOps team.</span>
            </p>
          </div>
        </section>

        {/* Core Features */}
        <section className="py-32 bg-slate-900 text-white px-6 relative overflow-hidden">
          <div className="max-w-7xl mx-auto grid md:grid-cols-3 gap-12 relative z-10">
            <div className="p-8 border border-slate-800 bg-slate-800/50 hover:border-indigo-500 transition-colors">
              <div className="w-12 h-12 bg-indigo-600 flex items-center justify-center mb-8">
                <Zap className="text-white" />
              </div>
              <h3 className="text-xl font-black uppercase tracking-tight mb-4">VPC-Native & Agentless</h3>
              <p className="text-slate-400 font-medium leading-relaxed">
                Deploys directly in your cloud environment. Your sensitive data never 
                leaves your infrastructure. Period.
              </p>
            </div>

            <div className="p-8 border border-slate-800 bg-slate-800/50 hover:border-orange-500 transition-colors">
              <div className="w-12 h-12 bg-orange-600 flex items-center justify-center mb-8">
                <Search className="text-white" />
              </div>
              <h3 className="text-xl font-black uppercase tracking-tight mb-4">Deep Data Observability</h3>
              <p className="text-slate-400 font-medium leading-relaxed">
                Map data flows and trace every database query back to the exact user, 
                service, or AI agent that triggered it.
              </p>
            </div>

            <div className="p-8 border border-slate-800 bg-slate-800/50 hover:border-emerald-500 transition-colors">
              <div className="w-12 h-12 bg-emerald-600 flex items-center justify-center mb-8">
                <Lock className="text-white" />
              </div>
              <h3 className="text-xl font-black uppercase tracking-tight mb-4">Automated Compliance</h3>
              <p className="text-slate-400 font-medium leading-relaxed">
                Turn complex privacy laws into actionable engineering policies. Generate 
                audit-ready reports with one click.
              </p>
            </div>
          </div>
        </section>

        {/* How It Works */}
        <section className="py-32 px-6 bg-white">
          <div className="max-w-7xl mx-auto">
            <div className="text-center mb-24">
              <h2 className="text-4xl font-black text-slate-900 uppercase tracking-tight mb-4">From deployment to compliance in minutes.</h2>
              <div className="w-20 h-1.5 bg-indigo-600 mx-auto" />
            </div>

            <div className="grid md:grid-cols-3 gap-8">
              {[
                { 
                  step: "01", 
                  title: "Deploy", 
                  desc: "Drop the Aurva engine into your cloud environment via Helm or Terraform.",
                  icon: Terminal
                },
                { 
                  step: "02", 
                  title: "Discover", 
                  desc: "Aurva automatically maps your databases, identities, and real-time data flows.",
                  icon: Activity
                },
                { 
                  step: "03", 
                  title: "Secure", 
                  desc: "Set privacy guardrails, monitor anomalies, and export compliance logs with one click.",
                  icon: Shield
                }
              ].map((item, idx) => (
                <div key={idx} className="relative p-10 bg-slate-50 border border-slate-100 group hover:bg-white hover:border-indigo-200 transition-all shadow-[8px_8px_0px_0px_rgba(15,23,42,0.03)]">
                  <span className="absolute top-6 right-6 text-6xl font-black text-slate-200 group-hover:text-indigo-100 transition-colors">{item.step}</span>
                  <item.icon className="w-10 h-10 text-indigo-600 mb-8 relative z-10" />
                  <h3 className="text-2xl font-black text-slate-900 uppercase tracking-tight mb-4 relative z-10">{item.title}</h3>
                  <p className="text-slate-500 font-medium leading-relaxed relative z-10">{item.desc}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* FAQ */}
        <section className="py-32 px-6 bg-slate-50">
          <div className="max-w-4xl mx-auto">
            <h2 className="text-4xl font-black text-slate-900 uppercase tracking-tight mb-16 text-center">Frequently Asked Questions</h2>
            <div className="space-y-6">
              {[
                {
                  q: "Does Aurva process or store my actual database payloads?",
                  a: "No. Aurva is strictly VPC-native. We analyze telemetry, query metadata, and access patterns. Your raw data stays yours."
                },
                {
                  q: "Will this slow down my databases?",
                  a: "Not at all. We use lightweight, modern technologies to monitor traffic out-of-band, ensuring zero performance impact on your production workloads."
                },
                {
                  q: "Does it integrate with my current stack?",
                  a: "Yes. Aurva seamlessly pipes high-context security alerts into your existing SIEM, Slack, or webhook of choice."
                }
              ].map((item, idx) => (
                <div key={idx} className="bg-white border border-slate-200 p-8 shadow-[4px_4px_0px_0px_rgba(15,23,42,0.05)]">
                  <h4 className="text-lg font-black text-slate-900 uppercase tracking-tight mb-4 flex items-start gap-4">
                    <span className="text-indigo-600 font-mono">Q:</span> {item.q}
                  </h4>
                  <p className="text-slate-500 font-medium leading-relaxed flex items-start gap-4">
                    <span className="text-orange-600 font-mono font-black">A:</span> {item.a}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Footer */}
        <footer className="bg-slate-900 text-white pt-24 pb-12 px-6">
          <div className="max-w-7xl mx-auto">
            <div className="grid grid-cols-2 md:grid-cols-5 gap-12 mb-24">
              <div className="col-span-2">
                <div className="flex items-center gap-2 mb-8">
                  <div className="w-8 h-8 bg-indigo-600 flex items-center justify-center">
                    <Shield size={18} />
                  </div>
                  <span className="text-xl font-black uppercase tracking-tighter">Aurva</span>
                </div>
                <p className="text-slate-400 font-medium leading-relaxed mb-8 max-w-xs">
                  India's first VPC-native security engine for modern data privacy and DPDP Act compliance.
                </p>
                <div className="flex gap-4">
                  <Globe className="w-5 h-5 text-slate-500 hover:text-white transition-colors cursor-pointer" />
                  <MessageSquare className="w-5 h-5 text-slate-500 hover:text-white transition-colors cursor-pointer" />
                  <Shield className="w-5 h-5 text-slate-500 hover:text-white transition-colors cursor-pointer" />
                </div>
              </div>
              
              {[
                { 
                  title: "Product", 
                  links: ["Integrations", "Pricing", "Changelog", "Roadmap"] 
                },
                { 
                  title: "Resources", 
                  links: ["Documentation", "Blog", "Security", "Status"] 
                },
                { 
                  title: "Company", 
                  links: ["About Us", "Careers", "Privacy Policy", "Terms"] 
                }
              ].map((group, idx) => (
                <div key={idx}>
                  <h4 className="text-[10px] font-black uppercase tracking-[0.2em] text-indigo-400 mb-6">{group.title}</h4>
                  <ul className="space-y-4">
                    {group.links.map((link, lIdx) => (
                      <li key={lIdx} className="text-[11px] font-bold text-slate-400 hover:text-white transition-colors cursor-pointer uppercase tracking-widest">{link}</li>
                    ))}
                  </ul>
                </div>
              ))}
            </div>
            
            <div className="border-t border-slate-800 pt-12 flex flex-col md:flex-row items-center justify-between gap-6">
              <p className="text-[10px] font-bold text-slate-500 uppercase tracking-widest">
                © 2026 Aurva Inc. All rights reserved.
              </p>
              <div className="flex items-center gap-6 text-[10px] font-bold text-slate-500 uppercase tracking-widest">
                <span>support@aurva.in</span>
                <div className="flex items-center gap-2 text-emerald-500">
                  <CheckCircle2 size={12} /> System Status: Normal
                </div>
              </div>
            </div>
          </div>
        </footer>
      </main>
    </div>
  );
}
