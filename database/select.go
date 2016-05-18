package database

import (
	"database/sql"
	"deadrop/api"
	_ "github.com/mattn/go-sqlite3"
)

// Läs en stash från databasen.

func getStash(db *sql.DB, token string) Stash {
	var s Stash

	rows, err := db.Query("SELECT Lifetime FROM stashes WHERE token=?", token)
	if err != nil {
		return s
	}
	var Lifetime int
	var StashName string
	for rows.Next() {
		err = rows.Scan(&Lifetime)
		if err != nil {
			return s
		}
	}

	rows, err = db.Query("SELECT StashName FROM stashes WHERE token=?", token)
	if err != nil {
		return s
	}
	for rows.Next() {
		err = rows.Scan(&StashName)
		if err != nil {
			return s
		}
	}

	s.Token = token
	s.StashName = StashName
	s.Lifetime = Lifetime
	return s
}

func getRows(db *sql.DB, token string) []StashFile {
	var count int
	SFile := make([]StashFile, count)
	error := db.QueryRow("SELECT COUNT(*) FROM " + token).Scan(&count)
	if error != nil {
		return SFile
	}

	SFile = make([]StashFile, count)
	rows, error := db.Query("SELECT ID, Fname, Size, Type, numD FROM " + token)
	if error != nil {
		return SFile
	}
	return getStashFiles(rows, SFile)
}

func getStashFiles(rows *sql.Rows, SFile []StashFile) []StashFile {
	i := 0
	for rows.Next() {
		var Id, Size, Download int
		var Fname, Type string
		err := rows.Scan(&Id, &Fname, &Size, &Type, &Download)
		if err != nil {
			return SFile
		}
		SFile[i] = StashFile{Id, Fname, Size, Type, Download}
		i++
	}
	return SFile
}

func SelectStash(db *sql.DB, token string) Stash {
	myStash := getStash(db, token)
	myStash.Files = getRows(db, token)
	return myStash
}

func SelectLifetime(db *sql.DB, token string) int {
	var tok int
	db.QueryRow("SELECT Lifetime FROM stashes WHERE token=?", token).Scan(&tok)
	return tok
}
