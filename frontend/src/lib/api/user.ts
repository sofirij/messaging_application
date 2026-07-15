import { User, UserAvatarRequest, UserUsernameRequest } from "@/types/http/user"
import { ApiResult } from "@/types/http/response"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/users"

export async function getUser(): Promise<ApiResult<User>> {
    const res = await fetch(`${api_url}/me`, {
        method: "GET",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
    })

    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null}
    }

    return {error: null, data: json.data}
}

export async function disableAccount(): Promise<ApiResult> {
    const res = await fetch(`${api_url}/me`, {
        method: "DELETE",
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

export async function updateUserAvatar(avatar_url: string): Promise<ApiResult> {
    const req: UserAvatarRequest = { avatar_url }

    const res = await fetch(`${api_url}/me/avatar`, {
        method: "PUT",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(req),
        
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message}
    }

    return {error: null}
}

export async function updateUsername(username: string): Promise<ApiResult> {
    const req: UserUsernameRequest = {username}

    const res = await fetch(`${api_url}/me/username`, {
        method: "PUT",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },
        body: JSON.stringify(req),
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message}
    }

    return {error: null}
}

export async function searchUsername(query: string): Promise<ApiResult<User[]>> {
    const res = await fetch(`${api_url}/search?q=${query}`, {
        method: "PUT",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },
    })

    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null}
    }

    return {error: null, data: json.data}
}

