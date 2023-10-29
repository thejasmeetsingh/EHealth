package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thejasmeetsingh/EHealth/config"
)

func TestSignUpAPI(t *testing.T) {
	router := config.GetRouter(true)

	payload := []byte(`{"email": "testing-user@example", "password": "12345678Aa@", "is_end_user": true}`)

	req := httptest.NewRequest(http.MethodPost, "/v1/signup/", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Response: %v", w.Body.String())
	}
}
