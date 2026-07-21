import { Ticket } from "@/types/ws/ticket"
import { OutboundTypingStartPayload, OutboundTypingStopPayload } from "@/types/ws/typing"
import { MessageReadPayload } from "@/types/ws/message"
import { Event, EventTypingStart, EventTypingStop } from "@/types/ws/event"
import { fetchWithAuth } from "@/lib/api/fetch"

const api_version = "/api/v1"
const ws_url = "ws://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/ws"
const ticket_url = "http://"+process.env.NEXT_PUBLIC_API_URL+api_version+"/ws/ticket"

export async function getWSTicket(): Promise<Ticket> {
    const url = `${ticket_url}/ticket`
    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    const res = await fetchWithAuth(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new Error(json.error.message)
    }

    return json.data
}

export async function getWSConn(): Promise<WebSocket> {
    // get ticket for creating websocket connection
    const ticket = await getWSTicket()
    
    const url = new URL(ws_url)
    url.searchParams.set("ticket", ticket.ticket)

    const ws = new WebSocket(url)

    ws.onopen = () => { console.log("opened ws conn") }
    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data)
        console.log(msg.Type)
    }
    ws.onclose = () => { console.log("conn closed") }
    ws.onerror = (error) => { console.log(error) }

    return ws
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