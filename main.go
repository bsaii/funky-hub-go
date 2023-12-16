package main

import (
	"log"
	"os"

	"github.com/bsaii/funky-hub-go/database"
	"github.com/bsaii/funky-hub-go/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	database.InitDb(dbUrl)

	gin.ForceConsoleColor()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to FunkyHub.",
		})
	})
	r.GET("/feeds", handlers.GetFeeds)
	r.POST("/feed", handlers.CreateFeed)
	r.GET("/feed/:id", handlers.GetFeedById)
	r.Run() // listen and serve on 0.0.0.0:8080
}
