package server

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "strings"
    "os"
    "api/dbMod"
    "time"
)

type fil struct {
	hash, size, time int
	format string
}

func add(db * sql.DB, fil fille) {
    stmt, err := db.Prepare("INSERT INTO deaddrops(hash, time, size, format) values(?,?,?,?)")
    checkErr(err)

   	stmt.Exec(fille.hash, fille.size, fille.time, fille.format)
}

 func skrivUt(db * sql.DB) {
 	rows, err := db.Query("SELECT * FROM deaddrops")
    checkErr(err)

        for rows.Next() {
	        var hash int
	        var size int
	        var time int
	        var format string
	        err = rows.Scan(&hash, &size, &time, &format)
	        checkErr(err)
	        fmt.Println(hash , " | ", size, " | ", time, " | ", format)
    }
}


func supervisorUp(db * sql.DB, token string, c chan string, cm *ChanMap) {

    select {
    case fname := <-c:
        fmt.Printf("received filename: %s\n", fname)
    case <-time.After(time.Second * 10):
        fmt.Println("timeout")
    }

}