package main

import (
    "log"
    "encoding/json"
    "code.google.com/p/go.net/websocket"
    "github.com/bitly/go-simplejson"
)

func queryCharInfo(this *player) int {
    rows, err := gamedb.Query(`SELECT cname, clevel, vitality, scene, stage, cards
    FROM charinfo WHERE cid=?`, this.charID)
    if err != nil {
        return 2
    }

    var cards string

    for rows.Next() {
        err = rows.Scan(&this.character.CharName, &this.character.Level,
        &this.character.Vitality,&this.character.Scene,
        &this.character.Stage, &cards)
        if err != nil {
            return 2
        }
        // decode cards
        err = json.Unmarshal([]byte(cards), &this.character.cards)
        if err != nil {
            panic(err)
        }
    }

    return 0
}

func cmCharCreateHander(this *player, command string, param *simplejson.Json) {
    rtnCode := 0

    defer func() {
        var rtnMsg interface{}
        if 0 == rtnCode {
            rtnMsg = this.character
        } else {
            rtnMsg, _ = errCodes[rtnCode]
        }
        rtnJson := responseJson(command, rtnCode, rtnMsg)
        if err := websocket.Message.Send(this.ws, rtnJson); err != nil {
            log.Printf("Send fail for cmCharCreateHander")
        }
    }()

    //check if login
    if this.userID == 0 {
        rtnCode = 6
        return
    }

    // get charname from json
    charname, err := param.Get("CharName").String()
    if err != nil {
        rtnCode = 1
        return
    }

    // insert to table charinfo
    stmt, err := gamedb.Prepare("INSERT charinfo SET uid=?, cname=?, cards=?")
    if err != nil {
        rtnCode = 2
        return
    }

    // TODO: read default vaule from config
    cards := `[{"ID":1, "Level":1, "Status":1},{"ID":2, "Level":1, "Status":1}]`//json array

    res, err := stmt.Exec(this.userID, charname, cards)
    if err != nil {
        log.Print(err)
        rtnCode = 5
        return
    }

    cid, err := res.LastInsertId()
    if err != nil {
        rtnCode = 2
        log.Printf("Error: database not support LastInsertId()")
        return
    }

    // safe cast for mysql INT
    this.charID = uint32(cid)
    // query that we jsut insert
    rtnCode = queryCharInfo(this)
}

func cmCharGetHander(this *player, command string, param *simplejson.Json) {
    rtnCode := 0

    defer func() {
        var rtnMsg interface{}
        if 0 == rtnCode {
            rtnMsg = this.character
        } else {
            rtnMsg, _ = errCodes[rtnCode]
        }
        rtnJson := responseJson(command, rtnCode, rtnMsg)
        if err := websocket.Message.Send(this.ws, rtnJson); err != nil {
            log.Printf("Send fail for cmCharGetHander")
        }
    }()

    //check if login
    if this.userID == 0 || this.charID == 0{
        rtnCode = 6
        return
    }

    rtnCode = queryCharInfo(this)
}
