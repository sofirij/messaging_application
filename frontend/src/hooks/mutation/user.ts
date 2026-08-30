"use client"
import { useRouter } from "next/navigation"
import { disableAccount, updateUserAvatar, updateUsername } from "@/lib/api/user"
import { validatePassword, validateUsername } from "@/utils/validation"
import { useMutation } from "@tanstack/react-query"
import { upload } from "@/lib/api/upload"
import { UserAuthRequest, UserUsernameRequest } from "@/types/http/user"
import { login, logout, register } from "@/lib/api/auth"
import { useAuthenticated, useUnauthenticated, useUserAvatarUpdated, useUsernameUpdated } from "@/hooks/query/user"

export function useLogin() {
    const router = useRouter()
    const handleAuthenticated = useAuthenticated()
    
    const mutation = useMutation({
        mutationFn: async (req: UserAuthRequest) => {
            validateLogin(req.username, req.password)
            return await login(req)
        },
        onSuccess: (user) => {
            handleAuthenticated(user)
            router.replace("/home")
        }
    })

    function validateLogin(username: string, password: string) {
        validateUsername(username)
        validatePassword(password)
    }

    async function handleLogin(username: string, password: string) {
        mutation.mutate({username, password})
    }

    return { handleLogin, loading: mutation.isPending }
}

export function useRegister() {
    const router = useRouter()
    const handleAuthenticated = useAuthenticated()

    const mutation = useMutation({
        mutationFn: async ({username, password, confirmPassword} : {username: string, password: string, confirmPassword: string}) => {
            const req : UserAuthRequest = {username, password}
            validateRegister(username, password, confirmPassword)
            await register(req)
            return await login(req)
        },
        onSuccess: (user) => {
            handleAuthenticated(user)
            router.replace("/home")
        }
    })

    function validateRegister(username: string, password: string, confirmPassword: string) {
        validateUsername(username)
        validatePassword(password)
        if (password !== confirmPassword) throw new Error("passwords don't match")
    }

    async function handleRegister(username: string, password: string, confirmPassword: string) {
        mutation.mutate({username, password, confirmPassword})
    }

    return { handleRegister, loading: mutation.isPending }
}

export function useLogout() {
    const router = useRouter()
    const handleUnauthenticated = useUnauthenticated()

    const mutation = useMutation({
        mutationFn: async () => {
            await logout()
        },
        onSuccess: () => {
            handleUnauthenticated()
            router.replace("/login")
        }
    })

    async function handleLogout() {
        mutation.mutate()
    }

    return { handleLogout, loading: mutation.isPending }
}

export function useDisableAccount() {
    const router = useRouter()
    const handleUnauthenticated = useUnauthenticated()

    const mutation = useMutation({
        mutationFn: async () => {
            await disableAccount()
        },
        onSuccess: () => {
            handleUnauthenticated()
            router.replace("/login")
        },
    })

    async function handleDisableAccount() {
        mutation.mutate()
    }

    return { handleDisableAccount, loading: mutation.isPending }
}

export function useUpdateUsername() {
    const handleUsernameUpdated = useUsernameUpdated()

    const mutation = useMutation({
        mutationFn: async (req: UserUsernameRequest) => {
            validateUsername(req.username)
            await updateUsername(req)
            return req.username
        },
        onSuccess: (username) => {
            handleUsernameUpdated(username)
        }
    })
    

    async function handleUpdateUsername(username: string) {
        mutation.mutate({username})
    }

    return { handleUpdateUsername, loading: mutation.isPending }
}

export function useUpdateAvatarURL() {
    const handleUserAvatarUpdated = useUserAvatarUpdated()

    const mutation = useMutation({
        mutationFn: async({file}: {file: File}) => {
            const data = await upload(file)
            await updateUserAvatar({avatar_url: data.url})
            return data.url
        },
        onSuccess: (url) => {
            handleUserAvatarUpdated(url)
        }
    })

    async function handleUpdateAvatarURL(file: File) {
        mutation.mutate({file})
    }

    return {handleUpdateAvatarURL, loading: mutation.isPending}
}

export function useClearAvatarURL() {
    const handleUserAvatarUpdated = useUserAvatarUpdated()

    const  mutation = useMutation({
        mutationFn: async () => {
            await updateUserAvatar({avatar_url: null})
        },
        onSuccess: () => {
            handleUserAvatarUpdated(null)
        }
    })

    async function handleClearAvatarURL() {
        mutation.mutate()
    }

    return { handleClearAvatarURL, loading: mutation.isPending}
}