"use client"
import { MessageEditRequest } from "@/types/http/message"
import { fetchWithCredentials } from "@/lib/api/fetch"

const apiVersion = "/api/v1"
const apiURL = "http://"+process.env.NEXT_PUBLIC_API_URL+apiVersion+"/messages"

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
        throw new Error(json.error.message)
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
        throw new Error(json.error.message)
    }
}