import { Attachment } from "@/types/http/message"
import Image from "next/image"
import { useState } from "react"

type InputProps = {
    attachments: Attachment[]
}

const serverURL = "http://"+process.env.NEXT_PUBLIC_API_URL+"/uploads/"

export default function AttachmentList({attachments}: InputProps) {
    const [fullDisplay, setFullDisplay] = useState(false)
    const [displayAttachment, setDisplayAttachment] = useState<Attachment|null>(null)

    function renderPreview(attachment: Attachment, height: number, width: number) {
        switch (attachment.type) {
            case "image":
                return <Image src={attachment.url} alt={attachment.filename} height={height} width={width} onClick={() => {setFullDisplay(true); setDisplayAttachment(attachment)}}/>
            case "audio":
                return <audio src={serverURL+attachment.url} onClick={() => {setFullDisplay(true); setDisplayAttachment(attachment)}}/>
            case "video":
                return <video disablePictureInPicture={true} src={serverURL+attachment.url} height={height} width={width} onClick={() => {setFullDisplay(true); setDisplayAttachment(attachment)}}/>
            case "application":
                return (
                    <div className="relative" onClick={() => {setFullDisplay(true); setDisplayAttachment(attachment)}}>
                        <iframe className="overflow-hidden" src={serverURL+attachment.url} height={height} width={width}></iframe>
                        <div className="absolute inset-0 z-10" />
                    </div>
                )
        }
    }

    function renderFullDisplay(attachment: Attachment, height: number, width: number) {
        switch (attachment.type) {
            case "image":
                return <Image src={attachment.url} alt={attachment.filename} height={height} width={width}/>
            case "audio":
                return <audio controls src={serverURL+attachment.url}/>
            case "video":
                return <video controls src={serverURL+attachment.url} height={height} width={width}/>
            case "application":
                return <iframe scrolling="no" src={serverURL+attachment.url} height={height} width={width}></iframe>
        }
    }

    return (
        <div>
            {attachments.map(attachment => (
                <div key={attachment.id} >
                    {renderPreview(attachment, 70, 70)}
                </div>
            ))}
            {fullDisplay && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60" onClick={() => setFullDisplay(false)}>
                    <div className="relative bg-white rounded-lg p-6" onClick={(e) => e.stopPropagation()}>
                        {displayAttachment ? renderFullDisplay(displayAttachment, 500, 700) : null}
                    </div>
                </div>
            )}
        </div>
    )
}