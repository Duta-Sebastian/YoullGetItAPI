package util

import (
	"encoding/json"
	"log"
	"net/http"
)

// DecodeRequestBody decodeRequestBody simplifies JSON decoding from request body
func DecodeRequestBody(r *http.Request, target interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	return decoder.Decode(target)
}

// RespondWithJSON Helper functions for standardized responses
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"message": message})
}
