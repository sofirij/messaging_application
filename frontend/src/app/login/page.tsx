"use client"
import Input from "@/components/ui/input"
import Form from "@/components/ui/form"
import { useLogin, useRegister } from "@/hooks/auth"
import { useState } from "react"

export default function Page() {
    const { handleLogin } = useLogin()
    const { handleRegister } = useRegister()
    const [mode, setMode] = useState<"login" | "register">("login")
    const [username, setUsername] = useState<string>("")
    const [password, setPassword] = useState<string>("")
    const [confirmPassword, setConfirmPassword] = useState<string>("")

    return (
        <main>
            {/* tab toggle */}
            <div>
                <button onClick={() => setMode("login")}>LOGIN</button> 
                <p>/</p>
                <button onClick={() => setMode("register")}>REGISTER</button>
            </div>

            {/* conditionally show one or the other */}
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