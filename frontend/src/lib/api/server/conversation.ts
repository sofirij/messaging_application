import { Conversation } from "@/types/http/conversation"
import { HTTPError, ErrorDetail } from "@/types/http/error"
import { httpProtocol, serverURL } from "@/config/defaults"
import { fetchAction } from "@/lib/api/server/fetch"
import { MessageQueryParams } from "../client/conversation"
import { PaginatedMessage } from "@/types/http/message"
import { User } from "@/types/http/user"

const apiVersion = "/api/v1"
const apiURL = httpProtocol+serverURL+apiVersion+"/conversations"

export async function getConversations(): Promise<Conversation[]> {
    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        },
    }
    
    const res = await fetchAction(apiURL, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
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

    const res = await fetchAction(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
}

export async function getMembers(conversationID: number): Promise<User[]> {
    const url = `${apiURL}/${conversationID}/members`
    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    const res = await fetchAction(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
}