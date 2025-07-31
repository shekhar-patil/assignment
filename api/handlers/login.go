package handlers

import (
	"encoding/json"
	"net/http"
	"shekhar-patil/assignment/api/storage"
)

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
		json.NewEncoder(w).Encode(map[string]string{"token": storage.ValidToken})
	} else {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
	}
}
