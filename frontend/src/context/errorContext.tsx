"use client"
import { createContext, useContext, useRef, useState } from "react"

const toastTimeout = 5000

type Toast = {
    id: number
    message: string
}

export type ErrorContextType = {
    toasts: Toast[]
    addError: (message: string) => void
    removeError: (id: number) => void
}

const ErrorContext = createContext<ErrorContextType|null>(null)

export function ErrorProvider({ children }: { children: React.ReactNode }) {
    const [toasts, setToasts] = useState<Toast[]>([])
    const ref = useRef(0) 

    function addError(message: string) {
        const id = ref.current
        ref.current += 1

        setToasts(prev => [...prev, { id: id, message }])
        setTimeout(() => removeError(id), toastTimeout)
    }

    function removeError(id: number) {
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