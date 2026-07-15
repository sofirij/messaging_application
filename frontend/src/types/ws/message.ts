export type MessageReadPayload = {
    conversation_id: number
    message_id: number
}

export type MessageEditedPayload = {
    conversation_id: number
    message_id: number
    body: string
}

export type MessageSeenPayload = {
    user_id: number
    conversation_id: number
    message_id: number
}

export type MessageDeletedPayload = {
    conversation_id: number
    message_id: number
}

