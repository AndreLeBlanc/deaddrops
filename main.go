package main

import (
	"database/sql"
	"deadrop/api"
	"deadrop/database"
	"deadrop/server"
	"fmt"
	"time"
	"sync"
)

func initstash() api.Stash {
	stash := api.NewEmptyStash()
	stash.Token = server.GenerateToken()
	stash.Lifetime = 999
	file := api.NewEmptyStashFile()
	file.Fname = "test1.txt"
	file.Download = 1
	file.Size = 100
	file.Type = "txt"
	stash.Files = append(stash.Files, file)
	return stash
}

func dbsupervisor(c chan int, wg sync.WaitGroup) {
	db := database.Init()
	defer database.Close(db)
	defer wg.Done()
	
	for {
		select {
		case i := <-c:
			stash := initstash()
			fmt.Printf("[%d]: ", i)
			fmt.Println(stash)
			err := database.InsertStash(db, &stash)
			database.CheckErr(err)
			if i >= 9 {
				return
			}
		}

	}
}

func insertstash(db *sql.DB, i int, wg sync.WaitGroup) {
	defer wg.Done()
	stash := initstash()
	fmt.Printf("[%d]: ", i)
	fmt.Println(stash)
	err := database.InsertStash(db, &stash)
	database.CheckErr(err)
}

func multithread() {
	var wg sync.WaitGroup
	wg.Add(10)

	db := database.Init()
	defer database.Close(db)

	for i := 0; i < 10; i++ {
		go insertstash(db, i, wg)
		time.Sleep(time.Second * 1)
	}

	wg.Wait()
}

func singlethread() {
	var wg sync.WaitGroup
	wg.Add(1)

	c := make(chan int)
	go dbsupervisor(c, wg)
	for i := 0; i < 10; i++ {
		c <- i
		time.Sleep(time.Second * 1)
	}

	wg.Wait()

}

func main() {
	//server.StartServer(server.InitServer())

	singlethread()
	// multithread()
}
