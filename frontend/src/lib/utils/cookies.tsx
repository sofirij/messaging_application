"use server"
import { cookies } from "next/headers"
import { cache } from "react"
import { jwtVerify, decodeJwt } from "jose"
import { jwtPublicKey } from "@/config/auth"

export const verifySession = cache(async () => {
    const accessToken = (await cookies()).get("access_token")?.value ?? ""
    try {
        await jwtVerify(accessToken, jwtPublicKey)
        const claims = decodeJwt(accessToken)
        
        if (claims.sub) return { sub: Number(claims.sub), err: null }
        return { sub: null, err: new Error("missing sub from jwt")}
    } catch {
        return {sub: null, err: new Error("invalid access token")}
    }
})

/*
export async function setCookie(res: Response) {
    const cookieStore = await cookies()
    const setCookies = parseSetCookie(res)

    for (const c of setCookies) {
        cookieStore.set(config(c))
    }
}*/