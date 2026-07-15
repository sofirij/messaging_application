type FormProps = {
    onSubmit: () => void
    children?: React.ReactNode
}

export default function Form({ children, onSubmit }: FormProps) {
    function handleSubmit(e: React.SubmitEvent<HTMLFormElement>) {
        e.preventDefault()
        onSubmit()
    }

    return (
        <form onSubmit={handleSubmit}>
            {children}
        </form>
    )
}