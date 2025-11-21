import { createContext, useContext, useEffect, useState } from 'react';
import { Spinner } from '../components/ui/spinner';
import { CoreService } from '@openvault/openvault';

export interface LockState {
  isLocked: boolean;
  unlock: (password: string) => Promise<void>;
  lock: () => Promise<void>;
}

const LockContext = createContext<LockState | undefined>(undefined);

export function LockProvider({ children }: { children: React.ReactNode }) {
  const [isLocked, setIsLocked] = useState<boolean>(true);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  useEffect(() => {
    CoreService.IsLocked()
      .then((locked) => {
        setIsLocked(locked);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  if (isLoading) {
    return (
      <div className="h-full w-full items-center justify-center">
        <Spinner className="size-8" />
      </div>
    );
  }

  const unlock = async (password: string) => {
    try {
      await CoreService.TryUnlock(password);
      setIsLocked(false);
    } catch (error) {
      throw new Error(`Failed to unlock: ${error}`);
    }
  };

  const lock = async () => {
    try {
      await CoreService.Lock();
      setIsLocked(true);
      window.location.replace('/');
    } catch (error) {
      throw new Error(`Failed to lock: ${error}`);
    }
  };

  return (
    <LockContext.Provider value={{ isLocked, unlock, lock }}>
      {children}
    </LockContext.Provider>
  );
}

export function useLock() {
  const context = useContext(LockContext);
  if (context === undefined) {
    throw new Error('useLock must be used within a LockProvider');
  }
  return context;
}
