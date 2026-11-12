import { SetState } from "@/types/react"
import { ReactNode } from "react"

type InputProps = {
    popup: boolean
    setPopup: SetState<boolean>
    children: ReactNode
    wide?: boolean
}

export default function Popup({setPopup, popup, children, wide}: InputProps) {
    return (
        popup && (
            <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" onClick={() => setPopup(false)}>
                <div className={`relative flex max-h-[90dvh] w-full flex-col gap-2 overflow-y-auto rounded-2xl bg-white p-6 shadow-xl ${wide ? "max-w-3xl" : "max-w-md"}`} onClick={(e) => e.stopPropagation()}>
                    <button className="btn-icon absolute top-3 right-3" onClick={() => setPopup(false)} aria-label="Close">✕</button>
                    {children}
                </div>
            </div>
        )
    )
}
