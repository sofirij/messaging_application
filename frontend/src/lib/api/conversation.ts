import { Conversation, ConversationAddMemberRequest, ConversationAvatarRequest, ConversationCreateRequest, ConversationRenameRequest } from "@/types/http/conversation"
import { AttachmentRequest, Message, MessageCreateRequest } from "@/types/http/message"
import { ApiResult } from "@/types/http/response"
import { fetchWithAuth } from "@/lib/api/fetch"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/conversations"

export async function getConversationsByUserID(): Promise<ApiResult<Conversation[]>> {
    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    return fetchWithAuth(api_url, init, true)
}

export async function createConversation(type: string, name: string | null, user_ids: number[]): Promise<ApiResult<Conversation>> {
    const req : ConversationCreateRequest = { type, name, user_ids }
    const init = {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }
    
    return fetchWithAuth(api_url, init, true)
}

export async function getConversationByID(id: number): Promise<ApiResult<Conversation>> {
    const url = `${api_url}/${id}`

    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    return fetchWithAuth(url, init, true)
}

export async function deleteConversation(id: number): Promise<ApiResult> {
    const url = `${api_url}/${id}`
    const init = {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        }
    }

    return fetchWithAuth(url, init, false)
}

export async function getMessages(conversationID: number, before: number | null, limit: number): Promise<ApiResult<{messages: Message[], next_cursor: number | null, has_more: boolean }>> {
    const url = new URL(`${api_url}/${conversationID}/messages`)
    url.searchParams.set("limit", String(limit))
    if (before !== null) {
        url.searchParams.set("before", String(before))
    }

    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }

    return fetchWithAuth(url, init, true)
}

export async function createMessage(conversation_id: number, reply_to_id: number | null, body: string | null, attachments: AttachmentRequest[]): Promise<ApiResult<Message>> {
    const req : MessageCreateRequest = { reply_to_id, body, attachments }
    const url = `${api_url}/${conversation_id}/messages`
    const init = {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    return fetchWithAuth(url, init, true)
}

export async function addMembers(conversation_id: number, user_ids: number[]): Promise<ApiResult> {
    const req : ConversationAddMemberRequest = { user_ids }
    const url = `${api_url}/${conversation_id}/members`
    const init = {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    return fetchWithAuth(url, init, false)
}

export async function updateConversationName(conversationID: number, name: string): Promise<ApiResult> {
    const req : ConversationRenameRequest = {name}
    const url = `${api_url}/${conversationID}/members`
    const init = {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    return fetchWithAuth(url, init, false)
}

export async function updateConversationAvatar(conversation_id: number, avatar_url: string): Promise<ApiResult> {
    const req : ConversationAvatarRequest = { avatar_url }
    const url = `${api_url}/${conversation_id}/avatar`
    const init = {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    return fetchWithAuth(url, init, false)
}

export async function clearMessages(conversation_id: number): Promise<ApiResult> {
    const url = `${api_url}/${conversation_id}/messages`
    const init ={
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    return fetchWithAuth(url, init, false)
}

export async function removeMember(conversationID: number, userID: number): Promise<ApiResult> {
    const url = `${api_url}/${conversationID}/members/${userID}`
    const init = {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    return fetchWithAuth(url, init, false)
}