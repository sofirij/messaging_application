export const conversationGroup = "group"
export const conversationDirect = "direct"
export type ConversationType = typeof conversationGroup  | typeof conversationDirect

export type Conversation = {
    id: number
    name: string
    avatar_url: string | null
    type: ConversationType
    last_message_id: number | null
    created_at: string
    created_by: number
    last_message_read: number | null
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

export type ConversationAddMemberRequest = {
    user_ids: number[]
}