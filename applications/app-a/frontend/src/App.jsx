import { useEffect, useState } from 'react'

const getCookie = (name) => {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
  return null;
};

function App() {
  const [userData, setUserData] = useState(null)
  const handleLogin = () => {
    fetch('/backend/login', {
      redirect: 'manual',
      origin: 'http://localhost:8692',
      credentials: 'include'
    })
    .then(response => {
      if (response.redirected) {
        console.log('The request was redirected!');
      }
      
      console.log('Final URL 1:', response.url); 
      window.location.href = response.url
    })
    .catch(err => console.error(err));
  }
  
  var pagecontent
  if (!userData) {
    pagecontent = <button className="loginbutton" onClick={handleLogin}>Login</button>
  } else {
    pagecontent = <h1 className="contenttext">Hello, {userData.name}</h1>
  }

  useEffect(() => {
    fetch('/backend/users', {
      method: 'GET',
      redirect: 'manual',
      origin: 'http://localhost:8692',
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
        {pagecontent}
      </div>
    </div>
  )
}

export default App
