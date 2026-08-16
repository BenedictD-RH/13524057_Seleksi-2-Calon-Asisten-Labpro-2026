import './ControlPanel.css'
import { useEffect, useState, useRef } from 'react'
import { stringify } from 'uuid';

function isUpdatableField(fieldKey) {
    if (fieldKey.substring(fieldKey.length - 2, fieldKey.length).toLowerCase() == 'id') {
        return false
    } else if (fieldKey.substring(fieldKey.length - 2, fieldKey.length).toLowerCase() == 'at') {
        return false
    }
    return true
}

function getSubsectionPaths() {
    const path = window.location.pathname
    if (path == "/users") {
        return ["/groups"]
    } else if (path == "/groups") {
        return ["/users"]
    }
    return []
}

function getFieldClass(field) {
    if (field.substring(field.length - 2, field.length).toLowerCase() == 'id') {
        return "idField"
    } else if (field.substring(field.length - 2, field.length).toLowerCase() == 'at') {
        return "dateField"
    }
    return field.toLowerCase() + "Field"
}

function formatEntryData(entry, dataKey) {
    if (dataKey.substring(dataKey.length - 2, dataKey.length).toLowerCase() == 'id') {
        return stringify(entry[dataKey])
    }
    return entry[dataKey]
}

// 0 == Cannot be filled
// 1 == Optional
// 2 == Required
function getRequirementCategory(fieldKey) {
    const required = ['Name', 'Email', 'PasswordHash', 'ClientId', 'ClientSecretHash']
    const optional = ['Description', 'LaunchUrl', 'LogoutNotificationUrl']
    for (const key of required) {
        if (fieldKey == key) return 2
    }
    for (const key of optional) {
        if (fieldKey == key) return 1
    }
    return 0
}

function toPayloadFieldKey(fieldKey) {
    var newStr = fieldKey
    if (isHashedField(newStr)) {
        newStr = newStr.substring(0, newStr.length - 4)
    }
    var newerStr = ""
    for (var i = 0; i < newStr.length; i++) {
        if (i > 0 && isCapital(newStr[i]) && newStr != "ID") {
            newerStr += "_" + newStr[i]
        } else {
            newerStr += newStr[i]
        }
    }
    return newerStr.toLowerCase()
}

function isHashedField(fieldKey) {
    return fieldKey.substring(fieldKey.length - 4, fieldKey.length) == 'Hash'
}

function isCapital(char) {
  return char == char.toUpperCase() && char != char.toLowerCase();
}

function seperateFieldName(fieldName) {
    var newStr = ""
    for (var i = 0; i < fieldName.length; i++) {
        if (i > 0 && isCapital(fieldName[i]) && fieldName != "ID") {
            newStr += " " + fieldName[i]
        } else {
            newStr += fieldName[i]
        }
    }
    return newStr
}

function DataContainer({entry, dataKey, remount}) {
    const data = formatEntryData(entry, dataKey)
    const [inputValue, setInputValue] = useState(data)
    const handleOnChange = (event) => {
        setInputValue(event.target.value)
    }
    const handleOnFocus = () => {
        setInputValue('')
    }
    const handleOnBlur = () => {
        setInputValue(data)
    }
    
    const handleDataUpdate = (event) => {
        if (event.key != 'Enter') return
        const id = entry['id']
        const path = window.location.pathname
        fetch('/server' + path, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({
                "id": id, 
                [dataKey]: inputValue
            }),
        })
        .then((res) => {
            if (res.ok) {
                return res.json();
            } else {
                throw new Error('Network response was not ok');
            }
            return null
        }).then((data) => {
            console.log(data)
            remount()
            setInputValue(entry[dataKey])
        }).catch((err) => {
            console.log(err)
            remount()
            setInputValue(data)
        })
    }

    return (
        <div className={'dataCont ' + getFieldClass(dataKey) + "Cont"}>
            {isUpdatableField(dataKey) ? 
                <input className={'dataInput ' + getFieldClass(dataKey)}
                       value={inputValue} onChange={handleOnChange} onKeyDown={handleDataUpdate} onFocus={handleOnFocus} onBlur={handleOnBlur}></input> : 
                <div className={'dataLabel ' + getFieldClass(dataKey)}>{data}</div>
            }
        </div>
    )
}

function EntryRow({entry, remount, subsection = false}) {
    const [expanded, setExpanded] = useState(false)
    const keys = entry ? Object.keys(entry) : []
    const handleDelete = () => {
        const path = window.location.pathname
        fetch('/server' + path, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({
                "id": entry['id'], 
            }),
        }).then((res) => {
            remount()
        }).catch((err) => {
            remount()
        })
    }

    const handleExpand = () => {
        setExpanded(!expanded)
    }

    return (
        <>
            <div className='dataRowCont'>
                {!subsection ? <button className='rowBtn expandBtn' onClick={handleExpand}>{expanded ? 'V' : '>'}</button> : 
                               <div className='rowBtn'></div>
                }
                <ul className='dataRow'>
                    {keys.map((key, index) => (
                        <li key={index}>
                            <DataContainer entry={entry} dataKey={key} remount={remount}></DataContainer>
                        </li>
                    ))}
                </ul>
                <button className='rowBtn deleteBtn' onClick={handleDelete}>X</button>
            </div>
            {expanded ? <SubsectionGroup uuid={entry['id']} remount={remount}></SubsectionGroup> : <></>}
        </>
    )
}

