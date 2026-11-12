import { useUpdateMessage } from "@/hooks/mutation/message"
import { SetState } from "@/types/react"
import { useState } from "react"

type InputProps = {
    conversationID: number
    messageID: number
    setEditClicked: SetState<boolean>
}

export default function EditInput({conversationID, messageID, setEditClicked}: InputProps) {
    const { handleUpdateMessage } = useUpdateMessage()
    const [body, setBody] = useState("")

    return (
        <div className="flex flex-col gap-2">
            <input type="text" className="text-field" placeholder="Edit message" onChange={(e) => setBody(e.target.value)} />
            <div className="flex justify-end gap-2">
                <button className="btn-ghost px-3 py-1" onClick={() => setEditClicked(false)}>Cancel</button>
                <button className="btn-primary px-3 py-1" onClick={() => {setEditClicked(false); handleUpdateMessage(conversationID, messageID, body)}}>Update</button>
            </div>
        </div>

    )
}
