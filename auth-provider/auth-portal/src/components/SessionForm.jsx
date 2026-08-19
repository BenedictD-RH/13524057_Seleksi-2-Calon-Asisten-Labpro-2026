import "./SessionForm.css"
import { useState, useEffect } from 'react'


function SessionForm() {
    const [sessionData, setSessionData] = useState(null)

    const handleExistingSession = () => {
        fetch('/server/session/use?' + new URLSearchParams(window.location.search) , {
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            origin: 'http://localhost:8081',
            credentials: 'include',
            redirect: 'follow'
        }).then(response => {
            return response.json()
        }).then(data => {
            if (data != null) {
                window.location.href = data['redirect']
            }
        }).catch(err => console.error(err));
    }

    const handleNewSession = () => {
        window.location.pathname = "/login"
    }

    useEffect(() => {
        fetch('/server/session', {
            method: 'GET',
            origin: 'http://localhost:8081',
            credentials: 'include'
        }).then(response => {
            if (response.ok) {
                return response.json()
            } else {
                throw new Error('Network response was not ok');
            }
        }).then(data => {
            console.log(data)
            setSessionData(data)
        }).catch(err => console.error(err));
    }, [])

    return (
        <div className="sessionContainer">
            <h1 className="sessionHeader">Session Found</h1>
            {
                sessionData != null ? 
                <div className="sessionChoice">
                    <div className="userCred">
                        <h1 className="userName">{sessionData['name']}</h1>
                        <h1 className="userEmail">{sessionData['email']}</h1>
                    </div>
                    <button className="useBtn" onClick={handleExistingSession}>{">"}</button>
                </div> : 
                <></>
            }
            <button className="anotherSessionBtn" onClick={handleNewSession}>Use Another Session</button>
        </div>
    )
}

export default SessionForm