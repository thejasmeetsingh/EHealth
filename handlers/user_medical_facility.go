package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/models"
)

// A common function for converting lat, lng string value to float
func parseLatLng(rawLatitude, rawLongitude string) (float64, float64, error) {
	if rawLatitude == "" || rawLongitude == "" {
		return 0, 0, fmt.Errorf("invalid location coordinates")
	}

	latitude, err := strconv.ParseFloat(rawLatitude, 64)

	if err != nil {
		return 0, 0, err
	}

	longitude, err := strconv.ParseFloat(rawLongitude, 64)

	if err != nil {
		return 0, 0, err
	}

	return latitude, longitude, nil
}

// API for fetching medical facilities sorted by distance, Nearest on top
func (apiCfg *ApiCfg) MedicalFacilityListing(c *gin.Context) {
	latitude, longitude, err := parseLatLng(c.Query("lat"), c.Query("lng"))

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
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

// API for fetching medical facility details
func (apiCfg *ApiCfg) MedicalFacilityDetail(c *gin.Context) {
	latitude, longitude, err := parseLatLng(c.Query("lat"), c.Query("lng"))

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// parse medical facility ID
	medicalFacilityIDStr := c.Param("id")
	if medicalFacilityIDStr == "" {
		ErrorResponse(c, http.StatusBadRequest, "Invalid facility ID")
		return
	}

	medicalFacilityID, err := uuid.Parse(medicalFacilityIDStr)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid facility ID")
		return
	}

	// fetch medical facility details
	medicalFacility, err := apiCfg.DB.MedicalFacilityDetail(c, database.MedicalFacilityDetailParams{
		StMakepoint:   latitude,
		StMakepoint_2: longitude,
		ID:            medicalFacilityID,
	})

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	SuccessResponse(c, http.StatusOK, "", models.DatabaseMedicalFacilityDetailToMedicalFacilityDetail(medicalFacility))
}
