"use client"
import { Conversation, conversationDirect } from "@/types/http/conversation"
import Image from "next/image"
import { useState } from "react"
import { useDeleteConversation, useLeaveGroup } from "@/hooks/conversation"
import { defaultAvatarURL } from "@/constants/defaults"
import { userQueryOptions } from "@/query/user"
import { useSuspenseQuery } from "@tanstack/react-query"
import { useRouter } from "next/navigation"

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
    
    return (
        <div>
                <Image src={avatarURL} alt="Conversation avatar" width={20} height={20}/>
                <p>{conversation.name}</p>
                <p onClick={() => {router.push(`/conversation/messages?id=${conversation.id}`)}}>{actionMessage}</p>
                <button onClick={() => setClickedOptions(true)}>Options</button>
                {clickedOptions && (
                    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60" onClick={() => setClickedOptions(false)}>
                        <div className="relative bg-white rounded-lg p-6" onClick={(e) => e.stopPropagation()}>
                            {conversation.type === conversationDirect ? (
                                <button onClick={() => handleDeleteConversation(conversation.id)}>Delete</button>
                            ): (
                                <button onClick={() => handleLeaveGroup(me.id, conversation.id)}>Leave Group</button>
                            )}
                        </div>
                    </div>
                )}
        </div>
    )
}