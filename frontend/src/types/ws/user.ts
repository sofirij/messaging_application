export type UserOnlinePayload = {
    user_id: number
}

export type UserOfflinePayload = {
    user_id: number
    last_seen_at: string
}