import './ControlPanel.css'
import { useEffect, useState, useRef } from 'react'
import { stringify } from 'uuid';

function isUpdatableField(fieldKey) {
    const unupdatableFields = ['id', 'created_at', 'updated_at']
    for (const key of unupdatableFields) {
        if (fieldKey == key) return false
    }
    return true
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
    const data = dataKey != 'id' ? entry[dataKey] : stringify(entry[dataKey])
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
        <div className='dataCont'>
            {isUpdatableField(dataKey) ? 
                <input className={'dataInput ' + dataKey + 'Field'}
                       value={inputValue} onChange={handleOnChange} onKeyDown={handleDataUpdate} onFocus={handleOnFocus} onBlur={handleOnBlur}></input> : 
                <div className={'dataLabel ' + dataKey + 'Field'}>{data}</div>
            }
        </div>
    )
}

function EntryRow({entry, remount}) {
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

    return (
        <div className='dataRowCont'>
            <div className='rowBtn'></div>
            <ul className='dataRow'>
                {keys.map((key, index) => (
                    <li key={index}>
                        <DataContainer entry={entry} dataKey={key} remount={remount}></DataContainer>
                    </li>
                ))}
            </ul>
            <button className='rowBtn deleteBtn' onClick={handleDelete}>X</button>
        </div> 
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
        <div className='createCont'>
            {getRequirementCategory(fieldKey) != 0 ? 
                <input className='createInput' value={value} onChange={handleValueChange}></input> : 
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
        console.log(copy)
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
                setCreateForm([])
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
                            <div className='fieldCont'>{seperateFieldName(key)}</div>
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

function ControlPanel({pageData, remount}) {
    return (
        <div>
            <TableHeader remount={remount}></TableHeader>
            {pageData.map((entry, index) => (
                <li key={entry.id}><EntryRow entry={entry} remount={remount}></EntryRow></li>
            ))}
        </div>
    )
}


export default ControlPanel