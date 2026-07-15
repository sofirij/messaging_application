// outbound
export const EventTypingStart = "typing.start"
export const EventTypingStop = "typing.stop"
export const EventMessageRead = "message.read"

// inbound
export const EventMessageNew      = "message.new"
export const EventMessageDeleted  = "message.deleted"
export const EventUserTypingStart = "user.typing.start"
export const EventUserTypingStop  = "user.typing.stop"
export const EventMessageSeen     = "message.seen"
export const EventUserOnline      = "user.online"
export const EventUserOffline     = "user.offline"
export const EventConversationNew = "conversation.new"
export const EventMessageEdited   = "message.edited"
export const EventMemberAdded     = "member.added"
export const EventMemberRemoved   = "member.removed"
export const EventBroadcastUserStatus = "broadcast_user_status"

export const EventError = "error"

export type Event = {
    Type: string
    Payload: any
}