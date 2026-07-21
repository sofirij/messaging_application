import { UserAuthRequest, User } from "@/types/http/user"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/auth"

export async function register(username: string, password: string): Promise<void> {
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
        throw new Error(json.error.message)
    }
}

export async function login(username: string, password: string): Promise<User> {
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
        throw new Error(json.error.message)
    }

    return json.data
}

export async function logout(): Promise<void> {
    const res = await fetch(`${api_url}/logout`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (!res.ok) {
        const json = await res.json()
        throw new Error(json.error.message)
    }
}

export async function refreshToken(): Promise<void> {
    const res = await fetch(`${api_url}/refresh`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },   
    })

    if (!res.ok) {
        const json = await res.json()
        throw new Error(json.error.message)
    }
}