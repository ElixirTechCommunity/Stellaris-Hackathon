import { create } from 'zustand';
import { PIIFinding } from './api';

interface AppState {
  selectedFinding: PIIFinding | null;
  setSelectedFinding: (finding: PIIFinding | null) => void;
  
  riskFilter: string[];
  setRiskFilter: (risks: string[]) => void;
  
  piiTypeFilter: string[];
  setPIITypeFilter: (types: string[]) => void;
  
  sortColumn: string;
  sortDirection: 'asc' | 'desc';
  setSorting: (column: string, direction: 'asc' | 'desc') => void;
  
  isDrawerOpen: boolean;
  setDrawerOpen: (open: boolean) => void;
}

export const useStore = create<AppState>((set) => ({
  selectedFinding: null,
  setSelectedFinding: (finding) => set({ selectedFinding: finding, isDrawerOpen: !!finding }),
  
  riskFilter: [],
  setRiskFilter: (risks) => set({ riskFilter: risks }),
  
  piiTypeFilter: [],
  setPIITypeFilter: (types) => set({ piiTypeFilter: types }),
  
  sortColumn: 'discovered_at',
  sortDirection: 'desc',
  setSorting: (column, direction) => set({ sortColumn: column, sortDirection: direction }),
  
  isDrawerOpen: false,
  setDrawerOpen: (open) => set({ isDrawerOpen: open }),
}));
