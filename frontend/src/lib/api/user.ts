import { User, UserAvatarRequest, UserUsernameRequest } from "@/types/http/user"
import { fetchWithAuth } from "@/lib/api/fetch"

const apiVersion = "/api/v1"
const apiURL = "http://"+process.env.NEXT_PUBLIC_API_URL+apiVersion+"/users"

export async function getUser(): Promise<User> {
    const url = `${apiURL}/me`
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
    const url = `${apiURL}/me`
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

export async function updateUserAvatar(req: UserAvatarRequest): Promise<void> {
    const url = `${apiURL}/me/avatar`
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

export async function updateUsername(req: UserUsernameRequest): Promise<void> {
    const url = `${apiURL}/me/username`
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
    if (query === "") return []

    const url = `${apiURL}/search?q=${query}`
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



