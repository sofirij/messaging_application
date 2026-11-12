"use client"

import { defaultAvatarURL } from "@/config/defaults"
import { User } from "@/types/http/user"
import Image from "next/image"

type InputProps = {
    member: User,
    children: React.ReactNode
}

export default function Toast({member, children}: InputProps) {
    return (
        <div className="flex items-center gap-3 py-2">
            <Image src={member.avatar_url ?? defaultAvatarURL} alt="User profile picture" width={32} height={32} className="avatar h-8 w-8"/>
            <p className="flex-1 text-sm">{member.username}</p>
            {children}
        </div>
    )
}