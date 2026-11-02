"use client"

import Input from "@/components/ui/input";
import { useState } from "react";
import { useCreateConversation } from "@/hooks/mutation/conversation";
import { User } from "@/types/http/user";
import { ConversationType, conversationGroup, conversationDirect } from "@/types/http/conversation";
import Toast from "@/components/conversations/toast";
import Popup from "@/components/ui/popup";
import UserSearch from "@/components/ui/userSearch";
import { useSuspenseQuery } from "@tanstack/react-query"
import { conversationQueryOptions } from "@/query/client/conversation";

export function Conversations() {
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
        <main className="h-full overflow-y-auto">
            <div className="mx-auto flex min-h-full max-w-2xl flex-col bg-white shadow-sm">
                <div className="flex items-center justify-between border-b border-gray-200 px-4 py-3">
                    <h1 className="text-xl font-semibold">Chats</h1>
                    <button className="btn-primary" onClick={() => setCreatingConversation(true)}>+ New Conversation</button>
                </div>
                {conversations.order.length > 0 ? (
                    <div>
                        {conversations.order.map((id) => (
                            <Toast key={id} conversation={conversations.data[id]}/>
                        ))}
                    </div>
                ) : (
                    <div className="flex flex-1 flex-col items-center justify-center gap-3 p-8 text-center">
                        <div className="text-5xl">💬</div>
                        <p className="text-lg font-medium">Start your first conversation</p>
                        <p className="text-sm text-gray-500">Find someone by username and say hello.</p>
                        <button className="btn-primary mt-2" onClick={() => setCreatingConversation(true)}>New Conversation</button>
                    </div>
                )}
            </div>
            <Popup popup={creatingConversation} setPopup={setCreatingConversation}>
                <h2 className="popup-title">New conversation</h2>
                <div className="flex flex-col gap-4">
                    <label className="flex flex-col gap-1.5 text-sm font-medium text-gray-700">
                        Type
                        <select className="text-field font-normal" value={mustBeGroup ? conversationGroup : type} onChange={(e) => setIsGroup(e.target.value === conversationGroup)}>
                            <option value={conversationGroup}>Group</option>
                            <option value={conversationDirect}>Direct</option>
                        </select>
                    </label>
                    {isGroup && (
                        <Input label="Group Name" type="text" onChange={(e) => setGroupName(e.target.value)}/>
                    )}
                    <UserSearch selected={selected} setSelected={setSelected}/>
                </div>
                <button className="btn-primary mt-4 w-full" onClick={() => handleCreateConversation(type, groupName, userIDs, onSuccess)}>Create</button>
            </Popup>
        </main>
    )
}
