"use client"
import { defaultAvatarURL } from "@/config/defaults"
import { useLogout } from "@/hooks/mutation/user"
import { userQueryOptions } from "@/query/client/user"
import { useSuspenseQuery } from "@tanstack/react-query"
import Link from "next/link"
import Image from "next/image"

export default function NavBar() {
    const { handleLogout } = useLogout()
    const { data: user} = useSuspenseQuery(userQueryOptions)

    const avatarURL = user.avatar_url ?? defaultAvatarURL 
    const username = user.username

    return (
        <nav className="flex items-center gap-4 bg-brand-dark px-4 py-2.5 text-white shadow">
            <Link href="/conversations" className="text-lg font-semibold">💬 Chat</Link>
            <div className="ml-auto flex items-center gap-1 sm:gap-3">
                <Link href="/conversations" className="rounded-lg px-3 py-1.5 text-sm hover:bg-white/10">Conversations</Link>
                <Link href="/profile" className="flex items-center gap-2 rounded-lg px-2 py-1 text-sm hover:bg-white/10">
                    <Image src={avatarURL} alt="User profile picture" width={32} height={32} className="avatar h-8 w-8" />
                    <span className="hidden sm:inline">{username}</span>
                </Link>
                <button onClick={handleLogout} className="rounded-lg px-3 py-1.5 text-sm hover:bg-white/10">Logout</button>
            </div>
        </nav>
    )
}
