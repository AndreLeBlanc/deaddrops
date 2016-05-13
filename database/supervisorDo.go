package database

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "time"
)

func delete(db * sql.DB, hash int) {
    out, error := db.Prepare("delete from userinfo where hash=?")
    CheckErr(error)

    out.Exec(hash)
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