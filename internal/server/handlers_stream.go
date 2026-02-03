package server

import (
	"encoding/json"
	"net/http"

	"github.com/SecurityDo/fluency_api/model"
)

// StreamConnectSinkRequest connects a router to a sink.
type StreamConnectSinkRequest struct {
	RouterID string `json:"routerId"`
	SinkID   string `json:"sinkId"`
}

// StreamConnectRouterRequest connects a source to a router.
type StreamConnectRouterRequest struct {
	SourceID string `json:"sourceId"`
	RouterID string `json:"routerId"`
}

// StreamAddRouterRequest adds a simple router.
type StreamAddRouterRequest struct {
	Processor string `json:"processor"`
	RouterName string `json:"routerName"`
}

func (s *Server) streamConnectSink(w http.ResponseWriter, r *http.Request) {
	var req StreamConnectSinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.RouterID == "" || req.SinkID == "" {
		writeJSONError(w, http.StatusBadRequest, "routerId and sinkId are required")
		return
	}
	err := s.API.SetRouterSink(req.RouterID, req.SinkID)
	if err != nil {
		s.Log.Error("connect sink failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) streamConnectRouter(w http.ResponseWriter, r *http.Request) {
	var req StreamConnectRouterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.SourceID == "" || req.RouterID == "" {
		writeJSONError(w, http.StatusBadRequest, "sourceId and routerId are required")
		return
	}
	err := s.API.SetSourceRouter(req.SourceID, req.RouterID)
	if err != nil {
		s.Log.Error("connect router failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) streamAddRouter(w http.ResponseWriter, r *http.Request) {
	var req StreamAddRouterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Processor == "" {
		writeJSONError(w, http.StatusBadRequest, "processor is required")
		return
	}
	id, err := s.API.AddSimpleRouter(req.Processor, req.RouterName)
	if err != nil {
		s.Log.Error("add router failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (s *Server) streamAddSource(w http.ResponseWriter, r *http.Request) {
	var source model.DataSourceConfig
	if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if source.Type == "" || source.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "type and name are required")
		return
	}
	if source.Format == "" {
		source.Format = "json"
	}
	resp, err := s.API.AddDataSource(&source)
	if err != nil {
		s.Log.Error("add data source failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) streamDeleteSource(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter id is required")
		return
	}
	err := s.API.DeleteDataSource(id)
	if err != nil {
		s.Log.Error("delete data source failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) streamListSources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	entries, err := s.API.ListDataSource()
	if err != nil {
		s.Log.Error("list data sources failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": entries})
}

func (s *Server) streamAddSink(w http.ResponseWriter, r *http.Request) {
	var sink model.DataSinkConfig
	if err := json.NewDecoder(r.Body).Decode(&sink); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if sink.Type == "" || sink.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "type and name are required")
		return
	}
	resp, err := s.API.AddDataSink(&sink)
	if err != nil {
		s.Log.Error("add data sink failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) streamDeleteSink(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "query parameter id is required")
		return
	}
	err := s.API.DeleteDataSink(id)
	if err != nil {
		s.Log.Error("delete data sink failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) streamListSinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	entries, err := s.API.ListDataSink()
	if err != nil {
		s.Log.Error("list data sinks failed", "error", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"entries": entries})
}
