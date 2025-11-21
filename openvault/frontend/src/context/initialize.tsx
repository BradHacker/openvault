import { createContext, useContext, useEffect, useState } from 'react';
import { Spinner } from '../components/ui/spinner';
import { CoreService } from '@openvault/openvault';
import type { InitOptions } from '@openvault/openvault/internal/fs';

export interface InitializeState {
  isInitialized: boolean;
  initialize: (options: InitOptions) => Promise<void>;
}

const InitializeContext = createContext<InitializeState | undefined>(undefined);

export function InitializeProvider({
  children
}: {
  children: React.ReactNode;
}) {
  const [isInitialized, setIsInitialized] = useState<boolean>(true);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  useEffect(() => {
    CoreService.IsInitialized()
      .then((initialized) => {
        setIsInitialized(initialized);
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

  const initialize = async (options: InitOptions) => {
    try {
      await CoreService.Initialize(options);
      setIsInitialized(true);
    } catch (error) {
      throw new Error(`Failed to initialize: ${error}`);
    }
  };

  return (
    <InitializeContext.Provider value={{ isInitialized, initialize }}>
      {children}
    </InitializeContext.Provider>
  );
}

export function useInitialize() {
  const context = useContext(InitializeContext);
  if (context === undefined) {
    throw new Error('useInitialize must be used within a InitializeProvider');
  }
  return context;
}
