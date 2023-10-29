package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/thejasmeetsingh/EHealth/config"
)

func main() {
	// Load env varriables
	godotenv.Load()

	// Get the rounter
	router := config.GetRouter(false)

	// Added health check endpoint which can be called to ensure that server is up and running
	router.GET("/health-check/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Up and Running!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not configured in the enviorment")
	}

	// Run the server on the given port
	router.Run(":" + port)
}
