"use client"
import Input from "@/components/ui/input";
import Search from "@/components/conversations/userSearch";
import { conversationQueryOptions } from "@/query/conversation";
import { useSuspenseQuery } from "@tanstack/react-query";
import { useState } from "react";
import { useCreateConversation } from "@/hooks/conversation";
import { User } from "@/types/http/user";
import Image from "next/image";
import { toggleSelected } from "@/utils/conversation";
import { ConversationType, conversationGroup, conversationDirect } from "@/types/http/conversation";
import Toast from "@/components/conversations/toast";
import { defaultAvatarURL } from "@/constants/defaults";

export default function Conversations() {
    const { data: conversations } = useSuspenseQuery(conversationQueryOptions)
    const [creatingConversation, setCreatingConversation] = useState(false)
    const [isGroup, setIsGroup] = useState(false)
    const [groupName, setGroupName] = useState("")
    const [selected, setSelected] = useState(new Set<User>())
    const { handleCreateConversation } = useCreateConversation()

    const userIDs = Array.from(selected).map(user => user.id)
    const type : ConversationType = isGroup ? conversationGroup : conversationDirect
    const mustBeGroup = isGroup || selected.size > 1

    function createConversation() {
        handleCreateConversation(type, groupName, userIDs, () => {
            setCreatingConversation(false)
            setSelected(new Set())
        })
    }

    return (
        <main>
            {conversations.order.length > 0 ? (
                <div>
                    <button onClick={() => setCreatingConversation(true)}>New Conversation</button>
                    {conversations.order.map((id) => (
                        <Toast key={id} conversation={conversations.data[id]}/>
                    ))}
                </div>
            ) : (
                <div>                                                                    
                    <p>Start your first conversation</p>
                    <button onClick={() => setCreatingConversation(true)}>New Conversation</button>
                </div>
            )}
            {creatingConversation && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60" onClick={() => setCreatingConversation(false)}>
                    <div className="relative bg-white rounded-lg p-6" onClick={(e) => e.stopPropagation()}>
                        <label>
                            Type
                            <select value={mustBeGroup ? conversationGroup : type} onChange={(e) => setIsGroup(e.target.value === conversationGroup)}>
                                <option value={conversationGroup}>Group</option>
                                <option value={conversationDirect}>Direct</option>
                            </select>
                            {isGroup && (
                                <Input label="Group Name" type="text" onChange={(e) => setGroupName(e.target.value)}/>
                            )}
                            <Search selected={selected} setSelected={setSelected}/>
                        </label>
                        <button onClick={createConversation}>Create</button>
                        {Array.from(selected).map(user => {
                            const avatarURL = user.avatar_url ?? defaultAvatarURL
                            return (
                                <div key={user.id}>
                                    <Image src={avatarURL} alt="User profile picture" width={20} height={20} />
                                    <p>{user.username}</p>
                                    <button onClick={() => toggleSelected(user, setSelected)}>Deselect</button>
                                </div>
                            )
                        })}
                    </div>
                </div>
            )}
        </main>
    )
}