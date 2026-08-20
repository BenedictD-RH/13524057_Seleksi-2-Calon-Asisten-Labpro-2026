import { useState, useEffect } from 'react'
import './App.css'
import LoginForm from './components/LoginForm'
import SessionForm from './components/SessionForm'
import { stringify } from 'uuid'

interface UserData {
  user_id: any
  name: string
  email: string
}

function App() {
  const [userData, setUserData] = useState<UserData | null>(null)
  const path: string = window.location.pathname

  const handleLogin = () => {
    window.location.pathname = '/login'
  }

  const handleLogout = () => {
    fetch('/server/logout', {
            method: 'POST',
            credentials: 'include'
        })
        .then((response: Response) => {
          if (response.ok) {
            setUserData(null)
          } else {
            throw new Error('Network response was not ok');
          }
        })
        .catch((err: unknown) => console.error(err));
  }

  useEffect(() => {
    fetch('/server/userinfo/central', {
      method: 'GET',
      credentials: 'include'
    }).then((res) => {
      if (res.ok) {
        return res.json();
      } else {
        throw new Error('Network response was not ok');
      }
      return null
    }).then((data) => {
      setUserData(data);
    }).catch((err) => {
      setUserData(null);
    })
  }, [])

  return (
    <>
      <div className="page">
        <div className="header">
          <h1 className="headertitle">Auth Provider</h1>
        </div>
        <div className="pagebody">
          {path === "/login" ? <LoginForm /> :
          (path === "/session" ? <SessionForm /> :
          (userData === null ? 
            <button className="referLogin" onClick={handleLogin}>Login</button> :
            <div className='userInfoMain'>
              <div className='userCredMain'>
                <div className='userIDMain'>{stringify(userData.user_id)}</div>
                <div className='userNameMain'>{userData.name}</div>
                <div className='userEmailMain'>{userData.email}</div>
              </div>
              <button className='logoutBtn' onClick={handleLogout}>Logout</button>
            </div>
          ))}
        </div>
      </div>
    </>
  )
}

export default App
