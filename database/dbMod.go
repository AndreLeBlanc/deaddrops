package database

import (
	"database/sql"
)

func Init() *sql.DB {
	Db, err := sql.Open("sqlite3", "database/deadrops.db")
	CheckErr(err)
	return Db
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Close(db *sql.DB) {
	db.Close()
}

func (e DError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}
