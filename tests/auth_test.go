package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestSignUpAPI(t *testing.T) {
	payload := []byte(`{"email": "testing-user@example.com", "password": "12345678Aa@", "is_end_user": true}`)

	response := getResponseRecorder(http.MethodPost, "/v1/signup/", payload)

	if response.Code != http.StatusCreated {
		t.Errorf("Response: %v", response.Body.String())
	}
}

func TestLoginAPI(t *testing.T) {
	payload := []byte(`{"email": "testing-user@example.com", "password": "12345678Aa@"}`)

	response := getResponseRecorder(http.MethodPost, "/v1/login/", payload)

	if response.Code != http.StatusOK {
		t.Errorf("Response: %v", response.Body.String())
	}
}

func TestRefreshTokenAPI(t *testing.T) {
	// Login the user with given credentials and aquired the token
	credentials := []byte(`{"email": "testing-user@example.com", "password": "12345678Aa@"}`)

	loginResponse := getResponseRecorder(http.MethodPost, "/v1/login/", credentials)

	if loginResponse.Code != http.StatusOK {
		t.Errorf("Response: %v", loginResponse.Body.String())
	}

	// Parse the API response
	type Response struct {
		Message string `json:"message"`
		Data    struct {
			Access  string `json:"access"`
			Refresh string `json:"refresh"`
		} `json:"data"`
	}

	var responseBody Response

	decoder := json.NewDecoder(loginResponse.Body)
	if err := decoder.Decode(&responseBody); err != nil {
		t.Errorf("Error while parsing the response: %v", err)
		return
	}

	// Call refresh token API with the given refresh token and check API gives a success response or not
	payload := []byte(fmt.Sprintf(`{"refresh_token": "%s"}`, responseBody.Data.Refresh))

	response := getResponseRecorder(http.MethodPost, "/v1/refresh-token/", payload)

	if response.Code != http.StatusOK {
		t.Errorf("Response: %v", response.Body.String())
	}
}
