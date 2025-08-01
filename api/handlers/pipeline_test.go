package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shekhar-patil/assignment/api/models"
	"shekhar-patil/assignment/api/storage"
	"testing"
)

func TestPipelineHandler(t *testing.T) {
	// Clear previous records
	storage.Mu.Lock()
	storage.PipelineData = []models.PipelineRecord{}
	storage.Mu.Unlock()

	tests := []struct {
		name           string
		payload        string
		expectedCode   int
		expectedMsg    string
		expectedRecord bool
	}{
		{
			name: "valid payload",
			payload: `{
				"runSpecId":"org1-app1-pipelineA",
				"jobName":"org1-app1",
				"status":"success",
				"initiateAt":"2025-06-12T11:20:33+00:00",
				"finishAt":"2025-06-12T11:24:46+00:00"
			}`,
			expectedCode:   http.StatusCreated,
			expectedMsg:    "Pipeline data saved",
			expectedRecord: true,
		},
		{
			name: "invalid JSON",
			payload: `{
				"runSpecId":"org1-app1-pipelineA",
				"jobName": "test"`, // malformed
			expectedCode:   http.StatusBadRequest,
			expectedMsg:    "Invalid JSON",
			expectedRecord: false,
		},
		{
			name: "invalid runSpecId",
			payload: `{
				"runSpecId":"org1-app1",
				"jobName":"org1-app1",
				"status":"success",
				"initiateAt":"2025-06-12T11:20:33+00:00",
				"finishAt":"2025-06-12T11:24:46+00:00"
			}`,
			expectedCode:   http.StatusBadRequest,
			expectedMsg:    "Invalid runSpecId format",
			expectedRecord: false,
		},
		{
			name: "invalid datetime format",
			payload: `{
				"runSpecId":"org1-app1-pipelineA",
				"jobName":"org1-app1",
				"status":"success",
				"initiateAt":"invalid-date",
				"finishAt":"2025-06-12T11:24:46+00:00"
			}`,
			expectedCode:   http.StatusBadRequest,
			expectedMsg:    "Invalid datetime format",
			expectedRecord: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/pipeline", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			PipelineHandler(rec, req)

			// Check status code
			if rec.Code != tt.expectedCode {
				t.Errorf("expected code %d, got %d", tt.expectedCode, rec.Code)
			}

			// Check response body
			var body map[string]string
			_ = json.Unmarshal(rec.Body.Bytes(), &body)
			if msg, ok := body["message"]; tt.expectedCode == http.StatusCreated && (!ok || msg != tt.expectedMsg) {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, msg)
			}

			// Check pipeline storage
			if tt.expectedRecord {
				storage.Mu.Lock()
				if len(storage.PipelineData) == 0 {
					t.Error("expected record to be saved but storage is empty")
				}
				storage.Mu.Unlock()
			}
		})
	}
}
