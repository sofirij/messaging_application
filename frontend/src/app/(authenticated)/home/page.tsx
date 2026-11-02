import Link from "next/link"

export default function Home() {
    return (
        <main className="flex h-full flex-col items-center justify-center gap-4 p-6 text-center">
            <div className="text-5xl">💬</div>
            <h1 className="text-2xl font-semibold">Welcome to Chat</h1>
            <p className="max-w-sm text-gray-500">Pick up where you left off or start a new conversation.</p>
            <Link href="/conversations" className="btn-primary">Go to conversations</Link>
        </main>
    )
}
