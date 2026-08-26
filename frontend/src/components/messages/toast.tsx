import { defaultAvatarURL } from "@/constants/defaults"
import { Message } from "@/types/http/message"
import { User } from "@/types/http/user"
import Image from "next/image"
import AttachmentList from "@/components/messages/attachment"
import { Dispatch, SetStateAction, useState } from "react"
import { useDeleteMessage } from "@/hooks/message"
import EditInput from "@/components/messages/editInput"
import ReplyToast from "@/components/messages/replyToast"
import Popup from "@/components/ui/popup"

type InputProps = {
    sender: User,
    message: Message,
    setReply: Dispatch<SetStateAction<Message | null>>
}

export default function Toast({sender, message, setReply} : InputProps) {
    const avatarURL = sender.avatar_url ?? defaultAvatarURL
    const [optionsClicked, setOptionsClicked] = useState(false)
    const { handleDeleteMessage } = useDeleteMessage()
    const [editClicked, setEditClicked] = useState(false)

    return (
        <div>
            <Image src={avatarURL} alt="user avatar" width={20} height={20}/>
            <p>{sender.username}</p>
            {message.reply !== null && (
                <ReplyToast reply={message.reply}/>
            )}
            {editClicked ? (
                <EditInput conversationID={message.conversation_id} messageID={message.id} setEditClicked={setEditClicked}/>
            ): (
                <></>
            )}
            <p>{message.body}</p>
            <AttachmentList attachments={message.attachments}/>
            <button onClick={() => setOptionsClicked(true)}>Options</button>
            <Popup popup={optionsClicked} setPopup={setOptionsClicked}>
                <button onClick={() => {setOptionsClicked(false); setEditClicked(true)}}>Edit</button>
                <button onClick={() => {setOptionsClicked(false); handleDeleteMessage(message.conversation_id, message.id)}}>Delete</button>
                <button onClick={() => {setOptionsClicked(false); setReply(message)}}>Reply</button>
            </Popup>
        </div>
    )
}