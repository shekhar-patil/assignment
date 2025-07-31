package handlers

import (
	"encoding/json"
	"net/http"
)

// curl -X GET http://localhost:8080/health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
