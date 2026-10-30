
import { conversationQueryOptions } from "@/query/server/conversation";
import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query";
import { Conversations } from "@/components/conversations/page"

import { verifyAndRedirect } from "@/lib/utils/cookies";

export default async function Page() {
    await verifyAndRedirect()

    const queryClient = new QueryClient()

    await queryClient.query(conversationQueryOptions)

    return (
        <HydrationBoundary state={dehydrate(queryClient)}>
            <Conversations/>
        </HydrationBoundary>
    )
}

