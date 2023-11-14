package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/thejasmeetsingh/EHealth/internal/database"
)

const (
	TimeFormat = "3:04 PM"
)

type medicalFacilityTiming struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
	Weekday       string    `json:"weekday"`
	StartDatetime string    `json:"start_datetime"`
	EndDatetime   string    `json:"end_datetime"`
}

func DatabaseMedicalFacilityTimingToMedicalFacilityTiming(dbMedicalFacilityTiming database.MedicalFacilityTiming) medicalFacilityTiming {
	startDateTime := dbMedicalFacilityTiming.StartDatetime.Format(TimeFormat)
	endDateTime := dbMedicalFacilityTiming.EndDatetime.Format(TimeFormat)

	return medicalFacilityTiming{
		ID:            dbMedicalFacilityTiming.ID,
		CreatedAt:     dbMedicalFacilityTiming.CreatedAt,
		ModifiedAt:    dbMedicalFacilityTiming.ModifiedAt,
		Weekday:       string(dbMedicalFacilityTiming.Weekday),
		StartDatetime: startDateTime,
		EndDatetime:   endDateTime,
	}
}

func DatabaseMedicalFacilityTimingsToMedicalFacilityTimings(dbMedicalFacilityTimings []database.MedicalFacilityTiming) []medicalFacilityTiming {
	medicalFacilityTimings := []medicalFacilityTiming{}

	for _, dbMedicalFacilityTiming := range dbMedicalFacilityTimings {
		medicalFacilityTimings = append(medicalFacilityTimings, DatabaseMedicalFacilityTimingToMedicalFacilityTiming(dbMedicalFacilityTiming))
	}

	return medicalFacilityTimings
}
