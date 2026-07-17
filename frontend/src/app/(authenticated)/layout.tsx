"use client"
import { useUserContext } from "@/context/userContext"
import Image from "next/image"
import { useLogout } from "@/hooks/auth"
import { useRouter } from "next/navigation"
import { useEffect } from "react"

const defaultAvatarURL = "fb24fc90-5e53-4972-b880-3edd0f8ccc64.jpg"

function NavBar() {
    const { user } = useUserContext()
    const { handleLogout } = useLogout()

    const avatarURL = user ? user.avatar_url ? user.avatar_url : defaultAvatarURL : defaultAvatarURL
    const username = user?.username

    return (
        <nav>
            <Image src={avatarURL} alt="User profile picture" width={40} height={10} />
            <a href="/profile">{username}</a> {"|"}
            <a href="/conversations">Conversations</a> {"|"}
            <button onClick={async () => await handleLogout()}>Logout</button>
        </nav>
    )
}

export default function AuthedLayout({children}: {children: React.ReactNode}) {
    const router = useRouter()
    const { user, loading } = useUserContext()
    
    useEffect(() => {
        if (loading) return

        if (!user) {
            console.log("navigating to login page")
            router.push("/login")
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [loading])

    return (
        <>
            <NavBar />
            {children}
        </>
    )
}