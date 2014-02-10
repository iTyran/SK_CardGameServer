package main

import (
    "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// TODO: load db info from config file
var gamedb *sql.DB
func initDB() {
    var err error
    gamedb, err = sql.Open("mysql", "kapai:kapai123@/kapai?charset=utf8")
    if err != nil {
        panic(err)
    }
}
