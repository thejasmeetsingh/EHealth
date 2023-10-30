package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestGetUserProfileAPI(t *testing.T) {
	// Login the user with given credentials and aquired the token
	credentials := []byte(`{"email": "testing-user@example.com", "password": "12345678Aa@"}`)

	authResponse, err := loginUser(credentials)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Call GetUserProfile API
	response := getResponseRecorder(http.MethodGet, "/v1/profile/", authResponse.Data.Access, nil)

	if response.Code != http.StatusOK {
		t.Errorf("Response: %v", response.Body.String())
	}
}

func TestUpdateUserProfileAPI(t *testing.T) {
	credentials := []byte(`{"email": "testing-user@example.com", "password": "12345678Aa@"}`)

	authResponse, err := loginUser(credentials)
	if err != nil {
		t.Errorf(err.Error())
	}

	newEmail := "testing-user1@example.com"

	payload := []byte(fmt.Sprintf(`{"email": "%s"}`, newEmail))

	profileResponse := getResponseRecorder(http.MethodPatch, "/v1/profile/", authResponse.Data.Access, payload)

	type Response struct {
		Message string `json:"message"`
		Data    struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			IsEndUser bool   `json:"is_end_user"`
		} `json:"data"`
	}

	var response Response

	decoder := json.NewDecoder(profileResponse.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Errorf(err.Error())
	}

	if profileResponse.Code != http.StatusOK || response.Data.Email != newEmail {
		t.Errorf("Response: %v", profileResponse.Body.String())
	}
}
