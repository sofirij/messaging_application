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

    function renderPreview(attachment: Attachment, height: number, width: number) {
        switch (attachment.type) {
            case "image":
                return <Image src={attachment.url} alt={attachment.filename} height={height} width={width} onClick={() => {setFullDisplay(true); setDisplayAttachment(attachment)}}/>
            case "audio":
                return <audio src={uploadURL+attachment.url} onClick={() => {setFullDisplay(true); setDisplayAttachment(attachment)}}/>
            case "video":
                return <video disablePictureInPicture={true} src={uploadURL+attachment.url} height={height} width={width} onClick={() => {setFullDisplay(true); setDisplayAttachment(attachment)}}/>
            case "application":
                return (
                    <div className="relative" onClick={() => {setFullDisplay(true); setDisplayAttachment(attachment)}}>
                        <iframe className="overflow-hidden" src={uploadURL+attachment.url} height={height} width={width}></iframe>
                        <div className="absolute inset-0 z-10" />
                    </div>
                )
        }
    }

    function renderFullDisplay(attachment: Attachment, height: number, width: number) {
        switch (attachment.type) {
            case "image":
                return <Image src={attachment.url} alt={attachment.filename} height={1920} width={1080}/>
            case "audio":
                return <audio controls src={uploadURL+attachment.url}/>
            case "video":
                return <video controls src={uploadURL+attachment.url} height={height} width={width}/>
            case "application":
                return <iframe scrolling="no" src={uploadURL+attachment.url} height={height} width={width}></iframe>
        }
    }

    return (
        <div>
            {attachments.map(attachment => (
                <div key={attachment.id} >
                    {renderPreview(attachment, 70, 70)}
                </div>
            ))}
            <Popup popup={fullDisplay} setPopup={setFullDisplay}>
                {displayAttachment ? renderFullDisplay(displayAttachment, 500, 700) : null}
            </Popup>
        </div>
    )
}