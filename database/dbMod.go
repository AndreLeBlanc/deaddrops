package database

import (
	"database/sql"
	"time"
	"fmt"
)

//A struct for storing error messages from the database package
type DError struct{
	When time.Time
	What string
}


// Initialises a database connection.
func Init() *sql.DB {
	Db, err := sql.Open("sqlite3", "database/deadrops.db")
	CheckErr(err)
	return Db
}

//Crashes the server if an error is sent to it.
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Closes database connection. Extremely important. 
func Close(db *sql.DB) {
	db.Close()
}

//Prints error
func (e DError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}
