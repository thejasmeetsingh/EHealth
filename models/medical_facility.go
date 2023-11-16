package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/internal/database"
)

type Coordinates struct {
	Lat interface{} `json:"lat"`
	Lng interface{} `json:"lng"`
}

type medicalFacility struct {
	ID           uuid.UUID   `json:"id"`
	CreatedAt    time.Time   `json:"created_at"`
	ModifiedAt   time.Time   `json:"modified_at"`
	Type         string      `json:"type"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Email        string      `json:"email"`
	MobileNumber string      `json:"mobile_number"`
	Charges      string      `json:"charges"`
	Address      string      `json:"address"`
	Location     interface{} `json:"location"`
}

func DatabaseMedicalFacilityToMedicalFacility(dbMedicalFacility database.GetMedicalFacilityByUserIdRow) medicalFacility {
	return medicalFacility{
		ID:           dbMedicalFacility.ID,
		CreatedAt:    dbMedicalFacility.CreatedAt,
		ModifiedAt:   dbMedicalFacility.ModifiedAt,
		Type:         string(dbMedicalFacility.Type),
		Name:         dbMedicalFacility.Name,
		Description:  dbMedicalFacility.Description.String,
		Email:        dbMedicalFacility.Email,
		MobileNumber: dbMedicalFacility.MobileNumber,
		Charges:      dbMedicalFacility.Charges,
		Address:      dbMedicalFacility.Address,
		Location: Coordinates{
			Lat: dbMedicalFacility.Lat,
			Lng: dbMedicalFacility.Lng,
		},
	}
}

type medicalFacilityListing struct {
	ID       uuid.UUID `json:"id"`
	Type     string    `json:"type"`
	Name     string    `json:"name"`
	Charges  string    `json:"charges"`
	Address  string    `json:"address"`
	Distance string    `json:"distance"`
}

func DatabaseMedicalFacilitiesToMedicalFacilities(dbMedicalFacilities []database.MedicalFacilityListingRow) []medicalFacilityListing {
	var medicalFacilities []medicalFacilityListing

	for _, dbMedicalFacility := range dbMedicalFacilities {
		medicalFacilities = append(medicalFacilities, medicalFacilityListing{
			ID:       dbMedicalFacility.ID,
			Type:     string(dbMedicalFacility.Type),
			Name:     dbMedicalFacility.Name,
			Charges:  dbMedicalFacility.Charges,
			Address:  dbMedicalFacility.Address,
			Distance: fmt.Sprintf("%.2f", dbMedicalFacility.Distance),
		})
	}

	return medicalFacilities
}
