"use client"
import { MessageEditRequest } from "@/types/http/message"
import { fetchWithCredentials } from "@/lib/api/fetch"
import { HTTPError, ErrorDetail } from "@/types/http/error"
import { httpProtocol, serverURL } from "@/constants/defaults"

const apiVersion = "/api/v1"
const apiURL = httpProtocol+serverURL+apiVersion+"/messages"

export async function updateMessage(id: number, req: MessageEditRequest): Promise<void> {
    const url = `${apiURL}/${id}`
    const init = {
        method: "PATCH",
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

export async function deleteMessage(id: number): Promise<void> {
    const url = `${apiURL}/${id}`
    const init = {
        method: "PATCH",
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