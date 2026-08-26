"use client"
import NavBar from "@/components/conversations/navBar";
import { conversationMemberQueryOptions, conversationQueryOptions } from "@/query/conversation";
import { messageQueryOptions } from "@/query/message";
import { useSuspenseInfiniteQuery, useSuspenseQuery } from "@tanstack/react-query";
import { notFound, useSearchParams } from "next/navigation";
import Toast from "@/components/messages/toast";
import Sender from "@/components/messages/sender"
import { useEffect, useRef, useState } from "react";
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

    const { data: messages, fetchNextPage, fetchPreviousPage, hasNextPage, hasPreviousPage, isFetchingNextPage, isFetchingPreviousPage } = useSuspenseInfiniteQuery(messageQueryOptions(conversationID))
    const { data: members } = useSuspenseQuery(conversationMemberQueryOptions(conversationID))
    const [reply, setReply] = useState<Message|null>(null)

    const sentinelRef = useRef<HTMLDivElement>(null)
    const bottomRef = useRef<HTMLDivElement>(null)
    const scrollRef = useRef<HTMLDivElement>(null)

    useEffect(() => {
        if (sentinelRef.current === null) return
        if (scrollRef.current === null) return

        const observer = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting && hasPreviousPage && !isFetchingPreviousPage) {
                console.log("top div is intersecting")
                fetchPreviousPage()
            } 
        }, {rootMargin: "100px", root: scrollRef.current, threshold: 0})
        observer.observe(sentinelRef.current)
        return () => observer.disconnect()
    }, [hasPreviousPage, isFetchingPreviousPage, fetchPreviousPage])

    useEffect(() => {
        bottomRef.current?.scrollIntoView({ block: "end" })
    }, [])

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
                <Sender reply={reply} setReply={setReply} conversationID={conversationID}/>
            </div>
        </main> 
    )
}