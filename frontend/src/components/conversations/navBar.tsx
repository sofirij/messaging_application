"use client"
import { defaultAvatarURL } from "@/constants/defaults"
import { Conversation, conversationDirect } from "@/types/http/conversation"
import Image from "next/image"
import { useQuery } from "@tanstack/react-query"
import { conversationMemberQueryOptions } from "@/query/conversation"
import { useState } from "react"

type InputProps = {
    conversation: Conversation
}



export function ConverstionNavBar({conversation}: InputProps) {
    const avatarURL = conversation.avatar_url ?? defaultAvatarURL
    const { isPending, data: member } = useQuery(conversationMemberQueryOptions(conversation.id))
    const [optionsClicked, setOptionsClicked] = useState(false)

    if (isPending || !member) {
        return "Loading..."
    }

    return (
        <div>
            <Image src={avatarURL} alt="Conversation profile picture" width={40} height={40}/>
            <p>{conversation.name}</p>
            <p className="truncate whitespace-nowrap overflow-hidden">{member.map(m => m.username).join(", ")}</p>
            <button onClick={() => setOptionsClicked(true)}>options</button>
            {optionsClicked && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60" onClick={() => setOptionsClicked(false)}>
                    <div className="relative bg-white rounded-lg p-6" onClick={(e) => e.stopPropagation()}>
                        {conversation.type === conversationDirect ? (
                            <>
                                <button>Clear Messages</button>
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