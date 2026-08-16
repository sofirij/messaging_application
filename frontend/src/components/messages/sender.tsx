import AttachmentSelect from "@/components/messages/attachmentSelect"
import { useCreateMessage } from "@/hooks/message"
import { Dispatch, SetStateAction, useState } from "react"
import { Message } from "@/types/http/message"
import BodyInput from "@/components/messages/bodyInput"

type InputProps = {
    reply: Message | null
    setReply: Dispatch<SetStateAction<Message | null>>
    conversationID: number
}

export default function Sender({reply, setReply, conversationID}: InputProps) {
    const [text, setText] = useState("")
    const [attachments, setAttachments] = useState<File[]>([])
    const { handleCreateMessage } = useCreateMessage()

    function processText(text: string): string | null {
        if (text.trim().length === 0) {
            return null
        }
        return text
    }

    return (
        <div>
            {reply && (
                <div>
                    <p>{reply.body ? reply.body : "attachment"}</p>
                    <button onClick={() => setReply(null)}>X</button>
                </div>
            )}
            <AttachmentSelect setAttachments={setAttachments} attachments={attachments}/>
            <BodyInput setText={setText}/>
            <div>
                <button onClick={async () => handleCreateMessage(conversationID, processText(text), reply ? reply.id : null, attachments)}>Send</button>
            </div>
        </div>
    )
}