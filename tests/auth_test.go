package tests

import (
	"fmt"
	"net/http"
	"testing"
)

var (
	userCredentials = `{"email": "testing-user@example.com", "password": "12345678Aa@", "is_end_user": true}`
)

func TestSignUpAPI(t *testing.T) {
	payload := []byte(userCredentials)
	_, err := createTestingUser(payload)

	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func TestLoginAPI(t *testing.T) {
	credentials := []byte(userCredentials)
	_, err := loginUser(credentials)

	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func TestRefreshTokenAPI(t *testing.T) {
	// Login the user with given credentials and aquired the token
	credentials := []byte(userCredentials)

	authResponse, err := loginUser(credentials)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Call refresh token API with the given refresh token and check API gives a success response or not
	payload := []byte(fmt.Sprintf(`{"refresh_token": "%s"}`, authResponse.Data.Refresh))

	response := getResponseRecorder(http.MethodPost, "/v1/refresh-token/", "", payload)

	if response.Code != http.StatusOK {
		t.Errorf("Response: %v", response.Body.String())
	}
}
