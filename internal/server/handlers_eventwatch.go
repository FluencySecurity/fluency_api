package server

import (
	"encoding/json"
	"net/http"
	"time"
)

// EventWatchSearchRequest is the JSON body for summary and timeline search.
type EventWatchSearchRequest struct {
	Query string `json:"query"`
	From  int64  `json:"from,omitempty"`
	To    int64  `json:"to,omitempty"`
}

// EventWatchRuleSearchRequest is the JSON body for rule search.
type EventWatchRuleSearchRequest struct {
	Query string `json:"query"`
}

func (s *Server) eventwatchSummarySearch(w http.ResponseWriter, r *http.Request) {
	var req EventWatchSearchRequest
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
	resp, err := s.API.SummarySearch(req.Query, from, to)
	if err != nil {
		s.Log.Error("summary search failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) eventwatchTimelineSearch(w http.ResponseWriter, r *http.Request) {
	var req EventWatchSearchRequest
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
	resp, err := s.API.TimelineSearch(req.Query, from, to)
	if err != nil {
		s.Log.Error("timeline search failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) eventwatchRuleSearch(w http.ResponseWriter, r *http.Request) {
	var req EventWatchRuleSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	resp, err := s.API.RuleSearch(req.Query)
	if err != nil {
		s.Log.Error("rule search failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
