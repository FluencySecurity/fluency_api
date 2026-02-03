package server

import (
	"encoding/json"
	"net/http"
	"time"
)

// AuditSearchRequest is the JSON body for audit search.
type AuditSearchRequest struct {
	Query string `json:"query"`
	From  int64  `json:"from,omitempty"`
	To    int64  `json:"to,omitempty"`
}

func (s *Server) auditSearch(w http.ResponseWriter, r *http.Request) {
	var req AuditSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	from, to := req.From, req.To
	if from == 0 && to == 0 {
		now := time.Now().UnixMilli()
		to = now
		from = now - int64(time.Hour/time.Millisecond)
	}
	resp, err := s.API.AuditSearch(req.Query, from, to)
	if err != nil {
		s.Log.Error("audit search failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// auditDbStatus handles POST /api/ds/db_status with body {}.
func (s *Server) auditDbStatus(w http.ResponseWriter, r *http.Request) {
	resp, err := s.API.DbStatus()
	if err != nil {
		s.Log.Error("db_status failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
