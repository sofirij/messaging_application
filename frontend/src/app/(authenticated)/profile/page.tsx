"use client"
import Form from "@/components/ui/form"
import { useClearAvatarURL, useDisableAccount, useUpdateAvatarURL, useUpdateUsername } from "@/hooks/user"
import { useRef, useState } from "react"
import Image from "next/image"
import Input from "@/components/ui/input"
import { useUserContext } from "@/context/userContext"
import { defaultAvatarURL } from "@/constants/defaults"


export default function Page() {
    const { handleUpdateUsername } = useUpdateUsername()
    const { handleDisableAccount } = useDisableAccount()
    const [newUsername, setNewUsername] = useState("")
    const [isExpanded, setIsExpanded] = useState(false)
    const user = useUserContext()
    const fileInputRef = useRef<HTMLInputElement>(null)
    const { handleUpdateAvatarURL } = useUpdateAvatarURL()
    const { handleClearAvatarURL } = useClearAvatarURL()

    function handleEditClick() {
        fileInputRef.current?.click()
    }

    function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
        const file = e.target.files?.[0]
        if (!file) return
        handleUpdateAvatarURL(file)
    }

    const avatarURL = user.avatar_url ?? defaultAvatarURL
    const username = user.username
    
    return (
        <main>
            <Image src={avatarURL} alt="User profile picture" width={40} height={40} onClick={() => setIsExpanded(true)}/>
            <Form onSubmit={async () => await handleUpdateUsername(newUsername)}>
                <Input label="username" type="text" onChange={(e) => setNewUsername(e.target.value)} placeholder={username}/>
                <button>Update</button>
            </Form>
            <button onClick={handleDisableAccount}>Disable Account</button>
            {isExpanded && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60" onClick={() => setIsExpanded(false)}>
                    <div className="relative bg-white rounded-lg p-6" onClick={(e) => e.stopPropagation()}>
                        <Image src={avatarURL} alt="User profile picture" width={100} height={100}/>
                        <div className="absolute bottom-4 right-4 flex gap-2">
                            <button onClick={handleEditClick}>Edit</button>
                            {user.avatar_url && (
                                <button onClick={handleClearAvatarURL}>Reset</button>
                            )}
                            <input type="file" accept="image/*" ref={fileInputRef} onChange={handleFileChange} style={{display: "none"}}/>
                        </div>
                    </div>
                </div>
            )}
        </main>
    )
}