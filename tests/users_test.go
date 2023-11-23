package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestGetUserProfileAPI(t *testing.T) {
	// Login the user with given credentials and aquired the token
	credentials := []byte(userCredentials)

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
	// Login the user with given credentials and aquired the token
	credentials := []byte(userCredentials)

	authResponse, err := loginUser(credentials)
	if err != nil {
		t.Errorf(err.Error())
	}

	newEmail := "testing-user1@example.com"

	payload := []byte(fmt.Sprintf(`{"email": "%s"}`, newEmail))

	// Call the profile update API and parse the response for validation check
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

	userCredentials = `{"email": "testing-user1@example.com", "password": "12345678Aa@", "is_end_user": true}`
}

func TestChangePasswordAPI(t *testing.T) {
	// Login the user with given credentials and aquired the token
	credentials := []byte(userCredentials)

	authResponse, err := loginUser(credentials)
	if err != nil {
		t.Errorf(err.Error())
	}

	payload := []byte(`{"current_password": "12345678Aa@", "new_password": "12345678Ab@", "new_password_confirm": "12345678Ab@"}`)

	response := getResponseRecorder(http.MethodPut, "/v1/change-password/", authResponse.Data.Access, payload)

	if response.Code != http.StatusOK {
		t.Errorf("Response: %s", response.Body.String())
	}
}
