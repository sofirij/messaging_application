"use client"

import { conversationMemberQueryOptions, conversationMemberRefetchOptions, conversationQueryOptions, conversationRefetchOptions } from "@/query/conversation"
import { userQueryOptions } from "@/query/user"
import { User } from "@/types/http/user"
import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query"

export function useConversationNew() {
    const queryClient = useQueryClient() 

    function handleConversationNew() {
        queryClient.refetchQueries(conversationRefetchOptions)
    }

    return handleConversationNew
}

export function useConversationDeleted() {
    const queryClient = useQueryClient()

    function handleConversationDeleted(conversationID: number) {
        const oldConversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
            if (!oldConversations) return

            const remainingData = { ...oldConversations.data }
            delete remainingData[conversationID]

            queryClient.setQueryData(conversationQueryOptions.queryKey, {
                data: remainingData,
                order: oldConversations.order.filter(orderID => orderID !== conversationID)
            })
    }

    return handleConversationDeleted
}

export function useMembersAdded() {
    const queryClient = useQueryClient()
    
    function handleMembersAdded(conversationID: number) {
        queryClient.refetchQueries(conversationMemberRefetchOptions(conversationID))
    }

    return handleMembersAdded
}

export function useMemberAdded() {
    const queryClient = useQueryClient()
    const {data: me} = useSuspenseQuery(userQueryOptions)

    function handleMemberAdded(conversationID: number, member: User) {
        const userID = member.id

        if (userID === me.id) {
            queryClient.refetchQueries(conversationRefetchOptions)
            queryClient.refetchQueries(conversationMemberRefetchOptions(conversationID))
        } else {
            const oldMembers = queryClient.getQueryData(conversationMemberQueryOptions(conversationID).queryKey)

            if (!oldMembers) return

            queryClient.setQueryData(conversationMemberQueryOptions(conversationID).queryKey, {
                data: {...oldMembers.data, [userID]: member},
                order: [...oldMembers.order.filter(id => id !== userID), userID].sort((a, b) => a - b)
            })
        }
    }

    return handleMemberAdded
}

export function useMemberRemoved() {
    const queryClient = useQueryClient()
    const {data: me} = useSuspenseQuery(userQueryOptions)

    function handleMemberRemoved(conversationID: number, memberID: number) {
        if (memberID === me.id) {
            queryClient.refetchQueries(conversationRefetchOptions)
            queryClient.refetchQueries(conversationMemberRefetchOptions(conversationID))
        } else {
            const oldMembers = queryClient.getQueryData(conversationMemberQueryOptions(conversationID).queryKey)

            if (!oldMembers) return

            const remainingData = {...oldMembers.data}
            delete remainingData[memberID]

            queryClient.setQueryData(conversationMemberQueryOptions(conversationID).queryKey, {
                data: remainingData,
                order: oldMembers.order.filter(id => id !== memberID)
            })
        }
    }

    return handleMemberRemoved
}

export function useConversationNameUpdated() {
    const queryClient = useQueryClient()

    function handleConversationNameUpdated(conversationID: number, name: string) {
        const oldConversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
        if (!oldConversations) return

        queryClient.setQueryData(conversationQueryOptions.queryKey, {
            data: {...oldConversations.data, [conversationID]: {...oldConversations.data[conversationID], name}},
            order: oldConversations.order
        })
    }

    return handleConversationNameUpdated
}

export function useConversationAvatarUpdated() {
    const queryClient = useQueryClient()

    function handleConversationAvatarUpdated(conversationID: number, url: string | null) {
        const oldConversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
        if (!oldConversations) return
        queryClient.setQueryData(conversationQueryOptions.queryKey, {
            data: {...oldConversations.data, [conversationID]: {...oldConversations.data[conversationID], avatar_url: url}},
            order: oldConversations.order
        })
    }

    return handleConversationAvatarUpdated
}

export function useGroupLeft() {
    const queryClient = useQueryClient()

    function handleGroupLeft(conversationID: number) {
        const oldConversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
        if (!oldConversations) return

        const remainingData = {...oldConversations.data}
        delete remainingData[conversationID]
        queryClient.setQueryData(conversationQueryOptions.queryKey, {
            data: remainingData,
            order: oldConversations.order.filter(orderID => orderID !== conversationID)
        })

        queryClient.removeQueries({queryKey: conversationMemberQueryOptions(conversationID).queryKey, exact: true})
    }

    return handleGroupLeft
} 