function CreateInput({fieldKey, changeFieldValue, resetTrigger}) {
    const [value, setValue] = useState('')
    const handleValueChange = (event) => {
        setValue(event.target.value)
        changeFieldValue(fieldKey, event.target.value)
    }
    useEffect(() => {
        setValue('')
    }, [resetTrigger])
    return (
        <div className={'createCont '+ getFieldClass(fieldKey) + "Cont"}>
            {getRequirementCategory(fieldKey) != 0 ? 
                <input className={'createInput ' + getFieldClass(fieldKey)} value={value} onChange={handleValueChange}></input> : 
                <></>
            }
        </div>
    )
}

function TableHeader({remount}) {
    const [fields, setFields] = useState(null)
    const [createForm, setCreateForm] = useState({})
    const [resetFieldTrigger, setResetFieldTrigger] = useState(0)
    const triggerFieldReset = () => {
        setResetFieldTrigger(resetFieldTrigger + 1)
    }
    const handleCreateFormChange = (fieldKey, value) => {
        const copy = createForm
        copy[toPayloadFieldKey(fieldKey)] = value
        setCreateForm(copy)
    }
    const path = window.location.pathname

    const handleCreateData  = () => {
        fetch('/server' + path, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify(createForm),
        }).then((res) => {
            return res.json()
        }).then((data) => {
            if (data['status'] == 'success') {
                setCreateForm({})
                remount()
                triggerFieldReset()
            } else {
                console.log(data)
            }
        }).catch((err) => {
            console.log(err)
            remount()
        })
    }

    
    useEffect(() => {
        fetch('/server' + path + "/fields")
        .then((res) => {
            if (res.ok) {
                return res.json();
            } else {
                throw new Error('Network response was not ok');
            }
            return null
        })
        .then((data) => {
            setFields(data)
        })
        .catch((err) => {
            console.log(err)
        })
    }, [])
    console.log(fields)
    return (
        <div>
            <div className="dataRowCont">
                <div className='rowBtn'></div>
                <ul className='fieldRow'>
                    {fields != null ? fields.fields.map((key, index) => (
                        <li key={index}>
                            <div className={'fieldCont ' + getFieldClass(key) + "Cont"}>{seperateFieldName(key)}</div>
                        </li>
                    )) : <></>}
                </ul>
            </div>
            <div className="dataRowCont">
                <button className='rowBtn createBtn' onClick={handleCreateData}>+</button>
                <ul className='dataRow'>
                    {fields != null ? fields.fields.map((key, index) => (
                        <li key={index}>
                            <CreateInput fieldKey={key} 
                                         changeFieldValue={handleCreateFormChange} 
                                         resetTrigger={resetFieldTrigger}>
                            </CreateInput>
                        </li>
                    )) : <></>}
                </ul>
            </div>
        </div>
    )
}

function SubsectionAddList({subsectionComplement, remount}) {
    return (
        <div>
            <div className="dataRowCont">
                <button className='rowBtn createBtn'>+</button>
                <ul className='fieldRow'>
                    {fields != null ? fields.fields.map((key, index) => (
                        <li key={index}>
                            <div className={'fieldCont ' +getFieldClass(key) + "Cont"}>{seperateFieldName(key)}</div>
                        </li>
                    )) : <></>}
                </ul>
            </div>
            {subsectionComplement.map((entry, index) => (
                <li key={index}><EntryRow entry={entry} remount={remount} subsection={true}></EntryRow></li>
            ))}
        </div>
    )
}

function SubsectionHeader({subsection_path, subsection_uuid, remount}) {
    const path = window.location.pathname
    const [fields, setFields] = useState(null)

    useEffect(() => {
        fetch('/server' + path + subsection_path + "/fields")
        .then((res) => {
            if (res.ok) {
                return res.json();
            } else {
                throw new Error('Network response was not ok');
            }
            return null
        })
        .then((data) => {
            console.log(data)
            setFields(data)
        })
        .catch((err) => {
            console.log(err)
        })
    }, [])

    return (
        <div className="dataRowCont">
            <button className='rowBtn createBtn'>+</button>
            <ul className='fieldRow'>
                {fields != null ? fields.fields.map((key, index) => (
                    <li key={index}>
                        <div className={'fieldCont ' +getFieldClass(key) + "Cont"}>{seperateFieldName(key)}</div>
                    </li>
                )) : <></>}
            </ul>
        </div>
    )
}

function Subsection({subsection_path, subsection_uuid, remount}) {
    const path = window.location.pathname
    const [subsectionData, setSubsectionData] = useState([])

    useEffect(() => {
        fetch('/server' + path + subsection_path + "/query", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({'id' : subsection_uuid}),
        })
        .then((res) => {
            if (res.ok) {
                return res.json();
            } else {
                throw new Error('Network response was not ok');
            }
            return null
        })
        .then((data) => {
            console.log(data)
            setSubsectionData(data)
        })
        .catch((err) => {
            console.log(err)
        })
    }, [])
    
    return (
        <ul className='subsectionRowCont'>
            <SubsectionHeader subsection_path={subsection_path} subsection_uuid={subsection_uuid} remount={remount}></SubsectionHeader>
            {subsectionData.map((entry, index) => (
                <li key={index}><EntryRow entry={entry} remount={remount} subsection={true}></EntryRow></li>
            ))}
        </ul>
    )
}

function SubsectionGroup({uuid, remount}) {
    const path = window.location.pathname
    const subsection_paths = getSubsectionPaths()

    return (
        <ul>
            {subsection_paths.map((path, index) => (
                <li key={index}><Subsection subsection_path={path} subsection_uuid={uuid} remount={remount}></Subsection></li>
            ))}
        </ul>
    )
}

function ControlPanel({pageData, remount}) {
    return (
        <div>
            <TableHeader remount={remount}></TableHeader>
            {pageData.map((entry, index) => (
                <li key={entry.id}>
                    <EntryRow entry={entry} remount={remount}></EntryRow></li>
            ))}
        </div>
    )
}


export default ControlPanel