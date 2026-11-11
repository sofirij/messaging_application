"use client"
import { UserAuthRequest, User } from "@/types/http/user"
import { fetchWithRetry } from "@/lib/api/client/fetch"
import { ErrorDetail, HTTPError, NetworkError } from "@/types/http/error"
import { httpProtocol, serverURL } from "@/config/defaults"

const apiVersion = "/api/v1"
const apiURL = httpProtocol+serverURL+apiVersion+"/auth"

export async function register(req: UserAuthRequest): Promise<void> {
    const res = await fetch(`${apiURL}/register`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json", 
        },
        body: JSON.stringify(req),
        credentials: "include"
    }).catch(() => { throw new NetworkError() })

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
    }
}

export async function login(req: UserAuthRequest): Promise<User> {
    const res = await fetch(`${apiURL}/login`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(req)
    }).catch(() => { throw new NetworkError() })

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
}

export async function logout(): Promise<void> {
    const res = await fetchWithRetry(`${apiURL}/logout`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
    }
}