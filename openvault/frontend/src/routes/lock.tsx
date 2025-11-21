import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput
} from '@/components/ui/input-group';
import { Spinner } from '@/components/ui/spinner';
import { cn } from '@/lib/utils';
import {
  createFileRoute,
  redirect,
  useNavigate,
  useRouter
} from '@tanstack/react-router';
import { ArrowRight, Lock, Vault } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';

export const Route = createFileRoute('/lock')({
  validateSearch: (search) => ({
    redirect: (search.redirect as string) || '/'
  }),
  beforeLoad: ({ context, search }) => {
    // If already unlocked, redirect to the intended page
    if (context.lock.isLocked === false) {
      throw redirect({
        to: search.redirect
      });
    }
  },
  component: LockScreen
});

function LockScreen() {
  const { lock } = Route.useRouteContext();
  const { redirect } = Route.useSearch();
  const navigate = useNavigate();
  const router = useRouter();
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [shakeLock, setShakeLock] = useState<boolean>(false);

  const handleSubmit = async () => {
    if (password.length === 0) return;
    setLoading(true);
    try {
      await lock.unlock(password);
      console.log('Vault unlocked successfully');
      await router.invalidate({ sync: true });
      navigate({
        to: redirect,
        viewTransition: true
      });
    } catch (err) {
      console.error('Failed to unlock vault:', err);
      toast.error(`Failed to unlock vault: ${err}`);
      setShakeLock(true);
      setTimeout(() => setShakeLock(false), 500);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex h-full w-full flex-col items-center justify-center gap-4">
      {/* Invisible draggable area */}
      <div
        className="fixed top-0 left-0 h-10 w-screen"
        style={{ '--wails-draggable': 'drag' }}
      />
      <h1 className="flex items-center gap-2 text-2xl font-bold text-primary">
        <Vault className="size-8" /> OpenVault
      </h1>
      <p className="text-center text-sm text-gray-500">
        Please unlock to continue
      </p>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          handleSubmit();
        }}
        className="w-1/3"
      >
        <InputGroup>
          <InputGroupInput
            type="password"
            placeholder="Enter your password"
            value={password}
            disabled={loading}
            autoFocus
            onChange={(e) => setPassword(e.target.value)}
          />
          <InputGroupAddon>
            <Lock
              className={cn(
                'transition-all',
                shakeLock ? 'animate-shake text-destructive' : ''
              )}
            />
          </InputGroupAddon>
          <InputGroupAddon align="inline-end">
            <InputGroupButton
              type="submit"
              size="icon-xs"
              className="hover:text-primary"
              disabled={password.length === 0}
            >
              {loading ? <Spinner /> : <ArrowRight />}
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </form>
    </div>
  );
}
