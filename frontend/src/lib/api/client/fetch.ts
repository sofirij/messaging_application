"use client"
export async function fetchWithCredentials(url: string | URL | Request, init: RequestInit | undefined) {
    const res = await fetch(url, {...init, credentials: "include"})

    return res
}