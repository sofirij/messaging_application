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
import Toast from "@/components/conversation/memberToast"

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
    const [listMembers, setListMembers] = useState(false)
    const [membersToAdd, setMembersToAdd] = useState<Set<User>>(new Set())
    const { handleAddMembers } = useAddMembers()
    const [newGroupname, setNewGroupname] = useState("")
    const { handleUpdateConversationName } = useUpdateConversationName()
    const [updatePicture, setUpdatePicture] = useState(false)
    const fileInputRef = useRef<HTMLInputElement>(null)
    const { handleUpdateConversationAvatar } = useUpdateConversationAvatar()
    const { handleRemoveMember } = useRemoveMember()

    const members = member.order.map(id => member.data[id])
    const userIDs = Array.from(membersToAdd).map(member => member.id)

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
                <p className="truncate whitespace-nowrap overflow-hidden text-xs text-gray-500">{members.map(m => m.username).join(", ")}</p>
            </div>
            <button className="btn-icon" onClick={() => setOptionsClicked(true)} aria-label="Conversation options">⋮</button>
            <Popup popup={optionsClicked} setPopup={setOptionsClicked}>
                <h2 className="popup-title">Options</h2>
                <button className="menu-item" onClick={() => {setOptionsClicked(false); setListMembers(true)}}>Member list</button>
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
                <UserSearch selected={membersToAdd} setSelected={setMembersToAdd}/>
                <button className="btn-primary mt-4 w-full" onClick={() => {setAddMember(false); handleAddMembers(conversation.id, userIDs)}}>Add Members</button>
            </Popup>
            <Popup popup={removeMember} setPopup={setRemoveMember}>
                <h2 className="popup-title">Remove a member</h2>
                <div className="divide-y divide-gray-100">
                    {members.filter(m => m.id != me.id).map(m => {
                        return (
                            <Toast key={m.id} member={m}>
                                <button className="btn-danger px-3 py-1.5" onClick={() => {setRemoveMember(false); handleRemoveMember(m.id, conversation.id)}}>Remove</button>
                            </Toast>
                        )
                    })}
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
            <Popup popup={listMembers} setPopup={setListMembers}>
                <div>
                    <h2 className="popup-title">Members</h2>
                    {members.map(m => {
                        return (
                            <Toast key={m.id} member={m}>
                                <span className={`flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs ${m.is_online ? "bg-brand/10 text-brand-dark" : "bg-gray-100 text-gray-500"}`}>
                                    <span className={`h-2 w-2 rounded-full ${m.is_online ? "bg-brand" : "bg-gray-400"}`} />
                                    {m.is_online ? "Online" : "Offline"}
                                </span>
                            </Toast>
                        )
                    })}
                </div>
            </Popup>
        </div>
    )
}
