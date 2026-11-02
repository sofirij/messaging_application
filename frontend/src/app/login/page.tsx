import { verifySession } from "@/lib/utils/cookies"
import LoginForm from "@/components/login/form"
import { redirect } from "next/navigation"

export default async function LoginPage() {
    const {err} = await verifySession()

    if (!err) {
        redirect("/home")
    }

    return (
        <LoginForm/>
    )
}