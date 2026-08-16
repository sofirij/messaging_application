import { Message } from "@/types/http/message"

export const conversationGroup = "group"
export const conversationDirect = "direct"
export type ConversationType = typeof conversationGroup  | typeof conversationDirect

export type Conversation = {
    id: number
    name: string
    avatar_url: string | null
    type: ConversationType
    created_at: string
    created_by: number
    last_message_read_by_user: number | null
    last_message_read_in_conversation: number | null
    last_message_sent_in_conversation: Message | null
}

export type ConversationCreateRequest = {
    type: ConversationType
    name: string | null
    user_ids: number[]
}

export type ConversationRenameRequest = {
    name: string
}

export type ConversationAvatarRequest = {
    avatar_url: string | null
}

export type ConversationAddMembersRequest = {
    user_ids: number[]
}