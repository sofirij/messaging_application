"use client"
import { useErrorContext } from "@/context/errorContext"

export default function ToastContainer() {
    const { toasts, removeError } = useErrorContext()

    return (
        <div style={{position: "fixed", top: 16, right: 16}}>
            {toasts.map(toast => {
                return (
                    <div key={toast.id}>
                        <p>{toast.message}</p>
                        <button onClick={() => removeError(toast.id)}>x</button>
                    </div>
                )
            })}
        </div>
    )
}