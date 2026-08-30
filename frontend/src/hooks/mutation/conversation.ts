"use client"
import { addMembers, createConversation, deleteConversation, removeMember, updateConversationAvatar, updateConversationName } from "@/lib/api/conversation";
import { upload } from "@/lib/api/upload";
import { ConversationAddMembersRequest, ConversationAvatarRequest, ConversationCreateRequest, ConversationRenameRequest, ConversationType } from "@/types/http/conversation";
import { Upload } from "@/types/http/upload";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useConversationAvatarUpdated, useConversationDeleted, useConversationNameUpdated, useConversationNew, useGroupLeft, useMemberRemoved, useMembersAdded } from "@/hooks/query/conversation"


export function useCreateConversation() {
    const handleConversationNew = useConversationNew()

    const mutation = useMutation({
        mutationFn: async(req: ConversationCreateRequest) => {
            await createConversation(req)
        },
        onSuccess: () => {
            handleConversationNew()
        }
    })

    async function handleCreateConversation(type: ConversationType, name: string | null, user_ids: number[], onSuccess: () => void) {
        mutation.mutate({type, name, user_ids}, {onSuccess})
    }

    return { handleCreateConversation, loading: mutation.isPending }
}

export function useDeleteConversation() {
    const handleConversationDeleted = useConversationDeleted()

    const mutation = useMutation({
        mutationFn: async(id: number) => {
            await deleteConversation(id)
            return id
        },
        onSuccess: (id) => {
            handleConversationDeleted(id)
        }
    })

    async function handleDeleteConversation(id: number) {
        mutation.mutate(id)
    }

    return { handleDeleteConversation, loading: mutation.isPending }
}

export function useRemoveMember() {
    const handleMemberRemoved = useMemberRemoved()

    const mutation = useMutation({
        mutationFn: async({conversationID, userID}: {conversationID: number, userID: number}) => {
            await removeMember(conversationID, userID)
            return {conversationID, userID}
        },
        onSuccess: ({conversationID, userID}) => {
            handleMemberRemoved(conversationID, userID)
        }
    })

    async function handleRemoveMember(userID: number, conversationID: number) {
        mutation.mutate({conversationID, userID})
    }

    return { handleRemoveMember, loading: mutation.isPending }
}

export function useAddMembers() {
    const handleMembersAdded = useMembersAdded()

    const mutation = useMutation({
        mutationFn: async ({conversationID, req}: {req: ConversationAddMembersRequest, conversationID: number}) => {
            await addMembers(conversationID, req)
            return conversationID
        },
        onSuccess: async (id) => {
            handleMembersAdded(id)
        }
    })

    async function handleAddMembers(conversationID: number, user_ids: number[]) {
        mutation.mutate({conversationID, req: {user_ids}})
    }

    return { handleAddMembers, loading: mutation.isPending}
}

export function useLeaveGroup() {
    const handleGroupLeft = useGroupLeft()
    const router = useRouter()

    const mutation = useMutation({
        mutationFn: async({conversationID, userID}: {conversationID: number, userID: number}) => {
            await removeMember(conversationID, userID)
            return conversationID
        },
        onSuccess: (id) => {
            handleGroupLeft(id)
            router.push("/conversations")
        }
    })

    async function handleLeaveGroup(userID: number, conversationID: number) {
        mutation.mutate({conversationID, userID})
    }

    return { handleLeaveGroup, loading: mutation.isPending }
}

export function useUpdateConversationName() {
    const handleConversationNameUpdated = useConversationNameUpdated()

    const mutation = useMutation({
        mutationFn: async ({req, conversationID} : {conversationID: number, req: ConversationRenameRequest}) => {
            await updateConversationName(conversationID, req)
            return {conversationID, name: req.name}
        },
        onSuccess: ({name, conversationID}) => {
            handleConversationNameUpdated(conversationID, name)
        }
    })

    async function handleUpdateConversationName(conversationID: number, name: string) {
        mutation.mutate({conversationID, req: {name}})
    }

    return { handleUpdateConversationName, loading: mutation.isPending}
}

export function useUpdateConversationAvatar() {
    const handleConversationAvatarUpdated = useConversationAvatarUpdated()

    const mutation = useMutation({
        mutationFn: async ({conversationID, req}: {conversationID: number, req: ConversationAvatarRequest}) => {
            await updateConversationAvatar(conversationID, req)
            return {conversationID, avatarURL: req.avatar_url}
        },
        onSuccess: ({conversationID, avatarURL}) => {
            handleConversationAvatarUpdated(conversationID, avatarURL)
        }
    })

    async function handleUpdateConversationAvatar(conversationID: number, file: File | null) {
        let uploaded: Upload | null = null
        if (file) {
            uploaded = await upload(file)
        }
        mutation.mutate({conversationID: conversationID, req: {avatar_url: uploaded ? uploaded.url : null}})
    }

    return { handleUpdateConversationAvatar, loading: mutation.isPending }
}