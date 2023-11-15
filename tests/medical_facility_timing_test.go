package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAddMedicalFacilityTimingAPI(t *testing.T) {
	// login the testing user
	authResponse, err := loginUser([]byte(doctorCredentials))
	if err != nil {
		t.Error(err.Error())
	}

	payload := []byte(`{"weekday": "1", "start_datetime": "5:00 AM", "end_datetime": "6:00 AM"}`)

	// Call medical facility timing POST API and check the respone Code
	response := getResponseRecorder(http.MethodPost, "/v1/medical-facility-timing/", authResponse.Data.Access, payload)

	if response.Code != http.StatusCreated {
		t.Errorf("Response: %v", response.Body.String())
	}
}

func TestGetMedicalFacilityTimingsAPI(t *testing.T) {
	// login the testing user
	authResponse, err := loginUser([]byte(doctorCredentials))
	if err != nil {
		t.Error(err.Error())
	}

	// Call medical facility timing GET API and check the respone Code
	response := getResponseRecorder(http.MethodGet, "/v1/medical-facility-timing/", authResponse.Data.Access, nil)

	if response.Code != http.StatusOK {
		t.Errorf("Response: %v", response.Body.String())
	}
}

func TestUpdateMedicalFacilityTiming(t *testing.T) {
	// login the testing user
	authResponse, err := loginUser([]byte(doctorCredentials))
	if err != nil {
		t.Error(err.Error())
	}

	// Call medical facility timing GET API
	medicalFacilityTimingResponse := getResponseRecorder(http.MethodGet, "/v1/medical-facility-timing/", authResponse.Data.Access, nil)

	if medicalFacilityTimingResponse.Code != http.StatusOK {
		t.Errorf("Response: %v", medicalFacilityTimingResponse.Body.String())
	}

	type ResponseList struct {
		Message string `json:"message"`
		Data    []struct {
			ID            uuid.UUID `json:"id"`
			CreatedAt     time.Time `json:"created_at"`
			ModifiedAt    time.Time `json:"modified_at"`
			Weekday       string    `json:"weekday"`
			StartDatetime string    `json:"start_datetime"`
			EndDatetime   string    `json:"end_datetime"`
		} `json:"data"`
	}
	var responseList ResponseList

	// Decode the response list
	decoder := json.NewDecoder(medicalFacilityTimingResponse.Body)
	if err := decoder.Decode(&responseList); err != nil {
		t.Errorf(err.Error())
	}

	url := fmt.Sprintf("/v1/medical-facility-timing/%s/", responseList.Data[0].ID.String())

	weekDay := "2"
	payload := []byte(fmt.Sprintf(`{"weekday": "%s"}`, weekDay))

	// Call medical facility timing update API
	response := getResponseRecorder(http.MethodPatch, url, authResponse.Data.Access, payload)

	type Response struct {
		Message string `json:"message"`
		Data    struct {
			ID            uuid.UUID `json:"id"`
			CreatedAt     time.Time `json:"created_at"`
			ModifiedAt    time.Time `json:"modified_at"`
			Weekday       string    `json:"weekday"`
			StartDatetime string    `json:"start_datetime"`
			EndDatetime   string    `json:"end_datetime"`
		} `json:"data"`
	}

	var updatedResponse Response

	// Decode the updated response
	decoder = json.NewDecoder(response.Body)
	if err := decoder.Decode(&updatedResponse); err != nil {
		t.Errorf(err.Error())
	}

	if response.Code != http.StatusOK || updatedResponse.Data.Weekday != weekDay {
		t.Errorf("Response: %v", response.Body.String())
	}
}
