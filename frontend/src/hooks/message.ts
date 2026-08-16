"use client"
import { clearMessages, createMessage } from "@/lib/api/conversation";
import { deleteMessage, updateMessage } from "@/lib/api/message";
import { uploadMany } from "@/lib/api/upload";
import { messageQueryOptions } from "@/query/message";
import { AttachmentRequest, MessageCreateRequest, MessageEditRequest } from "@/types/http/message";
import { useMutation, useQueryClient } from "@tanstack/react-query";


export function useUpdateMessage() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async ({messageID, conversationID, req}: {messageID: number, conversationID: number, req: MessageEditRequest}) => {
            await updateMessage(messageID, req)
            return {conversationID, messageID, body: req.body}
        },
        onSuccess: ({conversationID, messageID, body}) => {
            const old = queryClient.getQueryData(messageQueryOptions(conversationID).queryKey)
            if (!old) return

            queryClient.setQueryData(messageQueryOptions(conversationID).queryKey, {
                ...old,
                pages: old.pages.map(page => {
                    return {
                        ...page,
                        messages: page.messages.map(message => {
                            if (message.id === messageID) {
                                return {...message, body}
                            }
                            return message
                        })
                    }
                })
            })
        }
    })

    function handleUpdateMessage(conversationID: number, messageID: number, body: string) {
        mutation.mutate({messageID, conversationID, req: {body}})
    }

    return {handleUpdateMessage, loading: mutation.isPending}
}

export function useDeleteMessage() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async ({messageID, conversationID}: {conversationID: number, messageID: number}) => {
            await deleteMessage(messageID)
            return {conversationID, messageID}
        },
        onSuccess: ({conversationID, messageID}) => {
            const old = queryClient.getQueryData(messageQueryOptions(conversationID).queryKey)
            if (!old) return

            queryClient.setQueryData(messageQueryOptions(conversationID).queryKey, {
                ...old,
                pages: old.pages.map(page => {
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
    })

    function handleDeleteMessage(conversationID: number, messageID: number) {
        mutation.mutate({messageID, conversationID})
    }

    return {handleDeleteMessage, loading: mutation.isPending}
}

export function useCreateMessage() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async ({conversationID, req}: {conversationID: number, req: MessageCreateRequest}) => {
            const message = await createMessage(conversationID, req)
            return {message, conversationID}
        },
        onSuccess: ({message, conversationID}) => {
            const old = queryClient.getQueryData(messageQueryOptions(conversationID).queryKey)
            if (!old) return

            queryClient.setQueryData(messageQueryOptions(conversationID).queryKey, {
                ...old,
                pages: old.pages.map((page, i) => {
                    if (i === 0) {
                        return {...page, messages: [...page.messages, message]}
                    }
                    return page
                })
            })
        }
    })

    async function handleCreateMessage(conversationID: number, body: string | null, reply_to_id: number | null, files: File[]) {
        const uploads = files.length > 0 ? await uploadMany(files) : []
        const attachments: AttachmentRequest[] = uploads.map(upload => {
            return {
                url: upload.url,
                filename: upload.filename,
                type: upload.type,
                size: upload.size,
            }
        })
        mutation.mutate({conversationID, req: {body, reply_to_id, attachments}})
    }

    return {handleCreateMessage, loading: mutation.isPending}
}

export function useClearMessages() {
    const queryClient = useQueryClient()

    const mutation = useMutation({
        mutationFn: async({conversationID}: {conversationID: number}) => {
            await clearMessages(conversationID)
            return conversationID
        },
        onSuccess: (conversationID: number) => {
            queryClient.refetchQueries({queryKey: messageQueryOptions(conversationID).queryKey, exact: true})
        }
    })

    function handleClearMessages(conversationID: number) {
        mutation.mutate({conversationID})
    }

    return { handleClearMessages, loading: mutation.isPending }
}