package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/config"
	"github.com/thejasmeetsingh/EHealth/handlers"
	"github.com/thejasmeetsingh/EHealth/middlewares"
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
	v1.POST("/login/", apiCfg.Login)
	v1.GET("/profile/", middlewares.JWTAuth(apiCfg, apiCfg.GetUserProfile))
	v1.PATCH("/profile/", middlewares.JWTAuth(apiCfg, apiCfg.UpdateUserProfile))
	v1.DELETE("/profile/", middlewares.JWTAuth(apiCfg, apiCfg.DeleteUserProfile))
	v1.PUT("/change-password/", middlewares.JWTAuth(apiCfg, apiCfg.ChangePassword))
	v1.POST("/reset-password/", apiCfg.ResetPassword)
	return router
}
