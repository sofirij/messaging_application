"use client"
import { Conversation, ConversationAddMembersRequest, ConversationAvatarRequest, ConversationCreateRequest, ConversationRenameRequest } from "@/types/http/conversation"
import { Message, MessageCreateRequest, PaginatedMessage } from "@/types/http/message"
import { User } from "@/types/http/user"
import { fetchWithCredentials } from "@/lib/api/fetch"
import { HTTPError, ErrorDetail } from "@/types/http/error"
import { httpProtocol, serverURL } from "@/constants/defaults"

const apiVersion = "/api/v1"
const apiURL = httpProtocol+serverURL+apiVersion+"/conversations"

export async function getConversations(): Promise<Conversation[]> {
    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        },
    }
    
    const res = await fetchWithCredentials(apiURL, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
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
    
    const res = await fetchWithCredentials(apiURL, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
}

export async function getConversation(id: number): Promise<Conversation> {
    const url = `${apiURL}/${id}`

    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    const res = await fetchWithCredentials(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
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

    const res = await fetchWithCredentials(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
    }
}

export type MessageQueryParams = {
    limit: number
    before: number | null
    at: number | null
}

export async function getMessages(conversationID: number, params: MessageQueryParams): Promise<PaginatedMessage> {
    const url = new URL(`${apiURL}/${conversationID}/messages`)

    url.searchParams.set("limit", String(params.limit))
    if (params.before !== null) {
        url.searchParams.set("before", String(params.before))
    }
    if (params.at !== null) {
        url.searchParams.set("at", String(params.at))
    }

    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }

    const res = await fetchWithCredentials(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
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

    const res = await fetchWithCredentials(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
}

export async function addMembers(conversationID: number, req: ConversationAddMembersRequest): Promise<void> {
    const url = `${apiURL}/${conversationID}/members`
    const init = {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    const res = await fetchWithCredentials(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
    }
}

export async function updateConversationName(conversationID: number, req: ConversationRenameRequest): Promise<void> {
    const url = `${apiURL}/${conversationID}/name`
    const init = {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    const res = await fetchWithCredentials(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
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

    const res = await fetchWithCredentials(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
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
    
    const res = await fetchWithCredentials(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
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
    
    const res = await fetchWithCredentials(url, init)

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
    }
}

export async function getMembers(conversationID: number): Promise<User[]> {
    const url = `${apiURL}/${conversationID}/members`
    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    const res = await fetchWithCredentials(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
}