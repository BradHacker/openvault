import { createFileRoute, Navigate } from '@tanstack/react-router';
import { GetAccounts, GetVaultMetadatas } from '@/wailsjs/go/main/App.js';
import { useAccounts } from '@/state/account';
import { useEffect } from 'react';

export const Route = createFileRoute('/_authenticated/')({
  component: App,
  // beforeLoad: async ({ location }) => {
  //   const isInitialized = await IsInitialized();
  //   // console.log(`beforeLoad: isInitialized = ${isInitialized}`);
  //   if (!isInitialized) {
  //     throw redirect({
  //       to: '/initialize'
  //     });
  //   }
  //   const isLocked = await IsLocked();
  //   // console.log(`beforeLoad: isLocked = ${isLocked}`);
  //   if (isLocked) {
  //     throw redirect({
  //       to: '/lock',
  //       search: {
  //         redirect: location.href
  //       }
  //     });
  //   }
  // },
  loader: async () => {
    try {
      const accounts = await GetAccounts();
      // console.log(`loader: accounts = ${JSON.stringify(accounts)}`);
      const vaultMetas = await GetVaultMetadatas(accounts.map((a) => a.id));
      // console.log(`loader: vaultMetas = ${JSON.stringify(vaultMetas)}`);
      return {
        accounts,
        vaultMetas
      };
    } catch (error) {
      console.error('Error in loader:', error);
      throw new Error(`Failed to load data: ${error}`);
    }
  }
  // errorComponent: ({ error }) => {
  //   console.error(error);
  //   return (
  //     <ErrorScreen
  //       message={error.message}
  //       description={
  //         error.stack ? (
  //           <code className="max-h-48 max-w-9/12 overflow-auto rounded-md bg-card p-2">
  //             {error.stack}
  //           </code>
  //         ) : undefined
  //       }
  //     />
  //   );
  // }
});

function App() {
  const { accounts, vaultMetas } = Route.useLoaderData();
  const { setAccounts } = useAccounts();

  useEffect(() => {
    setAccounts(
      accounts.map((account) => {
        return {
          id: account.id,
          name: `${account.user_first_name} ${account.user_last_name}`,
          email: account.user_email
        };
      })
    );
  }, [accounts, setAccounts]);

  return <Navigate to="/vaults/$id" params={{ id: 'all' }} replace />;

  // return (
  //   <div className="flex justify-center">
  //     <Card className="p-4">
  //       <div className="align-end flex justify-between gap-4">
  //         <div className="flex items-center gap-x-2">
  //           <div
  //             className={cn(
  //               'size-2 animate-pulse rounded-full drop-shadow-glow',
  //               'bg-green-500 drop-shadow-green-500'
  //               // : 'bg-red-500 drop-shadow-red-500'
  //             )}
  //           />
  //           <span className="font-mono text-xs">
  //             {/* {isInitialized ? 'Initialized' : 'Not Initialized'} */}
  //             Initialized
  //           </span>
  //         </div>
  //         <ModeToggle />
  //       </div>
  //       <div className="flex max-w-64 flex-col break-all">
  //         {accounts.map((account) => (
  //           <div key={account.id} className="mb-2">
  //             <span className="font-bold">Account ID:</span> {account.id}
  //           </div>
  //         ))}
  //       </div>
  //       <h2 className="font-bold">Vaults</h2>
  //       <Separator />
  //       <div className="flex max-w-64 flex-col break-all">
  //         {vaultMetas.map((vault) => (
  //           <div key={vault.name} className="mb-2">
  //             <span className="font-bold">{vault.name}</span>
  //           </div>
  //         ))}
  //       </div>
  //       {/* <div className="align-center mb-12 flex justify-center">
  //         <AspectRatio ratio={16 / 9}>
  //           <img src={Logo} alt="Logo" />
  //         </AspectRatio>
  //       </div>
  //       <div className="text-md text-center font-bold">{resultText}</div>
  //       <Input
  //         id="name"
  //         onChange={updateName}
  //         autoComplete="off"
  //         name="input"
  //         type="text"
  //         placeholder="Enter your name"
  //         className="w-[20rem]"
  //       />
  //       <Button variant="outline" onClick={greet}>
  //         Greet
  //       </Button> */}
  //     </Card>
  //   </div>
  // );
}
