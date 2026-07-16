import { User, UserAvatarRequest, UserUsernameRequest } from "@/types/http/user"
import { ApiResult } from "@/types/http/response"
import { fetchWithAuth } from "@/lib/api/fetch"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/users"

export async function getUser(): Promise<ApiResult<User>> {
    const url = `${api_url}/me`
    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }

    return fetchWithAuth<User>(url, init, true)
}

export async function disableAccount(): Promise<ApiResult> {
    const url = `${api_url}/me`
    const init = {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    return fetchWithAuth(url, init, false)
}

export async function updateUserAvatar(avatar_url: string): Promise<ApiResult> {
    const req: UserAvatarRequest = { avatar_url }
    const url = `${api_url}/me/avatar`
    const init = {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    return fetchWithAuth(url, init, false)
}

export async function updateUsername(username: string): Promise<ApiResult> {
    const req: UserUsernameRequest = {username}
    const url = `${api_url}/me/username`
    const init = {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
    }

    return fetchWithAuth(url, init, false)
}

export async function searchUsername(query: string): Promise<ApiResult<User[]>> {
    const url = `${api_url}/search?q=${query}`
    const init = {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    return fetchWithAuth(url, init, true)
}

