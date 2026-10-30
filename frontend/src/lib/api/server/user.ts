import { httpProtocol, serverURL } from "@/config/defaults"
import { fetchAction } from "@/lib/api/server/fetch"
import { HTTPError, ErrorDetail } from "@/types/http/error"
import { User } from "@/types/http/user"

const apiVersion = "/api/v1"
const apiURL = httpProtocol+serverURL+apiVersion+"/users"

export async function getUser(): Promise<User> {
    const res = await fetchAction(`${apiURL}/me`, {
        method: "GET",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
    })

    if (!res.ok) {
        const json = await res.json()
        throw new HTTPError(json.error as ErrorDetail)
    }

    const json = await res.json()
    return json.data
}