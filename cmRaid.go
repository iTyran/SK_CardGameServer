package main

import (
    "log"
    "code.google.com/p/go.net/websocket"
    "github.com/bitly/go-simplejson"
)

type raidCombatData struct {
    Attacker string
    Beattacked string
    Skill int
    Damage int
}

type raidEndData struct {
    Result int
    Drop []int
    Experience int
}

type raidResponseData struct {
    My []combatCardData
    Monster []combatCardData
    Combat []raidCombatData
    End raidEndData
}

func cmRaidHander(this *player, command string, param *simplejson.Json) {
    rtnCode := 0
    var raidResponse raidResponseData

    defer func() {
        var rtnMsg interface{}
        if 0 == rtnCode {
            rtnMsg = raidResponse
        } else {
            rtnMsg, _ = errCodes[rtnCode]
        }
        rtnJson := responseJson(command, rtnCode, rtnMsg)
        if err := websocket.Message.Send(this.ws, rtnJson); err != nil {
            log.Printf("Send fail for cmRaidHander")
        }
    }()

    //check for charInfo
    if this.charID == 0 {
        rtnCode = 6
        return
    }

    scene, err := param.Get("Scene").Int()
    if err != nil {
        rtnCode = 1
        return
    }

    stage, err := param.Get("Stage").Int()
    if err != nil {
        rtnCode = 1
        return
    }

    // prepare data
    var myCards []combatCardData
    for index, value := range this.character.cards {
        card := gGameScript.getCombatCard(value.ID, value.Level, value.Status, index + 1)
        myCards = append(myCards, *card)
    }

    stageInfo := gGameScript.getStage(scene, stage)
    if stageInfo == nil {
        rtnCode = 7
        return
    }

    var monsterCards []combatCardData
    for _, value := range stageInfo.Monster {
        card := gGameScript.getCombatCard(value.ID, value.Level, 0, value.Position)
        monsterCards = append(monsterCards, *card)
    }

    // copy to keep HP full in json
    raidResponse.My = make([]combatCardData, len(myCards))
    copy(raidResponse.My, myCards)
    raidResponse.Monster = make([]combatCardData, len(monsterCards))
    copy(raidResponse.Monster, monsterCards)

    curMy := 1
    aliveMy := len(myCards)
    curMonster := 1
    aliveMonster := len(monsterCards)
    turn := true // true for my, false for monster
    for {
        //check end
        if aliveMy == 0 || aliveMonster == 0 {
            break
        }

        var attArr []combatCardData
        var beAttArr []combatCardData
        var curAttacker *int
        var alive *int

        if turn {
            attArr = myCards
            beAttArr = monsterCards
            curAttacker = &curMy
            alive = &aliveMonster
        } else {
            attArr = monsterCards
            beAttArr = myCards
            curAttacker = &curMonster
            alive = &aliveMy
        }
        // change turn
        turn = !turn

        // get attacker
        attacker := getLiveCard(attArr, curAttacker, len(attArr))
        // point to next attacker
        *curAttacker++
        if *curAttacker > len(attArr) {
            *curAttacker = 1
        }

        // choose skill
        // normal attack + special skill
        totolSkill := 1 + len(attacker.Skill)
        skillIndex := randInRange(0, totolSkill)

        attackRange := 1 //normal attack range
        if skillIndex > 0 {
            // special skill atack range
            attackRange = gGameScript.getSkill(skillIndex).Range
        }

        var attackPos int
        if attacker.Pos < 3 { // attacker in first line
            if attackRange < 2 {
                attackPos = attacker.Pos
            } else {
                attackPos = attacker.Pos + 2
            }
        } else { // attack in second line
            if attackRange < 2 {
                attackPos = 1
            } else {
                attackPos = attacker.Pos
            }
        }
        // XXX:Fix attackPos by max card counts
        if attackPos > len(beAttArr) {
            attackPos = len(beAttArr)
        }

        // get BeAttacked
        beAttacked := getLiveCard(beAttArr, &attackPos, len(beAttArr))

        // count hurtPoint
        attackPower := 1.0
        if skillIndex > 1 {
            attackPower = float64(gGameScript.getSkill(skillIndex).Power) / 100.0
        }
        hurtPoint := int((float64(attacker.Level) * 2 / 5) * float64(attacker.Attack) * attackPower / (float64(beAttacked.Defence) / 10))

        beAttacked.HP -= hurtPoint
        if beAttacked.HP <= 0 {
            *alive--
        }

        cbInfo := raidCombatData{attacker.Hash, beAttacked.Hash, skillIndex, hurtPoint}
        raidResponse.Combat = append(raidResponse.Combat, cbInfo)
    }

    // the end info
    raidResponse.End.Experience = 0
    if aliveMy > 0 {
        raidResponse.End.Result = 1
        // count exp
        for _, card := range monsterCards {
            raidResponse.End.Experience += (card.Experience * card.Level) / (7 * len(myCards))
        }
        // drop
        for _, award := range stageInfo.Award {
            card := gGameScript.getCard(award)
            rand := randInRange(0, 100)
            if rand < card.Drop {
                raidResponse.End.Drop = append(raidResponse.End.Drop, award)
            }
        }
        // TODO save progress
        // TODO add drop to database
    } else {
        raidResponse.End.Result = 0
    }
}

// pos [1, 5]
func getCardOfPos(cards []combatCardData, pos int) *combatCardData {
    for index := 0; index < len(cards); index++ {
        if cards[index].Pos == pos {
            return &cards[index]
        }
    }
    log.Println("Error:getCardOfPos fail")
    return nil
}

// maxPos [1, 5]
func getLiveCard(cards []combatCardData, startPos *int, maxPos int) *combatCardData {
    pos := *startPos
    for {
        card := getCardOfPos(cards, pos)
        if card.HP > 0 {
            *startPos = pos
            return card
        }
        pos++
        if pos > maxPos {
            pos = 1
        }

        if pos == *startPos {
            log.Println("Error:getLiveCard fail")
            break
        }
    }
    return nil
}
