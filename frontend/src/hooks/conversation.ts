import { createConversation } from "@/lib/api/conversation";
import { conversationQueryOptions } from "@/query/conversation";
import { Conversation, ConversationCreateRequest } from "@/types/http/conversation";
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

    function handleCreateConversation(type: string, name: string | null, user_ids: number[]) {
        mutation.mutate({type, name, user_ids})
    }

    return { handleCreateConversation, loading: mutation.isPending }
}