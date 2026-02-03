package server

import (
	"encoding/json"
	"net/http"

	"github.com/SecurityDo/fluency_api/model"
)

func (s *Server) applicationListTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	templates, err := s.API.ListAppTemplates()
	if err != nil {
		s.Log.Error("list app templates failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": templates})
}

// ApplicationInstallRequest is the JSON body for installing an app instance.
type ApplicationInstallRequest struct {
	App             string                 `json:"app"`
	Instance        string                 `json:"instance"`
	DisplayName     string                 `json:"displayName,omitempty"`
	Config          map[string]string     `json:"config,omitempty"`
	Secret          map[string]string     `json:"secret,omitempty"`
}

func (s *Server) applicationInstall(w http.ResponseWriter, r *http.Request) {
	var req ApplicationInstallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.App == "" || req.Instance == "" {
		writeJSONError(w, http.StatusBadRequest, "app and instance are required")
		return
	}
	var paras []*model.InputParameter
	for k, v := range req.Config {
		paras = append(paras, &model.InputParameter{Name: k, Value: v, Sensitive: false})
	}
	for k, v := range req.Secret {
		paras = append(paras, &model.InputParameter{Name: k, Value: v, Sensitive: true})
	}
	err := s.API.InstallAppInstance(req.App, req.Instance, req.DisplayName, paras)
	if err != nil {
		s.Log.Error("install app instance failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) applicationGetInstance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	app := r.URL.Query().Get("app")
	instance := r.URL.Query().Get("instance")
	if app == "" || instance == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameters app and instance are required")
		return
	}
	resp, err := s.API.GetAppInstance(app, instance)
	if err != nil {
		s.Log.Error("get app instance failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) applicationUninstall(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("app")
	instance := r.URL.Query().Get("instance")
	if app == "" || instance == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameters app and instance are required")
		return
	}
	err := s.API.UnInstallAppInstance(app, instance)
	if err != nil {
		s.Log.Error("uninstall app instance failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ApplicationAddTemplateRequest is the JSON body for adding a template.
type ApplicationAddTemplateRequest struct {
	Content string `json:"content"`
}

func (s *Server) applicationAddTemplate(w http.ResponseWriter, r *http.Request) {
	var req ApplicationAddTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Content == "" {
		writeJSONError(w, http.StatusBadRequest, "content is required")
		return
	}
	id, err := s.API.AddTemplate(req.Content)
	if err != nil {
		s.Log.Error("add template failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (s *Server) applicationDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("app")
	if app == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter app is required")
		return
	}
	err := s.API.DeleteTemplate(app)
	if err != nil {
		s.Log.Error("delete template failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// ApplicationUpdateTemplateRequest is the JSON body for updating a template.
type ApplicationUpdateTemplateRequest struct {
	Content string `json:"content"`
}

func (s *Server) applicationUpdateTemplate(w http.ResponseWriter, r *http.Request) {
	app := r.URL.Query().Get("app")
	if app == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter app is required")
		return
	}
	var req ApplicationUpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Content == "" {
		writeJSONError(w, http.StatusBadRequest, "content is required")
		return
	}
	err := s.API.UpdateTemplate(app, req.Content)
	if err != nil {
		s.Log.Error("update template failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
