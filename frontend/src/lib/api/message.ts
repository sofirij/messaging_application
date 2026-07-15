import { ApiResult } from "@/types/http/response"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/messages"

export async function editMessage(id: number): Promise<ApiResult> {
    const res = await fetch(`${api_url}/${id}`, {
        method: "PATCH",
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

export async function deleteMessage(id: number): Promise<ApiResult> {
    const res = await fetch(`${api_url}/${id}`, {
        method: "PATCH",
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