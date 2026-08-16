"use client"
import Input from "@/components/ui/input"
import Form from "@/components/ui/form"
import { useLogin, useRegister } from "@/hooks/auth"
import { useEffect, useState } from "react"
import { useQuery } from "@tanstack/react-query"
import { useRouter } from "next/navigation"
import { userQueryOptions } from "@/query/user"

export default function Page() {
    const router = useRouter()
    const { handleLogin } = useLogin()
    const { handleRegister } = useRegister()
    const [mode, setMode] = useState<"login" | "register">("login")
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")
    const { data, isStale } = useQuery(userQueryOptions)

    useEffect(() => {
        if (data && !isStale) {
            router.push("/home")
        }
    }, [data, isStale, router])

    return (
        <main>
            <div>
                <button onClick={() => setMode("login")}>LOGIN</button> 
                {"|"}
                <button onClick={() => setMode("register")}>REGISTER</button>
            </div>

            {mode === "login" ? (
                <Form onSubmit={async () => await handleLogin(username, password)}>
                    <Input label="username" type="text" onChange={(e) => setUsername(e.target.value)}/>
                    <Input label="password" type="password" onChange={(e) => setPassword(e.target.value)}/>
                    <button>Login</button>
                </Form>       
            ) : (
                <Form onSubmit={async () => await handleRegister(username, password, confirmPassword)}>
                    <Input label="username" type="text" onChange={(e) => setUsername(e.target.value)}/>
                    <Input label="password" type="password" onChange={(e) => setPassword(e.target.value)}/>
                    <Input label="confirm password" type="password" onChange={(e) => setConfirmPassword(e.target.value)}/>
                    <button>Register</button>
                </Form>
            )}
        </main>
    )
}