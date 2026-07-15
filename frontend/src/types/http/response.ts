export type ApiResult<T = void> = T extends void
    ? {error: string | null}
    : {error: string | null, data: T | null}