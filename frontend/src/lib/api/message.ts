import { ApiResult } from "@/types/http/response"
import { fetchWithAuth } from "@/lib/api/fetch"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/messages"

export async function editMessage(id: number): Promise<ApiResult> {
    const url = `${api_url}/${id}`
    const init = {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    return fetchWithAuth(url, init, false)
}

export async function deleteMessage(id: number): Promise<ApiResult> {
    const url = `${api_url}/${id}`
    const init = {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    return fetchWithAuth(url, init, false)
}