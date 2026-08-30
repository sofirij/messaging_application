"use client"

import { useConversationNew, useMemberAdded, useMemberRemoved } from "@/hooks/query/conversation"
import { useMessageDeleted, useMessageEdited, useMessageNew } from "@/hooks/query/message"
import { getWSConn, readMessage, startTyping, stopTyping } from "@/lib/api/ws"
import { userQueryOptions } from "@/query/user"
import { Event, EventConversationNew, EventMemberAdded, EventMemberRemoved, EventMessageDeleted, EventMessageEdited, EventMessageNew } from "@/types/ws/event"
import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query"
import { createContext, useContext, useEffect, useRef, useState } from "react"

export type WSContextType = {
    status: WebSocket["readyState"]
    readMessage: (ws: WebSocket, conversation_id: number, message_id: number) => void
    startTyping: (ws: WebSocket, conversation_id: number) => void
    stopTyping: (ws: WebSocket, conversation_id: number) => void
}

const WSContext = createContext<WSContextType|null>(null)

export function WSProvider({children}: {children: React.ReactNode}) {
    const ws = useRef<WebSocket|null>(null)
    const [status, setStatus] = useState<WebSocket["readyState"]>(0)
    const queryClient = useQueryClient()
    const { data: me } = useSuspenseQuery(userQueryOptions)

    const handleMessageNew = useMessageNew()
    const handleMessageDeleted = useMessageDeleted()
    const handleMessageEdited = useMessageEdited()
    const handleMemberAdded = useMemberAdded()
    const handleMemberRemoved = useMemberRemoved()
    const handleConversationNew = useConversationNew()

    useEffect(() => {
        let cancelled = false
        async function connect() {
            const conn = await getWSConn()

            if (cancelled) return

            conn.onopen = () => { console.log("opened ws conn") }

            conn.onmessage = (event) => {
            
                
                const msg: Event  = JSON.parse(event.data)
                console.log(msg.type)

                switch (msg.type) {
                    case EventMessageNew: {
                        handleMessageNew(msg.payload)
                        break
                    }
                    case EventMessageDeleted: {
                        handleMessageDeleted(msg.payload.conversation_id, msg.payload.message_id)
                        break
                    }
                    case EventMessageEdited: {
                        handleMessageEdited(msg.payload.conversation_id, msg.payload.message_id, msg.payload.body)
                        break
                    }
                    case EventMemberAdded: {
                        handleMemberAdded(msg.payload.conversation_id, msg.payload.user)
                        break
                    }
                    case EventMemberRemoved: {
                        handleMemberRemoved(msg.payload.conversation_id, msg.payload.user_id)
                        break
                    }
                    case EventConversationNew: {
                        handleConversationNew()
                        break
                    }
                }
            } 

            ws.current = conn
            setStatus(conn.readyState)
        }

        connect()

        return () => {
            cancelled = true
            if (ws.current) {
                ws.current.close()
                ws.current = null
            }
        }
    }, [handleConversationNew, handleMemberAdded, handleMemberRemoved, handleMessageDeleted, handleMessageEdited, handleMessageNew, me.id, queryClient])

    return (
        <WSContext value={{status, readMessage, startTyping, stopTyping}}>
            {children}
        </WSContext>
    )
}

export function useWSContext() {
    const ctx = useContext(WSContext)
    if (!ctx) throw new Error("ws context should be used within ws provider")
    return ctx
}