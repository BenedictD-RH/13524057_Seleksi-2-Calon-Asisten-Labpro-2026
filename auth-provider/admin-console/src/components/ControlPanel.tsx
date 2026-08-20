import './ControlPanel.css'
import { useEffect, useState } from 'react'
import { stringify } from 'uuid';

type Entry = Record<string, any>

interface FieldsResponse {
    fields: string[]
}

function isUpdatableInSubsection(fieldKey: string, subsection_path: string): boolean {
    var updatableFields: string[] = []
    const path = window.location.pathname + subsection_path
    if ((path == "/apps/groups") || (path == "/groups/apps")) {
        updatableFields = ['effect']
    }
    for (const key of updatableFields) {
        if (fieldKey == key) return true
    }
    return false
}

function updateFieldKey(dataKey: string): string {
    if (dataKey.substring(dataKey.length - 5, dataKey.length).toLowerCase() == '_hash') {
        return dataKey.substring(0, dataKey.length - 5)
    }
    return dataKey
}

function isUpdatableField(fieldKey: string, subsection_path: string): boolean {
    if (fieldKey.substring(fieldKey.length - 2, fieldKey.length).toLowerCase() == 'id' && fieldKey != "client_id") {
        return false
    } else if (fieldKey.substring(fieldKey.length - 2, fieldKey.length).toLowerCase() == 'at') {
        return false
    } else if (subsection_path != "" && !isUpdatableInSubsection(fieldKey, subsection_path)) {
        return false
    }
    return true
}

function getIDField(path: string): string {
    return path.substring(1, path.length - 1) + "_id"
}

function getPayloadIDField(path: string): string {
    if (path = "/apps") {
        return "application_id"
    }
    return path.substring(1, path.length - 1) + "_id"
}

function getSubsectionTitle(subsection_path: string): string {
    const path = window.location.pathname
    if (path == "/users") {
        if (subsection_path == "/groups") return "Groups"
    }
    else if (path == "/groups") {
        if (subsection_path == "/users") return "Users"
        if (subsection_path == "/apps") return "Application Policies"
    }
    else if (path == "/apps") {
        if (subsection_path == "/groups") return "Group Policies"
        if (subsection_path == "/uri") return "Redirect URIs"
    }

    return ""
}

function getSubsectionPaths(): string[] {
    const path = window.location.pathname
    if (path == "/users") {
        return ["/groups"]
    } else if (path == "/groups") {
        return ["/users", "/apps"]
    } else if (path == "/apps") {
        return ["/groups", "/uri"]
    }
    return []
}

function getFieldClass(field: string): string {
    if (field.substring(field.length - 2, field.length).toLowerCase() == 'id' && field != "client_id") {
        return "idField"
    } else if (field.substring(field.length - 2, field.length).toLowerCase() == 'at') {
        return "dateField"
    }
    return field.toLowerCase() + "Field"
}

function formatEntryData(entry: Entry, dataKey: string): any {
    if (dataKey.substring(dataKey.length - 2, dataKey.length).toLowerCase() == 'id' && dataKey != "client_id") {
        return stringify(entry[dataKey])
    }
    return entry[dataKey]
}

// 0 == Cannot be filled
// 1 == Optional
// 2 == Required
function getRequirementCategory(fieldKey: string): number {
    const required = ['Name', 'Email', 'PasswordHash', 'ClientId', 'ClientSecretHash', 'RedirectUri', 'LogoutNotificationUrl']
    const optional = ['Description', 'LaunchUrl']
    for (const key of required) {
        if (fieldKey == key) return 2
    }
    for (const key of optional) {
        if (fieldKey == key) return 1
    }
    return 0
}

function getSubsectionType(subsection_path: string): number {
    if (subsection_path == "/uri") {
        return 2;
    } else {
        return 1;
    }
}

