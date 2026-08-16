import { useErrorContext } from "@/context/errorContext"
import { Dispatch, SetStateAction, useEffect, useRef } from "react"

type InputProps = {
    attachments: File[]
    setAttachments: Dispatch<SetStateAction<File[]>>
}

function RenderPreview({attachments, setAttachments}: InputProps) {
    const previews = attachments.map(attachment => URL.createObjectURL(attachment))

    useEffect(() => {
        return () => {
            previews.forEach(preview => URL.revokeObjectURL(preview))
        }
    }, [previews])

    function getFileCategory(type: string) {
        if (type.startsWith("image/")) return "image";
        if (type.startsWith("video/")) return "video";
        if (type.startsWith("audio/")) return "audio";
        return "application"
    }

    function renderPreview(file: File, index: number) {
        const category = getFileCategory(file.type)

        switch (category) {
            case "image":
                // eslint-disable-next-line @next/next/no-img-element
                return <img src={previews[index]} alt={file.name}/>
            case "audio":
                return <audio controls src={previews[index]}/>
            case "video":
                return <video controls src={previews[index]}/>
            default: 
                return <div>📄 {file.name}</div>
        }
    }

    return (
        <div className="grid grid-cols-3 gap-[10px]">
            {attachments.map((file, i) => {
                return (
                    <div key={i}>
                        <button onClick={() => setAttachments(attachments.filter((_attachment, index) => index !== i))}>X</button>
                        {renderPreview(file, i)}
                    </div>
                )
            })}
        </div>
    )
}

export default function AttachmentSelect({setAttachments, attachments}: InputProps) {
    const ref = useRef<HTMLInputElement>(null)
    const { addError }= useErrorContext()

    function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
        const filelist = e.target.files
        if (!filelist) {
            return
        }

        const files = Array.from(filelist)
        if (attachments.length + files.length > 10) {
            addError("cannot upload more than 10 attachments")
            return
        }

        for (let i = 0; i < files.length; i++) {
            if (files[i].size > 100 * 1024 * 1024) {
                addError("media limit is 100MB per file")
                return
            }
        }
        setAttachments([...attachments, ...files])
    }

    return (
        <div>
            <button onClick={() => ref.current?.click()}>+</button>
            <input type="file" accept="application/*,image/*,audio/*,video/*" ref={ref} style={{display: "none"}} onChange={(e) => handleFileChange(e)} multiple/>
            <p>{attachments.length} selected</p>
            <RenderPreview attachments={attachments} setAttachments={setAttachments}/>
        </div>
    )
}