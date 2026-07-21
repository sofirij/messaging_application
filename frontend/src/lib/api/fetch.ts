import { refreshToken } from "@/lib/api/auth";

// prevent concurrent refresh token requests
let refreshPromise: ReturnType<typeof refreshToken> | null = null

function refreshOnce(): ReturnType<typeof refreshToken> {
    if (!refreshPromise) {
        refreshPromise = refreshToken().finally(() => {
            refreshPromise = null
        })
    }

    return refreshPromise
}

export async function fetchWithAuth(url: string | URL, init: RequestInit = {}): Promise<Response> { 
    async function doFetch(): Promise<Response> {
        const res = await fetch(url, {
            ...init,
            credentials: "include"
        })

        return res
    }

    let res = await doFetch()

    if (res.status === 401) {
        await refreshOnce()
        res = await doFetch()
    }

    return res
}

