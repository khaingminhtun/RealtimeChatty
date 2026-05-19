package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/realtimechatty/db"
)

func main() {
	pool, err := db.Connect()
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	defer pool.Close()

	log.Println("Connected to database")

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	r.Run(":8080")
}
