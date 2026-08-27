"use client"
import { Ticket } from "@/types/ws/ticket"
import { OutboundTypingStartPayload, OutboundTypingStopPayload } from "@/types/ws/typing"
import { MessageReadPayload } from "@/types/ws/message"
import { Event, EventMessageRead, EventTypingStart, EventTypingStop } from "@/types/ws/event"
import { fetchWithCredentials } from "@/lib/api/fetch"
import { HTTPError, ErrorDetail } from "@/types/http/error"
import { httpProtocol, serverURL, wsProtocol } from "@/constants/defaults"

const apiVersion = "/api/v1"
const wsURL = wsProtocol+serverURL+apiVersion+"/ws"
const ticketURL = httpProtocol+serverURL+apiVersion+"/ws/ticket"

export async function getWSTicket(): Promise<Ticket> {
    const url = `${ticketURL}`
    const init = {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }
    
    const res = await fetchWithCredentials(url, init)

    const json = await res.json()

    if (!res.ok) {
        throw new HTTPError(json.error as ErrorDetail)
    }

    return json.data
}

export async function getWSConn(): Promise<WebSocket> {
    // get ticket for creating websocket connection
    const ticket = await getWSTicket()
    
    const url = new URL(wsURL)
    url.searchParams.set("ticket", ticket.ticket)

    return new WebSocket(url)
}

export async function readMessage(ws: WebSocket, conversation_id: number, message_id: number): Promise<void> {
    const payload: MessageReadPayload = {conversation_id, message_id}

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