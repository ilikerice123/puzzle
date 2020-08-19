import React from 'react';
import PuzzleClient from './client'
import { AxiosResponse } from 'axios'
import { UserObject } from './game'
import { Tooltip } from '@material-ui/core'

type IdentityProps = {
    changeUser: (user: UserObject) => void
    user: UserObject | null
}

export default function Identity(props: IdentityProps) {
    const userID = props.user ? props.user.id : null
    const name = props.user ? props.user.name : null

    const [newName, setNewName] = React.useState("")
    const [newUserID, setNewUserID] = React.useState("")
    const [show, setShow] = React.useState(false)
    const client = new PuzzleClient()
    const modalRef = React.useRef<HTMLDivElement>(null)

    const updateUserID = async () => {
        let res: AxiosResponse<UserObject> | null
        if (newName != "") {
            // create new user
            res = await client.postJson<UserObject>('/users', {name: newName})
        } else {
            // check if id valid
            res = await client.get<UserObject>(`/users/${newUserID}`)
        }
        if (res == null || res.data == null) {
            return
        }
        props.changeUser(res.data)
        setShow(false)
    }

    const checkClick = (event: any) => {
        if (modalRef == null || modalRef.current == null) {
            return
        }
        if (!modalRef.current.contains(event.target)) {
            setShow(false)
        }
    }

    React.useEffect(() => {
        document.addEventListener('mousedown', checkClick);
        return () => {
            document.removeEventListener('mousedown', checkClick);
        }
    })

    return (
        <div>
            <div onClick={() => setShow(true)}>
                {userID == null ? (
                    "Sign in"
                ) : (
                    <Tooltip title={userID}>
                        <div>{name}</div>
                    </Tooltip>
                )}
            </div>
            {show && (
                <div className="modal" ref={modalRef}>
                    <p className="modal-message">Enter a new name or an existing user ID</p>
                    New Name: 
                    <input
                        className="ipt"
                        type="text"
                        value={newName}
                        onChange={(e) => {
                            setNewName(e.target.value)
                            setNewUserID("")
                        }} 
                    />
                    ID: 
                    <input
                        className="ipt"
                        type="text"
                        value={newUserID}
                        onChange={(e) => {
                            setNewUserID(e.target.value)
                            setNewName("")
                        }} 
                    />
                    <br />
                    <div className="modal-options">
                        <button className="btn" onClick={updateUserID}>
                            GO
                        </button>
                    </div>
                </div>
            )}
        </div>
    )
}