export type Conversation = {
    id: number
    name: string
    avatar_url: string | null
    type: 'group' | 'direct'
    last_message_id: number
    created_at: string
    created_by: number
    last_message_read: number | null
    members?: Member[]
}

export type Member = {
    id: number
    username: string
    avatar_url: string | null
    joined_at: string
}

export type ConversationCreateRequest = {
    type: string
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