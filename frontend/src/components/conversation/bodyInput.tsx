import { Dispatch, SetStateAction } from "react"

type InputProps = {
    setText: Dispatch<SetStateAction<string>>
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
        <div>
            <input type="text" placeholder="send message" onChange={(e) => setText(e.target.value)} value={text} onKeyDown={keyHandler}/>
        </div>
    )
}