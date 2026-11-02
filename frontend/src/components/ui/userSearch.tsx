"use client"
import { userQueryOptions, userSearchQueryOptions } from "@/query/client/user"
import { useQuery, useSuspenseQuery } from "@tanstack/react-query"
import { Dispatch, SetStateAction, useEffect, useState } from "react"
import Input from "@/components/ui/input"
import Image from "next/image"
import { User } from "@/types/http/user"
import { toggleSelected } from "@/lib/utils/conversation"
import { defaultAvatarURL } from "@/config/defaults"

type InputProps = {
    selected: Set<User>
    setSelected: Dispatch<SetStateAction<Set<User>>>
}

const queryDelay = 300

export default function UserSearch({selected, setSelected}: InputProps) {
    const { data: me } = useSuspenseQuery(userQueryOptions)
    const [query, setQuery] = useState("")
    const [debouncedQuery, setDebouncedQuery] = useState("")

    useEffect(() => {
        const timeout = setTimeout(() => setDebouncedQuery(query), queryDelay)
        return () => clearTimeout(timeout)
    }, [query])

    const { data, isPending } = useQuery({...userSearchQueryOptions(debouncedQuery), })
    return (
        <div className="flex flex-col gap-3">
            <Input label="search users" type="text" placeholder="Search by username" onChange={(e) => setQuery(e.target.value)} />

            {selected.size > 0 && (
                <div className="flex flex-wrap items-center gap-2">
                    {Array.from(selected).map(user => {
                        const avatarURL = user.avatar_url ?? defaultAvatarURL
                        return (
                            <div key={user.id} className="flex items-center gap-1.5 rounded-full bg-brand/10 py-1 pr-2 pl-1 text-sm text-brand-dark">
                                <Image src={avatarURL} alt="User profile picture" width={20} height={20} className="avatar h-5 w-5" />
                                <span>{user.username}</span>
                                <button onClick={() => toggleSelected(user, setSelected)} aria-label={`Deselect ${user.username}`} className="hover:text-red-600">✕</button>
                            </div>
                        )
                    })}
                    <button className="text-xs text-gray-500 hover:text-gray-800 hover:underline" onClick={() => setSelected(new Set())}>Deselect all</button>
                </div>
            )}

            <div className="flex max-h-60 flex-col overflow-y-auto rounded-lg border border-gray-200">
                {isPending || !data ? (
                    <div className="flex justify-center p-4 opacity-0 animate-[fadeIn_150ms_ease-in_200ms_forwards]">
                        <div className="animate-spin h-6 w-6 border-2 border-brand border-t-transparent rounded-full" />
                    </div>
                ): (
                    data.length === 0 ? (
                        <p className="p-4 text-center text-sm text-gray-500">No results found</p>
                    ): (
                        <div className="divide-y divide-gray-100">
                            {data.filter(user => user.id !== me.id).map(user => {
                                const avatarURL = user.avatar_url ?? defaultAvatarURL
                                const isSelected = selected.has(user)
                                return (
                                    <button key={user.id} onClick={() => toggleSelected(user, setSelected)} className={`flex w-full items-center gap-3 px-3 py-2 text-left text-sm transition-colors ${isSelected ? "bg-brand/10" : "hover:bg-gray-50"}`}>
                                        <Image src={avatarURL} alt="User profile picture" width={32} height={32} className="avatar h-8 w-8" />
                                        <span className="flex-1">{user.username}</span>
                                        {isSelected && (
                                            <span className="text-brand">✓</span>
                                        )}
                                    </button>
                                )
                            })}
                        </div>
                    )
                )}
            </div>
        </div>
    )
}
