import { createFileRoute, redirect, Outlet } from '@tanstack/react-router';

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: async ({ context, location }) => {
    console.log('_authenticated: isLocked = ', context.lock.isLocked);
    if (context.lock.isLocked) {
      throw redirect({
        to: '/lock',
        search: {
          // Save current location for redirect after login
          redirect: location.href
        }
      });
    }
  },
  component: () => <Outlet />
});
