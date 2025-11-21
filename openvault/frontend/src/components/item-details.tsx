import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem
} from '@/components/ui/dropdown-menu';
import { DropdownMenuTrigger } from '@radix-ui/react-dropdown-menu';
import { Clipboard } from '@wailsio/runtime';
import { ChevronDown, Circle } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Avatar, AvatarFallback } from './ui/avatar';
import { Separator } from './ui/separator';
import Markdown from 'react-markdown';

interface VaultDropdownProps {
  accountName: string;
  vaultName: string;
}

export function VaultDropdown({ accountName, vaultName }: VaultDropdownProps) {
  return (
    <Button
      variant="ghost"
      className="flex h-6 items-center gap-x-2 px-2! py-1 text-xs"
    >
      <Avatar className="size-4 text-xs">
        <AvatarFallback>{accountName.charAt(0)}</AvatarFallback>
      </Avatar>
      <span className="font-medium">{accountName}</span>
      <Separator orientation="vertical" />
      <Avatar className="size-4 text-xs">
        <AvatarFallback>{vaultName.charAt(0)}</AvatarFallback>
      </Avatar>
      <span className="font-medium">{vaultName}</span>
      <ChevronDown className="size-4 text-muted-foreground" />
    </Button>
  );
}

interface DetailsButtonProps {
  label: string;
  content: string;
  conceal?: boolean;
}

export function DetailsButton({ label, content, conceal }: DetailsButtonProps) {
  const [isRevealed, setIsRevealed] = useState(false);

  const copyToClipboard = () => {
    Clipboard.SetText(content).then(
      () => {
        toast.success(`Copied ${label.toLowerCase()} to clipboard`);
      },
      (err) => {
        toast.error(`Failed to copy ${label.toLowerCase()} to clipboard`, {
          description: err.message
        });
      }
    );
  };

  return (
    <div className="group/button flex h-14 items-center overflow-hidden border first-of-type:rounded-t-md last-of-type:rounded-b-md">
      <Button
        variant="default"
        className="h-full grow rounded-none bg-transparent px-4 text-foreground hover:bg-primary/25 active:bg-primary/50"
        onClick={copyToClipboard}
      >
        <div className="flex flex-col items-start gap-y-0">
          <span className="mr-4 text-xs font-medium text-primary">{label}</span>
          <span className="font-normal">
            {conceal && !isRevealed ? (
              <span className="my-1.5 inline-flex items-center justify-start gap-x-1">
                {Array(12)
                  .fill(0)
                  .map((_, i) => (
                    <Circle
                      key={`${label}-conceal-${i}`}
                      className="size-1.5 fill-foreground"
                    />
                  ))}
              </span>
            ) : (
              content
            )}
          </span>
        </div>
        <div className="ml-auto flex h-full items-center justify-center text-primary opacity-0 group-hover/button:opacity-100">
          Copy
        </div>
      </Button>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="default"
            className="group/dropdown h-14 items-center justify-center rounded-none bg-transparent hover:bg-primary/50 active:bg-primary/70 data-open:bg-primary/70"
          >
            <ChevronDown className="size-4 text-foreground opacity-0 transition-opacity group-hover/button:opacity-100 group-data-open/dropdown:opacity-100 focus-within:opacity-100" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent>
          {conceal && (
            <DropdownMenuItem onClick={() => setIsRevealed((prev) => !prev)}>
              {isRevealed ? 'Hide' : 'Reveal'}
            </DropdownMenuItem>
          )}
          <DropdownMenuItem>Show in Large Type</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  );
}

export function DetailsMarkdown({
  label,
  content
}: {
  label: string;
  content: string;
}) {
  return (
    <div className="flex flex-col gap-y-1 px-4">
      <span className="mr-4 text-xs font-medium text-primary">{label}</span>
      <Markdown
        components={{
          h1: (props) => <h1 className="text-xl font-bold" {...props} />,
          h2: (props) => <h2 className="text-lg font-bold" {...props} />,
          h3: (props) => <h3 className="text-base font-bold" {...props} />,
          p: (props) => <p className="mb-2 text-sm last:mb-0" {...props} />,
          a: (props) => <a className="text-primary underline" {...props} />
        }}
      >
        {content}
      </Markdown>
    </div>
  );
}
