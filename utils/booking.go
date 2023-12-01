package utils

import (
	"fmt"
	"time"

	"github.com/thejasmeetsingh/EHealth/internal/database"
)

// Check weather or not booking datetime is in range with the medical facility timings
func IsBookingTimeInRange(medicalFacilityTimings []database.MedicalFacilityTiming, startDateTime time.Time, endDateTime time.Time) bool {
	timeFormat := "15:04:05"

	// Run a loop and find the start and end time of medical facility based on the booking start datetime weekday
	for _, medicalFacilityTiming := range medicalFacilityTimings {
		if medicalFacilityTiming.Weekday == database.WeekdayType(fmt.Sprintf("%d", int(startDateTime.Weekday()))) {
			// Fetch the time value from the datetime object
			medicalFacilityStartTime := medicalFacilityTiming.StartDatetime.Format(timeFormat)
			medicalFacilityEndTime := medicalFacilityTiming.EndDatetime.Format(timeFormat)

			bookingStartTime := startDateTime.Format(timeFormat)
			bookingEndTime := endDateTime.Format(timeFormat)

			return bookingStartTime >= medicalFacilityStartTime && bookingEndTime <= medicalFacilityEndTime
		}
	}
	return false
}
