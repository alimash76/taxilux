package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}

func HealthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, healthResponse{Status: "degraded", DB: "unavailable"})
			return
		}
		writeJSON(w, http.StatusOK, healthResponse{Status: "ok", DB: "connected"})
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
