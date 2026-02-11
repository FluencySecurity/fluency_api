package server

import (
	"encoding/json"
	"net/http"
)

// CollectorStatusRequest is the JSON body for collector status (collector only; cargs are always el9).
type CollectorStatusRequest struct {
	Collector string `json:"collector"`
}

func (s *Server) collectorList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	entries, err := s.API.CollectorList()
	if err != nil {
		s.Log.Error("collector list failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": entries})
}

func (s *Server) collectorStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req CollectorStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Collector == "" {
		writeJSONError(w, http.StatusBadRequest, "collector is required")
		return
	}
	cargs := map[string]interface{}{"status": "el9"}
	out, err := s.API.CollectorStatus(req.Collector, cargs)
	if err != nil {
		s.Log.Error("collector status failed", "collector", req.Collector, "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, out)
}
