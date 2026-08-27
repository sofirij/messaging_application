export type ErrorDetail = {
    message: string
    code: number
}

export class HTTPError extends Error {
    code: number
    message: string

    constructor(public error: ErrorDetail) {
        super(error.message)
        this.code = error.code
        this.message = error.message
    }
}