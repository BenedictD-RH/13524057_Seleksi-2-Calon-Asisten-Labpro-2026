import './LogPopup.css'
import { useState, useEffect } from 'react'


function LogPopup() {
    const [logData, setLogData] = useState<string[] | null>(null)
    const [logType, setLogType] = useState("/activity")

    useEffect(() => {
        fetch('/backend/logs' + logType)
        .then((response) => response.json())
        .then((data) => setLogData(data['log']))
        .catch((error) => console.log(error))
    }, [logType])

    return (
        <div className='logPopup'>
            <div className='logNav'>
                <button onClick={() => setLogType("/activity")}>Activity</button>
                <button onClick={() => setLogType("/event")}>Events</button>
            </div>
            {logData?.map((entry, index) => (
                <li key={index}>
                    <div className='logEntry'>{"> "}{entry}</div>
                </li>
            ))}
        </div>
    )
}

export default LogPopup