package database

import (
	"database/sql"
	"deadrop/api"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"fmt"
)

func addToHashD(db *sql.DB, tok string, StashNom string, lifeTime int) error {
	_, err := db.Exec("INSERT INTO stashes(token, StashName, Lifetime) values(?,?,?)", tok, StashNom, lifeTime)
	return err
}

func createNewTable(db *sql.DB, token string, fil []api.StashFile) error {
    if fil == nil {
        return DError{time.Now(), "No stashfile!"}
    }
    fmt.Println(token)
    ex, err1 := db.Prepare("CREATE TABLE IF NOT EXISTS " + token + " (ID INTEGER PRIMARY KEY AUTOINCREMENT, Fname STRING, size INT, type STRING, numD INT, sqltime TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL)")
    if err1 != nil {
        return DError{time.Now(), "Couldn't create table!"}
    }
    _, err := ex.Exec()
    if err != nil {
        return DError{time.Now(), "Couldn't create table!"}
    }
    for _, element := range fil {
        ex, _ := db.Prepare("INSERT INTO "+ token +" (Fname, Size, Type, numD) values(?,?,?,?)")
	_, error := ex.Exec(element.Fname, element.Size, element.Type, element.Download)
        if error != nil {
            return DError{time.Now(), "Couldn't insert into table!"}
        }
    }
    return nil
}

//Inserts a new deadrop into the database. 
func InsertStash(db *sql.DB, s *api.Stash) error {
	error := addToHashD(db, s.Token, s.StashName, s.Lifetime)
	if error == nil {
		error = createNewTable(db, s.Token, s.Files)
	}
	return error
}
