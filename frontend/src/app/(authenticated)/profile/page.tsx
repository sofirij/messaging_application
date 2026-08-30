"use client"
import Form from "@/components/ui/form"
import { useClearAvatarURL, useDisableAccount, useUpdateAvatarURL, useUpdateUsername } from "@/hooks/mutation/user"
import { useRef, useState } from "react"
import Image from "next/image"
import Input from "@/components/ui/input"
import { defaultAvatarURL } from "@/constants/defaults"
import { userQueryOptions } from "@/query/user"
import { useSuspenseQuery } from "@tanstack/react-query"
import Popup from "@/components/ui/popup"


export default function Page() {
    const { handleUpdateUsername } = useUpdateUsername()
    const { handleDisableAccount } = useDisableAccount()
    const [newUsername, setNewUsername] = useState("")
    const [isExpanded, setIsExpanded] = useState(false)
    const { handleUpdateAvatarURL } = useUpdateAvatarURL()
    const { handleClearAvatarURL } = useClearAvatarURL()

    
    const { data: user } = useSuspenseQuery(userQueryOptions)
    const fileInputRef = useRef<HTMLInputElement>(null)
    
    const avatarURL = user.avatar_url ?? defaultAvatarURL
    const username = user.username

    function handleEditClick() {
        fileInputRef.current?.click()
    }

    function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
        const file = e.target.files?.[0]
        if (!file) return
        handleUpdateAvatarURL(file)
    }
    
    return (
        <main>
            <Image src={avatarURL} alt="User profile picture" width={40} height={40} onClick={() => setIsExpanded(true)}/>
            <Form onSubmit={() => handleUpdateUsername(newUsername)}>
                <Input label="username" type="text" onChange={(e) => setNewUsername(e.target.value)} placeholder={username}/>
                <button>Update</button>
            </Form>
            <button onClick={handleDisableAccount}>Disable Account</button>
            <Popup popup={isExpanded} setPopup={setIsExpanded}>
                <Image src={avatarURL} alt="User profile picture" width={100} height={100}/>
                <div className="absolute bottom-4 right-4 flex gap-2"> 
                    <button onClick={handleEditClick}>Edit</button>
                    {user.avatar_url && (
                        <button onClick={handleClearAvatarURL}>Reset</button>
                    )}
                    <input type="file" accept="image/*" ref={fileInputRef} onChange={handleFileChange} style={{display: "none"}}/>
                </div>
            </Popup>
        </main>
    )
}