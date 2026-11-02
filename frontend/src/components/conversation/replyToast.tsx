import { ReplyMetadata } from "@/types/http/message"

type InputProps = {
    reply: ReplyMetadata
}

export default function ReplyToast({reply}: InputProps) {
    return (
        <div className="rounded-md border-l-4 border-brand bg-black/5 px-2 py-1 text-xs text-gray-600">
            {reply.deleted ? (
                <p className="italic">This message was deleted</p>
            ): (
                <p className="line-clamp-2">{reply.body ? reply.body : "🖼️ Attachment"}</p>
            )}
        </div>
    )
}
