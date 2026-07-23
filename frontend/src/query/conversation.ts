import { getConversationsByUserID } from "@/lib/api/conversation";
import { queryOptions } from "@tanstack/react-query";

export const conversationQueryOptions = queryOptions({
    queryKey: ["conversations"],
    queryFn: async () => await getConversationsByUserID(),
    staleTime: Infinity,
    retry: false
})