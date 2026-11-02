"use client"
import { defaultAvatarURL } from "@/config/defaults"
import { Conversation, conversationDirect } from "@/types/http/conversation"
import Image from "next/image"
import { useSuspenseQuery } from "@tanstack/react-query"
import { conversationMemberQueryOptions } from "@/query/client/conversation"
import { ChangeEvent, useRef, useState } from "react"
import { useClearMessages } from "@/hooks/mutation/message"
import Popup from "@/components/ui/popup"
import UserSearch from "@/components/ui/userSearch"
import { User } from "@/types/http/user"
import { useAddMembers, useRemoveMember, useUpdateConversationAvatar, useUpdateConversationName } from "@/hooks/mutation/conversation"
import { userQueryOptions } from "@/query/client/user"
import Input from "@/components/ui/input"
import Link from "next/link"

type InputProps = {
    conversation: Conversation
}


export default function NavBar({conversation}: InputProps) {
    const avatarURL = conversation.avatar_url ?? defaultAvatarURL
    const {data: me} = useSuspenseQuery(userQueryOptions)
    const { data: member } = useSuspenseQuery(conversationMemberQueryOptions(conversation.id))
    const [optionsClicked, setOptionsClicked] = useState(false)
    const {handleClearMessages} = useClearMessages()
    const [addMember, setAddMember] = useState(false)
    const [removeMember, setRemoveMember] = useState(false)
    const [updateName, setUpdateName] = useState(false)
    const [members, setMembers] = useState<Set<User>>(new Set())
    const { handleAddMembers } = useAddMembers()
    const { handleRemoveMember } = useRemoveMember()
    const [newGroupname, setNewGroupname] = useState("")
    const { handleUpdateConversationName } = useUpdateConversationName()
    const [updatePicture, setUpdatePicture] = useState(false)
    const fileInputRef = useRef<HTMLInputElement>(null)
    const { handleUpdateConversationAvatar } = useUpdateConversationAvatar()

    const userIDs = Array.from(members).map(member => member.id)

    function handleProfilePicChange(e: ChangeEvent<HTMLInputElement>) {
        const file = e.target.files?.[0]
        if (!file) return
        handleUpdateConversationAvatar(conversation.id, file)
    }

    return (
        <div className="sticky top-0 z-10 flex items-center gap-3 border-b border-gray-200 bg-white px-3 py-2 shadow-sm">
            <Link href="/conversations" className="btn-icon" aria-label="Back to conversations">←</Link>
            <Image src={avatarURL} alt="Conversation profile picture" width={40} height={40} className="avatar h-10 w-10"/>
            <div className="min-w-0 flex-1">
                <p className="truncate font-medium">{conversation.name}</p>
                <p className="truncate whitespace-nowrap overflow-hidden text-xs text-gray-500">{member.order.map(memberID => member.data[memberID].username).join(", ")}</p>
            </div>
            <button className="btn-icon" onClick={() => setOptionsClicked(true)} aria-label="Conversation options">⋮</button>
            <Popup popup={optionsClicked} setPopup={setOptionsClicked}>
                <h2 className="popup-title">Options</h2>
                {conversation.type === conversationDirect ? (
                    <>
                        <button className="menu-item text-red-600" onClick={() => {setOptionsClicked(false); handleClearMessages(conversation.id)}}>Clear Messages</button>
                    </>
                ) : (
                    <>
                        <button className="menu-item" onClick={() => {setOptionsClicked(false); setAddMember(true)}}>Add Members</button>
                        <button className="menu-item" onClick={() => {setOptionsClicked(false); setRemoveMember(true)}}>Remove Member</button>
                        <button className="menu-item" onClick={() => {setOptionsClicked(false); setUpdateName(true)}}>Update Name</button>
                        <button className="menu-item" onClick={() => {setOptionsClicked(false); setUpdatePicture(true)}}>Update Conversation Picture</button>
                    </>
                )}
            </Popup>
            <Popup popup={addMember} setPopup={setAddMember}>
                <h2 className="popup-title">Add members</h2>
                <UserSearch selected={members} setSelected={setMembers}/>
                <button className="btn-primary mt-4 w-full" onClick={() => {setAddMember(false); handleAddMembers(conversation.id, userIDs)}}>Add Members</button>
            </Popup>
            <Popup popup={removeMember} setPopup={setRemoveMember}>
                <h2 className="popup-title">Remove a member</h2>
                <div className="divide-y divide-gray-100">
                    {member.order.filter(id => id !== me.id).map(id => (
                        <div key={id} className="flex items-center gap-3 py-2">
                            <Image src={member.data[id].avatar_url ?? defaultAvatarURL} alt="User profile picture" width={32} height={32} className="avatar h-8 w-8"/>
                            <p className="flex-1 text-sm">{member.data[id].username}</p>
                            <button className="btn-danger px-3 py-1.5" onClick={() => {setRemoveMember(false); handleRemoveMember(id, conversation.id)}}>Remove</button>
                        </div>
                    ))}
                </div>
            </Popup>
            <Popup popup={updateName} setPopup={setUpdateName}>
                <h2 className="popup-title">Rename group</h2>
                <Input label="group name" type="text" placeholder={conversation.name} onChange={(e) => setNewGroupname(e.target.value)}/>
                <button className="btn-primary mt-4 self-end" onClick={() => {setUpdateName(false); handleUpdateConversationName(conversation.id, newGroupname)}}>Update</button>
            </Popup>
            <Popup popup={updatePicture} setPopup={setUpdatePicture}>
                <h2 className="popup-title">Group picture</h2>
                <Image src={avatarURL} alt="Conversation profile picture" width={160} height={160} className="avatar mx-auto h-40 w-40"/>
                <div className="mt-4 flex justify-end gap-2">
                    {conversation.avatar_url && (
                        <button className="btn-secondary" onClick={() => {setUpdatePicture(false); handleUpdateConversationAvatar(conversation.id, null)}}>Remove picture</button>
                    )}
                    <button className="btn-primary" onClick={() => fileInputRef.current?.click()}>Select</button>
                </div>
                <input type="file" accept="image/*" ref={fileInputRef} onChange={handleProfilePicChange} style={{display: "none"}}/>
            </Popup>
        </div>
    )
}
