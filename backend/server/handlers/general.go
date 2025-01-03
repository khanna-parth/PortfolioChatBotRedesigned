package handlers

import "net/http"

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// SendResponse(w, map[string]string{"Status": "Alive"})
}