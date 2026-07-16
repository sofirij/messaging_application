import { Upload } from "@/types/http/upload"
import { ApiResult } from "@/types/http/response"
import { fetchWithAuth } from "@/lib/api/fetch"

const api_url = "http://"+process.env.NEXT_PUBLIC_API_URL+"/api"

export async function upload(file: File): Promise<ApiResult<Upload>> {
    const form = new FormData()
    form.append("file", file)

    const url = `${api_url}/upload`
    const init = {
        method: "POST",
        body: form
    }

    return fetchWithAuth(url, init, true)
}

export async function uploadMany(files: File[]): Promise<ApiResult<Upload[]>> {
    const form = new FormData()
    for (const file of files) {
        form.append("file", file)
    }

    const url = `${api_url}/upload-many`
    const init = {
        method: "POST",
        body: form
    }

    return fetchWithAuth(url, init, true)
}