"use client"
import { Conversation, conversationDirect } from "@/types/http/conversation"
import Image from "next/image"
import { useState } from "react"
import { useDeleteConversation, useRemoveMember } from "@/hooks/conversation"
import { useUserContext } from "@/context/userContext"
import { defaultAvatarURL } from "@/constants/defaults"

type InputProps = {
    conversation: Conversation
}

export function ConversationBar({conversation}: InputProps) {
    const avatarURL = conversation.avatar_url ?? defaultAvatarURL
    const [clickedOptions, setClickedOptions] = useState(false)
    const { handleDeleteConversation } = useDeleteConversation()
    const { handleRemoveMember } = useRemoveMember()
    const user = useUserContext()

    let actionMessage: string
    if (!conversation.last_message_id) {
        actionMessage = "Start the conversation"
    } else {
        if (conversation.last_message_read === conversation.last_message_id) {
            actionMessage = "seen"
        } else {
            actionMessage = "sent or received" // todo: add attributes to display from backend
        }
    }
    
    return (
        <div>
            <Image src={avatarURL} alt="Conversation avatar" width={20} height={20}/>
            <p>{conversation.name}</p>
            <p>{actionMessage}</p>
            <button onClick={() => setClickedOptions(true)}>Options</button>
            {clickedOptions && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60" onClick={() => setClickedOptions(false)}>
                    <div className="relative bg-white rounded-lg p-6" onClick={(e) => e.stopPropagation()}>
                        {conversation.type === conversationDirect ? (
                            <button onClick={() => handleDeleteConversation(conversation.id)}>Delete</button>
                        ): (
                            <button onClick={() => handleRemoveMember(user.id, conversation.id)}>Leave Group</button>
                        )}
                    </div>
                </div>
            )}
        </div>
    )
}