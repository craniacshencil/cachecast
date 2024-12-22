package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body != nil {
		fmt.Println("Empty body for request")
	}
	return json.NewDecoder(r.Body).Decode(v)
}
