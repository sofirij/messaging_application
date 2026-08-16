import { Attachment } from "@/types/http/message"
import Image from "next/image"

type InputProps = {
    attachments: Attachment[]
}

export default function AttachmentList({attachments}: InputProps) {
    return (
        <div>
            {attachments.map(attachment => {
                return <Image key={attachment.id} alt="Sent Image" src={attachment.url} width={40} height={40}/>
            })}
        </div>
    )
}