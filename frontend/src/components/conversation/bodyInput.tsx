import { useWSContext } from "@/context/wsProvider"
import { SetState } from "@/types/react"
import { useRef } from "react"

type InputProps = {
    setText: SetState<string>
    text: string
    submitHandler: () => void
    conversationID: number
}

export default function BodyInput({conversationID, setText, text, submitHandler}: InputProps) {
    const { startTyping, ws } = useWSContext()
    const lastSignal = useRef(0)
    const typingInterval = 1000
    
    function keyHandler(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter") {
            submitHandler()
        }
    }

    function signalTyping() {
        if (!ws.current) return
        const now = Date.now()
        if (now - lastSignal.current < typingInterval) return

        lastSignal.current = now
        startTyping(ws.current, conversationID)
    }



    return (
        <div className="flex-1">
            <input type="text" className="w-full rounded-full border border-gray-200 bg-white px-4 py-2.5 text-sm outline-none focus:border-brand" placeholder="Type a message" onChange={(e) => { setText(e.target.value); signalTyping()}} value={text} onKeyDown={keyHandler}/>
        </div>
    )
}
