package health

import (
	"encoding/json"
	"net/http"
)

// Handler is a HTTP Server Handler implementation
type Handler struct {
	CompositeChecker
}

// NewHandler returns a new Handler
func NewHandler() Handler {
	return Handler{}
}

// ServeHTTP returns a json encoded Health
// set the status to http.StatusServiceUnavailable if the check is down
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	health := h.CompositeChecker.Check(r.Context())

	if health.IsDown() {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(health)
}
