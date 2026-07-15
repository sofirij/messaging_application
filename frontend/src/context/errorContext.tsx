"use client"
import { createContext, useContext, useState } from "react"

type Toast = {
    id: number
    message: string
}

export type ErrorContextType = {
    toasts: Toast[]
    addError: (message: string) => void
    removeError: (id: number) => void
}

export const ErrorContext = createContext<ErrorContextType|null>(null)

export function ErrorProvider({ children }: { children: React.ReactNode }) {
    const [toasts, setToasts] = useState<Toast[]>([])

    function addError(message: string) {
        console.log("adding error")
        const id = Date.now()
        setToasts(prev => [...prev, { id, message }])

        const toastTimeout = 5000
        setTimeout(() => removeError(id), toastTimeout)
    }

    function removeError(id: number) {
        console.log("removing error")
        setToasts(prev => prev.filter(t => t.id !== id))
    }
    
    return (
        <ErrorContext value={{ toasts, addError, removeError }}>
            {children}
        </ErrorContext>
    )
}

export function useErrorContext() {
    const ctx = useContext(ErrorContext)
    if (!ctx) throw new Error("error context should be used within error provider")
    return ctx
}