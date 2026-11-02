import { useErrorContext } from "@/context/errorContext"
import { getFileCategory } from "@/lib/utils/file"
import { Dispatch, SetStateAction, useEffect, useRef } from "react"

type InputProps = {
    attachments: File[]
    setAttachments: Dispatch<SetStateAction<File[]>>
}

export function AttachmentPreview({attachments, setAttachments}: InputProps) {
    const previews = attachments.map(attachment => URL.createObjectURL(attachment))

    useEffect(() => {
        return () => {
            previews.forEach(preview => URL.revokeObjectURL(preview))
        }
    }, [previews])



    function renderPreview(file: File, index: number) {
        const category = getFileCategory(file.type)

        switch (category) {
            case "image":
                // eslint-disable-next-line @next/next/no-img-element
                return <img src={previews[index]} alt={file.name} className="h-full w-full object-cover"/>
            case "audio":
                return <div className="flex h-full w-full flex-col items-center justify-center gap-1 p-2 text-center text-xs"><span className="text-2xl">🎵</span><span className="w-full truncate">{file.name}</span></div>
            case "video":
                return <video src={previews[index]} className="h-full w-full object-cover"/>
            case "application":
                return <div className="flex h-full w-full flex-col items-center justify-center gap-1 p-2 text-center text-xs"><span className="text-2xl">📄</span><span className="w-full truncate">{file.name}</span></div>
        }
    }

    if (attachments.length === 0) return null

    return (
        <div className="flex gap-2 overflow-x-auto px-3 pt-3">
            {attachments.map((file, i) => {
                return (
                    <div key={i} className="relative h-20 w-20 shrink-0 overflow-hidden rounded-lg border border-gray-200 bg-white">
                        {renderPreview(file, i)}
                        <button className="absolute top-1 right-1 flex h-5 w-5 items-center justify-center rounded-full bg-black/60 text-[10px] text-white hover:bg-black/80" onClick={() => setAttachments(attachments.filter((_attachment, index) => index !== i))} aria-label={`Remove ${file.name}`}>✕</button>
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
        <div className="relative">
            <button className="btn-icon text-xl" onClick={() => ref.current?.click()} aria-label="Add attachment">+</button>
            <input type="file" accept="application/*,image/*,audio/*,video/*" ref={ref} style={{display: "none"}} onChange={(e) => handleFileChange(e)} multiple/>
            {attachments.length > 0 && (
                <span className="absolute -top-1 -right-1 flex h-4 min-w-4 items-center justify-center rounded-full bg-brand px-1 text-[10px] text-white">{attachments.length}</span>
            )}
        </div>
    )
}
