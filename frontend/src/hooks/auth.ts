"use client"
import { login, register, logout } from "@/lib/api/auth"
import { validateUsername, validatePassword } from "@/utils/validation"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { userQueryOptions, userRemoveQueryOptions } from "@/context/queryProvider"
import { User } from "@/types/http/user"
import { useRouter } from "next/navigation"

export function useLogin() {
    const queryClient = useQueryClient()
    const router = useRouter()
    
    const mutation = useMutation({
        mutationFn: async ({username, password} : {username: string, password: string}) => {
            validateLogin(username, password)
            return await login(username, password)
        },
        onSuccess: (user: User) => {
            queryClient.setQueryData(userQueryOptions.queryKey, user)
            router.push("/home")
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
    const queryClient = useQueryClient()
    const router = useRouter()

    const mutation = useMutation({
        mutationFn: async ({username, password, confirmPassword} : {username: string, password: string, confirmPassword: string}) => {
            validateRegister(username, password, confirmPassword)
            await register(username, password)
            return await login(username, password)
        },
        onSuccess: (user: User) => {
            queryClient.setQueryData(userQueryOptions.queryKey, user)
            router.push("/home")
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
    const queryClient = useQueryClient()
    const router = useRouter()

    const mutation = useMutation({
        mutationFn: async () => {
            await logout()
        },
        onSuccess: () => {
            queryClient.removeQueries(userRemoveQueryOptions)
            router.push("/login")
        }
    })

    async function handleLogout() {
        mutation.mutate()
    }

    return { handleLogout, loading: mutation.isPending }
}