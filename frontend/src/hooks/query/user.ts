"use client"

import { userQueryOptions } from "@/query/user"
import { User } from "@/types/http/user"
import { useQueryClient } from "@tanstack/react-query"

export function useUnauthenticated() {
    const queryClient = useQueryClient()

    function handleUnauthenticated() {
        queryClient.invalidateQueries({refetchType: "none"})
    }

    return handleUnauthenticated
}

export function useAuthenticated() {
    const queryClient = useQueryClient()

    function handleAuthenticated(user: User) {
        queryClient.setQueryData(userQueryOptions.queryKey, user)
    }

    return handleAuthenticated
}

export function useUsernameUpdated() {
    const queryClient = useQueryClient()

    function handleUsernameUpdated(username: string) {
        const oldUser = queryClient.getQueryData(userQueryOptions.queryKey)
        if (!oldUser) return

        queryClient.setQueryData(userQueryOptions.queryKey, {...oldUser, username})
    }

    return handleUsernameUpdated
}

export function useUserAvatarUpdated() {
    const queryClient = useQueryClient()

    function handleUserAvatarUpdated(url: string | null) {
        const oldUser = queryClient.getQueryData(userQueryOptions.queryKey)
        if (!oldUser) return

        queryClient.setQueryData(userQueryOptions.queryKey, {...oldUser, avatar_url: url})
    }

    return handleUserAvatarUpdated
}