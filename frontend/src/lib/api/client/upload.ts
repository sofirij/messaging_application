"use client"
import { Upload } from "@/types/http/upload"
import { fetchWithRetry } from "@/lib/api/client/fetch"
import { HTTPError, ErrorDetail } from "@/types/http/error"
import { httpProtocol, serverURL } from "@/config/defaults"

const apiVersion = "/api/v1"
const apiURL = httpProtocol+serverURL+apiVersion

export async function upload(file: File): Promise<Upload> {
    const form = new FormData()
    form.append("file", file)

    const url = `${apiURL}/upload`
    const init = {
        method: "POST",
        body: form
    }

    const res = await fetchWithRetry(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
}

export async function uploadMany(files: File[]): Promise<Upload[]> {
    const form = new FormData()
    for (const file of files) {
        form.append("file", file)
    }

    const url = `${apiURL}/upload-many`
    const init = {
        method: "POST",
        body: form
    }

    const res = await fetchWithRetry(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
}