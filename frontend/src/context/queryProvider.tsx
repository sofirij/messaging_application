"use client"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { useState } from "react"
import { useErrorContext } from "@/context/errorContext"

export function QueryProvider({children}: {children: React.ReactNode}) {
    const { addError } = useErrorContext()

    const [queryClient] = useState(() => new QueryClient({
        defaultOptions: {
            mutations: {
                onError: (e: Error) => {
                    addError(e.message)
                },
                throwOnError: false
            },
            queries: {
                throwOnError: false,
                retry: false
            }
        },
    }))

    return (
        <QueryClientProvider client={queryClient}>
            {children}
        </QueryClientProvider>
    )
}