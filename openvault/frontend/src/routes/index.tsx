import { createFileRoute, Navigate } from '@tanstack/react-router';

export const Route = createFileRoute('/')({
  component: RouteComponent
});

function RouteComponent() {
  return <Navigate to="/$itemFilter" params={{ itemFilter: 'all' }} replace />;
}
