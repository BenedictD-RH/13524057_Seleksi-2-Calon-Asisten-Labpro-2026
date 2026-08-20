import { useState } from 'react'
import './App.css'
import LoginForm from './components/LoginForm'
import SessionForm from './components/SessionForm'

function App() {
  const [unauthorizedMsg, setUnauthorizedMsg] = useState<string>("")
  const path: string = window.location.pathname

  return (
    <>
      <div className="page">
        <div className="header">
          <h1 className="headertitle">Auth Provider</h1>
        </div>
        <div className="pagebody">
          {path === "/login" ? <LoginForm /> :
          (path === "/session" ? <SessionForm /> :
          (path === "/unauthorized" ? 
          <>
            <div className="unauthorizedTitle">Unauthorized</div>
            <div className="unauthorizedMsg">{unauthorizedMsg}</div>
          </> : 
          null))}
        </div>
      </div>
    </>
  )
}

export default App
