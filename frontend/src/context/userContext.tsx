"use client"
import { createContext, Dispatch, SetStateAction, useContext, useEffect, useState } from "react"
import { User } from "@/types/http/user"
import { useErrorContext } from "./errorContext"
import { getUser } from "@/lib/api/user"

type UserContextType = {
    user: User | null
    setUser: Dispatch<SetStateAction<User | null>>
}

const UserContext = createContext<UserContextType | null>(null)

export function UserProvider({ children }: { children: React.ReactNode }) {
	const [ user, setUser ] = useState<User|null>(null)
	const { addError } = useErrorContext()

	useEffect(() => {
		async function loadUser() {
			const res = await getUser()

			if (res.error) {
				addError(res.error)
				setUser(null)
				return
			}

			setUser(res.data)
		}

		loadUser()
	}, [])

	return (
		<UserContext value={{ user, setUser }}>
			{children}
		</UserContext>
	)
}

export function useUserContext() {
	const ctx = useContext(UserContext)
	if (!ctx) throw new Error("user context should be used within user provider")
	return ctx
}