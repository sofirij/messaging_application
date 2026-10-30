import { httpProtocol, serverURL } from "@/config/defaults"
import { fetchAction } from "@/lib/api/server/fetch"
import { HTTPError, ErrorDetail } from "@/types/http/error"

const apiVersion = "/api/v1"
const apiURL = httpProtocol+serverURL+apiVersion+"/auth"

export async function refreshToken() {
    const res = await fetchAction(`${apiURL}/refresh`, {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
    }

    return res
}