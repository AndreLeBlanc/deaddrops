package database

import (
	"database/sql"
	//"deadrop/api"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// Ta bort en fil från en stash i databasen. Funktionen behöver även

// uppdatera json fältet, men det har Joel och Erik en funktion för

// säkert.

func RemoveFile(db *sql.DB, token string, fname string) error {
	var numD int
	error := db.QueryRow("SELECT numD FROM "+token+" WHERE Fname=?", fname).Scan(&numD)
	if error != nil {
		return DError{time.Now(), "Couldn't find File"}
	}

	if numD < 1 {
		RemoveFileHard(db, token, fname)
	} else {
		update, error := db.Prepare("UPDATE " + token + " SET numD=? WHERE Fname=?")
		if error != nil {
			return DError{time.Now(), "Couldn't update File"}
		}
		update.Exec(numD-1, fname)
	}
	return nil
}

func RemoveFileHard(db *sql.DB, token string, fname string) error {
	ut, error := db.Prepare("delete from " + token + " where Fname=?")
	if error != nil {
		return DError{time.Now(), "Couldn't remove File"}
	}
	ut.Exec(fname)
	return nil
}

// Ta bort en stash från databasen.

func RemoveStash(db *sql.DB, token string) error {
	ut, error := db.Prepare("delete from stashes where token=?")
	if error != nil {
		return DError{time.Now(), "No such stash"}
	}

	ut.Exec(token)
	if error != nil {
		return DError{time.Now(), "Couldn't delete Stash"}
	}

	db.Exec("DROP TABLE " + token)

	if error != nil {
		return DError{time.Now(), "Couldn't prepare tabledrop"}
	}

	return nil
}
