package main

import (
    "code.google.com/p/go.net/websocket"
)

type cardsOfChar struct {
    Status int
    ID     int
    Level  int
}

// charinfo table data
type charInfo struct {
    CharName string
    Level int
    Vitality int
    Scene int
    Stage int
    cards []cardsOfChar
}

type player struct {
    ws *websocket.Conn // ws of this player
    userID uint32 // uid of table userinfo
    charID uint32 // cid of table character

    character charInfo
}
