"use client"

import { conversationQueryOptions } from "@/query/client/conversation"
import { useQueryClient } from "@tanstack/react-query"

export function useReconnected() {
    const queryClient = useQueryClient()

    function handleReconnected() {
        queryClient.refetchQueries({ queryKey: conversationQueryOptions.queryKey })
    }

    return handleReconnected
}