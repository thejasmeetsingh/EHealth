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
		IsTest            bool      `json:"is_test"`
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

	// Check if requested DateTime overlaps with any other accepted booking for the same medical facility
	acceptedOverlappingBookingCount, err := apiCfg.DB.OverlappingAcceptedBookingCount(c, database.OverlappingAcceptedBookingCountParams{
		Overlaps:          params.StartDateTime,
		Overlaps_2:        params.EndDateTime,
		MedicalFacilityID: medicalFacilityID,
	})

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if acceptedOverlappingBookingCount > 0 {
		ErrorResponse(c, http.StatusBadRequest, "This slot is already book. Please select other time")
		return
	}

	medicalFacilityTimings, err := apiCfg.DB.GetMedicalFacilityTimingDetails(c, medicalFacilityID)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !utils.IsBookingTimeInRange(medicalFacilityTimings, params.StartDateTime, params.EndDateTime) {
		ErrorResponse(c, http.StatusBadRequest, "No booking slot available at this time. Please select other time")
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

	if !params.IsTest {
		_, err = emails.SendBookingCreationEmail(dbUser, dbBooking, *c.Request)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	SuccessResponse(c, http.StatusCreated, "Booking Created Successfully!", models.DatabaseBookingToBookingMedicalFacility(dbBooking, dbMedicalFacility))
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
		ErrorResponse(c, http.StatusBadRequest, "Invalid Booking ID")
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

// Cancel multiple bookings and set their status to rejected
func (apiCfg *ApiCfg) CancelBookings(c *gin.Context, bookings []database.OverlappingPendingBookingsRow, isTest bool) {
	for _, booking := range bookings {
		apiCfg.DB.UpdateBookingStatus(c, database.UpdateBookingStatusParams{
			Status: database.BookingStatusR,
			ID:     booking.ID,
		})

		if !isTest {
			go emails.SendBookingRejectedEmail(map[string]string{
				"name":       booking.Name,
				"address":    booking.Address,
				"user_email": booking.Email,
				"start_dt":   booking.StartDatetime.Format(time.RFC822),
				"end_dt":     booking.EndDatetime.Format(time.RFC822),
			}, *c.Request)
		}
	}
}

func (apiCfg *ApiCfg) UpdateBookingStatus(c *gin.Context) {
	dbUser, err := getDBUser(c)
	if err != nil {
		ErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	// Parse booking ID
	bookingIDStr := c.Param("id")
	bookingID, err := uuid.Parse(bookingIDStr)

	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid Booking ID")
		return
	}

	// Fetch booking object
	dbBooking, err := apiCfg.DB.GetBooking(c, bookingID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while fetching booking details: %v", err.Error()))
		return
	}

	// parse booking status coming in request data
	type Parameters struct {
		Status string `json:"status" binding:"required"`
		IsTest bool   `json:"is_test"`
	}
	var params Parameters

	if err = c.ShouldBindJSON(&params); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if status is pending in the request data
	if params.Status == string(database.BookingStatusP) {
		ErrorResponse(c, http.StatusBadRequest, "Invalid booking status")
		return
	}

	// Fetch user object
	bookingUser, err := apiCfg.DB.GetUserById(c, dbBooking.UserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while fetching booking user: %v", err.Error()))
		return
	}

	// Fetch medical facility object
	dbMedicalFacility, err := apiCfg.DB.GetMedicalFacilityById(c, dbBooking.MedicalFacilityID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while fetching booking medical facility: %v", err.Error()))
		return
	}

	if dbUser.ID != dbMedicalFacility.UserID {
		ErrorResponse(c, http.StatusForbidden, "You cannot access this resource")
		return
	}

	// Update the booking status
	dbBooking, err = apiCfg.DB.UpdateBookingStatus(c, database.UpdateBookingStatusParams{
		Status: database.BookingStatus(params.Status),
		ID:     bookingID,
	})

	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Error while updating the booking status: %v", err.Error()))
		return
	}

	// Check booking status
	// if booking is accepted, then reject the overlapping pending bookings
	// Otherwise send the booking rejected email to the user
	if params.Status == string(database.BookingStatusA) {

		// Fetch overlapping pending bookings, exlcuding the given booking
		pendingOverlappingBookings, err := apiCfg.DB.OverlappingPendingBookings(c, database.OverlappingPendingBookingsParams{
			MedicalFacilityID: dbBooking.MedicalFacilityID,
			Overlaps:          dbBooking.StartDatetime,
			Overlaps_2:        dbBooking.EndDatetime,
			ID:                bookingID,
		})

		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		// Cancel the overlapping bookings which have status pending
		go apiCfg.CancelBookings(c, pendingOverlappingBookings, params.IsTest)

		if !params.IsTest {
			// Send booking accepted email to end user
			_, err = emails.SendBookingAcceptedEmail(map[string]string{
				"name":       string(dbMedicalFacility.Name),
				"address":    dbMedicalFacility.Address,
				"user_email": bookingUser.Email,
				"start_dt":   dbBooking.StartDatetime.Format(time.RFC822),
				"end_dt":     dbBooking.EndDatetime.Format(time.RFC822),
			}, *c.Request)

			if err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
	} else {
		if !params.IsTest {
			// Send booking rejected email to end user
			_, err = emails.SendBookingRejectedEmail(map[string]string{
				"name":       string(dbMedicalFacility.Name),
				"address":    dbMedicalFacility.Address,
				"user_email": bookingUser.Email,
				"start_dt":   dbBooking.StartDatetime.Format(time.RFC822),
				"end_dt":     dbBooking.EndDatetime.Format(time.RFC822),
			}, *c.Request)

			if err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	SuccessResponse(c, http.StatusOK, "Booking status updated successfully!", models.DatabaseBookingToBookingUser(dbBooking, bookingUser))
}
