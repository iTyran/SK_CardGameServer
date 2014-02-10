package main

import (
    "log"
    "code.google.com/p/go.net/websocket"
    "github.com/bitly/go-simplejson"
)

func cmRegisterHander(this *player, command string, param *simplejson.Json) {
    rtnCode := 0

    defer func() {
        var rtnMsg interface{}
        if 0 == rtnCode {
            // client need display create character UI
            rtnMsg = "CreateCharacter"
        } else {
            rtnMsg, _ = errCodes[rtnCode]
        }
        rtnJson := responseJson(command, rtnCode, rtnMsg)
        if err := websocket.Message.Send(this.ws, rtnJson); err != nil {
            log.Printf("Send fail for cmRegisterHander")
        }
    }()

    // parse json
    username, err := param.Get("Username").String()
    if err != nil {
        rtnCode = 1
        return
    }

    passwd, err := param.Get("Password").String()
    if err != nil {
        rtnCode = 1
        return
    }

    email, err := param.Get("Email").String()
    if err != nil {
        rtnCode = 1
        return
    }

    // insert to db
    stmt, err := gamedb.Prepare("INSERT userinfo SET username=?, passwd=?, email=?, registered=?")
    if err != nil {
        rtnCode = 2
        return
    }

    res, err := stmt.Exec(username, sumSha1(passwd), email, nowToDateTime())
    if err != nil {
        rtnCode = 3
        return
    }

    uid, err := res.LastInsertId()
    if err != nil {
        rtnCode = 2
        log.Printf("Error: database not support LastInsertId()")
        return
    }

    // safe cast for mysql INT
    this.userID = uint32(uid)
}

func cmLoginHander(this *player, command string, param *simplejson.Json) {
    rtnCode := 0

    defer func() {
        var rtnMsg interface{}
        if 0 == rtnCode && 0 == this.charID {
            // client need display create character UI
            rtnMsg = "CreateCharacter"
        } else {
            rtnMsg, _ = errCodes[rtnCode]
        }
        rtnJson := responseJson(command, rtnCode, rtnMsg)
        if err := websocket.Message.Send(this.ws, rtnJson); err != nil {
            log.Printf("Send fail for cmLoginHander")
        }
    }()

    // parse json
    username, err := param.Get("Username").String()
    if err != nil {
        rtnCode = 1
        return
    }

    passwd, err := param.Get("Password").String()
    if err != nil {
        rtnCode = 1
        return
    }

    // query password of the user
    rows, err := gamedb.Query("SELECT uid,passwd FROM userinfo WHERE username=?", username)
    if err != nil {
        rtnCode = 2
        return
    }

    var uid uint32
    var dbPasswd string
    for rows.Next() {
        if err := rows.Scan(&uid, &dbPasswd); err != nil {
            rtnCode = 2
            return
        }
    }

    if sumSha1(passwd) != dbPasswd {
        rtnCode = 4
        return
    }

    this.userID = uid

    // check if have game character
    rows, err = gamedb.Query("SELECT cid FROM charinfo WHERE uid=?", uid)
    if err != nil {
        rtnCode = 2
        return
    }

    var cid uint32
    for rows.Next() {
        if err := rows.Scan(&cid); err != nil {
            rtnCode = 2
            return
        }
    }

    this.charID = cid
}
