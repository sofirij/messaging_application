
import { conversationQueryOptions } from "@/query/server/conversation";
import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query";
import { Conversations } from "@/components/conversations/page"
import { verifySession } from "@/lib/utils/cookies";
import { redirect } from "next/navigation";

export default async function Page() {
    const {err} = await verifySession()
    
    if (err) {
        console.log("redirect to login from conversations")
        redirect("/login")
    }
    
    const queryClient = new QueryClient()

    await queryClient.query(conversationQueryOptions)

    return (
        <HydrationBoundary state={dehydrate(queryClient)}>
            <Conversations/>
        </HydrationBoundary>
    )
}

