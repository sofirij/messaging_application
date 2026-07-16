export type ApiResult<T = void> = T extends void
    ? {error: string | null, status: number}
    : {error: string | null, data: T | null, status: number}