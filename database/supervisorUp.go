package database

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "time"
)

type fil struct {
	hash, size, time int
	format string
}

func add(db * sql.DB, fille fil) {
    stmt, err := db.Prepare("INSERT INTO userinfo(hash, time, size, format) values(?,?,?,?)")
    CheckErr(err)

   	stmt.Exec(fille.hash, fille.size, fille.time, fille.format)
}

 func skrivUt(db * sql.DB) {
 	rows, err := db.Query("SELECT * FROM userinfo")
    CheckErr(err)

        for rows.Next() {
	        var hash int
	        var size int
	        var time int
	        var format string
	        err = rows.Scan(&hash, &size, &time, &format)
	        CheckErr(err)
	        fmt.Println(hash , " | ", size, " | ", time, " | ", format)
    }
}


func SupervisorUp(db * sql.DB, token string, c chan string) {

    select {
    case fname := <-c:
        fmt.Printf("received filename: %s\n", fname)
        add(db, fil{2131231, 10, 333, ".exe"})
    case <-time.After(time.Second * 10):
        fmt.Println("timeout")
    }

}