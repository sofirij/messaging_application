"use client"
import { QueryCache, QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { useState } from "react"
import { useErrorContext } from "@/context/errorContext"

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