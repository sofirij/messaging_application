function showSection(name) {
    document.querySelectorAll('.section').forEach(s => s.classList.remove('active'))
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'))
    document.getElementById(name).classList.add('active')
    event.target.classList.add('active')
}

function login() {
    const username = document.getElementById("login-username").value
    const password = document.getElementById("login-password").value

    console.log(username)
    console.log(password)
}

function main() {
    document.getElementById("login-submit").addEventListener("on-click")
}
