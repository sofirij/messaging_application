import type { NextProxy } from "next/server"
import {NextResponse} from "next/server"
import { jwtPublicKey } from "@/config/auth"
import { jwtVerify } from "jose"
import { refreshToken } from "@/lib/api/server/auth"
import { parseSetCookie } from "set-cookie-parser"

export const proxy: NextProxy = async (request) => {
    console.log("proxy running on " + request.url)
    const jwt = request.cookies.get("access_token")?.value ?? ""
    console.log("access token = " + jwt)
    
    const refreshMargin = 30000 //30secs

    try {
        await jwtVerify(jwt, jwtPublicKey, {
            currentDate: new Date(Date.now() + refreshMargin)
        })
        console.log("jwt valid in proxy")
    } catch {
        try {
            console.log("attempting to refresh jwt in proxy")
            const res = await refreshToken()
            console.log("jwt refresh in proxy successful")

            const setCookies = parseSetCookie(res, {map: true})

            const accToken = setCookies["access_token"]
            const refToken = setCookies["refresh_token"]

            request.cookies.set("access_token", accToken.value)
            request.cookies.set("refresh_token", refToken.value)

            console.log("new access token is " + accToken.value)
            console.log("new refresh token is " + refToken.value)

            const response = NextResponse.next({ request: { headers: request.headers }})

            for (const c of res.headers.getSetCookie()) {
                response.headers.append("Set-Cookie", c)
            }

            return response
        } catch {
            console.log("jwt refresh failed in proxy")
            if (request.nextUrl.pathname !== "/login") return NextResponse.redirect(new URL(`/login?nav=${request.nextUrl.pathname}`, request.url))
        }
    }

    return NextResponse.next()
}

export const config = {
    matcher: [
        "/((?!_next/|favicon\.ico|\\login|.*\\.(?:png|jpg|jpeg|gif|webp|svg|ico|css|js|map|woff|woff2|ttf)|\.well-known.*$).*)"
    ]
}