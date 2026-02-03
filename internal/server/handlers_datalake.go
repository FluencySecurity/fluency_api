package server

import (
	"encoding/json"
	"net/http"
)

// DatalakeAddRequest is the JSON body for adding a datalake.
type DatalakeAddRequest struct {
	Name         string `json:"name"`
	Managed      bool   `json:"managed"`
	Integration  string `json:"integration,omitempty"`
}

// DatalakeIndexAddRequest is the JSON body for adding a datalake index.
type DatalakeIndexAddRequest struct {
	Datalake string `json:"datalake"`
	Index    string `json:"index"`
	Schema   string `json:"schema,omitempty"`
}

func (s *Server) datalakeList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	entries, err := s.API.ListDatalakes()
	if err != nil {
		s.Log.Error("list datalakes failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": entries})
}

func (s *Server) datalakeAdd(w http.ResponseWriter, r *http.Request) {
	var req DatalakeAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "name is required")
		return
	}
	err := s.API.AddDatalake(req.Name, req.Managed, req.Integration)
	if err != nil {
		s.Log.Error("add datalake failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) datalakeListIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	lake := r.URL.Query().Get("lake")
	if lake == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter lake is required")
		return
	}
	entries, err := s.API.ListDatalakeIndex(lake)
	if err != nil {
		s.Log.Error("list datalake index failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": entries})
}

func (s *Server) datalakeAddIndex(w http.ResponseWriter, r *http.Request) {
	var req DatalakeIndexAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Datalake == "" {
		req.Datalake = "managed"
	}
	if req.Index == "" {
		writeJSONError(w, http.StatusBadRequest, "index is required")
		return
	}
	if req.Schema == "" {
		req.Schema = "fluency default"
	}
	err := s.API.AddDatalakeIndex(req.Datalake, req.Index, req.Schema)
	if err != nil {
		s.Log.Error("add datalake index failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) datalakeDeleteIndex(w http.ResponseWriter, r *http.Request) {
	lake := r.URL.Query().Get("lake")
	index := r.URL.Query().Get("index")
	if lake == "" || index == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameters lake and index are required")
		return
	}
	err := s.API.DeleteDatalakeIndex(lake, index)
	if err != nil {
		s.Log.Error("delete datalake index failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
