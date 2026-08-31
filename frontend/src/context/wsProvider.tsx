"use client"

import { useConversationNew, useMemberAdded, useMemberRemoved } from "@/hooks/query/conversation"
import { useMessageDeleted, useMessageEdited, useMessageNew } from "@/hooks/query/message"
import { useReconnected } from "@/hooks/query/ws"
import { getWSConn, readMessage, startTyping, stopTyping } from "@/lib/api/ws"
import { HTTPError } from "@/types/http/error"
import { Event, EventConversationNew, EventMemberAdded, EventMemberRemoved, EventMessageDeleted, EventMessageEdited, EventMessageNew } from "@/types/ws/event"
import { useRouter } from "next/navigation"
import { createContext, useContext, useEffect, useRef, useState } from "react"
import { useErrorContext } from "@/context/errorContext"

export type WSContextType = {
    status: WebSocket["readyState"]
    readMessage: (ws: WebSocket, conversation_id: number, message_id: number) => void
    startTyping: (ws: WebSocket, conversation_id: number) => void
    stopTyping: (ws: WebSocket, conversation_id: number) => void
}

const WSContext = createContext<WSContextType|null>(null)

export function WSProvider({children}: {children: React.ReactNode}) {
    const ws = useRef<WebSocket|null>(null)
    const [status, setStatus] = useState<WebSocket["readyState"]>(WebSocket.CONNECTING)
    const router = useRouter()

    const handleMessageNew = useMessageNew()
    const handleMessageDeleted = useMessageDeleted()
    const handleMessageEdited = useMessageEdited()
    const handleMemberAdded = useMemberAdded()
    const handleMemberRemoved = useMemberRemoved()
    const handleConversationNew = useConversationNew()
    const handleReconnected = useReconnected()
    const { addError } = useErrorContext()

    const baseDelay = 2000 // 2secs
    const maxDelay = 60000 // 1min
    const attempt = useRef(0)
    const retryTimer = useRef<ReturnType<typeof setTimeout>>(null)
    const reconnected = useRef(false)

    useEffect(() => {
        let cancelled = false

        function clearRetry() {
            if (retryTimer.current) {
                clearTimeout(retryTimer.current)
                retryTimer.current = null
            }
        }

        function scheduleRetry() {
            clearRetry()
            const d = Math.min(baseDelay * 2 ** attempt.current, maxDelay)
            const jittered = Math.random() * d
            attempt.current++
            reconnected.current = true
            retryTimer.current = setTimeout(connect, jittered)
        }

        function reconnectNow() {
            if (document.visibilityState !== "visible" || navigator.onLine === false) return

            const rs = ws.current?.readyState

            if (rs === WebSocket.OPEN || rs === WebSocket.CONNECTING) return

            clearRetry()
            attempt.current = 0
            reconnected.current = true
            connect()
        }

        function onVisibility() {
            if (document.visibilityState === "visible") {
                reconnectNow()
            } else {
                clearRetry()
            }
        }

        async function connect() {
            if (document.visibilityState !== "visible" || navigator.onLine === false) return

            const rs = ws.current?.readyState

            if (rs === WebSocket.OPEN || rs === WebSocket.CONNECTING) return

            try {
                const conn = await getWSConn()

                if (cancelled) return

                conn.onopen = () => {
                    if (ws.current) {
                        setStatus(ws.current.readyState)
                    }
                    attempt.current = 0
                }

                conn.onclose = (closeEvent) => {
                    if (cancelled) {
                        return
                    }
                    if (ws.current) {
                        setStatus(ws.current.readyState)
                    }
                    if (!closeEvent.wasClean) {
                        scheduleRetry()
                    }
                }

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
                if (reconnected.current) {
                    handleReconnected()
                }
            } catch (e) {
                if (e instanceof Error) {
                    addError(e.message)
                }

                if (e instanceof HTTPError) {
                    if (e.code === 401 || e.code === 403) {
                        router.replace("/login")
                        return
                    }
                } 

                scheduleRetry()
            }
        }

        document.addEventListener("visibilitychange", onVisibility)
        window.addEventListener("online", reconnectNow)
        window.addEventListener("offline", clearRetry)

        connect()

        return () => {
            cancelled = true
            clearRetry()
            document.removeEventListener("visibilitychange", onVisibility)
            window.removeEventListener("online", reconnectNow)
            window.removeEventListener("offline", clearRetry)
            if (ws.current) {
                ws.current.close()
                ws.current = null
            }
        }
    }, [addError, handleConversationNew, handleMemberAdded, handleMemberRemoved, handleMessageDeleted, handleMessageEdited, handleMessageNew, handleReconnected, router])

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