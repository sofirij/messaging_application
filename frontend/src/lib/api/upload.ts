import { Upload } from "@/types/http/upload"
import { fetchWithAuth } from "@/lib/api/fetch"

const apiVersion = "/api/v1"
const apiURL = "http://"+process.env.NEXT_PUBLIC_API_URL+apiVersion

export async function upload(file: File): Promise<Upload> {
    const form = new FormData()
    form.append("file", file)

    const url = `${apiURL}/upload`
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

    const url = `${apiURL}/upload-many`
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