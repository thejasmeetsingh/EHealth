package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/config"
	"github.com/thejasmeetsingh/EHealth/handlers"
)

func getRouter() *gin.Engine {
	mode := os.Getenv("GIN_MODE")

	if mode == "" || mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	dbConn := config.GetDBConn()
	apiCfg := handlers.ApiCfg{
		DB: dbConn,
	}

	router := gin.Default()

	v1 := router.Group("/v1")
	v1.POST("/signup/", apiCfg.Singup)

	return router
}
