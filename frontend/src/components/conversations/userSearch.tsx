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

export default function Search({selected, setSelected}: InputProps) {
    const { data: me } = useSuspenseQuery(userQueryOptions)
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
                    {data.filter(user => user.id !== me.id).map(user => {
                        const avatarURL = user.avatar_url ?? defaultAvatarURL
                        return (
                            <div key={user.id} onClick={() => toggleSelected(user, setSelected)}>
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