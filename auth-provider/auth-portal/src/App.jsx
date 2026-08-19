import { useState } from 'react'
import './App.css'
import LoginForm from './components/LoginForm.jsx'
import SessionForm from './components/SessionForm.jsx'

function App() {
  const [count, setCount] = useState(0)
  const [unauthorizedMsg, setUnauthorizedMsg] = useState("")
  const path = window.location.pathname

  return (
    <>
      <div class="page">
        <div class="header">
          <h class="headertitle">Auth Provider</h>
        </div>
        <div class="pagebody">
          {path == "/login" ? <LoginForm></LoginForm> :
          (path == "/session" ? <SessionForm></SessionForm> :
          (path == "/unauthorized" ? 
          <>
            <div className='unauthorizedTitle'>Unauthorized</div>
            <div className='unauthorizedMsg'>{unauthorizedMsg}</div>
          </> : 
          <></>))}
        </div>
      </div>
    </>
  )
}

export default App
