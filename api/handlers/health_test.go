package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCHeck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	HealthCheck(rec, req)

	// Validate status code
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", rec.Code)
	}

	// Validate response body
	var body map[string]string
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	expected := "ok"
	if body["status"] != expected {
		t.Errorf("expected status %q, got %q", expected, body["status"])
	}
}
