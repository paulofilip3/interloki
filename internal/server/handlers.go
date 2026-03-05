package server

import (
	"encoding/json"
	"net/http"
)

// HandleStatus writes the current server status as JSON.
func (m *ClientManager) HandleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m.Status())
}
