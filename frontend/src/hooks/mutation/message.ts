"use client"
import { clearMessages, createMessage } from "@/lib/api/conversation";
import { deleteMessage, updateMessage } from "@/lib/api/message";
import { uploadMany } from "@/lib/api/upload";
import { AttachmentRequest, MessageCreateRequest, MessageEditRequest } from "@/types/http/message";
import { useMutation } from "@tanstack/react-query";
import { useMessageDeleted, useMessageEdited, useMessageNew, useMessagesCleared } from "@/hooks/query/message"


export function useUpdateMessage() {
    const handleMessageEdited = useMessageEdited()

    const mutation = useMutation({
        mutationFn: async ({messageID, conversationID, req}: {messageID: number, conversationID: number, req: MessageEditRequest}) => {
            await updateMessage(messageID, req)
            return {conversationID, messageID, body: req.body}
        },
        onSuccess: ({conversationID, messageID, body}) => {
            handleMessageEdited(conversationID, messageID, body)
        }
    })

    async function handleUpdateMessage(conversationID: number, messageID: number, body: string) {
        mutation.mutate({messageID, conversationID, req: {body}})
    }

    return {handleUpdateMessage, loading: mutation.isPending}
}

export function useDeleteMessage() {
    const handleMessageDeleted = useMessageDeleted()

    const mutation = useMutation({
        mutationFn: async ({messageID, conversationID}: {conversationID: number, messageID: number}) => {
            await deleteMessage(messageID)
            return {conversationID, messageID}
        },
        onSuccess: ({conversationID, messageID}) => {
            handleMessageDeleted(conversationID, messageID)
        }
    })

    async function handleDeleteMessage(conversationID: number, messageID: number) {
        mutation.mutate({messageID, conversationID})
    }

    return {handleDeleteMessage, loading: mutation.isPending}
}

export function useCreateMessage() {
    const handleMessageNew = useMessageNew()

    const mutation = useMutation({
        mutationFn: async ({conversationID, req}: {conversationID: number, req: MessageCreateRequest}) => {
            const message = await createMessage(conversationID, req)
            return message
        },
        onSuccess: (message) => {
            handleMessageNew(message)
        }
    })

    async function handleCreateMessage(conversationID: number, body: string | null, reply_to_id: number | null, files: File[], onSuccess: () => void) {
        const uploads = files.length > 0 ? await uploadMany(files) : []
        const attachments: AttachmentRequest[] = uploads.map(upload => {
            return {
                url: upload.url,
                filename: upload.filename,
                type: upload.type,
                size: upload.size,
            }
        })
        mutation.mutate({conversationID, req: {body, reply_to_id, attachments}}, {onSuccess})
    }

    return {handleCreateMessage, loading: mutation.isPending}
}

export function useClearMessages() {
    const handleMessagesCleared = useMessagesCleared()

    const mutation = useMutation({
        mutationFn: async({conversationID}: {conversationID: number}) => {
            await clearMessages(conversationID)
            return conversationID
        },
        onSuccess: (conversationID: number) => {
            handleMessagesCleared(conversationID)
        }
    })

    async function handleClearMessages(conversationID: number) {
        mutation.mutate({conversationID})
    }

    return { handleClearMessages, loading: mutation.isPending }
}