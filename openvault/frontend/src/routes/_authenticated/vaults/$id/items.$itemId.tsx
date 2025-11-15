import { GetVaultItemDetails } from '@/wailsjs/go/main/App';
import { createFileRoute } from '@tanstack/react-router';
import z from 'zod';

const routeSearchSchema = z.object({
  vaultId: z.string().optional()
});

export const Route = createFileRoute(
  '/_authenticated/vaults/$id/items/$itemId'
)({
  validateSearch: routeSearchSchema,
  loaderDeps: ({ search }) => ({ vaultId: search.vaultId }),
  loader: async ({ params, deps: { vaultId } }) => {
    const { id, itemId } = params;
    const details = await GetVaultItemDetails(vaultId || id, itemId);
    return { details };
  },
  component: RouteComponent
});

function RouteComponent() {
  const { details } = Route.useLoaderData();
  return (
    <div className="flex flex-col gap-4 p-4">
      {details.username}
      {details.password}
    </div>
  );
}
