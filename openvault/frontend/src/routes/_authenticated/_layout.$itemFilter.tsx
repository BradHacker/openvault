import { Button } from '@/components/ui/button';
import { Item, ItemContent, ItemMedia, ItemTitle } from '@/components/ui/item';
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup
} from '@/components/ui/resizable';
import { CoreService } from '@openvault/openvault';
import {
  createFileRoute,
  Link,
  Outlet,
  type ErrorComponentProps
} from '@tanstack/react-router';
import { StickyNote } from 'lucide-react';

export const Route = createFileRoute('/_authenticated/_layout/$itemFilter')({
  loader: async ({ params }) => {
    const { itemFilter } = params;
    let overviews = [];
    if (itemFilter === 'all') {
      overviews = await CoreService.ListAllItemOverviews();
    } else if (itemFilter === 'favorites') {
      // overviews = await CoreService.GetAllVaultItemOverviews();
      // TODO: implement favorites
      throw new Error('Favorites not implemented yet');
    } else {
      overviews = await CoreService.ListVaultItemOverviews(itemFilter);
    }
    return { overviews: overviews.filterNullish() };
  },
  component: RouteComponent,
  errorComponent: ErrorComponent
});

function RouteComponent() {
  const { itemFilter } = Route.useParams();
  const { overviews } = Route.useLoaderData();
  return (
    <ResizablePanelGroup direction="horizontal">
      <ResizablePanel className="min-w-64" defaultSize={33}>
        <div className="flex h-full w-full flex-col space-y-2 overflow-y-auto px-2 pt-2">
          {overviews.map((o) => (
            <Item key={o.item_id} variant="default" size="sm" asChild>
              <Link
                to="/$itemFilter/$itemId"
                params={{ itemFilter, itemId: o.item_id }}
                activeProps={{ 'data-active': true }}
                className="focus:bg-primary! focus:text-black data-[active=true]:bg-accent/50 hover:[&:not([data-active])]:bg-primary/50"
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
      <ResizableHandle withHandle />
      <ResizablePanel>
        <Outlet />
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}

function ErrorComponent({ error, info, reset }: ErrorComponentProps) {
  return (
    <div className="flex h-full w-full flex-col items-center justify-center gap-2">
      <p className="text-red-500">Error: {error.message}</p>
      {process.env.NODE_ENV === 'development' && info?.componentStack && (
        <pre className="text-sm whitespace-pre-wrap text-muted-foreground">
          {info.componentStack}
        </pre>
      )}
      <Button variant="outline" onClick={reset}>
        Try again
      </Button>
    </div>
  );
}
