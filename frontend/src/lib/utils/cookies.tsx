"use server"
import { cookies } from "next/headers"
import { cache } from "react"
import { jwtVerify, decodeJwt } from "jose"
import { jwtPublicKey } from "@/config/auth"
import { redirect } from "next/navigation"

export const verifySession = cache(async () => {
    const jwt = (await cookies()).get("access_token")?.value ?? ""

    try {
        await jwtVerify(jwt, jwtPublicKey)
        const claims = decodeJwt(jwt)
        
        if (claims.sub) return { sub: Number(claims.sub), err: null }
        return { sub: null, err: new Error("missing sub from jwt")}
    } catch {
        return {sub: null, err: new Error("invalid access token")}
    }
})

export async function verifyAndRedirect() {
    const {err} = await verifySession()

    if (err) {
        redirect("/login")
    }
}