import { Upload } from "@/types/http/upload"
import { fetchWithAuth } from "@/lib/api/fetch"

const api_version = "/api/v1"
const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version

export async function upload(file: File): Promise<Upload> {
    const form = new FormData()
    form.append("file", file)

    const url = `${api_url}/upload`
    const init = {
        method: "POST",
        body: form
    }

    const res = await fetchWithAuth(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new Error(json.error.message)
    }

    return json.data
}

export async function uploadMany(files: File[]): Promise<Upload[]> {
    const form = new FormData()
    for (const file of files) {
        form.append("file", file)
    }

    const url = `${api_url}/upload-many`
    const init = {
        method: "POST",
        body: form
    }

    const res = await fetchWithAuth(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new Error(json.error.message)
    }

    return json.data
}