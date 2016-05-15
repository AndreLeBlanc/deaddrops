package database

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "time"
)

type fil struct {
	hash, size, time, numD int
	format string
}

func add(db * sql.DB, fille fil) {
    stmt, err := db.Prepare("INSERT INTO userinfo(hash, time, size, format) values(?,?,?,?,?)")
    CheckErr(err)

   	stmt.Exec(fille.hash, fille.size, fille.time, fille.format, fille.numD)
}

 func skrivUt(db * sql.DB) {
 	rows, err := db.Query("SELECT * FROM userinfo")
    CheckErr(err)

        for rows.Next() {
	        var hash, size, time, numD int
	        var format string
	        err = rows.Scan(&hash, &size, &time, &format, &numD)
	        CheckErr(err)
	        fmt.Println(hash , " | ", size, " | ", time, " | ", format, " | ", numD)
    }
}


func SupervisorUp(db * sql.DB, token string, c chan string) {

    select {
    case fname := <-c:
        fmt.Printf("received filename: %s\n", fname)
        add(db, fil{2131231, 10, 333, 5, ".exe"})
    case <-time.After(time.Second * 10):
        fmt.Println("timeout")
    }

}