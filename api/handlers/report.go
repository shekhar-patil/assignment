package handlers

import (
	"encoding/json"
	"net/http"
	"shekhar-patil/assignment/api/storage"
)

//	curl -X GET http://localhost:8080/report \
//			-H "Authorization: Bearer s3cr3t-token"
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	storage.Mu.Lock()
	defer storage.Mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storage.PipelineData)
}
