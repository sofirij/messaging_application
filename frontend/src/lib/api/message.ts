import { fetchWithAuth } from "@/lib/api/fetch"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/messages"

export async function updateMessage(id: number): Promise<void> {
    const url = `${api_url}/${id}`
    const init = {
        method: "PATCH",
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

export async function deleteMessage(id: number): Promise<void> {
    const url = `${api_url}/${id}`
    const init = {
        method: "PATCH",
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