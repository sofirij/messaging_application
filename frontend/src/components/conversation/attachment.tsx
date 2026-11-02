import { serverURL } from "@/config/defaults"
import { Attachment } from "@/types/http/message"
import Image from "next/image"
import { useState } from "react"
import Popup from "@/components/ui/popup"

type InputProps = {
    attachments: Attachment[]
}

const uploadURL = "http://"+serverURL+"/uploads/"

export default function AttachmentList({attachments}: InputProps) {
    const [fullDisplay, setFullDisplay] = useState(false)
    const [displayAttachment, setDisplayAttachment] = useState<Attachment|null>(null)

    function open(attachment: Attachment) {
        setFullDisplay(true)
        setDisplayAttachment(attachment)
    }

    function renderPreview(attachment: Attachment, height: number, width: number) {
        switch (attachment.type) {
            case "image":
                return (
                    <button onClick={() => open(attachment)} className="block overflow-hidden rounded-md">
                        <Image src={attachment.url} alt={attachment.filename} height={height} width={width} className="h-40 w-40 object-cover transition-opacity hover:opacity-90"/>
                    </button>
                )
            case "audio":
                return <audio controls src={uploadURL+attachment.url} className="h-10 w-60 max-w-full"/>
            case "video":
                return (
                    <button onClick={() => open(attachment)} className="relative block overflow-hidden rounded-md bg-black">
                        <video disablePictureInPicture={true} src={uploadURL+attachment.url} height={height} width={width} className="h-40 w-40 object-cover"/>
                        <span className="absolute inset-0 flex items-center justify-center text-3xl text-white/90">▶</span>
                    </button>
                )
            case "application":
                return (
                    <button onClick={() => open(attachment)} className="flex w-60 max-w-full items-center gap-2 rounded-md bg-black/5 px-3 py-2 text-left text-sm hover:bg-black/10">
                        <span className="text-2xl">📄</span>
                        <span className="min-w-0 flex-1 truncate">{attachment.filename}</span>
                    </button>
                )
        }
    }

    function renderFullDisplay(attachment: Attachment, height: number, width: number) {
        switch (attachment.type) {
            case "image":
                return <Image src={attachment.url} alt={attachment.filename} height={1080} width={1920} className="max-h-[75dvh] w-full rounded-lg object-contain"/>
            case "audio":
                return <audio controls src={uploadURL+attachment.url} className="w-full"/>
            case "video":
                return <video controls src={uploadURL+attachment.url} height={height} width={width} className="max-h-[75dvh] w-full rounded-lg bg-black"/>
            case "application":
                return <iframe src={uploadURL+attachment.url} height={height} width={width} className="h-[75dvh] w-full rounded-lg border border-gray-200"></iframe>
        }
    }

    if (attachments.length === 0) return null

    return (
        <div className="flex flex-wrap gap-1.5">
            {attachments.map(attachment => (
                <div key={attachment.id}>
                    {renderPreview(attachment, 160, 160)}
                </div>
            ))}
            <Popup popup={fullDisplay} setPopup={setFullDisplay} wide>
                {displayAttachment && (
                    <p className="truncate pr-10 text-sm font-medium text-gray-700">{displayAttachment.filename}</p>
                )}
                {displayAttachment ? renderFullDisplay(displayAttachment, 500, 700) : null}
            </Popup>
        </div>
    )
}
