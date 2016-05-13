package api

func newConnect() {
	db, err := sql.Open("sqlite3", "database/deaddrops.db")
    checkErr(err)
    return db
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}