package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/khaingminhtun/realtimechatty/db"
	"github.com/khaingminhtun/realtimechatty/db/redis"
	"github.com/khaingminhtun/realtimechatty/internal/app"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		// If running from inside the 'cmd' directory, look one level up
		err = godotenv.Load("../.env")
		if err != nil {
			// If running from a deeper nested subdirectory, look two levels up
			err = godotenv.Load("../../.env")
		}
	}

	if err != nil {
		log.Println("Note: No .env file detected, using system environment variables instead")
	} else {
		log.Println("Successfully loaded configuration variables from .env file")
	}
	pool, err := db.Connect()
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	defer pool.Close()

	log.Println("Connected to database")

	rdb, err := redis.InitRedis()
	if err != nil {
		log.Fatal("failed to initialize redis: ", err)
	}
	defer rdb.Close()
	log.Println("Successfully connected to Redis cache")

	app := app.NewApp(pool,rdb)

	app.Run(":8080")
}
