import './ControlPanel.css'
import { stringify } from 'uuid';

function DataContainer({entry, dataKey}) {
    const data = dataKey != 'id' ? entry[dataKey] : stringify(entry[dataKey])
    return (
        <div className='dataCont'>{data}</div>
    )
}

function EntryRow({entry}) {
    const keys = entry ? Object.keys(entry) : []
    return (
        <ul className='dataRow'>
            {keys.map((key, index) => (
                <li key={index}>
                    <DataContainer entry={entry} dataKey={key}></DataContainer>
                </li>
            ))}
        </ul>
    )
}

function ControlPanel({pageData}) {
    return (
        <div>
            {pageData.map((entry, index) => (
                <li key={entry.id}><EntryRow entry={entry}></EntryRow></li>
            ))}
        </div>
    )
}


export default ControlPanel