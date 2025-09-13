package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/helper"
	"strings"
)

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// SendResponse(w, map[string]string{"Status": "Alive"})
}

func SuggestionHandler(w http.ResponseWriter, r *http.Request, config map[string]string) {
	uploads, exists := config["UPLOADS_PATH"]
	suggestionFile, suggestionExists := config["SUGGESTIONS_PROVIDER"]
	if !exists || !suggestionExists {
		http.Error(w, "Could not formulize suggestions path", http.StatusInternalServerError)
		return
	}

	file := filepath.Join(filepath.Join(uploads, "STATIC"), suggestionFile)
	data, err := os.ReadFile(file)
	if err != nil {
		http.Error(w, "Could not retrieve suggestions content", http.StatusInternalServerError)
		return
	}

	type suggestionsResponse struct {
		Suggestions []string `json:"suggestions"`
	}

	content := strings.Split(string(data), "\n")

	suggestions := &suggestionsResponse{
		Suggestions: content,
	}
	log.Printf("Sending suggestions response with %d suggestions\n", len(content))
	helper.SendResponse(w, suggestions)
}