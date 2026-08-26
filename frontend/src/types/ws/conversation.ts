import { User } from "@/types/http/user"

export type MemberAddedPayload = {
    conversation_id: number
    user: User
}

export type MemberRemovedPayload = {
    conversation_id: number
    user_id: number
}