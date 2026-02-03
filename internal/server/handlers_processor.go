package server

import (
	"encoding/json"
	"net/http"
)

// ProcessorAddRequest is the JSON body for adding a processor.
type ProcessorAddRequest struct {
	Name        string `json:"name"`
	Content     string `json:"content"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

func (s *Server) processorList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	entries, err := s.API.ListProcessor()
	if err != nil {
		s.Log.Error("list processors failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": entries})
}

func (s *Server) processorAdd(w http.ResponseWriter, r *http.Request) {
	var req ProcessorAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Content == "" {
		writeJSONError(w, http.StatusBadRequest, "name and content are required")
		return
	}
	if req.Type == "" {
		req.Type = "fpl_processor"
	}
	err := s.API.AddProcessor(req.Name, req.Content, req.Type, req.Description)
	if err != nil {
		s.Log.Error("add processor failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) processorDelete(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter name is required")
		return
	}
	err := s.API.DeleteProcessor(name)
	if err != nil {
		s.Log.Error("delete processor failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
