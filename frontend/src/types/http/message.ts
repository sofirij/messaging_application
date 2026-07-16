export type Message = {
    id: number
    conversation_id: number
    sender_id: number
    reply_to_id: number | null
    body: string | null
    deleted: boolean
    created_at: string
    attachments: Attachment[]
}

export type Attachment = {
    id: number
    type: string
    url: string
    filename: string
    size: number
}

export type PaginatedMessage = {
    messages: Message[]
    next_cursor: number | null
    has_more: boolean
}

export type MessageCreateRequest = {
    reply_to_id: number | null
    body: string | null
    attachments: AttachmentRequest[]
}

export type AttachmentRequest = {
    url: string
    filename: string
    type: string
    size: string
}

export type MessageEditRequest = {
    body: string
}