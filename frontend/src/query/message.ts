"use client"
import { getMessages, MessageQueryParams } from "@/lib/api/conversation"
import { Message } from "@/types/http/message"
import { infiniteQueryOptions } from "@tanstack/react-query"

const messageLimit = 20

export const messageQueryOptions = (conversationID: number) => {
    return infiniteQueryOptions({
        queryKey: ["conversations", [conversationID], "messages"],
        queryFn: async ({pageParam}) => {
            return await getMessages(conversationID, pageParam)  
        },
        initialPageParam: {before: null, limit: messageLimit, at: null} as MessageQueryParams,
        getNextPageParam: (lastPage) => {
            return lastPage.previous_cursor ? {limit: messageLimit, before: lastPage.previous_cursor, at: null} : undefined
        },
        getPreviousPageParam: (firstPage) => {
            return firstPage.next_cursor ? {limit: messageLimit, before: null, at: firstPage.next_cursor} : undefined
        },
        select: (result): Message[] => {
            return result.pages.flatMap(page => page.messages)
        },
        staleTime: Infinity
    })

}