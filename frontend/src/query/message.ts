"use client"
import { getMessages, MessageQueryParams } from "@/lib/api/conversation"
import { Message } from "@/types/http/message"
import { infiniteQueryOptions } from "@tanstack/react-query"
import { conversationQueryOptions } from "@/query/conversation"

const messageLimit = 20

export const messageQueryOptions = (conversationID: number) => {
    return infiniteQueryOptions({
        queryKey: [...conversationQueryOptions.queryKey, [conversationID], "messages"],
        queryFn: async ({pageParam}) => {
            return await getMessages(conversationID, pageParam)  
        },
        initialPageParam: {before: null, limit: messageLimit, at: null} as MessageQueryParams,
        getNextPageParam: (lastPage) => {
            return lastPage.next_cursor ? {limit: messageLimit, before: null, at: lastPage.next_cursor} : undefined
        },
        getPreviousPageParam: (firstPage) => {
            return firstPage.previous_cursor ? {limit: messageLimit, before: firstPage.previous_cursor, at: null} : undefined
        },
        select: (result): Message[] => {
            return result.pages.flatMap(page => page.messages)
        },
        staleTime: Infinity
    })

}