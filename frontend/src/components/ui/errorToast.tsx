"use client"
import { useErrorContext } from "@/context/errorContext"

export default function ToastContainer() {
    const { toasts, removeError } = useErrorContext()

    return (
        <div className="fixed top-4 right-4 z-[60] flex w-80 max-w-[calc(100vw-2rem)] flex-col gap-2">
            {toasts.map(toast => {
                return (
                    <div key={toast.id} className="flex items-start gap-3 rounded-lg border-l-4 border-red-500 bg-white px-4 py-3 shadow-lg">
                        <p className="flex-1 text-sm text-gray-800">{toast.message}</p>
                        <button className="text-gray-400 hover:text-gray-700" onClick={() => removeError(toast.id)} aria-label="Dismiss">✕</button>
                    </div>
                )
            })}
        </div>
    )
}
