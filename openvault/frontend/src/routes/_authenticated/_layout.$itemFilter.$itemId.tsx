import {
  DetailsButton,
  DetailsMarkdown,
  VaultDropdown
} from '@/components/item-details';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';

import { CoreService } from '@openvault/openvault';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute(
  '/_authenticated/_layout/$itemFilter/$itemId'
)({
  loader: async ({ params }) => {
    const { itemId } = params;
    const overview = await CoreService.GetItemOverview(itemId);
    const details = await CoreService.GetVaultItemDetails(itemId);
    if (!details || !overview) {
      throw new Error('Failed to load vault item');
    }
    const vaultMeta = await CoreService.GetVaultMetadata(overview.vault_id);
    if (!vaultMeta) {
      throw new Error('Failed to load vault item');
    }
    const account = await CoreService.GetAccount(vaultMeta.account_id);
    if (!account) {
      throw new Error('Failed to load account');
    }
    return { details, overview, vaultMeta, account };
  },
  component: RouteComponent
});

function RouteComponent() {
  const { details, overview, vaultMeta, account } = Route.useLoaderData();
  return (
    <div className="flex flex-col gap-4 p-4">
      <div className="flex items-center justify-between">
        <VaultDropdown
          accountName={account.user_email}
          vaultName={vaultMeta.name}
        />
      </div>
      <div className="flex items-center gap-x-4 px-4 py-2">
        <Avatar className="size-14 rounded-lg text-lg font-bold">
          <AvatarFallback className="rounded-lg">
            {overview.title.charAt(0)}
          </AvatarFallback>
        </Avatar>
        <h1 className="text-lg font-bold">{overview.title}</h1>
      </div>
      <div className="flex flex-col">
        <DetailsButton label="Username" content={details.username} />
        <DetailsButton label="Password" content={details.password} conceal />
      </div>
      <div className="flex flex-col">
        <DetailsMarkdown label="Notes" content={details.notes} />
      </div>
    </div>
  );
}
