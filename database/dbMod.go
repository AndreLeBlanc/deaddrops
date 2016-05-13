package database

import (
	"database/sql"
)

func NewConnect() * sql.DB {
	Db, err := sql.Open("sqlite3", "database/deaddrops.db")
    CheckErr(err)
    return Db
}

func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}