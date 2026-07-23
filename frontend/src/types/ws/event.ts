import { Message } from "@/types/http/message"
import { Conversation } from "@/types/http/conversation"
import { MessageDeletedPayload, MessageEditedPayload, MessageReadPayload, MessageSeenPayload } from "@/types/ws/message"
import { InboundTypingStartPayload, InboundTypingStopPayload, OutboundTypingStartPayload, OutboundTypingStopPayload } from "@/types/ws/typing"
import { UserOfflinePayload, UserOnlinePayload } from "@/types/ws/user"
import { MemberAddedPayload, MemberRemovedPayload } from "@/types/ws/conversation"
import { ErrorPayload } from "@/types/ws/error"

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

export const EventError = "error"

export type Event =
    // outbound
    | { type: typeof EventTypingStart, payload: OutboundTypingStartPayload }
    | { type: typeof EventTypingStop, payload: OutboundTypingStopPayload }
    | { type: typeof EventMessageRead, payload: MessageReadPayload}

    // inbound
    | { type: typeof EventMessageNew, payload: Message }
    | { type: typeof EventMessageDeleted, payload: MessageDeletedPayload }
    | { type: typeof EventTypingStart, payload: InboundTypingStartPayload }
    | { type: typeof EventTypingStop, payload: InboundTypingStopPayload }
    | { type: typeof EventMessageSeen, payload: MessageSeenPayload }
    | { type: typeof EventUserOnline, payload: UserOnlinePayload }
    | { type: typeof EventUserOffline, payload: UserOfflinePayload }
    | { type: typeof EventConversationNew, payload: Conversation }
    | { type: typeof EventMessageEdited, payload: MessageEditedPayload }
    | { type: typeof EventMemberAdded, payload: MemberAddedPayload }
    | { type: typeof EventMemberRemoved, payload: MemberRemovedPayload }
    | { type: typeof EventError, payload: ErrorPayload }
