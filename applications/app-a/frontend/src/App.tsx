import { useEffect, useState } from 'react'

interface UserData {
  name: string
  [key: string]: any
}

function App() {
  const [userData, setUserData] = useState<UserData | null>(null)
  const [unauthorizedMsg, setUnauthorizedMsg] = useState("")
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
  
  var pagecontent
  if (!userData) {
    pagecontent = <button className="loginbutton" onClick={handleLogin}>Login</button>
  } else {
    pagecontent = <div className="contenttext">Hello, {userData.name}</div>
  }

  useEffect(() => {
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
  }, [])
  


  return (
    <div className="page">
      <div className="pageheader">
        <h1 className="appname">App-A</h1>\
      </div>
      <div className="pagebody">
        {unauthorizedMsg != "" ? <div className='unauthorizedMsg'>{unauthorizedMsg}</div> : <></>}
        {pagecontent}
      </div>
    </div>
  )
}

export default App