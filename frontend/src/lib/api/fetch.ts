export function fetchWithCredentials(url: string | URL | Request, init: RequestInit | undefined): Promise<Response>{
    return fetch(url, {...init, credentials: "include"})
}