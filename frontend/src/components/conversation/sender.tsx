import AttachmentSelect, { AttachmentPreview } from "@/components/conversation/attachmentSelect"
import { useCreateMessage } from "@/hooks/mutation/message"
import { Dispatch, SetStateAction, useState } from "react"
import { Message } from "@/types/http/message"
import BodyInput from "@/components/conversation/bodyInput"

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
        <div className="flex flex-col">
            {reply && (
                <div className="mx-3 mt-3 flex items-center gap-2 rounded-lg border-l-4 border-brand bg-white px-3 py-2">
                    <div className="min-w-0 flex-1">
                        <p className="text-xs font-semibold text-brand-dark">Replying to</p>
                        <p className="truncate text-sm text-gray-600">{reply.body ? reply.body : "attachment"}</p>
                    </div>
                    <button className="btn-icon h-7 w-7" onClick={() => setReply(null)} aria-label="Cancel reply">✕</button>
                </div>
            )}
            <AttachmentPreview setAttachments={setAttachments} attachments={attachments}/>
            <div className="flex items-center gap-2 p-3">
                <AttachmentSelect setAttachments={setAttachments} attachments={attachments}/>
                <BodyInput setText={setText} text={text} submitHandler={createMessage}/>
                <button className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-brand text-white transition-colors hover:bg-brand-dark" onClick={createMessage} aria-label="Send message">➤</button>
            </div>
        </div>
    )
}
