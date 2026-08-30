import AttachmentSelect from "@/components/messages/attachmentSelect"
import { useCreateMessage } from "@/hooks/mutation/message"
import { Dispatch, SetStateAction, useState } from "react"
import { Message } from "@/types/http/message"
import BodyInput from "@/components/messages/bodyInput"

type InputProps = {
    reply: Message | null
    setReply: Dispatch<SetStateAction<Message | null>>
    conversationID: number
    bottomRef: React.RefObject<HTMLDivElement | null>
}

export default function Sender({reply, setReply, conversationID, bottomRef}: InputProps) {
    const [text, setText] = useState("")
    const [attachments, setAttachments] = useState<File[]>([])
    const { handleCreateMessage } = useCreateMessage()

    function processText(text: string): string | null {
        if (text.trim().length === 0) {
            return null
        }
        return text
    }

    function onSuccess() {
        console.log("on success")
        setAttachments([])
        setText("")
        setReply(null)
        bottomRef.current?.scrollIntoView({block: "end"})
    }

    function createMessage() {
        handleCreateMessage(conversationID, processText(text), reply ? reply.id : null, attachments, onSuccess)
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
            <BodyInput setText={setText} text={text} submitHandler={createMessage}/>
            <div>
                <button onClick={createMessage}>Send</button>
            </div>
        </div>
    )
}