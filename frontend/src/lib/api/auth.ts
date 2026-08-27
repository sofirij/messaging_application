"use client"
import { UserAuthRequest, User } from "@/types/http/user"
import { fetchWithCredentials } from "@/lib/api/fetch"
import { ErrorDetail, HTTPError } from "@/types/http/error"
import { httpProtocol, serverURL } from "@/constants/defaults"

const apiVersion = "/api/v1"
const apiURL = httpProtocol+serverURL+apiVersion+"/auth"

export async function register(req: UserAuthRequest): Promise<void> {
    const res = await fetchWithCredentials(`${apiURL}/register`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json", 
        },
        body: JSON.stringify(req)
    })

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
    }
}

export async function login(req: UserAuthRequest): Promise<User> {
    const res = await fetchWithCredentials(`${apiURL}/login`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(req)
    })

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
}

export async function logout(): Promise<void> {
    const res = await fetchWithCredentials(`${apiURL}/logout`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
    }
}

export async function refreshToken(): Promise<void> {
    const res = await fetchWithCredentials(`${apiURL}/refresh`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        },   
    })

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
    }
}