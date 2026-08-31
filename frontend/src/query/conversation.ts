"use client"
import { getConversations, getMembers } from "@/lib/api/conversation";
import { queryOptions } from "@tanstack/react-query";
import { ConversationRecord, MemberRecord } from "@/types/internal/conversation";
import { Conversation } from "@/types/http/conversation";
import { User } from "@/types/http/user";

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
    staleTime: Infinity
})

export const conversationRefetchOptions = {
    queryKey: conversationQueryOptions.queryKey,
    exact: true,
    throwOnError: true,
    type: "all"
} as const

export const conversationMemberQueryOptions = (id: number) => queryOptions({
    queryKey: [...conversationQueryOptions.queryKey, id, "members"],
    queryFn: async (): Promise<MemberRecord> => {
        const members = await getMembers(id)
        return {
            data: members.reduce((aggregate: Record<number, User>, member) => {
                aggregate[member.id] = member
                return aggregate
            }, {}),
            order: members.map(member => member.id)
        }
    },
    staleTime: Infinity,
    retry: false,
})

export const conversationMemberRefetchOptions = (id: number) => {
    return {
        queryKey: conversationMemberQueryOptions(id).queryKey,
        exact: true,
        throwOnError: true,
        type: "all"
    } as const
}