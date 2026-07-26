import { Conversation } from "@/types/http/conversation"

export type ConversationRecord = {
    data: Record<number, Conversation>,
    order: number[]
}