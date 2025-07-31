package handlers

import (
	"encoding/json"
	"net/http"
	"shekhar-patil/assignment/api/models"
	"shekhar-patil/assignment/api/storage"
	"strings"
	"time"
)

// // POST /pipeline
//
//	curl -X POST http://localhost:8080/pipeline \
//	  -H "Content-Type: application/json" \
//	  -H "Authorization: Bearer s3cr3t-token" \
//	  -d '{"runSpecId":"org1-app1-pipelineA","finishAt":"2025-06-12T11:24:46+00:00","jobName":"org1-app1","status":"success","initiateAt":"2025-06-12T11:20:33+00:00"}'
func PipelineHandler(w http.ResponseWriter, r *http.Request) {
	var req models.PipelineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse runSpecId: expected format org-app-pipeline
	parts := strings.Split(req.RunSpecID, "-")
	if len(parts) != 3 {
		http.Error(w, "Invalid runSpecId format", http.StatusBadRequest)
		return
	}

	initiateAt, err1 := time.Parse(time.RFC3339, req.InitiateAt)
	finishAt, err2 := time.Parse(time.RFC3339, req.FinishAt)
	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid datetime format", http.StatusBadRequest)
		return
	}

	record := models.PipelineRecord{
		Org:        parts[0],
		App:        parts[1],
		Pipeline:   parts[2],
		JobName:    req.JobName,
		Status:     req.Status,
		Duration:   finishAt.Sub(initiateAt),
		InitiateAt: initiateAt,
		FinishAt:   finishAt,
	}

	storage.Mu.Lock()
	storage.PipelineData = append(storage.PipelineData, record)
	storage.Mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pipeline data saved"})
}
