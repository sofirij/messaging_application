"use client"
import { useRouter } from "next/navigation"
import { login, register, logout } from "@/lib/api/auth"
import { useErrorContext } from "@/context/errorContext"
import { useUserContext } from "@/context/userContext"
import { validateUsername, validatePassword } from "@/utils/validation"


export function useLogin() {
    const router = useRouter()
    const { addError } = useErrorContext()
    const { setUser } = useUserContext()

    function validateLogin(username: string, password: string): string | null {
        const usernameIssue = validateUsername(username)
        if (usernameIssue) {
            return usernameIssue
        }

        return validatePassword(password)
    }

    async function handleLogin(username: string, password: string) {
        const err = validateLogin(username, password)
        if (err) {
            addError(err)
            return
        }

        const res = await login(username, password)

        if (res.error) {
            addError(res.error)
            return
        }

        setUser(res.data)
        console.log("logged in user")
        console.log("navigating to profile page")
        router.push("/profile")
    }

    return { handleLogin }
}

export function useRegister() {
    const router = useRouter()
    const { setUser } = useUserContext()
    const { addError } = useErrorContext()

    function validateRegister(username: string, password: string, confirmPassword: string): string | null {
        const usernameIssue = validateUsername(username)
        if (usernameIssue) {
            return usernameIssue
        }

        const passwordIssue = validatePassword(password)
        if (passwordIssue) {
            return passwordIssue
        }

        if (password !== confirmPassword) {
            return "passwords don't match"
        }

        return null
    }

    async function handleRegister(username: string, password: string, confirmPassword: string) {
        const err = validateRegister(username, password, confirmPassword)
        if (err) {
            addError(err)
            return
        }

        const registerRes = await register(username, password)

        if (registerRes.error) {
            addError(registerRes.error)
            return
        }

        const loginRes = await login(username, password)

        if (loginRes.error) {
            addError(loginRes.error)
            return
        }

        setUser(loginRes.data)
        console.log("logged in user")
        console.log("navigating to the profile page")
        router.push("/profile")
    }

    return { handleRegister }
}

export function useLogout() {
    const router = useRouter()
    const { addError } = useErrorContext()
    const { setUser } = useUserContext()

    async function handleLogout() {
        const res = await logout()

        if (res.error) {
            addError(res.error)
            return
        }

        setUser(null)
        console.log("logged out user")
        console.log("navigating to the login page")
        router.push("/login")
    }

    return { handleLogout }
}