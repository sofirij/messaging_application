import { useUpdateMessage } from "@/hooks/mutation/message"
import { Dispatch, SetStateAction, useState } from "react"

type InputProps = {
    conversationID: number
    messageID: number
    setEditClicked: Dispatch<SetStateAction<boolean>>
}

export default function EditInput({conversationID, messageID, setEditClicked}: InputProps) {
    const { handleUpdateMessage } = useUpdateMessage()
    const [body, setBody] = useState("")

    return (
        <div>
            <input type="text" onChange={(e) => setBody(e.target.value)} />
            <button onClick={() => {setEditClicked(false); handleUpdateMessage(conversationID, messageID, body)}}>Update</button>
            <button onClick={() => setEditClicked(false)}>X</button>
        </div>
        
    )
}