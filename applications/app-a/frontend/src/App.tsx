import { useEffect, useState } from 'react'
import LogPopup from './components/LogPopup'

interface UserData {
  name: string
  [key: string]: any
}

function App() {
  const [userData, setUserData] = useState<UserData | null>(null)
  const [unauthorizedMsg, setUnauthorizedMsg] = useState((new URLSearchParams(window.location.search)).get("errorMsg"))
  const [localSSID, setlocalSSID] = useState((new URLSearchParams(window.location.search)).get("local_ssid"))
  const [logPopupActive, setLogPopupActive] = useState(false)


  const handleLogin = () => {
    fetch('/backend/login', {
      redirect: 'manual',
      credentials: 'include'
    })
    .then(response => {
      console.log('Final URL 1:', response.url); 
      window.location.href = response.url
      if (response.redirected) {
        console.log('The request was redirected!');
        
      } else {
        return response.json()
      }
      
      return null
    }).then(data => {
      if (data != null) {
        setUnauthorizedMsg(data['error']['code'] + " | " + data['error']['message'])
      }
    }).catch(err => {
      console.error(err)
    });
  }

  const handleLogout = () => {
    fetch('/backend/logout', {
      method: 'POST',
      credentials: 'include'
    })
    .then(response => {
      if (response.ok) {
        setUserData(null)
      } else {
        throw new Error('Network response was not ok');
      }
    }).catch(err => {
      console.error(err)
    });
  }
  
  var pagecontent
  if (!userData) {
    pagecontent = <button className="loginbutton" onClick={handleLogin}>Login</button>
  } else {
    pagecontent = <div className="contenttext">Hello, {userData.name}</div>
  }

  const clearAllParamsNative = () => {
    const url = new URL(window.location.href);
    url.search = ''; // Empty out the query string
    
    window.history.replaceState(null, '', url.toString());
  };

  useEffect(() => {
    if (localSSID != null) {
      fetch('/backend/session?' + new URLSearchParams(window.location.search), {
        method: 'GET',
        redirect: 'manual',
        credentials: 'include'
      }).then((res) => {
        if (res.ok) {
          return res.json();
        } else if (res.status == 204) {
          return null
        } else {
          throw new Error('Network response was not ok');
        }
        return null
      }).then((data) => {
        clearAllParamsNative()
        fetch('/backend/users', {
          method: 'GET',
          redirect: 'manual',
          credentials: 'include'
        }).then((res) => {
          if (res.ok) {
            return res.json();
          } else if (res.status == 204) {
            return null
          } else {
            throw new Error('Network response was not ok');
          }
          return null
        }).then((data) => {
          setUserData(data);
        }).catch((err) => {
          setUserData(null);
        })
      }).catch((err) => {
        clearAllParamsNative()
      })
    } 
    else {
      clearAllParamsNative()
      fetch('/backend/users', {
        method: 'GET',
        redirect: 'manual',
        credentials: 'include'
      }).then((res) => {
        if (res.ok) {
          return res.json();
        } else if (res.status == 204) {
          return null
        } else {
          throw new Error('Network response was not ok');
        }
        return null
      }).then((data) => {
        setUserData(data);
      }).catch((err) => {
        setUserData(null);
      })
    }
  }, [])
  


  return (
    <div className="page">
      <div className="pageheader">
        <h1 className="appname">App-A</h1>
        {userData !== null ? <button className='logoutBtn' onClick={handleLogout}>Logout</button>: <></>}
      </div>
      {logPopupActive ? <LogPopup></LogPopup> : <></>}
      <div className="pagebody">
        {unauthorizedMsg != null ? <div className='unauthorizedMsg'>{unauthorizedMsg}</div> : <></>}
        {pagecontent}
      </div>
      <button onClick={() => setLogPopupActive(!logPopupActive)} className='logButton'>{logPopupActive ? 'Close Logs' : 'See Logs'}</button>
    </div>
  )
}

export default App