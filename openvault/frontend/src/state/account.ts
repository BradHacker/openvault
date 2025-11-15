import { create } from 'zustand';

interface Account {
  id: string;
  email: string;
  name: string;
}

interface AccountsState {
  accounts: Account[];
  setAccounts: (accounts: Account[]) => void;
}

export const useAccounts = create<AccountsState>((set) => ({
  accounts: [],
  setAccounts: (accounts) => set((state) => ({ ...state, accounts }))
}));
