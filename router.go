package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/config"
	"github.com/thejasmeetsingh/EHealth/handlers"
	"github.com/thejasmeetsingh/EHealth/middlewares"
)

// Create and return a router instance
func getRouter() *gin.Engine {
	// Check the application mode
	mode := os.Getenv("GIN_MODE")

	if mode == "" || mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Get the DB connection instance
	dbConn := config.GetDBConn()
	apiCfg := handlers.ApiCfg{
		DB: dbConn,
	}

	// Create a default router with default logging and recovery middleware
	router := gin.Default()

	// Set the default templates path and endpoints
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/reset-password/:token/", apiCfg.RenderResetPassword)
	router.POST("/validate-password/", apiCfg.ValidateResetPassword)

	// Create a group depicting the API Version
	v1 := router.Group("/v1")

	// Add all the REST API endponts in the created group
	v1.POST("/signup/", apiCfg.Singup)
	v1.POST("/login/", apiCfg.Login)
	v1.GET("/profile/", middlewares.JWTAuth(apiCfg, apiCfg.GetUserProfile))
	v1.PATCH("/profile/", middlewares.JWTAuth(apiCfg, apiCfg.UpdateUserProfile))
	v1.DELETE("/profile/", middlewares.JWTAuth(apiCfg, apiCfg.DeleteUserProfile))
	v1.PUT("/change-password/", middlewares.JWTAuth(apiCfg, apiCfg.ChangePassword))
	v1.POST("/reset-password/", apiCfg.ResetPassword)
	v1.POST("/refresh-token/", apiCfg.RefreshAccessToken)

	return router
}
