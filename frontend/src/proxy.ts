import type { NextProxy } from "next/server"
import {NextResponse} from "next/server"
import { jwtPublicKey } from "@/config/auth"
import { jwtVerify } from "jose"
import { refreshToken } from "@/lib/api/server/auth"
import { setCookie } from "@/lib/utils/cookies"

export const proxy: NextProxy = async (request) => {
    console.log("proxy running on " + request.url)
    const jwt = request.cookies.get("access_token")?.value ?? ""

    try {
        await jwtVerify(jwt, jwtPublicKey)
        
    } catch {
        try {
            console.log("attempting to refresh jwt in proxy")
            const res = await refreshToken()
            await setCookie(res)
        } catch {
            console.log("jwt refresh failed in proxy")
            if (request.nextUrl.pathname !== "/login") return NextResponse.redirect(new URL("/login", request.url))
        }
        console.log("jwt refresh in proxy successful")
        return NextResponse.next()
    }

    console.log("jwt valid in proxy")
    return NextResponse.next()
}

export const config = {
    matcher: [
        "/((?!_next/|favicon\.ico|\\login|.*\\.(?:png|jpg|jpeg|gif|webp|svg|ico|css|js|map|woff|woff2|ttf)|\.well-known.*$).*)"
    ]
}