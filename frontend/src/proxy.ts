import type { NextProxy } from "next/server"
import {NextResponse} from "next/server"
import { jwtPublicKey } from "@/config/auth"
import { jwtVerify } from "jose"
import { refreshToken } from "@/lib/api/server/auth"
import { parseSetCookie } from "set-cookie-parser"
import { config as cookieConfig } from "@/config/cookie"

export const proxy: NextProxy = async (request) => {
    console.log("proxy running on " + request.url)
    const jwt = request.cookies.get("access_token")?.value ?? ""
    console.log("access token = " + jwt)

    try {
        await jwtVerify(jwt, jwtPublicKey)
    } catch {
        try {
            console.log("attempting to refresh jwt in proxy")
            const res = await refreshToken()

            const setCookies = parseSetCookie(res, {map: true})

            request.cookies.set("access_token", setCookies["access_token"].toString())
            request.cookies.set("refresh_token", setCookies["refresh_token"].toString())

            const response = NextResponse.next() 

            response.cookies.set(cookieConfig(setCookies["access_token"]))
            response.cookies.set(cookieConfig(setCookies["refresh_token"]))

            return response
        } catch {
            console.log("jwt refresh failed in proxy")
            if (request.nextUrl.pathname !== "/login") return NextResponse.redirect(new URL(`/login?nav=${request.nextUrl.pathname}`, request.url))
        }
        console.log("jwt refresh in proxy successful")
    }

    console.log("jwt valid in proxy")

    return NextResponse.next()
}

export const config = {
    matcher: [
        "/((?!_next/|favicon\.ico|\\login|.*\\.(?:png|jpg|jpeg|gif|webp|svg|ico|css|js|map|woff|woff2|ttf)|\.well-known.*$).*)"
    ]
}