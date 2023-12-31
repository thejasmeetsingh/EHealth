package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/handlers"
	"github.com/thejasmeetsingh/EHealth/middlewares"
)

// Create and return a router instance
func GetRouter(isTest bool) *gin.Engine {
	// Check the application mode
	mode := os.Getenv("GIN_MODE")

	if mode == "" || mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Get the DB connection instance
	dbConn := getDBConn(isTest)
	apiCfg := handlers.ApiCfg{
		DB: dbConn,
	}

	// Create a default router with default logging and recovery middleware
	router := gin.Default()

	htmlTemplatePath := "templates/*.html"

	if isTest {
		htmlTemplatePath = "../" + htmlTemplatePath
	}

	// Set the default templates path and endpoints
	router.LoadHTMLGlob(htmlTemplatePath)
	router.GET("/reset-password/:token/", apiCfg.RenderResetPassword)
	router.POST("/validate-password/", apiCfg.ValidateResetPassword)

	// Create groups and attach respected middlewares
	v1 := router.Group("/v1")

	authResources := v1.Group("")
	authResources.Use(middlewares.JWTAuth(apiCfg))

	nonEndUserResources := authResources.Group("")
	nonEndUserResources.Use(middlewares.NonEndUser)

	endUserResources := authResources.Group("")
	endUserResources.Use(middlewares.EndUser)

	// Non-auth endpoints
	v1.POST("/signup/", apiCfg.Singup)
	v1.POST("/login/", apiCfg.Login)
	v1.POST("/reset-password/", apiCfg.ResetPassword)
	v1.POST("/refresh-token/", apiCfg.RefreshAccessToken)

	// Auth Endpoints
	authResources.GET("/profile/", apiCfg.GetUserProfile)
	authResources.PATCH("/profile/", apiCfg.UpdateUserProfile)
	authResources.DELETE("/profile/", apiCfg.DeleteUserProfile)
	authResources.PUT("/change-password/", apiCfg.ChangePassword)
	authResources.GET("/booking/", apiCfg.BookingList)
	authResources.GET("/booking/:id/", apiCfg.GetBooking)

	// Non end user resources
	nonEndUserResources.POST("/medical-facility/", apiCfg.AddMedicalFacility)
	nonEndUserResources.GET("/medical-facility/", apiCfg.GetMedicalFacilityDetails)
	nonEndUserResources.PATCH("/medical-facility/", apiCfg.UpdateMedicalFacility)

	nonEndUserResources.POST("/medical-facility-timing/", apiCfg.AddMedicalFacilityTiming)
	nonEndUserResources.GET("/medical-facility-timing/", apiCfg.GetMedicalFacilityTimings)
	nonEndUserResources.PATCH("/medical-facility-timing/:id/", apiCfg.UpdateMedicalFacilityTiming)
	nonEndUserResources.PATCH("/booking/:id/", apiCfg.UpdateBookingStatus)

	// End user resources
	endUserResources.GET("/facility/", apiCfg.MedicalFacilityListing)
	endUserResources.GET("/facility/:id/", apiCfg.MedicalFacilityDetail)
	endUserResources.POST("/booking/", apiCfg.CreateBooking)

	return router
}
