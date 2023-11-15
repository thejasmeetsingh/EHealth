package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const doctorCredentials = `{"email": "testing-doc@example.com", "password": "12345678Aa@"}`

func TestAddMedicalFacilityAPI(t *testing.T) {
	// Create a testing doctor user
	authResponse, err := createTestingUser([]byte(doctorCredentials))

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// Prepare medical facility payload and call the API
	payload := []byte(`{"type": "ID", "name": "test facility", "description": "lorem ipsum", "email": "contact@example.com", "mobile_number": "+911234567890", "charges": 4.00, "address": "asdas", "location": {"lat": 37.7749, "lng": -122.4194}}`)

	medicalFacilityResponse := getResponseRecorder(http.MethodPost, "/v1/medical-facility/", authResponse.Data.Access, payload)

	if medicalFacilityResponse.Code != http.StatusCreated {
		t.Errorf("Response: %v", medicalFacilityResponse.Body.String())
	}
}

func TestGetMedicalFacilityDetails(t *testing.T) {
	// login the testing user
	authResponse, err := loginUser([]byte(doctorCredentials))

	if err != nil {
		t.Error(err.Error())
	}

	// Call medical facility details API and check the respone Code
	response := getResponseRecorder(http.MethodGet, "/v1/medical-facility/", authResponse.Data.Access, nil)

	if response.Code != http.StatusOK {
		t.Errorf("Response: %v", response.Body.String())
	}
}

func TestUpdateMedicalFacilityAPI(t *testing.T) {
	// login the testing user
	authResponse, err := loginUser([]byte(doctorCredentials))

	if err != nil {
		t.Error(err.Error())
	}

	newMedicalFacilityName := "test facility 1"

	// Prepare medical facility payload and call the API
	payload := []byte(fmt.Sprintf(`{"name": "%s"}`, newMedicalFacilityName))

	medicalFacilityResponse := getResponseRecorder(http.MethodPatch, "/v1/medical-facility/", authResponse.Data.Access, payload)

	type Response struct {
		Message string `json:"message"`
		Data    struct {
			Type         string  `json:"type"`
			Name         string  `json:"name"`
			Description  string  `json:"description"`
			Email        string  `json:"email"`
			MobileNumber string  `json:"mobile_number"`
			Charges      float64 `json:"charges"`
			Address      string  `json:"address"`
			Location     struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"data"`
	}
	var response Response

	// Decode the response and check value is updated or not
	decoder := json.NewDecoder(medicalFacilityResponse.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Errorf(err.Error())
	}

	if medicalFacilityResponse.Code != http.StatusOK || response.Data.Name != newMedicalFacilityName {
		t.Errorf("Response: %v", medicalFacilityResponse.Body.String())
	}
}
