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

// API for creating booking record
func (apiCfg *ApiCfg) CreateBooking(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	type Parameters struct {
		MedicalFacilityID string    `json:"medical_facility_id"`
		StartDateTime     time.Time `json:"start_datetime"`
		EndDateTime       time.Time `json:"end_datetime"`
	}

	var params Parameters

	// Check if start datetime is less than end datetime or not
	if params.StartDateTime.Before(params.EndDateTime) || params.StartDateTime.Equal(params.EndDateTime) {
		ErrorResponse(c, http.StatusBadRequest, "End DateTime should be greater than Start DateTime")
		return
	}

	// parse medical facility ID
	medicalFacilityID, err := uuid.Parse(params.MedicalFacilityID)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid medical facility ID")
		return
	}

	// Create booking record
	dbBooking, err := apiCfg.DB.CreateBooking(c, database.CreateBookingParams{
		ID:                uuid.New(),
		CreatedAt:         time.Now().UTC(),
		ModifiedAt:        time.Now().UTC(),
		MedicalFacilityID: medicalFacilityID,
		UserID:            dbUser.ID,
		StartDatetime:     params.StartDateTime,
		EndDatetime:       params.EndDateTime,
		Status:            database.BookingStatusP,
	})

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while creating booking: %v", err.Error()))
		return
	}

	dbMedicalFacility, err := apiCfg.DB.GetMedicalFacilityById(c, medicalFacilityID)

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while fetching medical facility details: %v", err.Error()))
		return
	}

	SuccessResponse(c, http.StatusOK, "Booking Created Successfully!", models.DatabaseBookingToBookingMedicalFacility(dbBooking, dbMedicalFacility))
}

// API for getting booking details based on booking ID
func (apiCfg *ApiCfg) GetBooking(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	// Parse booking ID
	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid medical facility ID")
		return
	}

	dbBooking, err := apiCfg.DB.GetBooking(c, bookingID)

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while fetching booking details: %v", err.Error()))
		return
	}

	if dbUser.IsEndUser {
		dbMedicalFacility, err := apiCfg.DB.GetMedicalFacilityById(c, dbBooking.MedicalFacilityID)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while fetching medical facility details: %v", err.Error()))
			return
		}

		bookingResponse := models.DatabaseBookingToBookingMedicalFacility(dbBooking, dbMedicalFacility)
		SuccessResponse(c, http.StatusOK, "", bookingResponse)
	} else {
		bookingResponse := models.DatabaseBookingToBookingUser(dbBooking, dbUser)
		SuccessResponse(c, http.StatusOK, "", bookingResponse)
	}
}
