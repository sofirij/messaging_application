const minUsernameLength = 3
const maxUsernameLength = 15
const minPasswordLength = 8
const maxPasswordLength = 72
const usernameRegex = /^[a-zA-Z0-9_]+$/

export function validateUsername(username: string) {
    if (username.trim() === "") throw new Error("username is required")
    if (username.length < minUsernameLength) throw new Error("username too short")
    if (username.length > maxUsernameLength) throw new Error("username too long")
    if (!usernameRegex.test(username)) throw new Error("username can only contain letters and numbers") 
}

export function validatePassword(password: string) {
    if (password.trim() === "") throw new Error("password is required")
    if (password.length < minPasswordLength) throw new Error("password too short")
    if (password.length > maxPasswordLength) throw new Error("password too long")
}