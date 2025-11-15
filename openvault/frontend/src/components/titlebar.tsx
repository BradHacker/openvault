import { Maximize, Minimize, Minus, Vault, X } from 'lucide-react';
import { Button } from './ui/button';
import {
  Hide,
  Quit,
  WindowIsMaximised,
  WindowMaximise,
  WindowUnmaximise
} from '@/wailsjs/runtime/runtime';
import { useEffect, useState } from 'react';
import {
  Menubar,
  MenubarCheckboxItem,
  MenubarContent,
  MenubarItem,
  MenubarMenu,
  MenubarRadioGroup,
  MenubarRadioItem,
  MenubarSeparator,
  MenubarShortcut,
  MenubarSub,
  MenubarSubContent,
  MenubarSubTrigger,
  MenubarTrigger
} from './ui/menubar';
import { useAccounts } from '@/state/account';
import { useLock } from '@/lock';
import { cn } from '@/lib/utils';

export function TitleBar({ className }: { className?: string }) {
  const [isMaximized, setIsMaximized] = useState(false);

  useEffect(() => {
    const checkMaximized = async () => {
      const maximized = await WindowIsMaximised();
      setIsMaximized(maximized);
    };
    checkMaximized().catch((err) => {
      console.error('Error checking window maximized state:', err);
    });
  }, []);

  return (
    <>
      <div
        className={cn(
          'flex h-10 w-full cursor-default items-center gap-x-2 px-4',
          className
        )}
        style={{ '--wails-draggable': 'drag' }}
      >
        <Vault className="size-5 text-primary" />
        <h1 className="font-bold text-primary select-none">OpenVault</h1>
        <Menu />
      </div>
      {/* Spacer */}
      {/* <div className="h-10 w-full"></div> */}
    </>
  );
}

export function FloatingWindowControls() {
  const [isMaximized, setIsMaximized] = useState(false);

  useEffect(() => {
    const checkMaximized = async () => {
      const maximized = await WindowIsMaximised();
      setIsMaximized(maximized);
    };
    checkMaximized().catch((err) => {
      console.error('Error checking window maximized state:', err);
    });
  }, []);

  return (
    <div className="fixed top-0 right-0 flex h-10 items-center gap-x-2 px-2">
      <WindowControlButton onClick={() => Hide()}>
        <Minus className="size-4" />
      </WindowControlButton>
      <WindowControlButton
        onClick={() => {
          if (isMaximized) {
            WindowUnmaximise();
            setIsMaximized(false);
          } else {
            WindowMaximise();
            setIsMaximized(true);
          }
        }}
      >
        {isMaximized ? (
          <Minimize className="size-4" />
        ) : (
          <Maximize className="size-4" />
        )}
      </WindowControlButton>
      <WindowControlButton onClick={() => Quit()}>
        <X className="size-4" />
      </WindowControlButton>
    </div>
  );
}

function WindowControlButton(props: React.ComponentProps<typeof Button>) {
  return (
    <Button
      variant="ghost"
      size="icon"
      className="size-7 p-0"
      style={{ '--wails-draggable': 'no-drag' }}
      {...props}
    />
  );
}

function Menu() {
  const { accounts } = useAccounts();
  const { lock } = useLock();

  return (
    <Menubar className="border-none">
      <MenubarMenu>
        <MenubarTrigger>File</MenubarTrigger>
        <MenubarContent>
          <MenubarItem onClick={lock}>
            Lock <MenubarShortcut>⌘L</MenubarShortcut>
          </MenubarItem>
          <MenubarItem>
            New Tab <MenubarShortcut>⌘T</MenubarShortcut>
          </MenubarItem>
          <MenubarItem>
            New Window <MenubarShortcut>⌘N</MenubarShortcut>
          </MenubarItem>
          <MenubarItem disabled>New Incognito Window</MenubarItem>
          <MenubarSeparator />
          <MenubarSub>
            <MenubarSubTrigger>Share</MenubarSubTrigger>
            <MenubarSubContent>
              <MenubarItem>Email link</MenubarItem>
              <MenubarItem>Messages</MenubarItem>
              <MenubarItem>Notes</MenubarItem>
            </MenubarSubContent>
          </MenubarSub>
          <MenubarSeparator />
          <MenubarItem>
            Print... <MenubarShortcut>⌘P</MenubarShortcut>
          </MenubarItem>
        </MenubarContent>
      </MenubarMenu>
      <MenubarMenu>
        <MenubarTrigger>Edit</MenubarTrigger>
        <MenubarContent>
          <MenubarItem>
            Undo <MenubarShortcut>⌘Z</MenubarShortcut>
          </MenubarItem>
          <MenubarItem>
            Redo <MenubarShortcut>⇧⌘Z</MenubarShortcut>
          </MenubarItem>
          <MenubarSeparator />
          <MenubarSub>
            <MenubarSubTrigger>Find</MenubarSubTrigger>
            <MenubarSubContent>
              <MenubarItem>Search the web</MenubarItem>
              <MenubarSeparator />
              <MenubarItem>Find...</MenubarItem>
              <MenubarItem>Find Next</MenubarItem>
              <MenubarItem>Find Previous</MenubarItem>
            </MenubarSubContent>
          </MenubarSub>
          <MenubarSeparator />
          <MenubarItem>Cut</MenubarItem>
          <MenubarItem>Copy</MenubarItem>
          <MenubarItem>Paste</MenubarItem>
        </MenubarContent>
      </MenubarMenu>
      <MenubarMenu>
        <MenubarTrigger>View</MenubarTrigger>
        <MenubarContent>
          <MenubarCheckboxItem>Always Show Bookmarks Bar</MenubarCheckboxItem>
          <MenubarCheckboxItem checked>
            Always Show Full URLs
          </MenubarCheckboxItem>
          <MenubarSeparator />
          <MenubarItem inset>
            Reload <MenubarShortcut>⌘R</MenubarShortcut>
          </MenubarItem>
          <MenubarItem disabled inset>
            Force Reload <MenubarShortcut>⇧⌘R</MenubarShortcut>
          </MenubarItem>
          <MenubarSeparator />
          <MenubarItem inset>Toggle Fullscreen</MenubarItem>
          <MenubarSeparator />
          <MenubarItem inset>Hide Sidebar</MenubarItem>
        </MenubarContent>
      </MenubarMenu>
      <MenubarMenu>
        <MenubarTrigger>Accounts</MenubarTrigger>
        <MenubarContent>
          <MenubarRadioGroup value="benoit">
            {accounts.map((account) => (
              <MenubarRadioItem key={account.id} value={account.id}>
                {account.name} ({account.email})
              </MenubarRadioItem>
            ))}
          </MenubarRadioGroup>
          <MenubarSeparator />
          <MenubarItem inset>Edit...</MenubarItem>
          <MenubarSeparator />
          <MenubarItem inset>Add Account...</MenubarItem>
        </MenubarContent>
      </MenubarMenu>
    </Menubar>
  );
}
