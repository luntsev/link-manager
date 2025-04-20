package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func Json(w http.ResponseWriter, content any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(content)
	if err != nil {
		log.Printf("Error when responding to authorization: %s", err.Error())
	}
}
