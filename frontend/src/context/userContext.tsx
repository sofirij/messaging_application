"use client"
import { useQuery } from "@tanstack/react-query"
import { createContext, useContext } from "react"
import { userQueryOptions } from "@/context/queryProvider"
import { User } from "@/types/http/user"

const UserContext = createContext<User|null>(null)

export function UserProvider({children}: {children: React.ReactNode}) {
    const {data} = useQuery(userQueryOptions)

    if (!data) return

    return (
        <UserContext value={data}>
            {children}
        </UserContext>
    )
}

export function useUserContext(): User {
    const user = useContext(UserContext)
    if (!user) throw new Error("useUser must be used in an authenticated route")
    return user
}