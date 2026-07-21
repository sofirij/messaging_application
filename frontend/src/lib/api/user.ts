import { User, UserAvatarRequest, UserUsernameRequest } from "@/types/http/user"
import { fetchWithAuth } from "@/lib/api/fetch"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/users"

export async function getUser(): Promise<User> {
    const url = `${api_url}/me`
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

export async function disableAccount(): Promise<void> {
    const url = `${api_url}/me`
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

export async function updateUserAvatar(avatar_url: string | null): Promise<void> {
    const req: UserAvatarRequest = { avatar_url }
    const url = `${api_url}/me/avatar`
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

export async function updateUsername(username: string): Promise<void> {
    const req: UserUsernameRequest = {username}
    const url = `${api_url}/me/username`
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

export async function searchUsername(query: string): Promise<User[]> {
    const url = `${api_url}/search?q=${query}`
    const init = {
        method: "PUT",
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



