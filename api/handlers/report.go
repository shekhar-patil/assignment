package handlers

import (
	"encoding/json"
	"net/http"
	"shekhar-patil/assignment/api/models"
	"shekhar-patil/assignment/api/storage"
)

//	curl -X GET http://localhost:8080/report \
//			-H "Authorization: Bearer s3cr3t-token"
//			-H "Content-Type: application/json" \
//			-d '{"org":"org1", "app": "app1"}
func ReportHandler(w http.ResponseWriter, r *http.Request) {
	storage.Mu.RLock()
	defer storage.Mu.RUnlock()

	org := r.URL.Query().Get("org")
	app := r.URL.Query().Get("app")
	pipeline := r.URL.Query().Get("pipeline")

	var filtered []models.PipelineRecord
	for _, pr := range storage.PipelineData {
		if (org == "" || pr.Org == org) &&
			(app == "" || pr.App == app) &&
			(pipeline == "" || pr.Pipeline == pipeline) {
			filtered = append(filtered, pr)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}
