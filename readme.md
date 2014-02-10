# game server use golang and websocket

## Login module
Register
```
{
    "Command" : "CM_REGISTER",
        "Param":
        {
            "Username" : "xxx",
            "Password" : "aa",
            "Email" : "a@x.com"
        }
}

{
    "Command" : "CM_REGISTER",
        "Return" :
        {
            "Code" : 0,
            "Message" : "CreateCharacter"
        }
}
```

Login
```
{
    "Command" : "CM_LOGIN",
        "Param" :
        {
            "Username" : "xxx",
            "Password" : "aa"
        }
}

{
    "Command" : "CM_LOGIN",
        "Return" :
        {
            "Code" : 0,
            "Message" : "Success" // or "CreateCharacter"
        }
}
```

Create character (LOGIN FIRST)
```
{
    "Command" : "CM_CHAR_CREATE",
        "Param" :
        {
            "CharName" : "xxx",
        }
}

{
    "Command" : "CM_CHAR_CREATE",
        "Return" :
        {
            "Code" : 0,
            "Message" : {json of char info}
        }
}
```

Get character (LOGIN FIRST)
```
{
    "Command" : "CM_CHAR_GET",
        "Param" : ""
}

{
    "Command" : "CM_CHAR_GET",
        "Return" :
        {
            "Code" : 0,
            "Message" : {json of char info}
        }
}
```

json of char info
```
{
    "CharName": "a New Name ",
    "Level" : 1,
    "Vitality": 60,
    "Scene": 1,
    "Stage" : 1
}
```

Get card array
```
{
    "Command" : "CM_CARDS_GET",
        "Param" : ""
}

{
    "Command" : "CM_CARDS_GET",
        "Return" :
        {
            "Code" : 0,
            "Message" : {json of card array}
        }
}
```

json of card array
```
[
{
    "Name": "小花猫",
    "HP": 21,
    "Attack": 27,
    "Defence": 13,
    "Speed": 120,
    "Talent": 720,
    "Drop": 50,
    "Experience": 324,
    "Skill": [1,2],
    "ID": "002",
    "Level": 1,
    "Status": 1,
    "Hash": "0xc2000e9630",
}
]
```

Raid
```
{
    "Command" : "CM_RAID",
        "Param" : {
            "Scene" : 1,
            "Stage" : 1
        }
}

{
    "Command" : "CM_RAID",
        "Return" : 
        {
            "Code" : 0,
            "Message" : "{battle command order}"
        }
}
```
