import { useState, type ChangeEvent } from 'react'
import './LoginForm.css'

interface LoginError {
    code: "LOGIN_FAILED" | "INVALID_GRANT";
    message: string;
}

interface LoginResponse {
    error?: LoginError | null;
    redirect?: string;
    session_token?: string;
    status?: string;
}

interface ReturnURL {
    client_url?: string;
}

function LoginForm() {
    const [email, setEmail] = useState<string>("")
    const [password, setPassword] = useState<string>("")
    const [errorMessage, setErrorMessage] = useState<string>("")


    const handleEmailChange = (event: ChangeEvent<HTMLInputElement>): void => {
        setEmail(event.target.value)
    }

    const handlePasswordChange = (event: ChangeEvent<HTMLInputElement>): void => {
        setPassword(event.target.value)
    }

    const handleLogin = (): void => {
        fetch('/server/login?' + new URLSearchParams(window.location.search), {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({
                email: email,
                password: password
            }),
            redirect: 'follow'
        })
        .then((response: Response) => response.json())
        .then((data: LoginResponse) => {
            if (data.error != null) {
                const error = data.error
                console.log(error)
                if (error.code === "LOGIN_FAILED") {
                    setErrorMessage(error.message)
                } else if (error.code === "INVALID_GRANT") {
                    const searchParams = new URLSearchParams(window.location.search)
                    fetch('/server/return?client_id=' + searchParams.get('client_id'))
                    .then((response: Response) => response.json())
                    .then((data: ReturnURL) => {
                        if (data.client_url != "") {
                            window.location.href = data.client_url + "?errorMsg=" + error.message
                        }
                    })
                }
            } else if (data.redirect) {
                window.location.href = data.redirect + "?local_ssid=" + data.session_token
            } else if (data.status == "authorized") {
                window.location.pathname = ""
            }
        })
        .catch((err: unknown) => console.error(err));
    }

    return (
        <div className="container">
            <h1 className="loginheader">User Login</h1>
            
            <div className="inputcontainer">
                <span className="inputheader">Email:</span>
                <input 
                    className="logininput" 
                    placeholder="Enter Email" 
                    type="email" 
                    value={email}
                    onChange={handleEmailChange}
                />
            </div>
            
            <div className="inputcontainer">
                <span className="inputheader">Password:</span>
                <input 
                    className="logininput" 
                    placeholder="Enter Password" 
                    type="password" 
                    value={password}
                    onChange={handlePasswordChange}
                />
            </div>
            
            <div className="btncontainer">
                <button className="loginbutton" onClick={handleLogin}>Login</button>
            </div>
            
            <div className="errorMsg">{errorMessage}</div>
        </div>
    )
}

export default LoginForm
