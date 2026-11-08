import { verifySession } from "@/lib/utils/cookies"
import { redirect } from "next/navigation"
import ProfilePage from "@/components/profile/page"

export default async function Page() {
    const {err} = await verifySession()
        
    if (err) {
        console.log("redirect to login from profile")
        redirect("/login")
    }

    return <ProfilePage/>
}