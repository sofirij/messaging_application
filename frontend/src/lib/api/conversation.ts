import { Conversation, ConversationAddMemberRequest, ConversationAvatarRequest, ConversationCreateRequest, ConversationRenameRequest } from "@/types/http/conversation"
import { AttachmentRequest, Message, MessageCreateRequest } from "@/types/http/message"
import { ApiResult } from "@/types/http/response"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/conversations"

export async function getConversationsByUserID(): Promise<ApiResult<Conversation[]>> {
    const res = await fetch(`${api_url}`, {
        method: "GET",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },
    })

    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null}
    }

    return {error: null, data: json.data}
}

export async function createConversation(type: string, name: string | null, user_ids: number[]): Promise<ApiResult<Conversation>> {
    const req : ConversationCreateRequest = { type, name, user_ids }

    const res = await fetch(`${api_url}`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },
        body: JSON.stringify(req)
    })

    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null}
    }

    return {error: null, data: json.data}
}

export async function getConversationByID(id: number): Promise<ApiResult<Conversation>> {
    const res = await fetch(`${api_url}/${id}`, {
        method: "GET",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
    })

    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null}
    }

    return {error: null, data: json.data}
}

export async function deleteConversation(id: number): Promise<ApiResult> {
    const res = await fetch(`${api_url}/${id}`, {
        method: "DELETE",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message}
    }

    return {error: null}
}

export async function getMessages(conversationID: number, before: number | null, limit: number): Promise<ApiResult<{messages: Message[], next_cursor: number | null, has_more: boolean }>> {
    const url = new URL(`${api_url}/${conversationID}/messages`)
    url.searchParams.set("limit", String(limit))
    if (before !== null) {
        url.searchParams.set("before", String(before))
    }

    const res = await fetch(url, {
        method: "GET",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },    
    })

    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null}
    }

    return { error: null, data: {messages: json.data, next_cursor: json.next_cursor, has_more: json.has_more} }
}

export async function createMessage(conversation_id: number, reply_to_id: number | null, body: string | null, attachments: AttachmentRequest[]): Promise<ApiResult<Message>> {
    const req : MessageCreateRequest = { reply_to_id, body, attachments }

    const res = await fetch(`${api_url}/${conversation_id}/messages`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },
        body: JSON.stringify(req)
    })

    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null}
    }

    return {error: null, data: json.data}
}

export async function addMember(conversation_id: number, user_ids: number[]): Promise<ApiResult> {
    const req : ConversationAddMemberRequest = { user_ids }

    const res = await fetch(`${api_url}/${conversation_id}/members`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },
        body: JSON.stringify(req)
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message}
    }

    return {error: null}
}

export async function updateConversationName(conversationID: number, name: string): Promise<ApiResult> {
    const req : ConversationRenameRequest = {name}

    const res = await fetch(`${api_url}/${conversationID}/members`, {
        method: "PUT",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },
        body: JSON.stringify(req)
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message}
    }

    return {error: null}
}

export async function updateConversationAvatar(conversation_id: number, avatar_url: string): Promise<ApiResult> {
    const req : ConversationAvatarRequest = { avatar_url }

    const res = await fetch(`${api_url}/${conversation_id}/avatar`, {
        method: "PUT",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(req)
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message}
    }

    return {error: null}
}

export async function clearMessages(conversation_id: number): Promise<ApiResult> {
    const res = await fetch(`${api_url}/${conversation_id}/messages`, {
        method: "DELETE",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message}
    }

    return {error: null}
}

export async function removeMember(conversationID: number, userID: number): Promise<ApiResult> {
    const res = await fetch(`${api_url}/${conversationID}/members/${userID}`, {
        method: "DELETE",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message}
    }

    return {error: null}
}