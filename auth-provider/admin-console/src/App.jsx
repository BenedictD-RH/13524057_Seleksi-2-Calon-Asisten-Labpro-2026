import { useEffect, useState } from 'react'
import NavBar from './components/NavBar.jsx'
import ControlPanel from './components/ControlPanel.jsx'
import './App.css'

function App() {
  const [pathData, setPathData] = useState([])
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
  }, [])

  return (
    <>
      <div className='pageheader'>
        <NavBar></NavBar>
      </div>
      <div className='pagebody'>
        {pathData.length > 0 ? <ControlPanel pageData={pathData}></ControlPanel> : <></>}
      </div>
    </>
  )
}

export default App
