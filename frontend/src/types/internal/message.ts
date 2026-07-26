import { Message } from "@/types/http/message"

export type MessageRecord = {
    data: Record<number, Message>,
    order: number[]
}