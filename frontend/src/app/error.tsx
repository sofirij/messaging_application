"use client"
import ErrorPage from "@/components/ui/errorPage"

export default function Error({ reset }: { error: Error; reset: () => void }) {
    return <ErrorPage reset={reset} fullScreen />
}
