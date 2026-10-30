import Conversation from "@/components/conversation/page"
import { verifyAndRedirect } from "@/lib/utils/cookies"
import { conversationQueryOptions, conversationMemberQueryOptions } from "@/query/server/conversation"
import { messageQueryOptions } from "@/query/server/message"
import { QueryClient } from "@tanstack/react-query"
import { notFound } from "next/navigation"


export default async function Page({params} : {params: Promise<{ id: string }>}) {
    const { id } = await params
    
    const conversationID = Number(id)

    if (!Number.isInteger(conversationID)) notFound()

    await verifyAndRedirect()

    const queryClient = new QueryClient()

    await queryClient.query(conversationQueryOptions)
    await queryClient.infiniteQuery(messageQueryOptions(conversationID))
    await queryClient.query(conversationMemberQueryOptions(conversationID))


    return (
        <Conversation conversationID={conversationID}/>
    )
}


