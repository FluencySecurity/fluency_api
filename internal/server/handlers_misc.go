package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) statusCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	err := s.API.CheckStatus()
	if err != nil {
		s.Log.Error("status check failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ImportProcessorRequest is the JSON body for importing processors from a repo.
type ImportProcessorRequest struct {
	RepoName string `json:"repoName,omitempty"`
	Type     string `json:"type,omitempty"`
}

// ImportApplicationRequest is the JSON body for importing application templates from a repo.
type ImportApplicationRequest struct {
	RepoName string `json:"repoName,omitempty"`
}

func (s *Server) importProcessor(w http.ResponseWriter, r *http.Request) {
	var req ImportProcessorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Type == "" {
		req.Type = "fpl_processor"
	}
	err := s.API.ImportProcessor(req.Type, req.RepoName)
	if err != nil {
		s.Log.Error("import processor failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) importApplication(w http.ResponseWriter, r *http.Request) {
	var req ImportApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err := s.API.ImportAppTemplate(req.RepoName)
	if err != nil {
		s.Log.Error("import application failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
