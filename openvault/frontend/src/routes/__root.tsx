import { createRootRouteWithContext, Outlet } from '@tanstack/react-router';
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools';
import { TanStackDevtools } from '@tanstack/react-devtools';
import { ThemeProvider } from '@/components/theme-provider';
import { Toaster } from '@/components/ui/sonner';
import type { LockState } from '@/context/lock';
import { SidebarProvider } from '@/components/ui/sidebar';
import { FloatingWindowControls } from '@/components/titlebar';
import type { InitializeState } from '@/context/initialize';
import type { AccountFilterState } from '@/context/account-filter';

interface RootRouterContext {
  lock: LockState;
  initialize: InitializeState;
  accountFilter: AccountFilterState;
}

export const Route = createRootRouteWithContext<RootRouterContext>()({
  component: Root,
  notFoundComponent: NotFound
});

function NotFound() {
  const match = Route.useMatch();
  console.error('No route found for location:', match);
  return <div>404 - Page Not Found</div>;
}

function Providers({ children }: React.PropsWithChildren) {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <SidebarProvider>
        <Toaster />
        {children}
      </SidebarProvider>
    </ThemeProvider>
  );
}

function Root() {
  return (
    <Providers>
      <div className="flex h-screen w-screen flex-col">
        <Outlet />
        <FloatingWindowControls />
      </div>
      <TanStackDevtools
        config={{
          position: 'bottom-right'
        }}
        plugins={[
          {
            name: 'Tanstack Router',
            render: <TanStackRouterDevtoolsPanel />
          }
        ]}
      />
    </Providers>
  );
}
