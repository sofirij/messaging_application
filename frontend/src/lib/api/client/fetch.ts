"use client"

import { httpProtocol, serverURL } from "@/config/defaults"

const refreshURL = httpProtocol+serverURL+"/api/v1/auth/refresh"

let refreshing: Promise<boolean> | null = null

function refresh() {
    refreshing ??= fetch(refreshURL, {
        method: "POST",
        credentials: "include",
    })
    .then(res => res.ok)
    .catch(() => false)
    .finally(() => { refreshing = null})

    return refreshing
}

export async function fetchWithRetry(url: string | URL | Request, init: RequestInit | undefined) {
    const res = await fetch(url, {...init, credentials: "include"})

    if (res.status !== 401) return res

    const ok = await refresh()

    if (!ok) return res

    return fetch(url, {...init, credentials: "include"})
}