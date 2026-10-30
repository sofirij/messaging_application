import { cookies } from "next/headers"
import { setCookie } from "@/lib/utils/cookies"


export async function fetchAction(url: string | URL | Request, init: RequestInit | undefined) {
    const cookieStore = await cookies()
    const cookie = cookieStore.toString()

    const res = await fetch(url, {...init, credentials: "include", headers: {...init?.headers, "Cookie": cookie}})

    await setCookie(res)
    
    return res
}