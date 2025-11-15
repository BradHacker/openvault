import { Item, ItemContent, ItemMedia, ItemTitle } from '@/components/ui/item';
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup
} from '@/components/ui/resizable';
import {
  GetAllVaultItemOverviews,
  GetVaultItemOverviews
} from '@/wailsjs/go/main/App';
import { createFileRoute, Link, Outlet } from '@tanstack/react-router';
import { StickyNote } from 'lucide-react';

export const Route = createFileRoute('/_authenticated/vaults/$id')({
  loader: async ({ params }) => {
    const { id } = params;
    let overviews = [];
    if (id === 'all') {
      overviews = await GetAllVaultItemOverviews();
    } else {
      overviews = await GetVaultItemOverviews(id);
    }
    return { overviews };
  },
  component: RouteComponent
});

function RouteComponent() {
  const { id } = Route.useParams();
  const { overviews } = Route.useLoaderData();
  return (
    <ResizablePanelGroup direction="horizontal">
      <ResizablePanel>
        <div className="flex h-full w-full flex-col space-y-2 overflow-y-auto px-2 pt-2">
          {overviews.map((o) => (
            <Item key={o.item_id} variant="outline" size="sm" asChild>
              <Link
                to="/vaults/$id/items/$itemId"
                params={{ id, itemId: o.item_id }}
                search={{ vaultId: id === 'all' ? o.vault_id : undefined }}
                activeProps={{ 'data-active': true }}
                className="data-[active=true]:border-primary/50 data-[active=true]:bg-accent/50"
              >
                <ItemMedia>
                  <StickyNote className="size-5" />
                </ItemMedia>
                <ItemContent>
                  <ItemTitle>{o.title}</ItemTitle>
                </ItemContent>
              </Link>
            </Item>
          ))}
        </div>
      </ResizablePanel>
      <ResizableHandle />
      <ResizablePanel>
        <Outlet />
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}
