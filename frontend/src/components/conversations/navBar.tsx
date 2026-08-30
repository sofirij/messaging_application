"use client"
import { defaultAvatarURL } from "@/constants/defaults"
import { Conversation, conversationDirect } from "@/types/http/conversation"
import Image from "next/image"
import { useSuspenseQuery } from "@tanstack/react-query"
import { conversationMemberQueryOptions } from "@/query/conversation"
import { ChangeEvent, useRef, useState } from "react"
import { useClearMessages } from "@/hooks/mutation/message"
import Popup from "@/components/ui/popup"
import UserSearch from "@/components/ui/userSearch"
import { User } from "@/types/http/user"
import { useAddMembers, useRemoveMember, useUpdateConversationAvatar, useUpdateConversationName } from "@/hooks/mutation/conversation"
import { userQueryOptions } from "@/query/user"
import Input from "@/components/ui/input"

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
        <div className="sticky top-0 z-10">
            <Image src={avatarURL} alt="Conversation profile picture" width={40} height={40}/>
            <p>{conversation.name}</p>
            <p className="truncate whitespace-nowrap overflow-hidden">{member.order.map(memberID => member.data[memberID].username).join(", ")}</p>
            <button onClick={() => setOptionsClicked(true)}>options</button>
            <Popup popup={optionsClicked} setPopup={setOptionsClicked}>
                {conversation.type === conversationDirect ? (
                    <>
                        <button onClick={() => {setOptionsClicked(false); handleClearMessages(conversation.id)}}>Clear Messages</button>
                    </>
                ) : (
                    <>
                        <button onClick={() => {setOptionsClicked(false); setAddMember(true)}}>Add Members</button>
                        <button onClick={() => {setOptionsClicked(false); setRemoveMember(true)}}>Remove Member</button>
                        <button onClick={() => {setOptionsClicked(false); setUpdateName(true)}}>Update Name</button>
                        <button onClick={() => {setOptionsClicked(false); setUpdatePicture(true)}}>Update Conversation Picture</button>
                    </>
                )}
            </Popup>
            <Popup popup={addMember} setPopup={setAddMember}>
                <UserSearch selected={members} setSelected={setMembers}/>
                <button onClick={() => {setAddMember(false); handleAddMembers(conversation.id, userIDs)}}>Add Members</button>
            </Popup>
            <Popup popup={removeMember} setPopup={setRemoveMember}>
                {member.order.filter(id => id !== me.id).map(id => (
                    <div key={id}>
                        <p>{member.data[id].username}</p>
                        <button onClick={() => {setRemoveMember(false); handleRemoveMember(id, conversation.id)}}>Remove</button>
                    </div>
                ))}
            </Popup>
            <Popup popup={updateName} setPopup={setUpdateName}>
                <Input label="groupname" type="text" onChange={(e) => setNewGroupname(e.target.value)}/>
                <button onClick={() => {setUpdateName(false); handleUpdateConversationName(conversation.id, newGroupname)}}>Update</button>
            </Popup>
            <Popup popup={updatePicture} setPopup={setUpdatePicture}>
                <button onClick={() => {setUpdatePicture(false); handleUpdateConversationAvatar(conversation.id, null)}}>Remove picture</button>
                <button onClick={() => fileInputRef.current?.click()}>Select</button>
                <input type="file" accept="image/*" ref={fileInputRef} onChange={handleProfilePicChange} style={{display: "none"}}/>
            </Popup>
        </div>
    )
}