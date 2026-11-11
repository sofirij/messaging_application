"use client"
import { useQueryErrorResetBoundary } from "@tanstack/react-query"
import { useRouter } from "next/navigation"
import { useTransition } from "react"
import Link from "next/link"

type ErrorPageProps = {
    reset: () => void
    fullScreen?: boolean
}

export default function ErrorPage({ reset, fullScreen = false }: ErrorPageProps) {
    const { reset: resetQueries } = useQueryErrorResetBoundary()
    const router = useRouter()
    const [isPending, startTransition] = useTransition()

    function retry() {
        resetQueries()
        startTransition(() => {
            router.refresh()
            reset()
        })
    }

    return (
        <div className={`flex items-center justify-center p-4 ${fullScreen ? "min-h-dvh" : "h-full"}`}>
            <div className="w-full max-w-sm rounded-2xl bg-white p-8 text-center shadow-xl">
                <div className="mx-auto mb-3 flex h-14 w-14 items-center justify-center rounded-full bg-red-50 text-2xl font-semibold text-red-600">!</div>
                <h1 className="text-2xl font-semibold">Something went wrong</h1>
                <p className="mt-1 text-sm text-gray-500">We couldn&apos;t load this page. Check your connection and try again.</p>

                <div className="mt-6 flex flex-col gap-2">
                    <button className="btn-primary w-full py-2.5" onClick={retry} disabled={isPending}>
                        {isPending ? "Retrying…" : "Try again"}
                    </button>
                    <Link href="/home" className="btn-ghost w-full py-2.5">Go to home</Link>
                </div>
            </div>
        </div>
    )
}
