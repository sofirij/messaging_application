import { WSProvider } from "@/context/wsProvider"
import { verifySession } from "@/lib/utils/cookies"
import NavBar from "@/components/authedLayout/navBar"
import { userQueryOptions } from "@/query/server/user"
import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query"
import { redirect } from "next/navigation"



export default async function AuthedLayout({children}: {children: React.ReactNode}) {
    const {err} = await verifySession()

    if (err) {
        console.log("redirect to login from auth layout")
        redirect("/login")
    }

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