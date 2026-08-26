"use client"
import { userQueryOptions, userSearchQueryOptions } from "@/query/user"
import { useQuery, useSuspenseQuery } from "@tanstack/react-query"
import { Dispatch, SetStateAction, useEffect, useState } from "react"
import Input from "@/components/ui/input"
import Image from "next/image"
import { User } from "@/types/http/user"
import { toggleSelected } from "@/utils/conversation"
import { defaultAvatarURL } from "@/constants/defaults"

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
        <div>
            <div>
                <Input label="search" type="text" onChange={(e) => setQuery(e.target.value)} />
                <button onClick={() => setSelected(new Set())}>Deselect all</button>
            </div>
            <div className="flex flex-col">
                {isPending || !data ? (
                    <div className="opacity-0 animate-[fadeIn_150ms_ease-in_200ms_forwards]">
                        <div className="animate-spin h-8 w-8 border-2 border-current border-t-transparent rounded-full opacity-1" />
                    </div>
                ): (
                    data.length === 0 ? (
                        <div>
                            <p>No results found</p>
                        </div>
                    ): (
                        <div>
                            {data.filter(user => user.id !== me.id).map(user => {
                                const avatarURL = user.avatar_url ?? defaultAvatarURL
                                return (
                                    <div key={user.id} onClick={() => toggleSelected(user, setSelected)} className="flex flex-row">
                                        <Image src={avatarURL} alt="User profile picture" width={20} height={20} />
                                        <p>{user.username}</p>
                                        {selected.has(user) && (
                                            <p>✅</p>
                                        )}
                                    </div>
                                )
                            })}
                        </div>
                    ) 
                )}
                {Array.from(selected).map(user => {
                    const avatarURL = user.avatar_url ?? defaultAvatarURL
                    return (
                        <div key={user.id} className="flex flex-row">
                            <Image src={avatarURL} alt="User profile picture" width={20} height={20} />
                            <p>{user.username}</p>
                            <button onClick={() => toggleSelected(user, setSelected)}>Deselect</button>
                        </div>
                    )
                })}
            </div>
        </div>
    )
}