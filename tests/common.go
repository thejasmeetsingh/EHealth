package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/thejasmeetsingh/EHealth/config"
)

type AuthResponse struct {
	Message string `json:"message"`
	Data    struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	} `json:"data"`
}

// A common function for calling APIs and validating test cases based on API response
func getResponseRecorder(method, endpoint, accessToken string, payload []byte) *httptest.ResponseRecorder {
	router := config.GetRouter(true)

	req := httptest.NewRequest(method, endpoint, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	if accessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder
}

func parseAuthResponse(response *httptest.ResponseRecorder) (AuthResponse, error) {
	var authResponse AuthResponse

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&authResponse); err != nil {
		return AuthResponse{}, fmt.Errorf("error while parsing the response: %v", err)
	}

	return authResponse, nil
}

func createTestingUser(payload []byte) (AuthResponse, error) {
	response := getResponseRecorder(http.MethodPost, "/v1/signup/", "", payload)

	if response.Code != http.StatusCreated {
		return AuthResponse{}, fmt.Errorf("response: %v", response.Body.String())
	}

	return parseAuthResponse(response)
}

func loginUser(credentials []byte) (AuthResponse, error) {
	response := getResponseRecorder(http.MethodPost, "/v1/login/", "", credentials)

	if response.Code != http.StatusOK {
		return AuthResponse{}, fmt.Errorf("response: %v", response.Body.String())
	}

	return parseAuthResponse(response)
}
