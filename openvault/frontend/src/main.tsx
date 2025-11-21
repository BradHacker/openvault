import { StrictMode } from 'react';
import ReactDOM from 'react-dom/client';
import { RouterProvider, createRouter } from '@tanstack/react-router';

import '@wailsio/runtime';
import { Window } from '@wailsio/runtime';

window.onload = () => {
  Window.SetFrameless(true);
};

// Import array prototype extensions
import '@/lib/array-filter.ts';

// Import the generated route tree
import { routeTree } from './routeTree.gen';

import './styles.css';
import reportWebVitals from './reportWebVitals.ts';
import { LockProvider, useLock } from './context/lock.tsx';
import { InitializeProvider, useInitialize } from './context/initialize';
import {
  AccountFilterProvider,
  useAccountFilter
} from './context/account-filter';

// Create a new router instance
const router = createRouter({
  routeTree,
  context: {
    lock: undefined!, // Will be provided by LockProvider
    initialize: undefined!, // Will be provided by InitializeProvider
    accountFilter: undefined! // Will be provided by AccountProvider
  },
  defaultPreload: 'intent',
  scrollRestoration: true,
  defaultStructuralSharing: true,
  defaultPreloadStaleTime: 0
});

// Register the router instance for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

function InnerApp() {
  const lock = useLock();
  const initialize = useInitialize();
  const accountFilter = useAccountFilter();
  return (
    <RouterProvider
      router={router}
      context={{ lock, initialize, accountFilter }}
    />
  );
}

// Render the app
const rootElement = document.getElementById('app');
if (rootElement && !rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <StrictMode>
      <InitializeProvider>
        <LockProvider>
          <AccountFilterProvider>
            <InnerApp />
          </AccountFilterProvider>
        </LockProvider>
      </InitializeProvider>
    </StrictMode>
  );
}

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
