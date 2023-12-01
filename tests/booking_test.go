package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

var BookingID string

func TestBookingCreateAPI(t *testing.T) {
	// Create a testing doctor user
	authResponse, err := createTestingUser([]byte(`{"email": "testing-doc1@example.com", "password": "12345678Aa@"}`))

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// Prepare medical facility payload and call the API
	payload := []byte(`{"type": "ID", "name": "test facility", "description": "lorem ipsum", "email": "contact@example.com", "mobile_number": "+911234567890", "charges": 4.00, "address": "asdas", "location": {"lat": 37.7749, "lng": -122.4194}}`)

	medicalFacility := getResponseRecorder(http.MethodPost, "/v1/medical-facility/", authResponse.Data.Access, payload)

	if medicalFacility.Code != http.StatusCreated {
		t.Errorf("Response: %v", medicalFacility.Body.String())
	}

	type Response struct {
		Message string `json:"message"`
		Data    struct {
			ID string `json:"id"`
		}
	}
	var medicalFacilityResponse Response

	// Decode the medical facility create response
	decoder := json.NewDecoder(medicalFacility.Body)
	if err := decoder.Decode(&medicalFacilityResponse); err != nil {
		t.Errorf(err.Error())
	}

	// Login the user with given credentials and aquired the token
	credentials := []byte(userCredentials)

	authResponse, err = loginUser(credentials)
	if err != nil {
		t.Errorf(err.Error())
	}

	// prepare data for creating booking
	startDateTime := time.Now().Add(time.Hour).UTC()
	endDateTime := startDateTime.Add(time.Hour).UTC()

	payload = []byte(fmt.Sprintf(`{"medical_facility_id": "%s", "start_datetime": "%s", "end_datetime": "%s", "is_test": true}`, medicalFacilityResponse.Data.ID, startDateTime.Format(time.RFC3339), endDateTime.Format(time.RFC3339)))

	// call booking API and check response code
	booking := getResponseRecorder(http.MethodPost, "/v1/booking/", authResponse.Data.Access, payload)

	if booking.Code != http.StatusCreated {
		t.Errorf("Response: %s", booking.Body.String())
	}

	var bookingResponse Response

	// decode booking response and store booking ID
	decoder = json.NewDecoder(booking.Body)
	if err := decoder.Decode(&bookingResponse); err != nil {
		t.Errorf(err.Error())
	}

	BookingID = bookingResponse.Data.ID
}

func TestBookingListAPI(t *testing.T) {
	// Login the user with given credentials and aquired the token
	credentials := []byte(userCredentials)

	authResponse, err := loginUser(credentials)
	if err != nil {
		t.Errorf(err.Error())
	}

	response := getResponseRecorder(http.MethodGet, "/v1/booking/?status=P", authResponse.Data.Access, nil)

	if response.Code != http.StatusOK {
		t.Errorf("Response: %s", response.Body.String())
	}
}

func TestBookingDetailAPI(t *testing.T) {
	// Login the user with given credentials and aquired the token
	credentials := []byte(userCredentials)

	authResponse, err := loginUser(credentials)
	if err != nil {
		t.Errorf(err.Error())
	}

	url := fmt.Sprintf("/v1/booking/%s/", BookingID)
	response := getResponseRecorder(http.MethodGet, url, authResponse.Data.Access, nil)

	if response.Code != http.StatusOK {
		t.Errorf("Response: %s", response.Body.String())
	}
}
