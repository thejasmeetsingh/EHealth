package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/internal/database"
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

	SuccessResponse(c, http.StatusOK, "Booking Created Successfully!", dbBooking)
}
