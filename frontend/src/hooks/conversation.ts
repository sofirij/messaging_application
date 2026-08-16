"use client"
import { addMembers, createConversation, deleteConversation, removeMember, updateConversationAvatar, updateConversationName } from "@/lib/api/conversation";
import { conversationMemberMutationOptions, conversationMemberQueryOptions, conversationQueryOptions } from "@/query/conversation";
import { ConversationAddMembersRequest, ConversationAvatarRequest, ConversationCreateRequest, ConversationRenameRequest, ConversationType } from "@/types/http/conversation";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";


export function useCreateConversation() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async(req: ConversationCreateRequest) => {
            return await createConversation(req)
        },
        onSuccess: (conversation) => {
            const conversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
            if (!conversations) return
            queryClient.setQueryData(conversationQueryOptions.queryKey, {
                data: {...conversations.data, [conversation.id]: conversation},
                order: [conversation.id, ...conversations.order]
            })
        }
    })

    function handleCreateConversation(type: ConversationType, name: string | null, user_ids: number[], onSuccess: () => void) {
        mutation.mutate({type, name, user_ids}, {onSuccess: () => {
            onSuccess()
        }})
    }

    return { handleCreateConversation, loading: mutation.isPending }
}

export function useDeleteConversation() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async(id: number) => {
            await deleteConversation(id)
            return id
        },
        onSuccess: (id) => {
            const conversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
            if (!conversations) return

            const remainingData = { ...conversations.data }
            delete remainingData[id]

            queryClient.setQueryData(conversationQueryOptions.queryKey, {
                data: remainingData,
                order: conversations.order.filter(orderID => orderID !== id)
            })
        }
    })

    function handleDeleteConversation(id: number) {
        mutation.mutate(id)
    }

    return { handleDeleteConversation, loading: mutation.isPending }
}

export function useRemoveMember() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async({conversationID, userID}: {conversationID: number, userID: number}) => {
            await removeMember(conversationID, userID)
            return {conversationID, userID}
        },
        onSuccess: ({conversationID, userID}) => {
            const old = queryClient.getQueryData(conversationMemberQueryOptions(conversationID).queryKey)
            if (!old) return

            const remainingData = {...old.data}
            delete remainingData[userID]

            queryClient.setQueryData(conversationMemberQueryOptions(conversationID).queryKey, {
                data: remainingData,
                order: old.order.filter(id => id !== userID)
            })
        }
    })

    function handleRemoveMember(userID: number, conversationID: number) {
        mutation.mutate({conversationID, userID})
    }

    return { handleRemoveMember, loading: mutation.isPending }
}

export function useAddMember() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async ({conversationID, req}: {req: ConversationAddMembersRequest, conversationID: number}) => {
            await addMembers(conversationID, req)
            return conversationID
        },
        onSuccess: async (id) => {
            await queryClient.refetchQueries(conversationMemberMutationOptions(id))
        }
    })

    function handleAddMember(conversationID: number, user_ids: number[]) {
        mutation.mutate({conversationID, req: {user_ids}})
    }

    return { handleAddMember, loading: mutation.isPending}
}

export function useLeaveGroup() {
    const queryClient = useQueryClient()
    const router = useRouter()

    const mutation = useMutation({
        mutationFn: async({conversationID, userID}: {conversationID: number, userID: number}) => {
            await removeMember(conversationID, userID)
            return conversationID
        },
        onSuccess: (id) => {
            const conversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
            if (!conversations) return
            const remainingData = { ...conversations.data }
            delete remainingData[id]
            queryClient.setQueryData(conversationQueryOptions.queryKey, {
                data: remainingData,
                order: conversations.order.filter(orderID => orderID !== id)
            })
            router.push("/conversations")
        }
    })

    function handleLeaveGroup(userID: number, conversationID: number) {
        mutation.mutate({conversationID, userID})
    }

    return { handleLeaveGroup, loading: mutation.isPending }
}

export function useUpdateConversationName() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async ({req, conversationID} : {conversationID: number, req: ConversationRenameRequest}) => {
            await updateConversationName(conversationID, req)
            return {conversationID, name: req.name}
        },
        onSuccess: ({name, conversationID}) => {
            const conversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
            if (!conversations) return
            queryClient.setQueryData(conversationQueryOptions.queryKey, {
                data: {...conversations.data, [conversationID]: {...conversations.data[conversationID], name}},
                order: conversations.order
            })
        }
    })

    function handleUpdateConversationName(conversationID: number, name: string) {
        mutation.mutate({conversationID, req: {name}})
    }

    return { handleUpdateConversationName, loading: mutation.isPending}
}

export function useUpdateConversationAvatar() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async ({conversationID, req}: {conversationID: number, req: ConversationAvatarRequest}) => {
            await updateConversationAvatar(conversationID, req)
            return {conversationID, avatarURL: req.avatar_url}
        },
        onSuccess: ({conversationID, avatarURL}) => {
            const conversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
            if (!conversations) return
            queryClient.setQueryData(conversationQueryOptions.queryKey, {
                data: {...conversations.data, [conversationID]: {...conversations.data[conversationID], avatar_url: avatarURL}},
                order: conversations.order
            })
        }
    })

    function handleUpdateConversationAvatar(conversationID: number, req: ConversationAvatarRequest) {
        mutation.mutate({conversationID, req})
    }

    return { handleUpdateConversationAvatar, loading: mutation.isPending }
}