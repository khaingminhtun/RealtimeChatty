package main

import (
	"log"

	"github.com/khaingminhtun/realtimechatty/db"
	"github.com/khaingminhtun/realtimechatty/internal/app"
)

func main() {
	pool, err := db.Connect()
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	defer pool.Close()

	log.Println("Connected to database")

	app := app.NewApp(pool)

	app.Run(":8080")
}
