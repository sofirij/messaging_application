"use client"
import { defaultAvatarURL } from "@/constants/defaults"
import { Conversation, conversationDirect } from "@/types/http/conversation"
import Image from "next/image"
import { useSuspenseQuery } from "@tanstack/react-query"
import { conversationMemberQueryOptions } from "@/query/conversation"
import { useState } from "react"
import { useClearMessages } from "@/hooks/message"

type InputProps = {
    conversation: Conversation
}



export default function NavBar({conversation}: InputProps) {
    const avatarURL = conversation.avatar_url ?? defaultAvatarURL
    const { data: member } = useSuspenseQuery(conversationMemberQueryOptions(conversation.id))
    const [optionsClicked, setOptionsClicked] = useState(false)
    const {handleClearMessages} = useClearMessages()

    return (
        <div className="sticky top-0 z-10">
            <Image src={avatarURL} alt="Conversation profile picture" width={40} height={40}/>
            <p>{conversation.name}</p>
            <p className="truncate whitespace-nowrap overflow-hidden">{member.order.map(memberID => member.data[memberID].username).join(", ")}</p>
            <button onClick={() => setOptionsClicked(true)}>options</button>
            {optionsClicked && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60" onClick={() => setOptionsClicked(false)}>
                    <div className="relative bg-white rounded-lg p-6" onClick={(e) => e.stopPropagation()}>
                        {conversation.type === conversationDirect ? (
                            <>
                                <button onClick={() => {setOptionsClicked(false); handleClearMessages(conversation.id)}}>Clear Messages</button>
                            </>
                        ) : (
                            <>
                                <button>Add Member</button>
                                <button>Remove Member</button>
                                <button>Update Name</button>
                                <button>Update Conversation Picture</button>
                            </>
                        )}
                    </div>
                </div>
            )}
        </div>
    )
}