"use client"
import Input from "@/components/ui/input";
import { conversationQueryOptions } from "@/query/conversation";
import { useSuspenseQuery } from "@tanstack/react-query";
import { useState } from "react";
import { useCreateConversation } from "@/hooks/mutation/conversation";
import { User } from "@/types/http/user";
import { ConversationType, conversationGroup, conversationDirect } from "@/types/http/conversation";
import Toast from "@/components/conversations/toast";
import Popup from "@/components/ui/popup";
import UserSearch from "@/components/ui/userSearch";

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

    function onSuccess() {
        setCreatingConversation(false)
        setSelected(new Set())
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
            <Popup popup={creatingConversation} setPopup={setCreatingConversation}>
                <label>
                    Type
                    <select value={mustBeGroup ? conversationGroup : type} onChange={(e) => setIsGroup(e.target.value === conversationGroup)}>
                        <option value={conversationGroup}>Group</option>
                        <option value={conversationDirect}>Direct</option>
                    </select>
                    {isGroup && (
                        <Input label="Group Name" type="text" onChange={(e) => setGroupName(e.target.value)}/>
                    )}
                    <UserSearch selected={selected} setSelected={setSelected}/>
                </label>
                <button onClick={() => handleCreateConversation(type, groupName, userIDs, onSuccess)}>Create</button>
            </Popup>
        </main>
    )
}