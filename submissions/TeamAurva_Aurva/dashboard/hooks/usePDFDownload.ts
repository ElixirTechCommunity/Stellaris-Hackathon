import { api } from '@/lib/api';
import toast from 'react-hot-toast';

export function usePDFDownload() {
  const download = async () => {
    try {
      toast.loading('Generating PDF report...', { id: 'pdf-download' });
      const blob = await api.downloadReport();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `aurva-compliance-report-${Date.now()}.pdf`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
      toast.success('Report downloaded', { id: 'pdf-download' });
    } catch (error) {
      toast.error('Failed to download report', { id: 'pdf-download' });
    }
  };

  return { download };
}
