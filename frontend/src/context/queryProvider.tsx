"use client"
import { QueryCache, QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { useState } from "react"
import { useErrorContext } from "@/context/errorContext"

export function QueryProvider({children}: {children: React.ReactNode}) {
    const { addError } = useErrorContext()

    const [queryClient] = useState(() => new QueryClient({
        defaultOptions: {
            queries: {
                retry: false,
            },
            mutations: {
                onError: (error: Error) => {
                    addError(error.message)
                }
            }
        },
        queryCache: new QueryCache({
            onError: (error) => addError(error.message)
        })
    }))

    return (
        <QueryClientProvider client={queryClient}>
            {children}
        </QueryClientProvider>
    )
}