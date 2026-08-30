"use client"

import { conversationQueryOptions } from "@/query/conversation"
import { messageQueryOptions } from "@/query/message"
import { Message } from "@/types/http/message"
import { useQueryClient } from "@tanstack/react-query"

export function useMessageNew() {
    const queryClient = useQueryClient()

    function handleMessageNew(message: Message) {
        const conversationID = message.conversation_id
        const oldConversations = queryClient.getQueryData(conversationQueryOptions.queryKey)
        if (!oldConversations) return

        queryClient.setQueryData(conversationQueryOptions.queryKey, {
            data: {...oldConversations.data, [conversationID]: {...oldConversations.data[conversationID], last_message_sent_in_conversation: message}},
            order: [conversationID, ...oldConversations.order.filter(id => id !== conversationID)],
        })

        const oldMessages = queryClient.getQueryData(messageQueryOptions(conversationID).queryKey)
        if (!oldMessages) return
        
        queryClient.setQueryData(messageQueryOptions(conversationID).queryKey, {
            ...oldMessages,
            pages: oldMessages.pages.map((page, i) => {
                if (i === oldMessages.pages.length - 1) {
                    // sort in the event of network jitter
                    // ws message can come in before http response
                    return {...page, messages: [...page.messages.filter(msg => msg.id !== message.id), message].sort((a, b) => a.id - b.id)}
                }
                return page
            })
        })
    }

    return handleMessageNew
}

export function useMessageDeleted() {
    const queryClient = useQueryClient()

    function handleMessageDeleted(conversationID: number, messageID: number) {
        const oldMessages = queryClient.getQueryData(messageQueryOptions(conversationID).queryKey)
        if (!oldMessages) return

        queryClient.setQueryData(messageQueryOptions(conversationID).queryKey, {
            ...oldMessages,
            pages: oldMessages.pages.map(page => {
                return {
                    ...page,
                    messages: page.messages.map(message => {
                        if (message.id === messageID) {
                            return {...message, body: null, attachments: [], deleted: true}
                        }
                        return message
                    })
                }
            })
        })
    }

    return handleMessageDeleted
}

export function useMessageEdited() {
    const queryClient = useQueryClient()

    function handleMessageEdited(conversationID: number, messageID: number, body: string) {
        const oldMessages = queryClient.getQueryData(messageQueryOptions(conversationID).queryKey)
        if (!oldMessages) return

        queryClient.setQueryData(messageQueryOptions(conversationID).queryKey, {
            ...oldMessages,
            pages: oldMessages.pages.map(page => {
                return {
                    ...page,
                    messages: page.messages.map(message => {
                        if (message.id === messageID) {
                            return {...message, body: body}
                        }
                        return message
                    })
                }
            })
        })
    }

    return handleMessageEdited
}

export function useMessagesCleared() {
    const queryClient = useQueryClient()

    function handleMessagesCleared(conversationID: number) {
        queryClient.refetchQueries({queryKey: messageQueryOptions(conversationID).queryKey, exact: true})
    }

    return handleMessagesCleared
}