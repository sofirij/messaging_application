"use client"
import { conversationQueryOptions } from "@/context/queryProvider";
import { useQuery } from "@tanstack/react-query";

export default function Conversations() {
    const { data, isPending } = useQuery(conversationQueryOptions)

    if (isPending || !data) {
        return "Loading..."
    }

    return (
        <main>
            {data.length > 0 ? (
                <></>
            ) : (
                <div>
                    <p>Start your first conversation</p>
                    <button>New Conversation</button>
                </div>
            )}
        </main>
    )
}