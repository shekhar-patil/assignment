package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shekhar-patil/assignment/api/models"
	"shekhar-patil/assignment/api/storage"
	"testing"
	"time"
)

func TestReportHandler(t *testing.T) {
	// Prepare test data
	initiateAt := time.Date(2025, 6, 12, 11, 20, 33, 0, time.UTC)
	finishAt := time.Date(2025, 6, 12, 11, 24, 46, 0, time.UTC)

	testRecords := []models.PipelineRecord{
		{
			Org:        "org1",
			App:        "app1",
			Pipeline:   "pipelineA",
			JobName:    "jobA",
			Status:     "success",
			Duration:   finishAt.Sub(initiateAt),
			InitiateAt: initiateAt,
			FinishAt:   finishAt,
		},
		{
			Org:        "org2",
			App:        "app2",
			Pipeline:   "pipelineB",
			JobName:    "jobB",
			Status:     "failure",
			Duration:   finishAt.Sub(initiateAt),
			InitiateAt: initiateAt,
			FinishAt:   finishAt,
		},
	}

	storage.Mu.Lock()
	storage.PipelineData = testRecords
	storage.Mu.Unlock()

	tests := []struct {
		name          string
		queryParams   string
		expectedCount int
	}{
		{
			name:          "no filters returns all",
			queryParams:   "",
			expectedCount: 2,
		},
		{
			name:          "filter by org",
			queryParams:   "?org=org1",
			expectedCount: 1,
		},
		{
			name:          "filter by app",
			queryParams:   "?app=app2",
			expectedCount: 1,
		},
		{
			name:          "filter by org and app",
			queryParams:   "?org=org1&app=app1",
			expectedCount: 1,
		},
		{
			name:          "filter with no match",
			queryParams:   "?org=orgX",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/report"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			ReportHandler(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("expected 200 OK, got %d", rec.Code)
			}

			var result []models.PipelineRecord
			err := json.NewDecoder(rec.Body).Decode(&result)
			if err != nil {
				t.Fatalf("could not decode response: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d records, got %d", tt.expectedCount, len(result))
			}
		})
	}
}
