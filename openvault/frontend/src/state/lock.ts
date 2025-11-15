import { create } from 'zustand';

interface LockState {
  isLocked: boolean;
  setIsLocked: (isLocked: boolean) => void;
}

export const useLock = create<LockState>((set) => ({
  isLocked: false,
  setIsLocked: (isLocked) => set((state) => ({ ...state, isLocked }))
}));
