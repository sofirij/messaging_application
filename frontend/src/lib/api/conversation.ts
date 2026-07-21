import { Conversation, ConversationAddMemberRequest, ConversationAvatarRequest, ConversationCreateRequest, ConversationRenameRequest } from "@/types/http/conversation"
import { AttachmentRequest, Message, MessageCreateRequest, PaginatedMessage } from "@/types/http/message"
import { fetchWithAuth } from "@/lib/api/fetch"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/conversations"

export async function getConversationsByUserID(): Promise<Conversation[]> {
    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    const res = await fetchWithAuth(api_url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new Error(json.error.message)
    }

    return json.data
}

export async function createConversation(type: string, name: string | null, user_ids: number[]): Promise<Conversation> {
    const req : ConversationCreateRequest = { type, name, user_ids }
    const init = {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }
    
    const res = await fetchWithAuth(api_url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new Error(json.error.message)
    }

    return json.data
}

export async function getConversationByID(id: number): Promise<Conversation> {
    const url = `${api_url}/${id}`

    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    const res = await fetchWithAuth(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new Error(json.error.message)
    }

    return json.data
}

export async function deleteConversation(id: number): Promise<void> {
    const url = `${api_url}/${id}`
    const init = {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        }
    }

    const res = await fetchWithAuth(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new Error(json.error.message)
    }
}

export async function getMessages(conversationID: number, before: number | null, limit: number): Promise<PaginatedMessage> {
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

    const res = await fetchWithAuth(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new Error(json.error.message)
    }

    return json.data
}

export async function createMessage(conversation_id: number, reply_to_id: number | null, body: string | null, attachments: AttachmentRequest[]): Promise<Message> {
    const req : MessageCreateRequest = { reply_to_id, body, attachments }
    const url = `${api_url}/${conversation_id}/messages`
    const init = {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    const res = await fetchWithAuth(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new Error(json.error.message)
    }

    return json.data
}

export async function addMembers(conversation_id: number, user_ids: number[]): Promise<void> {
    const req : ConversationAddMemberRequest = { user_ids }
    const url = `${api_url}/${conversation_id}/members`
    const init = {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    const res = await fetchWithAuth(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new Error(json.error.message)
    }
}

export async function updateConversationName(conversationID: number, name: string): Promise<void> {
    const req : ConversationRenameRequest = {name}
    const url = `${api_url}/${conversationID}/members`
    const init = {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    const res = await fetchWithAuth(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new Error(json.error.message)
    }
}

export async function updateConversationAvatar(conversation_id: number, avatar_url: string): Promise<void> {
    const req : ConversationAvatarRequest = { avatar_url }
    const url = `${api_url}/${conversation_id}/avatar`
    const init = {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    const res = await fetchWithAuth(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new Error(json.error.message)
    }
}

export async function clearMessages(conversation_id: number): Promise<void> {
    const url = `${api_url}/${conversation_id}/messages`
    const init ={
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    const res = await fetchWithAuth(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new Error(json.error.message)
    }
}

export async function removeMember(conversationID: number, userID: number): Promise<void> {
    const url = `${api_url}/${conversationID}/members/${userID}`
    const init = {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    const res = await fetchWithAuth(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new Error(json.error.message)
    }
}