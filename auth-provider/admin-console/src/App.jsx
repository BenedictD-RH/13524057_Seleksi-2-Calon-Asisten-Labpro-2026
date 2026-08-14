import { useEffect, useState } from 'react'
import NavBar from './components/NavBar.jsx'
import ControlPanel from './components/ControlPanel.jsx'
import './App.css'

function App() {
  const [pathData, setPathData] = useState([])
  const [remountKey, setRemountKey] = useState(0)
  const remountData = () => {
    setRemountKey(remountKey + 1)
  }
  const path = window.location.pathname
  useEffect(() => {
    if (path != "") {
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
      }).catch((err) => {
        console.log(err)
      })
    }
  }, [remountKey])

  return (
    <>
      <div className='pageheader'>
        <NavBar></NavBar>
      </div>
      <div className='pagebody'>
        {<ControlPanel pageData={pathData} remount={remountData}></ControlPanel>}
      </div>
    </>
  )
}

export default App
