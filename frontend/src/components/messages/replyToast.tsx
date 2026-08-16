import { ReplyMetadata } from "@/types/http/message"

type InputProps = {
    reply: ReplyMetadata
}

export default function ReplyToast({reply}: InputProps) {
    return (
        <div>
            {reply.deleted ? (
                <p>This message was deleted</p>
            ): (
                <p>{reply.body ? reply.body : "🖼️"}</p>
            )}
        </div>
    )
}