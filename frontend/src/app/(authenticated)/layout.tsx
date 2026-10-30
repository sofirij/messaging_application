import { WSProvider } from "@/context/wsProvider"
import { verifyAndRedirect } from "@/lib/utils/cookies"
import NavBar from "@/components/authedLayout/navBar"
import { userQueryOptions } from "@/query/server/user"
import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query"



export default async function AuthedLayout({children}: {children: React.ReactNode}) {
    await verifyAndRedirect()

    const queryClient = new QueryClient()

    await queryClient.query(userQueryOptions)
    
    

    return (
        <HydrationBoundary state={dehydrate(queryClient)}>
            <WSProvider>
                <div className="h-dvh flex flex-col">
                    <NavBar  />
                    <div className="flex-1 min-h-0">
                        {children}
                    </div>
                </div>
            </WSProvider>
        </HydrationBoundary>
        
    )
}