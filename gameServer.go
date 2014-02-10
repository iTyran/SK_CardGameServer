package main

import (
    "log"
    "net/http"
    "code.google.com/p/go.net/websocket"
    "github.com/bitly/go-simplejson"
)

func wsHandler(ws *websocket.Conn) {
    var err error
    var this player // link ws with player

    this.ws = ws
    // need loop to keep socket connect
    for {
        var reply string

        if err = websocket.Message.Receive(ws, &reply); err != nil {
            log.Printf("connect closed!")
            break
        }

        js, err := simplejson.NewJson([]byte(reply))
        if err != nil {
            // TODO: Send error json back to client
            log.Printf("parse json error:", err);
            continue
        }

        commandDispatcher(&this, js)
    }
}

func main() {
    log.Print("initing database ...");
    initDB()
    log.Print("initing game script data ...");
    initGameData()

    log.Print("starting socket server ...");
    //TODO read port from config
    http.Handle("/", websocket.Handler(wsHandler))
    if err := http.ListenAndServe(":1234", nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
