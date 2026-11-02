"use client"

import NavBar from "@/components/conversations/navBar";
import { conversationMemberQueryOptions, conversationQueryOptions } from "@/query/client/conversation";
import { messageQueryOptions } from "@/query/client/message";
import { useSuspenseInfiniteQuery, useSuspenseQuery } from "@tanstack/react-query";
import Toast from "@/components/conversation/toast";
import Sender from "@/components/conversation/sender"
import { useEffect, useLayoutEffect, useRef, useState } from "react";
import { Message } from "@/types/http/message";
import { userQueryOptions } from "@/query/client/user";
import { useWSContext } from "@/context/wsProvider";

export default function Conversation({conversationID}: {conversationID: number}) {
    const { data: conversations } = useSuspenseQuery(conversationQueryOptions)

    const { data: messages, fetchPreviousPage, hasPreviousPage, isFetchingPreviousPage } = useSuspenseInfiniteQuery(messageQueryOptions(conversationID))
    const { data: members } = useSuspenseQuery(conversationMemberQueryOptions(conversationID))
    const { data: me } = useSuspenseQuery(userQueryOptions)
    const [reply, setReply] = useState<Message|null>(null)

    const sentinelRef = useRef<HTMLDivElement>(null)
    const bottomRef = useRef<HTMLDivElement>(null)
    const scrollRef = useRef<HTMLDivElement>(null)
    const shouldScroll = useRef(true)

    const { readMessage, ws } = useWSContext()

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
    }, [conversationID, conversations.data, me.id, messages, readMessage, ws])

    return (
        <main className="h-full flex flex-col bg-chat">
            <NavBar conversation={conversations.data[conversationID]}></NavBar>
            <div ref={scrollRef} className="flex flex-1 overflow-y-auto flex-col gap-1 px-3 py-4 sm:px-8">
                <div ref={sentinelRef} />
                {isFetchingPreviousPage && (
                    <div className="flex justify-center py-2">
                        <div className="animate-spin h-5 w-5 border-2 border-brand border-t-transparent rounded-full" />
                    </div>
                )}
                {messages.length === 0 && (
                    <p className="m-auto rounded-lg bg-white/80 px-4 py-2 text-sm text-gray-600 shadow-sm">No messages yet. Say hello 👋</p>
                )}
                {messages.filter(message => !message.deleted).map(message => {
                    return (
                        <Toast key={message.id} sender={members.data[message.sender_id]} message={message} setReply={setReply}/>
                    )
                })}
                <div ref={bottomRef}/>
            </div>
            <div className="border-t border-gray-200 bg-gray-50">
                <Sender reply={reply} setReply={setReply} conversationID={conversationID} bottomRef={bottomRef}/>
            </div>
        </main> 
    )
}