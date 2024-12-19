package handlers

import (
	"encoding/json"
	"hw3/internal/nodes"
	"net/http"
)

func MakeGetHandler(v *nodes.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := v.Get()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
