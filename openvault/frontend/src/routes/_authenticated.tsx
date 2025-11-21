import { useAccountFilter } from '@/context/account-filter';
import { CoreService } from '@openvault/openvault';
import { createFileRoute, redirect, Outlet } from '@tanstack/react-router';
import { useEffect } from 'react';

export const Route = createFileRoute('/_authenticated')({
  component: RouteComponent,
  loader: async ({ context, location }) => {
    // console.log(
    //   '_authenticated: isInitialized = ',
    //   context.initialize.isInitialized
    // );
    if (!context.initialize.isInitialized) {
      throw redirect({
        to: '/initialize'
      });
    }
    // console.log('_authenticated: isLocked = ', context.lock.isLocked);
    if (context.lock.isLocked) {
      throw redirect({
        to: '/lock',
        search: {
          // Save current location for redirect after login
          redirect: location.href
        }
      });
    }
    const accounts = (await CoreService.GetAccounts()).filterNullish();

    return {
      accounts
    };
  }
});

function RouteComponent() {
  const { accounts } = Route.useLoaderData();
  const { setAllAccounts } = useAccountFilter();

  useEffect(() => {
    setAllAccounts(accounts);
  }, [accounts, setAllAccounts]);

  return <Outlet />;
}
