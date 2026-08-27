"use client"
import { useRouter } from "next/navigation"
import { disableAccount, updateUserAvatar, updateUsername } from "@/lib/api/user"
import { validateUsername } from "@/utils/validation"
import { useQueryClient } from "@tanstack/react-query"
import { userQueryOptions } from "@/query/user"
import { useMutation } from "@tanstack/react-query"
import { upload } from "@/lib/api/upload"
import { UserUsernameRequest } from "@/types/http/user"

export function useDisableAccount() {
    const queryClient = useQueryClient()
    const router = useRouter()

    const mutation = useMutation({
        mutationFn: async () => {
            await disableAccount()
        },
        onSuccess: () => {
            queryClient.invalidateQueries({refetchType: 'none'})
            router.replace("/login")
        },
    })

    async function handleDisableAccount() {
        mutation.mutate()
    }

    return { handleDisableAccount, loading: mutation.isPending }
}

export function useUpdateUsername() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async (req: UserUsernameRequest) => {
            validateUsername(req.username)
            await updateUsername(req)
            return req.username
        },
        onSuccess: (username) => {
            const user = queryClient.getQueryData(userQueryOptions.queryKey)
            if (!user) return
            queryClient.setQueryData(userQueryOptions.queryKey, {...user, username})
        }
    })
    

    async function handleUpdateUsername(username: string) {
        mutation.mutate({username})
    }

    return { handleUpdateUsername, loading: mutation.isPending }
}

export function useUpdateAvatarURL() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async({file}: {file: File}) => {
            const data = await upload(file)
            await updateUserAvatar({avatar_url: data.url})
            return data.url
        },
        onSuccess: (url) => {
            const user = queryClient.getQueryData(userQueryOptions.queryKey)
            if (!user) return
            queryClient.setQueryData(userQueryOptions.queryKey, {...user, avatar_url: url})
        }
    })

    async function handleUpdateAvatarURL(file: File) {
        mutation.mutate({file})
    }

    return {handleUpdateAvatarURL, loading: mutation.isPending}
}

export function useClearAvatarURL() {
    const queryClient = useQueryClient()

    const  mutation = useMutation({
        mutationFn: async () => {
            await updateUserAvatar({avatar_url: null})
        },
        onSuccess: () => {
            const user = queryClient.getQueryData(userQueryOptions.queryKey)
            if (!user) return
            queryClient.setQueryData(userQueryOptions.queryKey, {...user, avatar_url: null})
        }
    })

    async function handleClearAvatarURL() {
        mutation.mutate()
    }

    return { handleClearAvatarURL, loading: mutation.isPending}
}