import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api, ScanRequest } from '@/lib/api';
import toast from 'react-hot-toast';

export function useScan(scanId: string, enabled: boolean = true) {
  return useQuery({
    queryKey: ['scan', scanId],
    queryFn: () => api.getScanStatus(scanId),
    enabled: enabled && !!scanId,
    refetchInterval: (query) => {
      const data = query.state.data;
      if (!data) return 2000;
      return data.status === 'completed' || data.status === 'failed' ? false : 2000;
    },
  });
}

export function useTriggerScan() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: ScanRequest) => api.triggerScan(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['findings'] });
      toast.success('Scan initiated successfully');
    },
    onError: (error: Error) => {
      toast.error(`Failed to start scan: ${error.message}`);
    },
  });
}

export function useFindings() {
  return useQuery({
    queryKey: ['findings'],
    queryFn: () => api.getFindings(),
    staleTime: 30 * 1000,
  });
}

export function useHealth() {
  return useQuery({
    queryKey: ['health'],
    queryFn: () => api.checkHealth(),
    refetchInterval: 10000,
    retry: false,
  });
}
