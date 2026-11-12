import { SetState } from "@/types/react"

type InputProps = {
    setText: SetState<string>
    text: string
    submitHandler: () => void
}

export default function BodyInput({setText, text, submitHandler}: InputProps) {
    function keyHandler(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key === "Enter") {
            submitHandler()
        }
    }

    return (
        <div className="flex-1">
            <input type="text" className="w-full rounded-full border border-gray-200 bg-white px-4 py-2.5 text-sm outline-none focus:border-brand" placeholder="Type a message" onChange={(e) => setText(e.target.value)} value={text} onKeyDown={keyHandler}/>
        </div>
    )
}
