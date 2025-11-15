export function ErrorScreen({
  message,
  description
}: {
  message?: string;
  description?: React.ReactNode;
}) {
  return (
    <div className="flex h-full w-full flex-col items-center justify-center">
      <p className="text-center text-red-500">
        {message || 'Unknown Error Occurred'}
      </p>
      {description}
    </div>
  );
}
