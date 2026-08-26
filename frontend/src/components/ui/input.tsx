type InputProps = {
    label: string
    type: "password" | "text"
    labelClassName?: string
    inputClassName?: string
    placeholder?: string
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void
}

export default function Input({label, type, labelClassName, inputClassName, placeholder, onChange}: InputProps) {
    return (
        <>
            <label className={labelClassName}>
                {label}
                <input type={type} className={inputClassName} placeholder={placeholder} onChange={onChange}></input>
            </label>
        </>
    )
}