package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shekhar-patil/assignment/api/storage"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]string
		expectedCode   int
		expectedResult string // part of body to check
	}{
		{
			name: "valid credentials",
			input: map[string]string{
				"username": "admin",
				"password": "Admin",
			},
			expectedCode:   http.StatusOK,
			expectedResult: storage.ValidToken,
		},
		{
			name: "invalid credentials",
			input: map[string]string{
				"username": "admin",
				"password": "wrong",
			},
			expectedCode:   http.StatusUnauthorized,
			expectedResult: "Authentication failed",
		},
		{
			name:           "invalid JSON",
			input:          nil, // will send invalid raw data
			expectedCode:   http.StatusBadRequest,
			expectedResult: "Invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.input != nil {
				body, _ := json.Marshal(tt.input)
				req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			} else {
				req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{invalid_json}"))
			}
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			LoginHandler(rec, req)

			if rec.Code != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, rec.Code)
			}

			body := rec.Body.String()
			if tt.expectedCode == http.StatusOK {
				var resp map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				if err != nil {
					t.Fatalf("could not decode response: %v", err)
				}
				if resp["token"] != tt.expectedResult {
					t.Errorf("expected token %q, got %q", tt.expectedResult, resp["token"])
				}
			} else if !bytes.Contains([]byte(body), []byte(tt.expectedResult)) {
				t.Errorf("expected response to contain %q, got %q", tt.expectedResult, body)
			}
		})
	}
}
