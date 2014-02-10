package main

var errCodes = map[int]string{
    0 : "Success",
    1 : "Json parse fail", // fail for parse json of client
    2 : "general database fail", // server fail for sql exec
    3 : "register fail:duplicate name",
    4 : "login fail: wrong password",
    5 : "create char fail: duplicate name",
    6 : "require login",
    7 : "stage out of range",
}
