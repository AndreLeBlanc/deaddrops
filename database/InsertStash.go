package database

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "time"
    "deadrop/api"
)

func addToHashD(db * sql.DB, tok string, lifeTime int) error {
    _, err := db.Exec("INSERT INTO stashes(token, Lifetime) values(?,?)", tok, lifeTime)
    return err
}

func createNewTable(db * sql.DB, token string, fil * []StashFile) error {
    if fil == nil {
        return DError{time.Now(), "No stashfile!"}
    }   

    _ , err := db.Exec("CREATE TABLE IF NOT EXISTS " + token + " (Fname STRING, size INT, type String, numD Int, sqltime TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL)")

    if err != nil {

        return DError{time.Now(), "Couldn't create table!"}
    }

    for _ ,element := range *fil {
        _, error := db.Exec("INSERT INTO " + token + " (Fname, Type, numD) values(?,?,?)", element.Fname, element.Type, element.Download)      
        if error != nil {
            return DError{time.Now(), "Couldn't insert into table!"}
        }
    }

    return nil
}

func skrivUt(db * sql.DB) {
 	rows, err := db.Query("SELECT * FROM userinfo")
    CheckErr(err)

        for rows.Next() {
	        var hash, size, numD int
	        var format string
            var time time.Time
	        err = rows.Scan(&hash, &size, &format, &numD, &time)
	        CheckErr(err)
	        fmt.Println(hash , " | ", size, " | ", time, " | ", format, " | ", numD)
    }
}


func InsertStash(db * sql.DB, s * Stash) error {
    error := addToHashD(db, s.Token, s.Lifetime)
    if error == nil {
        error = createNewTable(db, s.Token, &s.Files)
    }
    return error
}
