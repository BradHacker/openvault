import { Spinner } from './ui/spinner';

export function LoadingScreen() {
  return (
    <div className="flex h-full w-full items-center justify-center">
      <Spinner className="size-8" />
    </div>
  );
}
