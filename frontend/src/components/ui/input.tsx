type InputProps = {
    label: string
    type: "password" | "text"
    labelClassName?: string
    inputClassName?: string
    placeholder?: string
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void
    value?: string
}

export default function Input({label, type, labelClassName, inputClassName, placeholder, onChange, value}: InputProps) {
    return (
        <>
            <label className={labelClassName}>
                {label}
                <input type={type} className={inputClassName} placeholder={placeholder} onChange={onChange} value={value}></input>
            </label>
        </>
    )
}