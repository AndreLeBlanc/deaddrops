package database

import (
	"database/sql"
	"deadrop/api"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func addToHashD(db *sql.DB, tok string, StashNom string, lifeTime int) error {
	_, err := db.Exec("INSERT INTO stashes(token, StashName, Lifetime) values(?,?,?)", tok, StashNom, lifeTime)
	return err
}

func createNewTable(db *sql.DB, token string, fil *[]api.StashFile) error {
	if fil == nil {
		return DError{time.Now(), "No stashfile!"}
	}

	for _, element := range *fil {
		_, error := db.Exec("INSERT INTO files(Token, Fname, Size, Type, numD) values(?,?,?,?,?)", token, element.Fname, element.Size, element.Type, element.Download)
		if error != nil {
			return DError{time.Now(), "Couldn't insert into table!"}
		}
	}

	return nil
}

func InsertStash(db *sql.DB, s *api.Stash) error {
	error := addToHashD(db, s.Token, s.StashName, s.Lifetime)
	if error == nil {
		error = createNewTable(db, s.Token, &s.Files)
	}
	return error
}
