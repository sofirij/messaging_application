"use client"
import Image from "next/image"
import { useLogout } from "@/hooks/auth"
import { useSuspenseQuery } from "@tanstack/react-query"
import { userQueryOptions } from "@/query/user"
import Link from "next/link"
import { defaultAvatarURL } from "@/constants/defaults"

function NavBar() {
    const { handleLogout } = useLogout()
    const { data: user} = useSuspenseQuery(userQueryOptions)

    const avatarURL = user.avatar_url ?? defaultAvatarURL 
    const username = user.username

    return (
        <nav>
            <Image src={avatarURL} alt="User profile picture" width={40} height={10} />
            <Link href="/profile">{username}</Link> {"|"}
            <Link href="/conversations">Conversations</Link> {"|"}
            <button onClick={async () => await handleLogout()}>Logout</button>
        </nav>
    )
}

export default function AuthedLayout({children}: {children: React.ReactNode}) {
    return (
        <>
            <NavBar />
            {children}
        </>
    )
}