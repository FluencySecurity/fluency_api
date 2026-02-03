package server

import (
	"encoding/json"
	"net/http"

	"github.com/SecurityDo/fluency_api/model"
)

func (s *Server) integrationList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	entries, err := s.API.ListIntegration()
	if err != nil {
		s.Log.Error("list integrations failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": entries})
}

// IntegrationAddRequest matches the CLI: integration type, name, description, config (raw JSON object), secret (raw JSON object).
type IntegrationAddRequest struct {
	Integration string          `json:"integration"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Config      json.RawMessage `json:"config,omitempty"`
	Secret      json.RawMessage `json:"secret,omitempty"`
}

func (s *Server) integrationAdd(w http.ResponseWriter, r *http.Request) {
	var req IntegrationAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Integration == "" || req.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "integration and name are required")
		return
	}
	if req.Config == nil {
		req.Config = []byte("{}")
	}
	if req.Secret == nil {
		req.Secret = []byte("{}")
	}
	entry := &model.Integration{
		Integration: req.Integration,
		Name:        req.Name,
		Description: req.Description,
		Config:      req.Config,
		Secret:      req.Secret,
	}
	id, err := s.API.AddIntegration(entry)
	if err != nil {
		s.Log.Error("add integration failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (s *Server) integrationDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter id is required")
		return
	}
	err := s.API.DeleteIntegration(id)
	if err != nil {
		s.Log.Error("delete integration failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
