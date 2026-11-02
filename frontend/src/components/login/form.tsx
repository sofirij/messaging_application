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

    function tabClassName(tab: "login" | "register") {
        return `flex-1 rounded-md py-2 text-sm font-medium transition-colors ${mode === tab ? "bg-white text-gray-900 shadow-sm" : "text-gray-500 hover:text-gray-800"}`
    }

    return (
        <div className="flex min-h-dvh items-center justify-center bg-gradient-to-b from-brand-dark to-brand p-4">
            <div className="w-full max-w-sm rounded-2xl bg-white p-8 shadow-xl">
                <div className="mb-6 text-center">
                    <div className="mx-auto mb-3 flex h-14 w-14 items-center justify-center rounded-full bg-brand text-2xl text-white">💬</div>
                    <h1 className="text-2xl font-semibold">{mode === "login" ? "Welcome back" : "Create an account"}</h1>
                    <p className="mt-1 text-sm text-gray-500">{mode === "login" ? "Log in to continue chatting" : "Sign up to start chatting"}</p>
                </div>

                <div className="mb-6 flex rounded-lg bg-gray-100 p-1">
                    <button className={tabClassName("login")} onClick={() => setMode("login")}>Login</button>
                    <button className={tabClassName("register")} onClick={() => setMode("register")}>Register</button>
                </div>

                {mode === "login" ? (
                    <Form onSubmit={() => handleLogin(username, password)} key={mode}>
                        <Input label="username" type="text" onChange={(e) => setUsername(e.target.value)}/>
                        <Input label="password" type="password" onChange={(e) => setPassword(e.target.value)}/>
                        <button className="btn-primary mt-2 w-full py-2.5">Login</button>
                    </Form>
                ) : (
                    <Form onSubmit={() => handleRegister(username, password, confirmPassword)} key={mode}>
                        <Input label="username" type="text" onChange={(e) => setUsername(e.target.value)}/>
                        <Input label="password" type="password" onChange={(e) => setPassword(e.target.value)}/>
                        <Input label="confirm password" type="password" onChange={(e) => setConfirmPassword(e.target.value)}/>
                        <button className="btn-primary mt-2 w-full py-2.5">Register</button>
                    </Form>
                )}
            </div>
        </div>
    )
}
