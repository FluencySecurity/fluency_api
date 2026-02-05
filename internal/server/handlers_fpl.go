package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SecurityDo/fluency_api/model"
)

func (s *Server) fplRunReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req model.RunFPLV2Report
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := s.API.RunReport(&req)
	if err != nil {
		s.Log.Error("fpl run_report failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"id": id})
}

func (s *Server) fplGetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter id is required")
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "id must be a non-negative integer")
		return
	}
	resp, err := s.API.GetTaskByID(uint(id))
	if err != nil {
		s.Log.Error("fpl get_task failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) fplGetResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter id is required")
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "id must be a non-negative integer")
		return
	}
	resp, err := s.API.GetResultsByID(uint(id))
	if err != nil {
		s.Log.Error("fpl get_results failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
