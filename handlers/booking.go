package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/emails"
	"github.com/thejasmeetsingh/EHealth/internal/database"
	"github.com/thejasmeetsingh/EHealth/models"
	"github.com/thejasmeetsingh/EHealth/utils"
)

// API for creating booking record
func (apiCfg *ApiCfg) CreateBooking(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	type Parameters struct {
		MedicalFacilityID string    `json:"medical_facility_id" binding:"required"`
		StartDateTime     time.Time `json:"start_datetime" binding:"required"`
		EndDateTime       time.Time `json:"end_datetime" binding:"required"`
	}

	var params Parameters

	if err = c.ShouldBindJSON(&params); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if datetime is in the past or not
	if time.Now().After(params.StartDateTime) || time.Now().Equal(params.StartDateTime) {
		ErrorResponse(c, http.StatusBadRequest, "Start DateTime should be greater than current DateTime")
		return
	}

	// Check if start datetime is less than end datetime or not
	if params.StartDateTime.After(params.EndDateTime) || params.StartDateTime.Equal(params.EndDateTime) {
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

	_, err = emails.SendBookingCreationEmail(dbUser, dbBooking, *c.Request)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err.Error())
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

	// return booking detail response based on the user type
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

// API for fetching booking list
func (apiCfg *ApiCfg) BookingList(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	// Retrive booking status from query params
	bookingStatus := c.Query("status")

	if bookingStatus == "" {
		ErrorResponse(c, http.StatusBadRequest, "status parameter is required")
		return
	}

	// Check user type and return respected booking list response
	if dbUser.IsEndUser {
		bookings, err := apiCfg.DB.GetUserBookings(c, database.GetUserBookingsParams{
			UserID: dbUser.ID,
			Status: database.BookingStatus(bookingStatus),
			Limit:  10,
			Offset: utils.GetOffset(c),
		})

		if err != nil {
			ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while fetching booking records: %v", err.Error()))
			return
		}

		SuccessResponse(c, http.StatusOK, "", models.DatebaseBookingListingToBookingListing(bookings))
	} else {
		// Fetch medical facility object
		dbMedicalFacility, err := apiCfg.getMedicalFacilityObject(c)

		if err != nil {
			ErrorResponse(c, http.StatusBadRequest, "Error while fetching facility details or facility does not exists")
			return
		}

		bookings, err := apiCfg.DB.GetMedicalFacilityBookings(c, database.GetMedicalFacilityBookingsParams{
			MedicalFacilityID: dbMedicalFacility.ID,
			Status:            database.BookingStatus(bookingStatus),
			Limit:             10,
			Offset:            utils.GetOffset(c),
		})

		if err != nil {
			ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Error while fetching booking records: %v", err.Error()))
			return
		}

		SuccessResponse(c, http.StatusOK, "", models.DatebaseBookingListingToBookingListing(bookings))
	}
}

func (apiCfg *ApiCfg) UpdateBookingStatus(c *gin.Context) {

}