function toPayloadFieldKey(fieldKey: string): string {
    var newStr = fieldKey
    if (isHashedField(newStr)) {
        newStr = newStr.substring(0, newStr.length - 4)
    }
    var newerStr = ""
    for (var i = 0; i < newStr.length; i++) {
        if (i > 0 && isCapital(newStr.charAt(i)) && newStr != "ID") {
            newerStr += "_" + newStr[i]
        } else {
            newerStr += newStr[i]
        }
    }
    return newerStr.toLowerCase()
}

function isHashedField(fieldKey: string): boolean {
    return fieldKey.substring(fieldKey.length - 4, fieldKey.length) == 'Hash'
}

function isCapital(char: string): boolean {
  return char == char.toUpperCase() && char != char.toLowerCase();
}

function seperateFieldName(fieldName: string): string {
    var newStr = ""
    for (var i = 0; i < fieldName.length; i++) {
        if (i > 0 && isCapital(fieldName.charAt(i)) && fieldName != "ID") {
            newStr += " " + fieldName[i]
        } else {
            newStr += fieldName[i]
        }
    }
    return newStr
}

interface DataContainerProps {
    entry: Entry
    dataKey: string
    remount: () => void
    subsection_path?: string
    subsection_uuid?: any
}

function DataContainer({entry, dataKey, remount, subsection_path = "", subsection_uuid = []}: DataContainerProps) {
    const data = formatEntryData(entry, dataKey)
    const [inputValue, setInputValue] = useState(data)
    const handleOnChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setInputValue(event.target.value)
    }
    const handleOnFocus = () => {
        setInputValue('')
    }
    const handleOnBlur = () => {
        setInputValue(data)
    }
    
    const handleDataUpdate = (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (event.key != 'Enter') return
        const id = entry['id']
        const path = window.location.pathname
        console.log((subsection_path == "" ? JSON.stringify({"id": entry['id'], [updateFieldKey(dataKey)]: inputValue}) : 
                                           JSON.stringify({
                                                [getIDField(path)] : subsection_uuid,
                                                [getIDField(subsection_path)] : entry[getPayloadIDField(subsection_path)],
                                                [updateFieldKey(dataKey)]: inputValue
                                            })))
        fetch('/server' + path + subsection_path, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: (subsection_path == "" ? JSON.stringify({"id": entry['id'], [updateFieldKey(dataKey)]: inputValue}) : 
                                           JSON.stringify({
                                                [getIDField(path)] : subsection_uuid,
                                                [getIDField(subsection_path)] : entry[getPayloadIDField(subsection_path)],
                                                [updateFieldKey(dataKey)]: inputValue
                                            })),
                })
        .then((res) => {
            if (res.ok) {
                return res.json();
            } else {
                throw new Error('Network response was not ok');
            }
            return null
        }).then((data) => {
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
            {isUpdatableField(dataKey, subsection_path) ? 
                <input className={'dataInput ' + getFieldClass(dataKey)}
                       value={inputValue} onChange={handleOnChange} onKeyDown={handleDataUpdate} onFocus={handleOnFocus} onBlur={handleOnBlur}></input> : 
                <div className={'dataLabel ' + getFieldClass(dataKey)}>{data}</div>
            }
        </div>
    )
}

interface EntryRowProps {
    entry: Entry
    remount: () => void
    subsection?: string
    subsection_uuid?: any
}

