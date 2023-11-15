package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/models"
)

// A common function for fetching the Medical Facility DB object
// This is to be used for all the medical facility timings handler
func (apiCfg *ApiCfg) getMedicalFacilityObject(c *gin.Context) (database.GetMedicalFacilityByUserIdRow, error) {
	dbUser, err := getDBUser(c)
	if err != nil {
		return database.GetMedicalFacilityByUserIdRow{}, err
	}

	// Fetch medical facility details from the DB
	dbMedicalFacility, err := apiCfg.DB.GetMedicalFacilityByUserId(c, dbUser.ID)
	if err != nil {
		return database.GetMedicalFacilityByUserIdRow{}, err
	}

	return dbMedicalFacility, nil
}

// Get timings of a medical facility
func (apiCfg *ApiCfg) GetMedicalFacilityTimings(c *gin.Context) {
	// Fetch medical facility details from the DB
	dbMedicalFacility, err := apiCfg.getMedicalFacilityObject(c)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Error while fetching facility details or facility does not exists")
		return
	}

	dbMedicalFacilityTimings, err := apiCfg.DB.GetMedicalFacilityTimingDetails(c, dbMedicalFacility.ID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while fetching medical facility timings: %v", err.Error()))
		return
	}

	SuccessResponse(c, http.StatusOK, "", models.DatabaseMedicalFacilityTimingsToMedicalFacilityTimings(dbMedicalFacilityTimings))
}

// API for adding medical facility timing details
func (apiCfg *ApiCfg) AddMedicalFacilityTiming(c *gin.Context) {
	// Fetch medical facility details from the DB
	dbMedicalFacility, err := apiCfg.getMedicalFacilityObject(c)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Error while fetching facility details or facility does not exists")
		return
	}

	type Parameters struct {
		Weekday       string `json:"weekday" binding:"required"`
		StartDateTime string `json:"start_datetime" binding:"required"`
		EndDateTime   string `json:"end_datetime" binding:"required"`
	}
	var params Parameters

	if err := c.ShouldBindJSON(&params); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request: %v", err.Error()))
		return
	}

	// Parse time string to time object
	startDateTime, err := time.Parse(models.TimeFormat, params.StartDateTime)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid Start datetime format: %v", err.Error()))
		return
	}

	endDateTime, err := time.Parse(models.TimeFormat, params.EndDateTime)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid End datetime format: %v", err.Error()))
		return
	}

	// Check wheather the start datetime is less than equal to end datetime or not
	if endDateTime.Hour() <= startDateTime.Hour() {
		ErrorResponse(c, http.StatusBadRequest, "End datetime should be greater than Start datetime")
		return
	}

	// Create medical facility timing record
	dbMedicalFacilityTimings, err := apiCfg.DB.AddMedicalFacilityTimings(c, database.AddMedicalFacilityTimingsParams{
		ID:                uuid.New(),
		CreatedAt:         time.Now().UTC(),
		ModifiedAt:        time.Now().UTC(),
		MedicalFacilityID: dbMedicalFacility.ID,
		Weekday:           database.WeekdayType(params.Weekday),
		StartDatetime:     startDateTime,
		EndDatetime:       endDateTime,
	})

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while adding medical facility timing: %v", err.Error()))
		return
	}

	SuccessResponse(c, http.StatusCreated, "Medical Facility Timings Details Added Successfully!", models.DatabaseMedicalFacilityTimingToMedicalFacilityTiming(dbMedicalFacilityTimings))
}

// API for updating a medical facility timing record
func (apiCfg *ApiCfg) UpdateMedicalFacilityTiming(c *gin.Context) {
	// Convert the passed ID to UUID
	medicalFacilityTimingIdStr := c.Param("id")
	medicalFacilityTimingId, err := uuid.Parse(medicalFacilityTimingIdStr)

	if err != nil {
		ErrorResponse(c, http.StatusForbidden, "Invalid ID")
		return
	}

	// Fetch medical facility detail of a given ID from DB
	dbMedicalFacilityTiming, err := apiCfg.DB.GetMedicalFacilityTimingById(c, medicalFacilityTimingId)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, "Medical Facility Detail with the given ID does not exists")
		return
	}

	type Parameters struct {
		Weekday       string `json:"weekday"`
		StartDateTime string `json:"start_datetime"`
		EndDateTime   string `json:"end_datetime"`
	}
	var params Parameters

	if err := c.ShouldBindJSON(&params); err != nil {
		ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while parsing the request: %v", err.Error()))
		return
	}

	startDateTime := dbMedicalFacilityTiming.StartDatetime
	endDateTime := dbMedicalFacilityTiming.EndDatetime
	weekDay := dbMedicalFacilityTiming.Weekday

	if params.Weekday != "" {
		weekDay = database.WeekdayType(params.Weekday)
	}

	// Parse time string to time object
	if params.StartDateTime != "" {
		startDateTime, err = time.Parse(models.TimeFormat, params.StartDateTime)
		if err != nil {
			ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid Start datetime format: %v", err.Error()))
			return
		}
	}

	if params.EndDateTime != "" {
		endDateTime, err = time.Parse(models.TimeFormat, params.EndDateTime)
		if err != nil {
			ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid End datetime format: %v", err.Error()))
			return
		}
	}

	// Check wheather the start datetime is less than equal to end datetime or not
	if endDateTime.Hour() <= startDateTime.Hour() {
		ErrorResponse(c, http.StatusBadRequest, "End datetime should be greater than Start datetime")
		return
	}

	// Update the details in DB
	dbMedicalFacilityTiming, err = apiCfg.DB.UpdateMedicalFacilityTimings(c, database.UpdateMedicalFacilityTimingsParams{
		Weekday:       weekDay,
		StartDatetime: startDateTime,
		EndDatetime:   endDateTime,
		ID:            medicalFacilityTimingId,
	})

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while updating medical facility timing: %v", err.Error()))
		return
	}

	SuccessResponse(c, http.StatusOK, "Medical Facility Timings Details Updated Successfully!", models.DatabaseMedicalFacilityTimingToMedicalFacilityTiming(dbMedicalFacilityTiming))
}
