package controller

import (
	"encoding/json"
	"net/http"
)

func getLoginFromToken(r *http.Request) string {
	return r.URL.Query().Get("login") // login из query
}

func writeError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code": code,
			"text": err.Error(),
		},
	})
}

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
