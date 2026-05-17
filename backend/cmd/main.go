package main

import (
	"log"

	"github.com/khaingminhtun/realtimechatty/db"
)

func main() {
	pool, err := db.Connect()
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	defer pool.Close()

	log.Println("Connected to database")
}
