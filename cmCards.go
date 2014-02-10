package main

import (
    "code.google.com/p/go.net/websocket"
    "github.com/bitly/go-simplejson"
    "log"
)

// return all cards of this player, in troop and off troop
func cmCardsGetHander(this *player, command string, param *simplejson.Json) {
    rtnCode := 0
    //TODO make this vaule globle for this player, PUSH_CARDS_ORDER will use it
    var cardsArray []lvCardData

    defer func() {
        var rtnMsg interface{}
        if 0 == rtnCode {
            rtnMsg = cardsArray
        } else {
            rtnMsg, _ = errCodes[rtnCode]
        }
        rtnJson := responseJson(command, rtnCode, rtnMsg)
        log.Println(rtnJson)
        if err := websocket.Message.Send(this.ws, rtnJson); err != nil {
            log.Printf("Send fail for cmCardsGetHander")
        }
    }()

    //check for charInfo
    if this.charID == 0 {
        rtnCode = 6
        return
    }

    for _, value := range this.character.cards {
        card := gGameScript.getLvCard(value.ID, value.Level, value.Status)
        cardsArray = append(cardsArray, *card)
    }
}
