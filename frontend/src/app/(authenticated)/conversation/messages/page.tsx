"use client"
import NavBar from "@/components/conversations/navBar";
import { conversationMemberQueryOptions, conversationQueryOptions } from "@/query/conversation";
import { messageQueryOptions } from "@/query/message";
import { useSuspenseInfiniteQuery, useSuspenseQuery } from "@tanstack/react-query";
import { notFound, useSearchParams } from "next/navigation";
import Toast from "@/components/messages/toast";
import Sender from "@/components/messages/sender"
import { useEffect, useLayoutEffect, useRef, useState } from "react";
import { Message } from "@/types/http/message";
import { userQueryOptions } from "@/query/user";
import { useWSContext } from "@/context/wsProvider";

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

    const { data: messages, fetchPreviousPage, hasPreviousPage, isFetchingPreviousPage } = useSuspenseInfiniteQuery(messageQueryOptions(conversationID))
    const { data: members } = useSuspenseQuery(conversationMemberQueryOptions(conversationID))
    const { data: me } = useSuspenseQuery(userQueryOptions)
    const [reply, setReply] = useState<Message|null>(null)

    const sentinelRef = useRef<HTMLDivElement>(null)
    const bottomRef = useRef<HTMLDivElement>(null)
    const scrollRef = useRef<HTMLDivElement>(null)
    const shouldScroll = useRef(true)

    const { readMessage, ws, status } = useWSContext()

    useEffect(() => {
        if (sentinelRef.current === null) return
        if (scrollRef.current === null) return
        if (bottomRef.current === null) return

        const fetchPreviousObserver = new IntersectionObserver(([entry]) => {
            if (entry.isIntersecting && hasPreviousPage && !isFetchingPreviousPage) {
                console.log("top div is intersecting")
                fetchPreviousPage()
            } 
        }, {rootMargin: "100px 0px 0px 0px", root: scrollRef.current, threshold: 0})
        fetchPreviousObserver.observe(sentinelRef.current)

        const scrollObserver = new IntersectionObserver(([entry]) => {
            shouldScroll.current = entry.isIntersecting
            if (entry.isIntersecting) {
                console.log("bottom ref is intersecting")
    
            } else console.log("bottom ref is not intersecting")
        }, {rootMargin: "0px 0px 100px 0px", root: scrollRef.current, threshold: 0})
        scrollObserver.observe(bottomRef.current)

        return () => { 
            fetchPreviousObserver.disconnect()
            scrollObserver.disconnect()
        }
    }, [hasPreviousPage, isFetchingPreviousPage, fetchPreviousPage])

    useLayoutEffect(() => {
        if (shouldScroll.current) bottomRef.current?.scrollIntoView({block: "end"})
        if (messages.length === 0) return

        // read the message if it wasn't sent by you and its after the last message you've read
        const lastMessage = messages[messages.length - 1]
        const lastReadByUser = conversations.data[conversationID].last_message_read_by_user
        if (lastMessage.sender_id !== me.id && (!lastReadByUser || lastMessage.id > lastReadByUser)) {
            if (!ws.current) return
            readMessage(ws.current, conversationID, lastMessage.id)
        }
    }, [conversationID, conversations.data, me.id, messages, readMessage, status, ws])

    return (
        <main className="h-full flex flex-col">
            <NavBar conversation={conversations.data[conversationID]}></NavBar>
            <div ref={scrollRef} className="flex-1 overflow-y-auto">
                <div style={{outline: "1px solid red"}} ref={sentinelRef} />
                {messages.filter(message => !message.deleted).map(message => {
                    return (
                        <Toast key={message.id} sender={members.data[message.sender_id]} message={message} setReply={setReply}/>
                    )
                })}
                <div style={{outline: "1px solid red"}} ref={bottomRef}/>
            </div>
            <div>
                <Sender reply={reply} setReply={setReply} conversationID={conversationID} bottomRef={bottomRef}/>
            </div>
        </main> 
    )
}