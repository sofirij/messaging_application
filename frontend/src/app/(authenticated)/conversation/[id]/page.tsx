import Conversation from "@/components/conversation/page"
import { verifySession } from "@/lib/utils/cookies"
import { conversationQueryOptions, conversationMemberQueryOptions } from "@/query/server/conversation"
import { messageQueryOptions } from "@/query/server/message"
import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query"
import { notFound, redirect } from "next/navigation"


export default async function Page({params} : {params: Promise<{ id: string }>}) {
    const {err} = await verifySession()
    
    if (err) {
        console.log("redirect to login from conversation")
        redirect("/login")
    }

    
    const { id } = await params
    
    const conversationID = Number(id)

    if (!Number.isInteger(conversationID)) notFound()

    const queryClient = new QueryClient()

    await queryClient.query(conversationQueryOptions)
    await queryClient.infiniteQuery(messageQueryOptions(conversationID))
    await queryClient.query(conversationMemberQueryOptions(conversationID))


    return (
        <HydrationBoundary state={dehydrate(queryClient)}>
            <Conversation conversationID={conversationID}/>
        </HydrationBoundary>
    )
}


