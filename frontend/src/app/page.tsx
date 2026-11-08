import { redirect } from "next/navigation"

export default function Page() {
    console.log("redirect to login from index")
    redirect("/login")
}
