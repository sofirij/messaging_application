import { Dispatch, ReactNode, SetStateAction } from "react"

type InputProps = {
    popup: boolean
    setPopup: Dispatch<SetStateAction<boolean>>
    children: ReactNode
}

export default function Popup({setPopup, popup, children}: InputProps) {
    return (
        popup && (
            <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60" onClick={() => setPopup(false)}>
                <div className="relative bg-white rounded-lg p-6 flex flex-col" onClick={(e) => e.stopPropagation()}>
                   {children} 
                </div>
            </div>
        )
    )
}