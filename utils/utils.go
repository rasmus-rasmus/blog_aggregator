package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
	if statusCode > 499 {
		log.Printf("Responding with 5XX error: %s", errorMessage)
	}
	responseBody := struct {
		Msg string `json:"error"`
	}{errorMessage}
	RespondWithJSON(w, statusCode, responseBody)
}
