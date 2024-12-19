package handlers

import (
	"encoding/json"
	"hw3/internal/nodes"
	"io"
	"net/http"
)

func MakeUpdateHandler(v *nodes.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid Content-Type, expected application/json", http.StatusUnsupportedMediaType)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var data map[string]string
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		updates := make([]nodes.Update, 0, len(data))
		for key, value := range data {
			updates = append(updates, nodes.Update{
				Key:   nodes.Key(key),
				Value: nodes.Value(value),
			})
		}

		v.Broadcast(updates)

		w.WriteHeader(http.StatusOK)
	}
}
