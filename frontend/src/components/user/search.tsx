"use client"
import { userSearchQueryOptions } from "@/query/user"
import { useQuery } from "@tanstack/react-query"
import { Dispatch, SetStateAction, useEffect, useState } from "react"
import Input from "@/components/ui/input"
import Image from "next/image"
import { User } from "@/types/http/user"

type InputProps = {
    selected: Set<User>,
    setSelected: Dispatch<SetStateAction<Set<User>>>
}

const defaultAvatarURL = "fb24fc90-5e53-4972-b880-3edd0f8ccc64.jpg"
const queryDelay = 300

export function Search({selected, setSelected}: InputProps) {
    function toggleSelected(user: User) {
        setSelected(prev => {
            const next = new Set(prev)
            if (next.has(user)) {
                next.delete(user)
            } else {
                next.add(user)
            }
            return next
        })
    }

    const [query, setQuery] = useState("")
    const [debouncedQuery, setDebouncedQuery] = useState("")

    useEffect(() => {
        const timeout = setTimeout(() => setDebouncedQuery(query), queryDelay)
        return () => clearTimeout(timeout)
    }, [query])

    const { data, isPending } = useQuery(userSearchQueryOptions(debouncedQuery))
    return (
        <div>
            <div>
                <Input label="search" type="text" onChange={(e) => setQuery(e.target.value)} />
                <button onClick={() => setSelected(new Set())}>Deselect all</button>
            </div>
            {isPending || !data ? (
                <div>
                    <p>Loading...</p>
                </div>
            ): (
                <div>
                    {data.map(user => {
                        const avatarURL = user.avatar_url ?? defaultAvatarURL
                        return (
                            <div key={user.id} onClick={() => toggleSelected(user)}>
                                <Image src={avatarURL} alt="User profile picture" width={20} height={20} />
                                <p>{user.username}</p>
                                {selected.has(user) && (
                                    <p>✅</p>
                                )}
                            </div>
                        )
                    })}
                </div>
            )}
        </div>
    )
}