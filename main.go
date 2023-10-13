package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// Load env varriables
	godotenv.Load()

	router := getRouter()
	router.GET("/health-check/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Up and Running!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not configured in the enviorment")
	}

	router.Run(":" + port)
}
