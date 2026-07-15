import { Upload } from "@/types/http/upload"
import { ApiResult } from "@/types/http/response"

const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+"/api"

export async function upload(file: File): Promise<ApiResult<Upload>> {
    const form = new FormData()
    form.append("file", file)

    const res = await fetch(`${api_url}/upload`, {
        method: "POST",
        credentials: "include",
        body: form,
    })
    
    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null}
    }

    return {error: json.error.message, data: json.data}
}

export async function uploadMany(files: File[]): Promise<ApiResult<Upload[]>> {
    const form = new FormData()
    for (const file of files) {
        form.append("file", file)
    }

    const res = await fetch(`${api_url}/upload-many`, {
        method: "POST",
        credentials: "include",
        body: form,
    })

    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null}
    }

    return {error: null, data: json.data}
}