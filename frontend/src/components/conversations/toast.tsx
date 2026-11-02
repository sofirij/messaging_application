"use client"
import { Conversation, conversationDirect } from "@/types/http/conversation"
import Image from "next/image"
import { useState } from "react"
import { useDeleteConversation, useLeaveGroup } from "@/hooks/mutation/conversation"
import { defaultAvatarURL } from "@/config/defaults"
import { userQueryOptions } from "@/query/client/user"
import { useSuspenseQuery } from "@tanstack/react-query"
import { useRouter } from "next/navigation"
import Popup from "@/components/ui/popup"

type InputProps = {
    conversation: Conversation
}

export default function Toast({conversation}: InputProps) {
    const avatarURL = conversation.avatar_url ?? defaultAvatarURL
    const [clickedOptions, setClickedOptions] = useState(false)
    const { handleDeleteConversation } = useDeleteConversation()
    const { handleLeaveGroup } = useLeaveGroup()
    const { data: me } = useSuspenseQuery(userQueryOptions)
    const router = useRouter()

    let actionMessage: string
    if (!conversation.last_message_sent_in_conversation) {
        actionMessage = "Start the conversation"
    } else {
        if (conversation.last_message_sent_in_conversation.sender_id === me.id) {
            if (!conversation.last_message_read_in_conversation || conversation.last_message_read_in_conversation < conversation.last_message_sent_in_conversation.id) {
                actionMessage = "Sent"
            } else {
                actionMessage = "Opened"
            }
        } else {
            if (!conversation.last_message_read_by_user || conversation.last_message_read_by_user < conversation.last_message_sent_in_conversation.id) {
                actionMessage = "New Message"
            } else {
                actionMessage = "Received"
            }
        }
    }

    const isNew = actionMessage === "New Message"

    return (
        <div className="flex items-center gap-3 border-b border-gray-100 px-4 py-3 transition-colors hover:bg-gray-50">
            <button className="flex min-w-0 flex-1 items-center gap-3 text-left" onClick={() => {router.push(`/conversation/${conversation.id}`)}}>
                <Image src={avatarURL} alt="Conversation avatar" width={48} height={48} className="avatar h-12 w-12"/>
                <div className="min-w-0 flex-1">
                    <p className="truncate font-medium">{conversation.name}</p>
                    <p className={`truncate text-sm ${isNew ? "font-semibold text-brand" : "text-gray-500"}`}>{actionMessage}</p>
                </div>
                {isNew && (
                    <span className="h-2.5 w-2.5 shrink-0 rounded-full bg-brand" />
                )}
            </button>
            <button className="btn-icon" onClick={() => setClickedOptions(true)} aria-label="Conversation options">⋮</button>
            <Popup popup={clickedOptions} setPopup={setClickedOptions}>
                <h2 className="popup-title">{conversation.name}</h2>
                {conversation.type === conversationDirect ? (
                    <button className="btn-danger w-full" onClick={() => handleDeleteConversation(conversation.id)}>Delete conversation</button>
                ): (
                    <button className="btn-danger w-full" onClick={() => handleLeaveGroup(me.id, conversation.id)}>Leave Group</button>
                )}
            </Popup>
        </div>
    )
}
