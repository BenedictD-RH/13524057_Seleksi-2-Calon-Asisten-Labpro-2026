import { useEffect, useState } from 'react'
import NavBar from './components/NavBar'
import ControlPanel from './components/ControlPanel'
import './App.css'

function App() {
  const [pathData, setPathData] = useState<any[]>([])
  const [remountKey, setRemountKey] = useState(0)
  const [authorized, setAuthorized] = useState(false)
  const [unauthorizedMsg, setUnauthorizedMsg] = useState("")
  const remountData = () => {
    setRemountKey(remountKey + 1)
  }
  const path = window.location.pathname
  useEffect(() => {
    fetch('/server/authorize/administrator', {
      headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
      },
      credentials: 'include',
    }).then(response => {
      return response.json()
    }).then(data => {
      if (data['status'] == "authorized") {
        if (path != "/") {
          fetch('/server' + path)
          .then((res) => {
            if (res.ok) {
              return res.json();
            } else {
              throw new Error('Network response was not ok');
            }
            return null
          }).then((data) => {
            console.log(data)
            setPathData(data)
            setAuthorized(true)
          }).catch((err) => {
            console.log(err)
          })
        }
      } else {
        setUnauthorizedMsg(data['error']['message'])
      }
    }).catch(error => console.log(error)) 

    

  }, [remountKey])

  return (
    <>
      {authorized ? 
      <div>
        <div className='pageheader'>
          <NavBar></NavBar>
        </div>
        <div className='pagebody'>
          {<ControlPanel pageData={pathData} remount={remountData}></ControlPanel>}
        </div>
      </div> :
      <div>
        <div className='unauthorizedTitle'>Unauthorized</div>
        <div className='unauthorizedMsg'>{unauthorizedMsg}</div>
      </div>
      }
    </>
  )
}

export default App