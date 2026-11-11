import { NetworkError } from "@/types/http/error"
import { cookies } from "next/headers"


export async function fetchAction(url: string | URL | Request, init: RequestInit | undefined) {
    const cookieStore = await cookies()
    const cookie = cookieStore.toString()

    console.log("access token is " + cookieStore.get("access_token")?.value)

    const res = await fetch(url, {...init, credentials: "include", headers: {...init?.headers, "Cookie": cookie}}).catch(() => { throw new NetworkError() })
    
    return res
}