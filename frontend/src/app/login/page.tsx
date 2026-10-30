import { verifySession } from "@/lib/utils/cookies"
import { redirect } from "next/navigation"
import LoginForm from "@/components/login/form"

export default async function LoginPage() {
    const {err} = await verifySession()

    if (!err) {
        redirect("/home")
    }

    return (
        <LoginForm/>
    )
}