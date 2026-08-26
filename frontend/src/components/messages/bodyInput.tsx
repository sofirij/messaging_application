import { Dispatch, SetStateAction } from "react"

type InputProps = {
    setText: Dispatch<SetStateAction<string>>
    text: string
}

export default function BodyInput({setText, text}: InputProps) {
    return (
        <div>
            <input type="text" placeholder="send message" onChange={(e) => setText(e.target.value)} value={text}/>
        </div>
    )
}