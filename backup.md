```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type PipelineRequest struct {
	RunSpecID  string `json:"runSpecId"`
	FinishAt   string `json:"finishAt"`
	JobName    string `json:"jobName"`
	Status     string `json:"status"`
	InitiateAt string `json:"initiateAt"`
}

// Parsed version with org/app/pipeline split
type PipelineRecord struct {
	Org        string    `json:"org"`
	App        string    `json:"app"`
	Pipeline   string    `json:"pipeline"`
	JobName    string    `json:"jobName"`
	Status     string    `json:"status"`
	InitiateAt time.Time `json:"initiateAt"`
	FinishAt   time.Time `json:"finishAt"`
}

var (
	pipelineData []PipelineRecord
	mu           sync.Mutex
)

var validToken = "s3cr3t-token"

// curl -X GET http://localhost:8080/health
func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// --- Login Handler ---
//
//	curl -X POST http://localhost:8080/login \
//	  -H "Content-Type: application/json" \
//	  -d '{"username":"admin","password":"Admin"}'
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if creds.Username == "admin" && creds.Password == "Admin" {
		json.NewEncoder(w).Encode(map[string]string{"token": validToken})
	} else {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != validToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// // POST /pipeline
//
//	curl -X POST http://localhost:8080/pipeline \
//	  -H "Content-Type: application/json" \
//	  -H "Authorization: Bearer s3cr3t-token" \
//	  -d '{"runSpecId":"org1-app1-pipelineA","finishAt":"2025-06-12T11:24:46+00:00","jobName":"org1-app1","status":"success","initiateAt":"2025-06-12T11:20:33+00:00"}'
func PipelineHandler(w http.ResponseWriter, r *http.Request) {
	var req PipelineRequest
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

	record := PipelineRecord{
		Org:        parts[0],
		App:        parts[1],
		Pipeline:   parts[2],
		JobName:    req.JobName,
		Status:     req.Status,
		InitiateAt: initiateAt,
		FinishAt:   finishAt,
	}

	mu.Lock()
	pipelineData = append(pipelineData, record)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pipeline data saved"})
}

//	curl -X GET http://localhost:8080/report \
//			-H "Authorization: Bearer s3cr3t-token"
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pipelineData)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.Handle("/pipeline", AuthMiddleware(http.HandlerFunc(PipelineHandler))).Methods("POST")
	r.Handle("/report", AuthMiddleware(http.HandlerFunc(ReportHandler))).Methods("GET")
	r.HandleFunc("/login", LoginHandler).Methods("POST")

	// Start server
	port := ":8080"
	fmt.Println("Server running on http://localhost" + port)
	http.ListenAndServe(port, r)
}
```