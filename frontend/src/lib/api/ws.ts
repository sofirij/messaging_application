import { Ticket } from "@/types/ws/ticket"
import { OutboundTypingStartPayload, OutboundTypingStopPayload } from "@/types/ws/typing"
import { MessageReadPayload } from "@/types/ws/message"
import { Event, EventTypingStart, EventTypingStop } from "@/types/ws/event"
import { ApiResult } from "@/types/http/response"

const api_version = "/api/v1"
const ws_url = "ws://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/ws"
const ticket_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/ws/ticket"

export async function getWSTicket(): Promise<ApiResult<Ticket>> {
    const res = await fetch(`${ticket_url}/ticket`, {
        method: "GET",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
            
        }
    })

    const json = await res.json()

    if (!res.ok) {
        return {error: json.error.message, data: null}
    }

    return {error: null, data: json.data}
}

export async function getWSConn(): Promise<ApiResult<WebSocket>> {
    // get ticket for creating websocket connection
    const res = await getWSTicket()

    if (res.error) {
        return {error: res.error, data: null}
    }

    const ticket = res.data?.ticket
    
    const url = new URL(ws_url)
    url.searchParams.set("ticket", ticket ? ticket : "")

    const ws = new WebSocket(url)

    ws.onopen = () => { console.log("opened ws conn") }
    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data)
        console.log(msg.Type)
    }
    ws.onclose = () => { console.log("conn closed") }
    ws.onerror = (error) => { console.log(error) }

    return {error: null, data: ws}
}

export async function readMessage(ws: WebSocket, conversationID: number, messageID: number): Promise<void> {
    const req : MessageReadPayload = {
        conversation_id: conversationID,
        message_id: messageID
    }

    ws.send(JSON.stringify(req))
}

export async function startTyping(ws: WebSocket, conversationID: number): Promise<void> {
    const payload : OutboundTypingStartPayload = {
        conversation_id: conversationID,
    }
    
    const req : Event = {
        Type: EventTypingStart,
        Payload: payload
    }

    ws.send(JSON.stringify(req))
}

export async function stopTyping(ws: WebSocket, conversationID: number): Promise<void> {
    const payload : OutboundTypingStopPayload = {
        conversation_id: conversationID,
    }

    const req : Event = {
        Type: EventTypingStop,
        Payload: payload
    }

    ws.send(JSON.stringify(req))
}