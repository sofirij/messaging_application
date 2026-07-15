const minUsernameLength = 3
const maxUsernameLength = 15
const minPasswordLength = 8
const maxPasswordLength = 72
const usernameRegex = /^[a-zA-Z0-9_]+$/

export function validateUsername(username: string): string | null {
    if (username.trim() === "") return "username is required"
    if (username.length < minUsernameLength) return "username too short"
    if (username.length > maxUsernameLength) return "username too long"
    if (!usernameRegex.test(username)) return "username can only contain letters and numbers"
    return null 
}

export function validatePassword(password: string): string | null {
    if (password.trim() === "") return "password is required"
    if (password.length < minPasswordLength) return "password too short"
    if (password.length > maxPasswordLength) return "password too long"
    return null
}