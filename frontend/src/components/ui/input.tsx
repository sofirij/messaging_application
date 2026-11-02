"use client"
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
            <label className={labelClassName ?? "flex flex-col gap-1.5 text-sm font-medium text-gray-700 capitalize"}>
                {label}
                <input type={type} className={inputClassName ?? "text-field font-normal normal-case"} placeholder={placeholder} onChange={onChange}></input>
            </label>
        </>
    )
}
