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
        <nav>
            <Image src={avatarURL} alt="User profile picture" width={40} height={10} />
            <Link href="/profile">{username}</Link> {"|"}
            <Link href="/conversations">Conversations</Link> {"|"}
            <button onClick={handleLogout}>Logout</button>
        </nav>
    )
}