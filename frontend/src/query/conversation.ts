import { getConversations, getMembers } from "@/lib/api/conversation";
import { queryOptions } from "@tanstack/react-query";
import { ConversationRecord } from "@/types/internal/conversation";
import { Conversation } from "@/types/http/conversation";

export const conversationQueryOptions = queryOptions({
    queryKey: ["conversations"],
    queryFn: async (): Promise<ConversationRecord> => {
        const conversations = await getConversations()
        return {
            data: conversations.reduce((aggregate: Record<number, Conversation>, conversation) => {
                aggregate[conversation.id] = conversation
                return aggregate
            }, {}),
            order: conversations.map(conversation => conversation.id)
        }
    },
    staleTime: Infinity,
    retry: false
})

export const conversationMemberQueryOptions = (id: number) => queryOptions({
    queryKey: ["conversations", id, "members"],
    queryFn: async () => await getMembers(id),
    staleTime: Infinity,
    retry: false,
})

export const conversationMemberMutationOptions = (id: number) => {
    return {
        queryKey: conversationMemberQueryOptions(id).queryKey,
        exact: true,
        throwOnError: true,
        type: "all"
    } as const
}