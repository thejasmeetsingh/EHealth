package tests

import (
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
