"use client"
import Image from "next/image"
import { useLogout } from "@/hooks/auth"
import { useRouter } from "next/navigation"
import { useEffect } from "react"
import { useQuery } from "@tanstack/react-query"
import { userQueryOptions } from "@/query/user"
import { UserProvider, useUserContext } from "@/context/userContext"
import Link from "next/link"
import { defaultAvatarURL } from "@/constants/defaults"

function NavBar() {
    const { handleLogout } = useLogout()
    const user = useUserContext()

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
    const router = useRouter()
    const { isPending, isError, data} = useQuery(userQueryOptions)
    
    useEffect(() => {
        if (isPending) return

        if (isError || !data) {
            router.push("/login")
        }
    }, [isError, router, data, isPending])

    if (isPending || !data) {
        return "Loading..."
    }

    return (
        <UserProvider>
            <NavBar />
            {children}
        </UserProvider>
    )
}