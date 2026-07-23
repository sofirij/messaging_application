import { Ticket } from "@/types/ws/ticket"
import { OutboundTypingStartPayload, OutboundTypingStopPayload } from "@/types/ws/typing"
import { MessageReadPayload } from "@/types/ws/message"
import { Event, EventMessageRead, EventTypingStart, EventTypingStop } from "@/types/ws/event"
import { fetchWithAuth } from "@/lib/api/fetch"

const apiVersion = "/api/v1"
const wsURL = "ws://"+process.env.NEXT_PUBLIC_API_URL+apiVersion+"/ws"
const ticketURL = "http://"+process.env.NEXT_PUBLIC_API_URL+apiVersion+"/ws/ticket"

export async function getWSTicket(): Promise<Ticket> {
    const url = `${ticketURL}/ticket`
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
    
    const url = new URL(wsURL)
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

export async function readMessage(ws: WebSocket, payload: MessageReadPayload): Promise<void> {
    const req : Event = {
        type: EventMessageRead,
        payload: payload
    }

    ws.send(JSON.stringify(req))
}

export async function startTyping(ws: WebSocket, conversationID: number): Promise<void> {
    const payload : OutboundTypingStartPayload = {
        conversation_id: conversationID,
    }
    
    const req : Event = {
        type: EventTypingStart,
        payload: payload
    }

    ws.send(JSON.stringify(req))
}

export async function stopTyping(ws: WebSocket, conversationID: number): Promise<void> {
    const payload : OutboundTypingStopPayload = {
        conversation_id: conversationID,
    }

    const req : Event = {
        type: EventTypingStop,
        payload: payload
    }

    ws.send(JSON.stringify(req))
}