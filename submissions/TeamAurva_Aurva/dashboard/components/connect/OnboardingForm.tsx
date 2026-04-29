"use client";

import { useState } from 'react';
import { motion } from 'framer-motion';
import { useTriggerScan } from '@/hooks/useScan';
import { useRouter } from 'next/navigation';
import { Loader2, Terminal } from 'lucide-react';

export function OnboardingForm() {
  const [accountId, setAccountId] = useState('');
  const [roleArn, setRoleArn] = useState('');
  const [nickname, setNickname] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  
  const router = useRouter();
  const { mutate: triggerScan, isPending } = useTriggerScan();

  const validateAccountId = (value: string) => {
    if (!value) return 'Account ID is required';
    if (!/^\d{12}$/.test(value)) return 'Must be 12 digits';
    return '';
  };

  const validateRoleArn = (value: string) => {
    if (!value) return 'Role ARN is required';
    if (!value.startsWith('arn:aws:iam::')) return 'Must start with arn:aws:iam::';
    return '';
  };

  const validateNickname = (value: string) => {
    if (!value) return 'Nickname is required for the demo';
    if (value.length < 3) return 'Nickname must be at least 3 chars';
    return '';
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    const newErrors: Record<string, string> = {};
    const accountError = validateAccountId(accountId);
    const roleError = validateRoleArn(roleArn);
    const nicknameError = validateNickname(nickname);
    
    if (accountError) newErrors.accountId = accountError;
    if (roleError) newErrors.roleArn = roleError;
    if (nicknameError) newErrors.nickname = nicknameError;
    
    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    triggerScan(
      {
        account_id: accountId,
        role_arn: roleArn,
        account_nickname: nickname,
      },
      {
        onSuccess: (data) => {
          router.push(`/scan/${data.scan_id}`);
        },
      }
    );
  };

  return (
    <motion.form
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: 0.2 }}
      onSubmit={handleSubmit}
      className="space-y-8"
    >
      <div>
        <label htmlFor="accountId" className="block text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-2">
          AWS Account ID
        </label>
        <input
          id="accountId"
          type="text"
          value={accountId}
          onChange={(e) => {
            setAccountId(e.target.value);
            setErrors({ ...errors, accountId: '' });
          }}
          placeholder="123456789012"
          className="w-full bg-slate-50 border border-slate-200 rounded-none px-5 py-4 text-slate-900 placeholder-slate-400 focus:outline-none focus:border-indigo-600 focus:bg-white transition-all font-mono text-sm shadow-[4px_4px_0px_0px_rgba(15,23,42,0.05)]"
        />
        {errors.accountId && (
          <p className="text-rose-600 text-[10px] font-bold uppercase mt-2 tracking-wide">{errors.accountId}</p>
        )}
      </div>

      <div>
        <label htmlFor="roleArn" className="block text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-2">
          IAM Role ARN
        </label>
        <input
          id="roleArn"
          type="text"
          value={roleArn}
          onChange={(e) => {
            setRoleArn(e.target.value);
            setErrors({ ...errors, roleArn: '' });
          }}
          placeholder="arn:aws:iam::123456789012:role/AurvaReadOnly"
          className="w-full bg-slate-50 border border-slate-200 rounded-none px-5 py-4 text-slate-900 placeholder-slate-400 focus:outline-none focus:border-indigo-600 focus:bg-white transition-all font-mono text-xs shadow-[4px_4px_0px_0px_rgba(15,23,42,0.05)]"
        />
        {errors.roleArn && (
          <p className="text-rose-600 text-[10px] font-bold uppercase mt-2 tracking-wide">{errors.roleArn}</p>
        )}
      </div>

      <div>
        <label htmlFor="nickname" className="block text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-2">
          Account Nickname
        </label>
        <input
          id="nickname"
          type="text"
          value={nickname}
          onChange={(e) => {
            setNickname(e.target.value);
            setErrors({ ...errors, nickname: '' });
          }}
          placeholder="Prod-Cluster-01"
          className="w-full bg-slate-50 border border-slate-200 rounded-none px-5 py-4 text-slate-900 placeholder-slate-400 focus:outline-none focus:border-indigo-600 focus:bg-white transition-all text-sm shadow-[4px_4px_0px_0px_rgba(15,23,42,0.05)]"
        />
        {errors.nickname && (
          <p className="text-rose-600 text-[10px] font-bold uppercase mt-2 tracking-wide">{errors.nickname}</p>
        )}
      </div>

      <button
        type="submit"
        disabled={isPending}
        className="w-full h-14 bg-indigo-600 hover:bg-indigo-700 text-white font-black text-xs uppercase tracking-[0.2em] transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3 shadow-[6px_6px_0px_0px_rgba(79,70,229,0.3)]"
      >
        {isPending ? (
          <>
            <Loader2 className="w-5 h-5 animate-spin" />
            Connecting...
          </>
        ) : (
          <>
            <Terminal size={18} /> Establish Connection
          </>
        )}
      </button>

      <p className="text-slate-400 text-[10px] text-center font-bold uppercase tracking-widest leading-relaxed">
        System guarantees 100% data residency.<br/>Credentials never stored.
      </p>
    </motion.form>
  );
}
