"use client"
import NavBar from "@/components/conversations/navBar";
import { conversationMemberQueryOptions, conversationQueryOptions } from "@/query/conversation";
import { messageQueryOptions } from "@/query/message";
import { useSuspenseInfiniteQuery, useSuspenseQuery } from "@tanstack/react-query";
import { notFound, useSearchParams } from "next/navigation";
import Toast from "@/components/messages/toast";
import Sender from "@/components/messages/sender"
import { useState } from "react";
import { Message } from "@/types/http/message";

export default function Conversation() {
    const searchParams = useSearchParams()
    const id = searchParams.get("id")
    const conversationID = Number(id)

    if (!Number.isInteger(conversationID)) notFound()
    
    return <ConversationContent conversationID={conversationID}/>
}

function ConversationContent({conversationID}: {conversationID: number}) {
    const { data: conversations } = useSuspenseQuery(conversationQueryOptions)
    if (!conversations.data[conversationID]) notFound()

    const { data: messages } = useSuspenseInfiniteQuery(messageQueryOptions(conversationID))
    const { data: members } = useSuspenseQuery(conversationMemberQueryOptions(conversationID))
    const [reply, setReply] = useState<Message|null>(null)

    return (
        <main>
            <NavBar conversation={conversations.data[conversationID]}></NavBar>
            <div>
                {messages.filter(message => !message.deleted).map(message => {
                    return (
                        <Toast key={message.id} sender={members.data[message.sender_id]} message={message} setReply={setReply}/>
                    )
                })}
            </div>
            <div>
                <Sender reply={reply} setReply={setReply} conversationID={conversationID}/>
            </div>
        </main> 
    )
}