function EntryRow({entry, remount, subsection = "", subsection_uuid = []}: EntryRowProps) {
    const [expanded, setExpanded] = useState(false)
    const keys = entry ? Object.keys(entry) : []
    const handleDelete = () => {
        const path = window.location.pathname
        fetch('/server' + path + subsection, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: (subsection == "" || getSubsectionType(subsection) == 2  ? JSON.stringify({"id": entry['id']}) : 
                                      JSON.stringify({
                                        [getIDField(path)] : subsection_uuid,
                                        [getIDField(subsection)] : entry[getIDField(subsection)]
                                      }))
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
                {subsection == "" ? <button className='rowBtn expandBtn' onClick={handleExpand}>{expanded ? 'V' : '>'}</button> : 
                                    <div className='rowBtn'></div>
                }
                <ul className='dataRow'>
                    {keys.map((key, index) => (
                        <li key={index}>
                            <DataContainer entry={entry} dataKey={key} remount={remount} subsection_path={subsection} subsection_uuid={subsection_uuid}></DataContainer>
                        </li>
                    ))}
                </ul>
                <button className='rowBtn deleteBtn' onClick={handleDelete}>X</button>
            </div>
            {expanded ? <SubsectionGroup uuid={entry['id']}></SubsectionGroup> : <></>}
        </>
    )
}

interface CreateInputProps {
    fieldKey: string
    changeFieldValue: (fieldKey: string, value: string) => void
    resetTrigger: number
}

function CreateInput({fieldKey, changeFieldValue, resetTrigger}: CreateInputProps) {
    const [value, setValue] = useState('')
    const handleValueChange = (event: React.ChangeEvent<HTMLInputElement>) => {
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

interface TableHeaderProps {
    remount: () => void
}

function TableHeader({remount}: TableHeaderProps) {
    const [fields, setFields] = useState<FieldsResponse | null>(null)
    const [createForm, setCreateForm] = useState<Record<string, any>>({})
    const [resetFieldTrigger, setResetFieldTrigger] = useState(0)
    const triggerFieldReset = () => {
        setResetFieldTrigger(resetFieldTrigger + 1)
    }
    const handleCreateFormChange = (fieldKey: string, value: string) => {
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

interface SubsectionAddListProps {
    subsection_path: string
    subsection_uuid: any
    remount: () => void
    remountKey: number
}

function SubsectionAddList({subsection_path, subsection_uuid, remount, remountKey}: SubsectionAddListProps) {
    const [addList, setAddList] = useState<Entry[]>([])
    const path = window.location.pathname

    const handleAddToSubsection = (pointed_id: any) => {
        fetch('/server' + path + subsection_path, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({
                [getIDField(path)] : subsection_uuid,
                [getIDField(subsection_path)] : pointed_id
            }),
        }).then((res) => {
            return res.json()
        }).then((data) => {
            if (data['status'] == 'success') {
                remount()
            }
        }).catch((err) => {
            console.log(err)
        })
    }

    useEffect(() => {
        fetch('/server' + path + subsection_path + "/query/complement", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({'id' : subsection_uuid}),
        }).then((res) => {
            return res.json()
        }).then((data) => {
            setAddList(data)
        }).catch((err) => {
            console.log(err)
        })
    }, [remountKey])

    if (addList.length == 0) {
        return <></>
    }

    return (
        <div className='popup'>
            {addList.map((entry) => (
                <li key={entry['id']}>
                    <div className='dataRowCont'>
                        <div className='rowBtn'></div>
                        <ul className='dataRow addListRow'>
                            {Object.keys(entry).map((key, index) => (
                                <li key={index}>
                                    <DataContainer entry={entry} dataKey={key} remount={remount}></DataContainer>
                                </li>
                            ))}
                        </ul>
                        <button className='rowBtn createBtn' onClick={() => handleAddToSubsection(entry['id'])}>+</button>
                    </div>
                </li>
            ))}
        </div>
    )
}

interface SubsectionHeaderProps {
    subsection_path: string
    subsection_uuid: any
    remount: () => void
    remountKey: number
}

function SubsectionHeader({subsection_path, subsection_uuid, remount, remountKey}: SubsectionHeaderProps) {
    const path = window.location.pathname
    const [fields, setFields] = useState<FieldsResponse | null>(null)
    const [addListActive, setAddListActive] = useState(false)
    const [createForm, setCreateForm] = useState<Record<string, any>>({[getIDField(path)] : subsection_uuid})
    const [resetFieldTrigger, setResetFieldTrigger] = useState(0)
    const triggerFieldReset = () => {
        setResetFieldTrigger(resetFieldTrigger + 1)
    }
    const handleCreateFormChange = (fieldKey: string, value: string) => {
        const copy = createForm
        copy[toPayloadFieldKey(fieldKey)] = value
        console.log(copy)
        setCreateForm(copy)
    }
    const handleCreateData  = () => {
        fetch('/server' + path + subsection_path, {
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
                setCreateForm({[getIDField(path)] : subsection_uuid})
                remount()
                triggerFieldReset()
            } else {
                console.log(data)
            }
        }).catch((err) => {
            console.log(err)
            setCreateForm({[getIDField(path)] : subsection_uuid})
            remount()
            triggerFieldReset()
        })
    }

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
            setFields(data)
        })
        .catch((err) => {
            console.log(err)
        })
    }, [])

    return (
        <>
            <div className="dataRowCont">
                {getSubsectionType(subsection_path) == 1 ?
                    <button className={'rowBtn ' + (!addListActive ? 'createBtn' : 'deleteBtn')} 
                        onClick={() => setAddListActive(!addListActive)}>{!addListActive ? "+" : "X"}
                    </button> :
                    <div className='rowBtn'></div>
                }
                <ul className='fieldRow'>
                    {fields != null ? fields.fields.map((key, index) => (
                        <li key={index}>
                            <div className={'fieldCont ' +getFieldClass(key) + "Cont"}>{seperateFieldName(key)}</div>
                        </li>
                    )) : <></>}
                </ul>
            </div>
            {getSubsectionType(subsection_path) == 1 ?
                (addListActive ? <SubsectionAddList subsection_path={subsection_path} subsection_uuid={subsection_uuid} remount={remount} remountKey={remountKey}></SubsectionAddList> :
                                <></>) :
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
            }
        </>
    )
}

