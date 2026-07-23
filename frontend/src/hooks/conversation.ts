import { createConversation, deleteConversation } from "@/lib/api/conversation";
import { conversationQueryOptions } from "@/query/conversation";
import { Conversation, ConversationCreateRequest, ConversationType } from "@/types/http/conversation";
import { useMutation, useQueryClient } from "@tanstack/react-query";


export function useCreateConversation() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async(req: ConversationCreateRequest) => {
            return await createConversation(req)
        },
        onSuccess: (conversation: Conversation) => {
            const conversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
            if (!conversations) return
            queryClient.setQueryData(conversationQueryOptions.queryKey, [conversation, ...conversations])
        }
    })

    function handleCreateConversation(type: ConversationType, name: string | null, user_ids: number[]) {
        mutation.mutate({type, name, user_ids})
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
        onSuccess: (id: number) => {
            const conversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
            if (!conversations) return
            queryClient.setQueryData(conversationQueryOptions.queryKey, conversations.filter(conversation => conversation.id !== id))
        }
    })

    function handleDeleteConversation(id: number) {
        mutation.mutate(id)
    }

    return { handleDeleteConversation, loading: mutation.isPending }
}