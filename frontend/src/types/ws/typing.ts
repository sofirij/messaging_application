export type OutboundTypingStopPayload = {
    conversation_id: number
}

export type OutboundTypingStartPayload = {
    conversation_id: number
}

export type InboundTypingStopPayload = {
    conversation_id: number
    user_id: number
}

export type InboundTypingStartPayload = {
    conversation_id: number
    user_id: number
}