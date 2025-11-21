import { createContext, useContext } from 'react';
import type { VaultMetadata } from '@openvault/openvault/internal/structs';

export interface VaultState {
  vault: VaultMetadata | 'all' | 'favorites';
}

const VaultContext = createContext<VaultState | undefined>(undefined);

export function VaultProvider({
  children,
  vault
}: {
  children: React.ReactNode;
  vault: VaultMetadata | 'all' | 'favorites';
}) {
  return (
    <VaultContext.Provider
      value={{
        vault
      }}
    >
      {children}
    </VaultContext.Provider>
  );
}

export function useVault() {
  const context = useContext(VaultContext);
  if (context === undefined) {
    throw new Error('useVault must be used within a VaultProvider');
  }
  return context;
}
