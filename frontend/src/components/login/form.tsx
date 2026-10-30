"use client"
import Form from "@/components/ui/form"
import Input from "@/components/ui/input"
import { useLogin, useRegister } from "@/hooks/mutation/user"
import { useState } from "react"


export default function LoginForm() {
    const { handleLogin } = useLogin()
    const { handleRegister } = useRegister()
    const [mode, setMode] = useState<"login" | "register">("login")
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")

    return (
        <div>
            <div>
                <button onClick={() => setMode("login")}>LOGIN</button> 
                {"|"}
                <button onClick={() => setMode("register")}>REGISTER</button>
            </div>

            {mode === "login" ? (
                <Form onSubmit={() => handleLogin(username, password)} key={mode}>
                    <Input label="username" type="text" onChange={(e) => setUsername(e.target.value)}/>
                    <Input label="password" type="password" onChange={(e) => setPassword(e.target.value)}/>
                    <button>Login</button>
                </Form>       
            ) : (
                <Form onSubmit={() => handleRegister(username, password, confirmPassword)} key={mode}>
                    <Input label="username" type="text" onChange={(e) => setUsername(e.target.value)}/>
                    <Input label="password" type="password" onChange={(e) => setPassword(e.target.value)}/>
                    <Input label="confirm password" type="password" onChange={(e) => setConfirmPassword(e.target.value)}/>
                    <button>Register</button>
                </Form>
            )}
        </div>
    )
}