"use client"

import { getWSConn, readMessage, startTyping, stopTyping } from "@/lib/api/ws"
import { conversationMemberQueryOptions, conversationMemberRefetchOptions, conversationRefetchOptions } from "@/query/conversation"
import { messageQueryOptions } from "@/query/message"
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

    useEffect(() => {
        let cancelled = false
        async function connect() {
            const conn = await getWSConn()

            if (cancelled) return

            conn.onopen = () => { console.log("opened ws conn") }

            conn.onmessage = (event) => {
                console.log("message received")
                
                const msg: Event  = JSON.parse(event.data)

                switch (msg.type) {
                    case EventMessageNew: {
                        const conversationID = msg.payload.conversation_id
                        const old = queryClient.getQueryData(messageQueryOptions(conversationID).queryKey)
                        if (!old) return
            
                        queryClient.setQueryData(messageQueryOptions(conversationID).queryKey, {
                            ...old,
                            pages: old.pages.map((page, i) => {
                                if (i === old.pages.length - 1) {
                                    // sort in the event of network jitter
                                    // ws message can come in before http response
                                    return {...page, messages: [...page.messages.filter(message => message.id !== msg.payload.id), msg.payload].sort((a, b) => a.id - b.id)}
                                }
                                return page
                            })
                        })
                        break
                    }
                    case EventMessageDeleted: {
                        const conversationID = msg.payload.conversation_id
                        const old = queryClient.getQueryData(messageQueryOptions(conversationID).queryKey)
                        if (!old) return

                        queryClient.setQueryData(messageQueryOptions(conversationID).queryKey, {
                            ...old,
                            pages: old.pages.map(page => {
                                return {
                                    ...page,
                                    messages: page.messages.map(message => {
                                        if (message.id === msg.payload.message_id) {
                                            return {...message, body: null, attachments: [], deleted: true}
                                        }
                                        return message
                                    })
                                }
                            })
                        })
                        break
                    }
                    case EventMessageEdited: {
                        const conversationID = msg.payload.conversation_id
                        const old = queryClient.getQueryData(messageQueryOptions(conversationID).queryKey)
                        if (!old) return

                        queryClient.setQueryData(messageQueryOptions(conversationID).queryKey, {
                            ...old,
                            pages: old.pages.map(page => {
                                return {
                                    ...page,
                                    messages: page.messages.map(message => {
                                        if (message.id === msg.payload.message_id) {
                                            return {...message, body: msg.payload.body}
                                        }
                                        return message
                                    })
                                }
                            })
                        })
                        break
                    }
                    case EventMemberAdded: {
                        const conversationID = msg.payload.conversation_id
                        const userID = msg.payload.user.id

                        if (userID === me.id) {
                            queryClient.refetchQueries(conversationRefetchOptions)
                            queryClient.refetchQueries(conversationMemberRefetchOptions(conversationID))
                        } else {
                            const old = queryClient.getQueryData(conversationMemberQueryOptions(conversationID).queryKey)

                            if (!old) return

                            queryClient.setQueryData(conversationMemberQueryOptions(conversationID).queryKey, {
                                data: {...old.data, [userID]: msg.payload.user},
                                order: [...old.order.filter(id => id !== userID), userID].sort((a, b) => a - b)
                            })
                        }
        
                        break
                    }
                    case EventMemberRemoved: {
                        const conversationID = msg.payload.conversation_id
                        const userID = msg.payload.user_id

                        if (userID === me.id) {
                            queryClient.refetchQueries(conversationRefetchOptions)
                            queryClient.refetchQueries(conversationMemberRefetchOptions(conversationID))
                        } else {
                            const old = queryClient.getQueryData(conversationMemberQueryOptions(conversationID).queryKey)

                            if (!old) return


                            const remainingData = {...old.data}
                            delete remainingData[userID]

                            queryClient.setQueryData(conversationMemberQueryOptions(conversationID).queryKey, {
                                data: remainingData,
                                order: old.order.filter(id => id !== userID)
                            })
                        }
                        break
                    }
                    case EventConversationNew: {
                        queryClient.refetchQueries(conversationRefetchOptions)
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
    }, [me.id, queryClient])

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