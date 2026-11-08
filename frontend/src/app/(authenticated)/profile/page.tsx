"use client"
import Form from "@/components/ui/form"
import { useClearAvatarURL, useDisableAccount, useUpdateAvatarURL, useUpdateUsername } from "@/hooks/mutation/user"
import { useRef, useState } from "react"
import Image from "next/image"
import Input from "@/components/ui/input"
import { defaultAvatarURL } from "@/config/defaults"
import { userQueryOptions } from "@/query/client/user"
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
        <main className="h-full overflow-y-auto p-4">
            <div className="mx-auto flex max-w-md flex-col gap-4">
                <section className="flex flex-col items-center gap-3 rounded-2xl bg-white p-6 shadow-sm">
                    <button onClick={() => setIsExpanded(true)} className="group relative rounded-full" aria-label="View profile picture">
                        <Image src={avatarURL} alt="User profile picture" width={112} height={112} className="avatar h-28 w-28"/>
                        <span className="absolute inset-0 flex items-center justify-center rounded-full bg-black/40 text-sm text-white opacity-0 transition-opacity group-hover:opacity-100">View</span>
                    </button>
                    <h1 className="text-xl font-semibold">{username}</h1>
                </section>

                <section className="rounded-2xl bg-white p-6 shadow-sm">
                    <h2 className="mb-4 font-semibold">Change username</h2>
                    <Form onSubmit={() => handleUpdateUsername(newUsername)}>
                        <Input label="username" type="text" onChange={(e) => setNewUsername(e.target.value)} placeholder={username}/>
                        <button className="btn-primary self-end">Update</button>
                    </Form>
                </section>

                <section className="rounded-2xl bg-white p-6 shadow-sm">
                    <h2 className="font-semibold text-red-600">Danger zone</h2>
                    <p className="mt-1 mb-4 text-sm text-gray-500">Disabling your account will log you out.</p>
                    <button onClick={handleDisableAccount} className="btn-danger">Disable Account</button>
                </section>
            </div>
            <Popup popup={isExpanded} setPopup={setIsExpanded}>
                <h2 className="popup-title">Profile picture</h2>
                <Image src={avatarURL} alt="User profile picture" width={320} height={320} className="mx-auto aspect-square w-full max-w-xs rounded-xl object-cover"/>
                <div className="mt-4 flex justify-end gap-2">
                    {user.avatar_url && (
                        <button onClick={handleClearAvatarURL} className="btn-secondary">Reset</button>
                    )}
                    <button onClick={handleEditClick} className="btn-primary">Change</button>
                    <input type="file" accept="image/*" ref={fileInputRef} onChange={handleFileChange} style={{display: "none"}}/>
                </div>
            </Popup>
        </main>
    )
}
