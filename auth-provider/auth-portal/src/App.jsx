import { useState } from 'react'
import './App.css'
import LoginForm from './components/LoginForm.jsx'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <div class="page">
        <div class="header">
          <h class="headertitle">Auth Provider</h>
        </div>
        <div class="pagebody">
          <LoginForm></LoginForm>
        </div>
      </div>
    </>
  )
}

export default App
