package tests

import (
	"fmt"
	"net/http"
	"testing"
)

func TestSignUpAPI(t *testing.T) {
	payload := []byte(`{"email": "testing-user@example.com", "password": "12345678Aa@", "is_end_user": true}`)
	_, err := createTestingUser(payload)

	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func TestLoginAPI(t *testing.T) {
	credentials := []byte(`{"email": "testing-user@example.com", "password": "12345678Aa@"}`)
	_, err := loginUser(credentials)

	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func TestRefreshTokenAPI(t *testing.T) {
	// Login the user with given credentials and aquired the token
	credentials := []byte(`{"email": "testing-user@example.com", "password": "12345678Aa@"}`)

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
