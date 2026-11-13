import { defaultAvatarURL } from "@/config/defaults"
import { Message } from "@/types/http/message"
import { User } from "@/types/http/user"
import Image from "next/image"
import AttachmentList from "@/components/conversation/attachment"
import { useState } from "react"
import { useDeleteMessage } from "@/hooks/mutation/message"
import EditInput from "@/components/conversation/editInput"
import ReplyToast from "@/components/conversation/replyToast"
import Popup from "@/components/ui/popup"
import { useSuspenseQuery } from "@tanstack/react-query"
import { userQueryOptions } from "@/query/client/user"
import { SetState } from "@/types/react"
import { conversationQueryOptions } from "@/query/client/conversation"

type InputProps = {
    sender: User,
    message: Message,
    setReply: SetState<Message | null>
}

export default function Toast({sender, message, setReply} : InputProps) {
    const avatarURL = sender.avatar_url ?? defaultAvatarURL
    const [optionsClicked, setOptionsClicked] = useState(false)
    const { handleDeleteMessage } = useDeleteMessage()
    const {data: conversations} = useSuspenseQuery(conversationQueryOptions)
    const [editClicked, setEditClicked] = useState(false)
    const {data: me} = useSuspenseQuery(userQueryOptions)
    const isMine = me.id === message.sender_id
    const alignment = isMine ? "justify-end" : "justify-start"
    const bubbleColor = isMine ? "bg-bubble-out rounded-tr-none" : "bg-white rounded-tl-none"
    const time = new Date(message.created_at).toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"})
    const conversation = conversations.data[message.conversation_id]

    return (
        <div className={`group flex items-end gap-2 ${alignment}`}>
            {!isMine && (
                <Image src={avatarURL} alt="user avatar" width={28} height={28} className="avatar h-7 w-7"/>
            )}
            <div className={`relative flex max-w-[75%] flex-col gap-1 rounded-lg px-3 py-2 shadow-sm sm:max-w-[60%] ${bubbleColor}`}>
                <div className="flex items-start gap-2">
                    {!isMine && (
                        <p className="text-xs font-semibold text-brand-dark">{sender.username}</p>
                    )}
                    <button onClick={() => setOptionsClicked(true)} className="-mt-1 -mr-1 ml-auto rounded px-1 text-gray-400 opacity-60 transition-opacity hover:text-gray-700 group-hover:opacity-100" aria-label="Message options">⋮</button>
                </div>
                {message.reply !== null && (
                    <ReplyToast reply={message.reply}/>
                )}
                {editClicked ? (
                    <EditInput conversationID={message.conversation_id} messageID={message.id} setEditClicked={setEditClicked}/>
                ): (
                    <></>
                )}
                {message.body && (
                    <p className="whitespace-pre-wrap break-words text-sm">{message.body}</p>
                )}
                <AttachmentList attachments={message.attachments}/>
                <div className="flex items-center gap-1 self-end">
                    <span className="text-[10px] text-gray-500" suppressHydrationWarning>{time}</span>
                    {(isMine && conversation.last_message_sent_in_conversation && message.id === conversation.last_message_sent_in_conversation.id) && (
                        <>
                            {(conversation.last_message_read_in_conversation && conversation.last_message_read_in_conversation >= message.id) ? (
                                <svg viewBox="0 0 18 12" className="h-3 w-[18px] text-sky-500" fill="none" stroke="currentColor" strokeWidth={1.6} strokeLinecap="round" strokeLinejoin="round" role="img" aria-label="Seen">
                                    <path d="M1 6.5 4.5 10 12 2"/>
                                    <path d="M7.5 9 8.5 10 16 2"/>
                                </svg>
                            ): (
                                <svg viewBox="0 0 14 12" className="h-3 w-3.5 text-gray-500" fill="none" stroke="currentColor" strokeWidth={1.6} strokeLinecap="round" strokeLinejoin="round" role="img" aria-label="Sent">
                                    <path d="M2 6.5 5.5 10 13 2"/>
                                </svg>
                            )}
                        </>
                    )}
                </div>
                <Popup popup={optionsClicked} setPopup={setOptionsClicked}>
                    { isMine ? (
                        <>
                            <h2 className="popup-title">Message</h2>
                            <button className="menu-item" onClick={() => {setOptionsClicked(false); setReply(message)}}>Reply</button>
                            <button className="menu-item" onClick={() => {setOptionsClicked(false); setEditClicked(true)}}>Edit</button>
                            <button className="menu-item text-red-600" onClick={() => {setOptionsClicked(false); handleDeleteMessage(message.conversation_id, message.id)}}>Delete</button>
                        </>
                    ) : (
                        <>
                            <h2 className="popup-title">Message</h2>
                            <button className="menu-item" onClick={() => {setOptionsClicked(false); setReply(message)}}>Reply</button>
                        </>
                    )}
                    
                </Popup>
            </div>
        </div>
    )
}
