package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/models"
)

// API for fetching medical facilities sorted by distance, Nearest on top
func (apiCfg *ApiCfg) MedicalFacilityListing(c *gin.Context) {
	// Validate the location coordinates
	rawLatitude := c.Query("lat")
	rawLongitude := c.Query("lng")

	if rawLatitude == "" || rawLongitude == "" {
		ErrorResponse(c, http.StatusBadRequest, "Invalid location coordinates")
		return
	}

	latitude, err := strconv.ParseFloat(rawLatitude, 64)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid location coordinates")
		return
	}

	longitude, err := strconv.ParseFloat(rawLongitude, 64)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid location coordinates")
		return
	}

	// Fetch the medical facilities based on the coordinates sent in query params
	dbMedicalFacilities, err := apiCfg.DB.MedicalFacilityListing(c, database.MedicalFacilityListingParams{
		StMakepoint:   latitude,
		StMakepoint_2: longitude,
	})

	if err != nil {
		ErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Error while fetching medical facilities: %v", err.Error()))
		return
	}

	SuccessResponse(c, http.StatusOK, "", models.DatabaseMedicalFacilitiesToMedicalFacilities(dbMedicalFacilities))
}
