import { Conversation } from "@/types/http/conversation"
import { User } from "@/types/http/user"

export type ConversationRecord = {
    data: Record<number, Conversation>,
    order: number[]
}

export type MemberRecord = {
    data: Record<number, User>,
    order: number[]
}