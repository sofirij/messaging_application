"use client"
import Form from "@/components/ui/form"
import { useUserContext } from "@/context/userContext"
import { useDisableAccount, useEditUsername } from "@/hooks/user"
import { useState } from "react"
import Image from "next/image"
import Input from "@/components/ui/input"

const defaultAvatarURL = "fb24fc90-5e53-4972-b880-3edd0f8ccc64.jpg"

export default function Page() {

    const { user } = useUserContext()
    const { handleEditUsername } = useEditUsername()
    const { handleDisableAccount } = useDisableAccount()
    const [newUsername, setNewUsername] = useState("")
    const avatarURL = user ? user.avatar_url ? user.avatar_url : defaultAvatarURL : defaultAvatarURL
    const username = user ? user.username : "not set"
    
    return (
        <main>
            <Image src={avatarURL} alt="User profile picture" width={40} height={10} />
            <Form onSubmit={async () => await handleEditUsername(newUsername)}>
                <Input label="username" type="text" onChange={(e) => setNewUsername(e.target.value)} placeholder={username}/>
                <button>Update</button>
            </Form>
            <button onClick={handleDisableAccount}>Disable Account</button>
        </main>
    )
}