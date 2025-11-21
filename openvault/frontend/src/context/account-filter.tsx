import { createContext, useContext, useState } from 'react';
import { AccountWithUnlockStatus } from '@openvault/openvault';

export interface AccountFilterState {
  activeAccount: 'all' | AccountWithUnlockStatus;
  setActiveAccountId: (filter: 'all' | string) => void;
  allAccounts: AccountWithUnlockStatus[];
  setAllAccounts: (accounts: AccountWithUnlockStatus[]) => void;
}

const AccountFilterContext = createContext<AccountFilterState | undefined>(
  undefined
);

export function AccountFilterProvider({
  children
}: {
  children: React.ReactNode;
}) {
  const [allAccounts, setAllAccounts] = useState<AccountWithUnlockStatus[]>([]);
  const [activeAccountId, setActiveAccountId] = useState<'all' | string>('all');

  const activeAccount = allAccounts.find((a) => a.id === activeAccountId);

  if (activeAccount === undefined && activeAccountId !== 'all') {
    throw new Error(
      `AccountFilterProvider: activeAccountId "${activeAccountId}" not found in allAccounts`
    );
  }

  return (
    <AccountFilterContext.Provider
      value={{
        activeAccount: activeAccount || 'all',
        setActiveAccountId,
        allAccounts,
        setAllAccounts
      }}
    >
      {children}
    </AccountFilterContext.Provider>
  );
}

export function useAccountFilter() {
  const context = useContext(AccountFilterContext);
  if (context === undefined) {
    throw new Error(
      'useAccountFilter must be used within a AccountFilterProvider'
    );
  }
  return context;
}
