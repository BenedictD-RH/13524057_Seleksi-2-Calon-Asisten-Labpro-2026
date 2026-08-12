import { useState } from 'react'
import './LoginForm.css'

function LoginForm() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")

    const handleEmailChange = (event) => {
        setEmail(event.target.value)
    }
    const handlePasswordChange = (event) => {
        setPassword(event.target.value)
    }
    const handleLogin = () => {
        fetch('/server/login?' + new URLSearchParams(window.location.search), {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({
                "email": email,
                "password": password
            }),
            redirect: 'follow'
        })
        .then(response => {
        if (response.redirected) {
            console.log('The request was redirected!');
        }
        
        console.log('Final URL 1:', response.url); 
        window.location.href = response.url
        })
        .catch(err => console.error(err));
    }

    return (
        <div class="container">
            <h class="loginheader">User Login</h>
            <div class="inputcontainer">
                <h class="inputheader">Email:</h>
                <input class="logininput" placeholder="Enter Email" type="email" onChange={handleEmailChange}></input>
            </div>
            <div class="inputcontainer">
                <h class="inputheader">Password:</h>
                <input class="logininput" placeholder="Enter Password" type="password" onChange={handlePasswordChange}></input>
            </div>
            <div class="btncontainer">
                <button class="loginbutton" onClick={handleLogin}>Login</button>
            </div>
        </div>
    )
}

export default LoginForm