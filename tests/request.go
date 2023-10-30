package tests

import (
	"bytes"
	"net/http/httptest"

	"github.com/thejasmeetsingh/EHealth/config"
)

// A common function for calling APIs and validating test cases based on API response
func getResponseRecorder(method string, endpoint string, payload []byte) *httptest.ResponseRecorder {
	router := config.GetRouter(true)

	req := httptest.NewRequest(method, endpoint, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder
}
