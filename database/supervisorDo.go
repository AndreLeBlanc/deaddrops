package database

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "time"
)

func delete(db * sql.DB, hash int) {
    var numD int
    error := db.QueryRow("SELECT numD FROM userinfo WHERE hash=?", hash).Scan(&numD)
    CheckErr(error)
    
    if numD < 1 {
        ut, error := db.Prepare("delete from userinfo where hash=?")
        CheckErr(error)
        ut.Exec(hash)
    } else {
        update, error := db.Prepare("UPDATE userinfo SET numD=? WHERE hash=?")
        CheckErr(error)
        update.Exec(numD-1, hash)
    }
}

func SupervisorDo(db * sql.DB, token string, c chan string) {

    select {
    case fname := <-c:
        fmt.Printf("received filename: %s\n", fname)
        delete(db, 2131231)
    case <-time.After(time.Second * 10):
        fmt.Println("timeout")
    }
}