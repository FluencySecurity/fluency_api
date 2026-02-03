package server

import (
	"encoding/json"
	"net/http"
)

// AuthAddUserRequest is the JSON body for adding a user.
type AuthAddUserRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Role        string `json:"role"`
	Org         string `json:"org,omitempty"`
}

// AuthAddTokenRequest is the JSON body for adding a token.
type AuthAddTokenRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	Role        string `json:"role"`
}

func (s *Server) authAddUser(w http.ResponseWriter, r *http.Request) {
	var req AuthAddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Role == "" {
		writeJSONError(w, http.StatusBadRequest, "name and role are required")
		return
	}
	if req.Org == "" {
		req.Org = "fluency"
	}
	err := s.API.AddUser(req.Name, req.DisplayName, req.Role, req.Org)
	if err != nil {
		s.Log.Error("add user failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) authDeleteUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter name is required")
		return
	}
	err := s.API.DeleteUser(name)
	if err != nil {
		s.Log.Error("delete user failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) authListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	users, err := s.API.ListUser()
	if err != nil {
		s.Log.Error("list users failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": users})
}

func (s *Server) authAddToken(w http.ResponseWriter, r *http.Request) {
	var req AuthAddTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Role == "" {
		writeJSONError(w, http.StatusBadRequest, "name and role are required")
		return
	}
	token, err := s.API.AddToken(req.Name, req.DisplayName, req.Role)
	if err != nil {
		s.Log.Error("add token failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (s *Server) authDeleteToken(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter name is required")
		return
	}
	err := s.API.DeleteToken(name)
	if err != nil {
		s.Log.Error("delete token failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) authListTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	tokens, err := s.API.ListToken()
	if err != nil {
		s.Log.Error("list tokens failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": tokens})
}
