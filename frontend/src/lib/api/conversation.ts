import { Conversation, ConversationAddMemberRequest, ConversationAvatarRequest, ConversationCreateRequest, ConversationRenameRequest } from "@/types/http/conversation"
import { Message, MessageCreateRequest, PaginatedMessage } from "@/types/http/message"
import { fetchWithAuth } from "@/lib/api/fetch"

const apiVersion = "/api/v1"
const apiURL = "http://"+process.env.NEXT_PUBLIC_API_URL+apiVersion+"/conversations"

export async function getConversationsByUserID(): Promise<Conversation[]> {
    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    const res = await fetchWithAuth(apiURL, init)

    const json = await res.json()

    if (!res.ok) {
        throw new Error(json.error.message)
    }

    return json.data
}

export async function createConversation(req: ConversationCreateRequest): Promise<Conversation> {
    const init = {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }
    
    const res = await fetchWithAuth(apiURL, init)

    const json = await res.json()

    if (!res.ok) {
        throw new Error(json.error.message)
    }

    return json.data
}

export async function getConversationByID(id: number): Promise<Conversation> {
    const url = `${apiURL}/${id}`

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
    const url = `${apiURL}/${id}`
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
    const url = new URL(`${apiURL}/${conversationID}/messages`)
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

export async function createMessage(conversationID: number, req: MessageCreateRequest): Promise<Message> {
    const url = `${apiURL}/${conversationID}/messages`
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

export async function addMembers(conversationID: number, req: ConversationAddMemberRequest): Promise<void> {
    const url = `${apiURL}/${conversationID}/members`
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

export async function updateConversationName(conversationID: number, req: ConversationRenameRequest): Promise<void> {
    const url = `${apiURL}/${conversationID}/members`
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

export async function updateConversationAvatar(conversationID: number, req: ConversationAvatarRequest): Promise<void> {
    const url = `${apiURL}/${conversationID}/avatar`
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
    const url = `${apiURL}/${conversation_id}/messages`
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
    const url = `${apiURL}/${conversationID}/members/${userID}`
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