"use client"
import { QueryCache, QueryClient, QueryClientProvider, queryOptions } from "@tanstack/react-query"
import { getUser } from "@/lib/api/user"
import { useState } from "react"
import { useErrorContext } from "@/context/errorContext"
import { getConversationsByUserID } from "@/lib/api/conversation"

export const userQueryOptions = queryOptions({
    queryKey: ["user"],
    queryFn: async () => await getUser(),
    staleTime: Infinity,
    retry: false
})

export const userRemoveQueryOptions = {
    queryKey: userQueryOptions.queryKey,
    exact: true
}

export const conversationQueryOptions = queryOptions({
    queryKey: ["conversations"],
    queryFn: async () => await getConversationsByUserID(),
    staleTime: Infinity,
    retry: false
})

export function QueryProvider({children}: {children: React.ReactNode}) {
    const { addError } = useErrorContext()

    const [queryClient] = useState(() => new QueryClient({
        defaultOptions: {
            mutations: {
                onError: (e: Error) => {
                    addError(e.message)
                }
            }   
        },
        queryCache: new QueryCache({
            onError: (e: Error) => {
                addError(e.message)
            }
        })
    }))

    return (
        <QueryClientProvider client={queryClient}>
            {children}
        </QueryClientProvider>
    )
}