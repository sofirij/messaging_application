import { Cookie } from "set-cookie-parser"

export function config(c: Cookie) {
    return {
        name: c.name,
        value: c.value,
        expires: c.expires,
        maxAge: c.maxAge,
        secure: c.secure,
        httpOnly: c.httpOnly,
        partitioned: false,
        sameSite: "lax"
    } as const
}