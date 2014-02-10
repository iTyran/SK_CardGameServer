package main

import (
    "io/ioutil"
    "encoding/json"
    "fmt"
    "log"
)

//FIXME: change to your own data path.
//make sure end with '/'
const gDataPath = "/Users/u0u0/mygo/src/gameServer/"

// scene.json
type monsterData struct {
    ID int
    Level int
    Position int
}

type stageDate struct {
    Monster []monsterData
    Award []int
    Vitality int
}
// end of scene.json

// stages of a scene
type sceneData struct {
    stages []stageDate
}

// card.json
// base info
type cardData struct {
    Name string
    HP int
    Attack int
    Defence int
    Speed int
    Talent int
    Drop int
    Experience int
    Skill []int
}

// card info of level addition
type lvCardData struct {
    // anonymous filed
    cardData
    // additional info
    ID string
    Level int
    Status int // 1: in troop, 0: off troop
    Hash string
}

// combat card info
type combatCardData struct {
    lvCardData
    Pos int
}

// skill.json
type skillData struct {
    Name string
    Power int
    Hit int
    Range int
}

// all scene, card, skill data info
type gameScript struct {
    scenes []sceneData
    cards []cardData
    skills []skillData
}

// ======global var define
var gGameScript gameScript

// =========== method of gameScript =========
func (script *gameScript) getLvCard(id int, level int, status int) *lvCardData {
    idString := fmt.Sprintf("%03d", id)
    lvData := lvCardData{cardData:*script.getCard(id), ID:idString, Level:level, Status:status}
    lvData.Hash = fmt.Sprintf("%p", &lvData)// %p to print point address
    //corretive HP, Attack, Defence
    lvData.HP = (lvData.Talent * 2 + lvData.HP) * level / 100 + 5 + level
    lvData.Attack = (lvData.Talent * 2 + lvData.Attack) * level / 100 + 10
    lvData.Defence = (lvData.Defence * 2 + lvData.Defence) * level / 100 + 10
    return &lvData
}

func (script *gameScript) getCombatCard(id int, level int, status int, pos int) *combatCardData {
    return &combatCardData{*script.getLvCard(id, level, status), pos}
}

// id [1,n]
func (script *gameScript) getCard(id int) *cardData {
    if id < 1 || id > len(script.cards) {
        log.Println("Error: getCard out of range", id)
        return nil
    }
    return &script.cards[id - 1]
}

// scene [1,n]
// stage [1,n]
func (script *gameScript) getStage(scene int, stage int) *stageDate {
    if scene < 1 || scene > len(script.scenes) {
        log.Println("Error: scene out of range", scene)
        return nil
    }
    if stage < 1 || stage > len(script.scenes[scene - 1].stages) {
        log.Println("Error: stage out of range", stage)
        return nil
    }
    return &script.scenes[scene - 1].stages[stage - 1]
}

// id [1,n]
func (script *gameScript) getSkill(id int) *skillData {
    if id < 1 || id > len(script.skills) {
        log.Println("Error: getSkill out of range", id)
        return nil
    }
    return &script.skills[id - 1]
}

func (script *gameScript) initSceneData() {
    totalScene := 0
    for {
        filepath := fmt.Sprintf("%sdata/scene/scene%02d-01.json", gDataPath, totalScene + 1)
        if false == isExist(filepath) {
            break
        }
        totalScene++

        stageOfScene := 0
        var scene sceneData
        for {
            filepath = fmt.Sprintf("%sdata/scene/scene%02d-%02d.json",
            gDataPath, totalScene, stageOfScene + 1)
            if false == isExist(filepath) {
                break
            }
            // read whole file to buffer
            buf, err := ioutil.ReadFile(filepath)
            if err != nil {
                panic(err)
            }

            stageOfScene++

            var stage stageDate
            err = json.Unmarshal(buf, &stage)
            if err != nil {
                panic(err)
            }
            scene.stages = append(scene.stages, stage)
        }
        // important: add scene after the scene changed
        script.scenes = append(script.scenes, scene)
    }
}

func (script *gameScript) initCardData() {
    totalCard := 0
    for {
        filepath := fmt.Sprintf("%sdata/card/card%03d.json", gDataPath, totalCard + 1)
        if false == isExist(filepath) {
            break
        }
        // read whole file to buffer
        buf, err := ioutil.ReadFile(filepath)
        if err != nil {
            panic(err)
        }

        totalCard++

        var card cardData
        err = json.Unmarshal(buf, &card)
        if err != nil {
            panic(err)
        }
        script.cards = append(script.cards, card)
    }
}

func (script *gameScript) initSkillData() {
    totalSkill := 0
    for {
        filepath := fmt.Sprintf("%sdata/skill/skill%03d.json", gDataPath, totalSkill + 1)
        if false == isExist(filepath) {
            break
        }
        // read whole file to buffer
        buf, err := ioutil.ReadFile(filepath)
        if err != nil {
            panic(err)
        }

        totalSkill++

        var skill skillData
        err = json.Unmarshal(buf, &skill)
        if err != nil {
            panic(err)
        }
        script.skills = append(script.skills, skill)
    }
}

// script module init
func initGameData() {
    gGameScript.initSceneData()
    gGameScript.initCardData()
    gGameScript.initSkillData()
}
