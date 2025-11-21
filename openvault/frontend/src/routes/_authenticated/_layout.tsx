import { TitleBar } from '@/components/titlebar';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupAction,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSkeleton,
  SidebarTrigger
} from '@/components/ui/sidebar';
import { useAccountFilter } from '@/context/account-filter';
import { cn } from '@/lib/utils';
import { CoreService } from '@openvault/openvault';
import { VaultMetadata } from '@openvault/openvault/internal/structs';
import { createFileRoute, Link, Outlet } from '@tanstack/react-router';
import {
  Check,
  ChevronsUpDown,
  Lock,
  Plus,
  Star,
  Users,
  Vault,
  WalletCards
} from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

export const Route = createFileRoute('/_authenticated/_layout')({
  component: RouteComponent
});

function RouteComponent() {
  const { activeAccount, allAccounts, setActiveAccountId } = useAccountFilter();
  const [vaultMetadatas, setVaultMetadatas] = useState<VaultMetadata[]>([]);
  const [vaultMetadatasLoading, setVaultMetadatasLoading] = useState(false);

  useEffect(() => {
    setVaultMetadatasLoading(true);
    let accountIds: string[] = [];
    if (activeAccount === 'all') {
      accountIds = allAccounts.map((a) => a.id);
    } else {
      accountIds = [activeAccount.id];
    }
    CoreService.ListVaultMetadatas(accountIds)
      .then(
        (vaults) => setVaultMetadatas(vaults.filterNullish()),
        (err) => {
          toast.error(`Failed to load vaults: ${err.message}`);
        }
      )
      .finally(() => setVaultMetadatasLoading(false));
  }, [activeAccount, allAccounts]);

  return (
    <div className="grid h-full w-full grid-cols-[min-content_auto]">
      <Sidebar collapsible="icon">
        <SidebarContent>
          <AppSidebarHeader />
          <SidebarGroup>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild>
                    <Link
                      to="/$itemFilter"
                      params={{ itemFilter: 'all' }}
                      activeProps={{ 'data-active': true }}
                      className="data-[active=true]:bg-accent/50"
                    >
                      <WalletCards />
                      <span>All Vaults</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem>
                  <SidebarMenuButton asChild>
                    <Link
                      to="/$itemFilter"
                      params={{ itemFilter: 'favorites' }}
                      activeProps={{ 'data-active': true }}
                      className="data-[active=true]:bg-accent/50"
                    >
                      <Star />
                      <span>Favorites</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
          <SidebarGroup>
            <SidebarGroupLabel>Vaults</SidebarGroupLabel>
            <SidebarGroupAction title="New Vault">
              <Plus className="size-4" />
              <span className="sr-only">New Vault</span>
            </SidebarGroupAction>
            <SidebarGroupContent>
              <SidebarMenu>
                {vaultMetadatasLoading &&
                  Array.from({ length: 3 }).map((_, i) => (
                    <SidebarMenuItem key={`vault_skeleton_${i}`}>
                      <SidebarMenuSkeleton />
                    </SidebarMenuItem>
                  ))}
                {!vaultMetadatasLoading &&
                  vaultMetadatas.map((vault) => (
                    <SidebarMenuItem key={vault.vault_id}>
                      <SidebarMenuButton
                        tooltip={{
                          children: vault.name
                        }}
                        asChild
                      >
                        <Link
                          to="/$itemFilter"
                          params={{ itemFilter: vault.vault_id }}
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
      <div className="flex h-screen w-full flex-col">
        <TitleBar />
        <main className="w-full grow">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

function AppSidebarHeader() {
  const { allAccounts, activeAccount, setActiveAccountId } = useAccountFilter();

  return (
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <SidebarMenuButton
                size="lg"
                className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
              >
                <Avatar className="bg-sidebar-secondary text-sidebar-secondary-foreground flex aspect-square size-8 items-center justify-center">
                  <AvatarFallback>
                    {activeAccount === 'all' ? (
                      <Users className="size-4" />
                    ) : (
                      activeAccount.user_first_name[0] +
                      activeAccount.user_last_name[0]
                    )}
                  </AvatarFallback>
                </Avatar>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="font-medium">
                    {activeAccount === 'all'
                      ? 'All Accounts'
                      : `${activeAccount.user_first_name} ${activeAccount.user_last_name}`}
                  </span>
                  {activeAccount !== 'all' && (
                    <span className="text-xs text-muted-foreground">
                      {activeAccount.user_email}
                    </span>
                  )}
                </div>
                <ChevronsUpDown className="ml-auto" />
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              className="w-(--radix-dropdown-menu-trigger-width)"
              align="start"
            >
              <DropdownMenuItem onSelect={() => setActiveAccountId('all')}>
                <Avatar className="bg-sidebar-secondary text-sidebar-secondary-foreground flex aspect-square size-8 items-center justify-center">
                  <AvatarFallback>
                    <Users className="size-4" />
                  </AvatarFallback>
                </Avatar>
                <div className="flex flex-col">
                  <span>All Accounts</span>
                </div>
                {activeAccount === 'all' && <Check className="ml-auto" />}
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              {allAccounts.map((account) => (
                <DropdownMenuItem
                  key={account.id}
                  onSelect={() => setActiveAccountId(account.id)}
                >
                  <Avatar className="bg-sidebar-secondary text-sidebar-secondary-foreground flex aspect-square size-8 items-center justify-center">
                    <AvatarFallback>
                      {account.user_first_name[0]}
                      {account.user_last_name[0]}
                    </AvatarFallback>
                  </Avatar>
                  <div
                    className={cn(
                      'flex flex-col',
                      !account.is_unlocked && 'text-muted-foreground'
                    )}
                  >
                    <span>
                      {account.user_first_name} {account.user_last_name}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {account.user_email}
                    </span>
                  </div>
                  {!account.is_unlocked ? (
                    <Lock className="ml-auto text-muted-foreground" />
                  ) : (
                    activeAccount !== 'all' &&
                    account.id === activeAccount.id && (
                      <Check className="ml-auto" />
                    )
                  )}
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
  );
}
