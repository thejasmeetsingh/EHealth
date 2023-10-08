package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load env varriables
	godotenv.Load()

	mode := os.Getenv("GIN_MODE")

	if mode == "" || mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()
	r.GET("/health-check/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Up and Running!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not configured in the enviorment")
	}

	r.Run(":" + port)
}
