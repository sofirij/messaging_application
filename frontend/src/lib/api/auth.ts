import { UserAuthRequest, User } from "@/types/http/user"
import { ApiResult } from "@/types/http/response"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/auth"

export async function register(username: string, password: string): Promise<ApiResult> {
    const req: UserAuthRequest = {username, password}

    const res = await fetch(`${api_url}/register`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json", 
        },
        body: JSON.stringify(req)
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message, status: res.status}
    }

    return {error: null, status: res.status}
}

export async function login(username: string, password: string): Promise<ApiResult<User>> {
    const req: UserAuthRequest = {username, password}

    const res = await fetch(`${api_url}/login`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(req)
    })

    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null, status: res.status}
    }

    return {error: null, data: json.data, status: res.status}
}

export async function logout(): Promise<ApiResult> {
    const res = await fetch(`${api_url}/logout`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message, status: res.status}
    }

    return {error: null, status: res.status}
}

export async function refreshToken(): Promise<ApiResult> {
    const res = await fetch(`${api_url}/refresh`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },   
    })

    if (!res.ok) {
        const json = await res.json()
        return {error: json.error.message, status: res.status}
    }

    return {error: null, status: res.status}
}