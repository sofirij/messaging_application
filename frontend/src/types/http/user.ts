export type User = {
    id: number
    username: string
    avatar_url: string | null
    is_online: boolean
    last_seen_at: string | null
    created_at: string 
}

export type UserAuthRequest = {
    username: string
    password: string
}

export type UserAvatarRequest = {
    avatar_url: string | null
}

export type UserUsernameRequest = {
    username: string | null
}