interface SubsectionProps {
    subsection_path: string
    subsection_uuid: any
}

function Subsection({subsection_path, subsection_uuid}: SubsectionProps) {
    const path = window.location.pathname
    const [subsectionData, setSubsectionData] = useState<Entry[]>([])
    const [remountKey, setRemountKey] = useState(0)
    const [expanded, setExpanded] = useState(false)
    const handleExpand = () => {
        setExpanded(!expanded)
    }
    const triggerRemount = () => {
        setRemountKey(prev => prev + 1)
    }

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
    }, [remountKey])
    
    return (
        <ul className='subsectionRowCont'>
            <div className='subsectionPreview'>
                <button className='rowBtn expandBtn' onClick={handleExpand}>{expanded ? 'V' : '>'}</button>
                <div className={'subsectionTitle' + (expanded ? ' subsectionTitleExpanded': '')}>{getSubsectionTitle(subsection_path)}</div>
                <div className='rowBtn'></div>
            </div>
            {expanded ? 
            <div>
                <SubsectionHeader subsection_path={subsection_path} subsection_uuid={subsection_uuid} remount={triggerRemount} remountKey={remountKey}></SubsectionHeader>
                {(subsectionData.map((entry) => (
                    <li key={entry[getIDField(subsection_path)]}><EntryRow entry={entry} remount={triggerRemount} subsection={subsection_path} subsection_uuid={subsection_uuid}></EntryRow></li>
                )))}
            </div> : <></>}
        </ul>
    )
}

interface SubsectionGroupProps {
    uuid: any
}

function SubsectionGroup({uuid}: SubsectionGroupProps) {
    const path = window.location.pathname
    const subsection_paths = getSubsectionPaths()

    return (
        <ul>
            {subsection_paths.map((path, index) => (
                <li key={index}><Subsection subsection_path={path} subsection_uuid={uuid}></Subsection></li>
            ))}
        </ul>
    )
}

interface ControlPanelProps {
    pageData: Entry[]
    remount: () => void
}

function ControlPanel({pageData, remount}: ControlPanelProps) {
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
