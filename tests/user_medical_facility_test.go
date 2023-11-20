package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

const (
	latitude  = "27.644800"
	longitude = "77.216721"
)

func TestMedicalFacilityListing(t *testing.T) {
	// Login the user with given credentials and aquired the token
	credentials := []byte(userCredentials)

	authResponse, err := loginUser(credentials)
	if err != nil {
		t.Errorf(err.Error())
	}

	url := fmt.Sprintf("/v1/facility/?lat=%s&lng=%s", latitude, longitude)
	response := getResponseRecorder(http.MethodGet, url, authResponse.Data.Access, nil)

	if response.Code != http.StatusOK {
		t.Errorf("Response: %s", response.Body.String())
	}
}

func TestMedicalFacilityDetail(t *testing.T) {
	// Login the user with given credentials and aquired the token
	credentials := []byte(userCredentials)

	authResponse, err := loginUser(credentials)
	if err != nil {
		t.Errorf(err.Error())
	}

	medicalListingUrl := fmt.Sprintf("/v1/facility/?lat=%s&lng=%s", latitude, longitude)
	medicalListingRecorder := getResponseRecorder(http.MethodGet, medicalListingUrl, authResponse.Data.Access, nil)

	if medicalListingRecorder.Code != http.StatusOK {
		t.Errorf("Response: %s", medicalListingRecorder.Body.String())
	}

	type Response struct {
		Message string `json:"message"`
		Data    []struct {
			ID       uuid.UUID `json:"id"`
			Type     string    `json:"type"`
			Name     string    `json:"name"`
			Charges  string    `json:"charges"`
			Address  string    `json:"address"`
			Distance string    `json:"distance"`
		} `json:"data"`
	}
	var response Response

	decoder := json.NewDecoder(medicalListingRecorder.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Errorf(err.Error())
	}

	url := fmt.Sprintf("/v1/facility/%v/?lat=%s&lng=%s", response.Data[0].ID, latitude, longitude)
	medicalDetailRecorder := getResponseRecorder(http.MethodGet, url, authResponse.Data.Access, nil)

	if medicalDetailRecorder.Code != http.StatusOK {
		t.Errorf("Response: %s", medicalDetailRecorder.Body.String())
	}
}
