package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/internal/database"
)

type bookingMedicalFacility struct {
	ID           uuid.UUID `json:"id"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Email        string    `json:"email"`
	MobileNumber string    `json:"mobile_number"`
	Charges      string    `json:"charges"`
	Address      string    `json:"address"`
}

type bookingUser struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type bookingFacility struct {
	ID              uuid.UUID              `json:"id"`
	CreatedAt       time.Time              `json:"created_at"`
	ModifiedAt      time.Time              `json:"modified_at"`
	StartDatetime   time.Time              `json:"start_datetime"`
	EndDatetime     time.Time              `json:"end_datetime"`
	Status          database.BookingStatus `json:"status"`
	MedicalFacility bookingMedicalFacility `json:"medical_facility"`
}

type bookingEndUser struct {
	ID            uuid.UUID              `json:"id"`
	CreatedAt     time.Time              `json:"created_at"`
	ModifiedAt    time.Time              `json:"modified_at"`
	StartDatetime time.Time              `json:"start_datetime"`
	EndDatetime   time.Time              `json:"end_datetime"`
	Status        database.BookingStatus `json:"status"`
	User          bookingUser            `json:"user"`
}

type bookingListing struct {
	ID            uuid.UUID              `json:"id"`
	CreatedAt     time.Time              `json:"created_at"`
	ModifiedAt    time.Time              `json:"modified_at"`
	StartDatetime time.Time              `json:"start_datetime"`
	EndDatetime   time.Time              `json:"end_datetime"`
	Status        database.BookingStatus `json:"status"`
}

func DatabaseBookingToBookingMedicalFacility(dbBooking database.Booking, dbMedicalFacility database.GetMedicalFacilityByIdRow) bookingFacility {
	return bookingFacility{
		ID:            dbBooking.ID,
		CreatedAt:     dbBooking.CreatedAt,
		ModifiedAt:    dbBooking.ModifiedAt,
		StartDatetime: dbBooking.StartDatetime,
		EndDatetime:   dbBooking.EndDatetime,
		Status:        dbBooking.Status,
		MedicalFacility: bookingMedicalFacility{
			ID:           dbMedicalFacility.ID,
			Type:         string(dbMedicalFacility.Type),
			Name:         dbMedicalFacility.Name,
			Description:  dbMedicalFacility.Description.String,
			Email:        dbMedicalFacility.Email,
			MobileNumber: dbMedicalFacility.MobileNumber,
			Charges:      dbMedicalFacility.Charges,
			Address:      dbMedicalFacility.Address,
		},
	}
}

func DatabaseBookingToBookingUser(dbBooking database.Booking, dbUser database.User) bookingEndUser {
	return bookingEndUser{
		ID:            dbBooking.ID,
		CreatedAt:     dbBooking.CreatedAt,
		ModifiedAt:    dbBooking.ModifiedAt,
		StartDatetime: dbBooking.StartDatetime,
		EndDatetime:   dbBooking.EndDatetime,
		Status:        dbBooking.Status,
		User: bookingUser{
			ID:    dbUser.ID,
			Name:  dbUser.Name.String,
			Email: dbUser.Email,
		},
	}
}

func DatebaseBookingListingToBookingListing(dbBookingListing []database.Booking) []bookingListing {
	var bookings []bookingListing

	for _, booking := range dbBookingListing {
		bookings = append(bookings, bookingListing{
			ID:            booking.ID,
			CreatedAt:     booking.CreatedAt,
			ModifiedAt:    booking.ModifiedAt,
			StartDatetime: booking.StartDatetime,
			EndDatetime:   booking.EndDatetime,
			Status:        booking.Status,
		})
	}

	return bookings
}
