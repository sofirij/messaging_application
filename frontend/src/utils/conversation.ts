import { User } from "@/types/http/user"
import { Dispatch, SetStateAction } from "react"

export function toggleSelected(user: User, setSelected: Dispatch<SetStateAction<Set<User>>>) {
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