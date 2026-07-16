import { ApiResult } from "@/types/http/response";
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

export async function fetchWithAuth<T = void>(url: string | URL, init: RequestInit = {}, parseBody: boolean): Promise<ApiResult<T>> { 
    async function doFetch(): Promise<ApiResult<T>> {
        const res = await fetch(url, {
            ...init,
            credentials: "include"
        })

        if (!res.ok) {
            const json = await res.json()
            return {error: json.error.message, status: res.status} as ApiResult<T>
        }

        if (parseBody) {
            const json = await res.json()
            return {error: null, data: json.data, status: res.status} as ApiResult<T>
        }

        return {error: null, status: res.status} as ApiResult<T>
    }

    const res = await doFetch()

    if (res.status === 401) {
        const refreshRes = await refreshOnce()
        if (refreshRes.error) {
            return {...refreshRes, data: null} as ApiResult<T>
        }

        const redoRes = await doFetch()
        return redoRes as ApiResult<T>
    }

    return res
}

