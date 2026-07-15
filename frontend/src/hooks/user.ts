"use client"
import { useRouter } from "next/navigation"
import { disableAccount, updateUsername } from "@/lib/api/user"
import { useErrorContext } from "@/context/errorContext"
import { validateUsername } from "@/utils/validation"
import { useUserContext } from "@/context/userContext"
import { getUser as getUserApi } from "@/lib/api/user"

export function useDisableAccount() {
    const router = useRouter()
    const { setUser } = useUserContext()
    const { addError } = useErrorContext()

    async function handleDisableAccount() {
        const res = await disableAccount()

        if (res.error) {
            addError(res.error)
            return
        }

        setUser(null)
        console.log("disabled account")
        console.log("navigating to login page")
        router.push("/login")
    }

    return { handleDisableAccount }
}

export function useEditUsername() {
    const { addError } = useErrorContext()
    const { user, setUser } = useUserContext()

    async function handleEditUsername(username: string) {
        if (!user) {
            console.log("user is not set")
            return
        }

        if (user.username === username) {
            addError("username unchanged")
            return
        }

        const err = validateUsername(username)
        if (err) {
            addError(err)
            return
        }

        const res = await updateUsername(username)

        if (res.error) {
            addError(res.error)
            return
        }

        setUser(prev => prev ? {...prev, username: username}: null)
    }

    return { handleEditUsername }
}

export function useGetUser() {
    const { addError } = useErrorContext()
    const { setUser } = useUserContext()

    async function getUser() {
        const res = await getUserApi()

        if (res.error) {
            addError(res.error)
            return
        }

        setUser(res.data)
    }

    return { getUser }
}