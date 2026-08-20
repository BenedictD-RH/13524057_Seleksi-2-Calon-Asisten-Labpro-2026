import { useState } from 'react'
import './LoginForm.css'

function LoginForm() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [errorMessage, setErrorMessage] = useState("")

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
            return response.json()
        }).then(data => {
            if (data['error'] != null) {
                
                var error = data['error']
                console.log(error)
                if (error['code'] == "LOGIN_FAILED") {
                    setErrorMessage(error['message'])
                } else if (error['code'] == "INVALID_GRANT") {
                    //redirect back to page but with invalid grant msg
                }
            } else {
                window.location.href = data['redirect']
            }
        })
        .catch(err => console.error(err));
    }

    return (
        <div className="container">
            <h className="loginheader">User Login</h>
            <div className="inputcontainer">
                <h className="inputheader">Email:</h>
                <input className="logininput" placeholder="Enter Email" type="email" onChange={handleEmailChange}></input>
            </div>
            <div className="inputcontainer">
                <h className="inputheader">Password:</h>
                <input className="logininput" placeholder="Enter Password" type="password" onChange={handlePasswordChange}></input>
            </div>
            <div className="btncontainer">
                <button className="loginbutton" onClick={handleLogin}>Login</button>
            </div>
            <div className="errorMsg">{errorMessage}</div>
        </div>
    )
}

export default LoginForm