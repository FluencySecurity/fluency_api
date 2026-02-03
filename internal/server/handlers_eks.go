package server

import (
	"encoding/json"
	"net/http"
)

// EKSTestAssumedRoleRequest is the JSON body for testing an assumed role.
type EKSTestAssumedRoleRequest struct {
	RoleARN     string `json:"roleArn"`
	ExternalID  string `json:"externalId,omitempty"`
}

// EKSAddAssumedRoleRequest is the JSON body for adding an assumed role.
type EKSAddAssumedRoleRequest struct {
	Name        string `json:"name"`
	RoleARN     string `json:"roleArn"`
	ExternalID  string `json:"externalId,omitempty"`
}

func (s *Server) eksGetPodRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	role, arn, err := s.API.GetPodRole()
	if err != nil {
		s.Log.Error("get pod role failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"role": role, "arn": arn})
}

func (s *Server) eksTestAssumedRole(w http.ResponseWriter, r *http.Request) {
	var req EKSTestAssumedRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.RoleARN == "" {
		writeJSONError(w, http.StatusBadRequest, "roleArn is required")
		return
	}
	err := s.API.TestAssumedRole(req.RoleARN, req.ExternalID)
	if err != nil {
		s.Log.Error("test assumed role failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) eksAddAssumedRole(w http.ResponseWriter, r *http.Request) {
	var req EKSAddAssumedRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.RoleARN == "" {
		writeJSONError(w, http.StatusBadRequest, "name and roleArn are required")
		return
	}
	id, err := s.API.AddAssumedRole(req.Name, req.RoleARN, req.ExternalID)
	if err != nil {
		s.Log.Error("add assumed role failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (s *Server) eksDeleteAssumedRole(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter id is required")
		return
	}
	err := s.API.DeleteAssumedRole(id)
	if err != nil {
		s.Log.Error("delete assumed role failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) eksListAssumedRoles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	roles, err := s.API.ListAssumedRole()
	if err != nil {
		s.Log.Error("list assumed roles failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": roles})
}
