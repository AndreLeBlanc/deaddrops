package database

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "time"
    "deadrop/api"
)


// Läs en stash från databasen.

func SelectStash(db * sql.DB, token string) stash {
	rows, err := db.Query("SELECT * FROM userinfo")
}

​

// Se en stashs Lifetime.

func SelectLifetime(db * sql.DB, token string) int {
    var tok int
    db.QueryRow("SELECT Lifetime FROM stashes WHERE token=?", token).Scan(&tok)
    return tok
}