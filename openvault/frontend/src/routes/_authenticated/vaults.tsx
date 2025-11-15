import { TitleBar } from '@/components/titlebar';

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarTrigger
} from '@/components/ui/sidebar';
import { GetAccounts, GetVaultMetadatas } from '@/wailsjs/go/main/App';
import { createFileRoute, Link, Outlet } from '@tanstack/react-router';
import { Vault } from 'lucide-react';

export const Route = createFileRoute('/_authenticated/vaults')({
  component: RouteComponent,
  loader: async () => {
    const accounts = await GetAccounts();
    // Load the vault; metadata to display in sidebar
    const vaultMetadatas = await GetVaultMetadatas(
      accounts.filter((a) => a.is_unlocked).map((a) => a.id)
    );
    return { vaultMetadatas };
  }
});

function RouteComponent() {
  const { vaultMetadatas } = Route.useLoaderData();

  return (
    <div className="grid h-full w-full grid-cols-[min-content_auto]">
      <Sidebar collapsible="icon">
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Vaults</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild>
                    <Link
                      to="/vaults/$id"
                      params={{ id: 'all' }}
                      activeProps={{ 'data-active': true }}
                      className="data-[active=true]:bg-accent/50"
                    >
                      <Vault />
                      <span>All Vaults</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                {vaultMetadatas.map((vault) => (
                  <SidebarMenuItem key={vault.vault_id}>
                    <SidebarMenuButton asChild>
                      <Link
                        to="/vaults/$id"
                        params={{ id: vault.vault_id }}
                        activeProps={{ 'data-active': true }}
                        className="data-[active=true]:bg-accent/50"
                      >
                        <Vault />
                        <span>{vault.name}</span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <div className="flex w-full justify-end">
            <SidebarTrigger />
          </div>
        </SidebarFooter>
      </Sidebar>
      <div className="h-full w-full">
        <TitleBar className="col-span-2" />
        <Outlet />
      </div>
    </div>
  );
}
