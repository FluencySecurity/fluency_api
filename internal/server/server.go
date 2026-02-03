package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/SecurityDo/fluency_api/internal/api"
)

// Server runs the HTTP API backed by the Fluency API client.
type Server struct {
	API  *api.Client
	Addr string
	Log  *slog.Logger
}

// New builds a server that uses the given API client.
func New(apiClient *api.Client, addr string, log *slog.Logger) *Server {
	if log == nil {
		log = slog.Default()
	}
	return &Server{API: apiClient, Addr: addr, Log: log}
}

// Routes returns the HTTP handler for all API and health routes.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /health", s.healthHandler)

	// EventWatch
	mux.HandleFunc("POST /api/eventwatch/summary-search", s.eventwatchSummarySearch)
	mux.HandleFunc("POST /api/eventwatch/timeline-search", s.eventwatchTimelineSearch)
	mux.HandleFunc("POST /api/eventwatch/rule-search", s.eventwatchRuleSearch)

	// Audit
	mux.HandleFunc("POST /api/audit/search", s.auditSearch)

	// Datalake
	mux.HandleFunc("GET /api/datalake", s.datalakeList)
	mux.HandleFunc("POST /api/datalake", s.datalakeAdd)
	mux.HandleFunc("GET /api/datalake/index", s.datalakeListIndex)
	mux.HandleFunc("POST /api/datalake/index", s.datalakeAddIndex)
	mux.HandleFunc("DELETE /api/datalake/index", s.datalakeDeleteIndex)

	// Auth (users and tokens)
	mux.HandleFunc("POST /api/auth/users", s.authAddUser)
	mux.HandleFunc("DELETE /api/auth/users", s.authDeleteUser)
	mux.HandleFunc("GET /api/auth/users", s.authListUsers)
	mux.HandleFunc("POST /api/auth/tokens", s.authAddToken)
	mux.HandleFunc("DELETE /api/auth/tokens", s.authDeleteToken)
	mux.HandleFunc("GET /api/auth/tokens", s.authListTokens)

	// Application
	mux.HandleFunc("GET /api/application/templates", s.applicationListTemplates)
	mux.HandleFunc("POST /api/application/install", s.applicationInstall)
	mux.HandleFunc("GET /api/application/instance", s.applicationGetInstance)
	mux.HandleFunc("DELETE /api/application/instance", s.applicationUninstall)
	mux.HandleFunc("POST /api/application/templates", s.applicationAddTemplate)
	mux.HandleFunc("DELETE /api/application/templates", s.applicationDeleteTemplate)
	mux.HandleFunc("PUT /api/application/templates", s.applicationUpdateTemplate)

	// Processor
	mux.HandleFunc("GET /api/processor", s.processorList)
	mux.HandleFunc("POST /api/processor", s.processorAdd)
	mux.HandleFunc("DELETE /api/processor", s.processorDelete)

	// Integration
	mux.HandleFunc("GET /api/integration", s.integrationList)
	mux.HandleFunc("POST /api/integration", s.integrationAdd)
	mux.HandleFunc("DELETE /api/integration", s.integrationDelete)

	// Stream (sources, sinks, routers)
	mux.HandleFunc("POST /api/stream/connect-sink", s.streamConnectSink)
	mux.HandleFunc("POST /api/stream/connect-router", s.streamConnectRouter)
	mux.HandleFunc("POST /api/stream/router", s.streamAddRouter)
	mux.HandleFunc("POST /api/stream/source", s.streamAddSource)
	mux.HandleFunc("DELETE /api/stream/source", s.streamDeleteSource)
	mux.HandleFunc("GET /api/stream/source", s.streamListSources)
	mux.HandleFunc("POST /api/stream/sink", s.streamAddSink)
	mux.HandleFunc("DELETE /api/stream/sink", s.streamDeleteSink)
	mux.HandleFunc("GET /api/stream/sink", s.streamListSinks)

	// EKS / AWS assumed roles
	mux.HandleFunc("GET /api/eks/pod-role", s.eksGetPodRole)
	mux.HandleFunc("POST /api/eks/test-assumed-role", s.eksTestAssumedRole)
	mux.HandleFunc("POST /api/eks/assumed-role", s.eksAddAssumedRole)
	mux.HandleFunc("DELETE /api/eks/assumed-role", s.eksDeleteAssumedRole)
	mux.HandleFunc("GET /api/eks/assumed-role", s.eksListAssumedRoles)

	// Status
	mux.HandleFunc("GET /api/status", s.statusCheck)

	// Import
	mux.HandleFunc("POST /api/import/processor", s.importProcessor)
	mux.HandleFunc("POST /api/import/application", s.importApplication)

	return mux
}

// Run starts the HTTP server and blocks until it exits.
func (s *Server) Run() error {
	// #region agent log
	if f, err := os.OpenFile("/Users/yuanshentay/GitHub/fluency_api/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(f).Encode(map[string]interface{}{"timestamp": time.Now().UnixMilli(), "location": "server.go:Run", "message": "Run() entered", "hypothesisId": "H4", "sessionId": "debug-session", "runId": "run1", "data": map[string]interface{}{"addr": s.Addr, "log_nil": s.Log == nil}})
		f.Write([]byte("\n"))
		f.Close()
	}
	// #endregion
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		// #region agent log
		if f, e := os.OpenFile("/Users/yuanshentay/GitHub/fluency_api/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); e == nil {
			json.NewEncoder(f).Encode(map[string]interface{}{"timestamp": time.Now().UnixMilli(), "location": "server.go:Run", "message": "net.Listen failed", "hypothesisId": "H3", "sessionId": "debug-session", "runId": "run1", "data": map[string]interface{}{"err": err.Error()}})
			f.Write([]byte("\n"))
			f.Close()
		}
		// #endregion
		return err
	}
	// #region agent log
	if f, err := os.OpenFile("/Users/yuanshentay/GitHub/fluency_api/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(f).Encode(map[string]interface{}{"timestamp": time.Now().UnixMilli(), "location": "server.go:Run", "message": "listener created, about to Log.Info", "hypothesisId": "H2", "sessionId": "debug-session", "runId": "run1", "data": map[string]interface{}{"listening_on": listener.Addr().String()}})
		f.Write([]byte("\n"))
		f.Close()
	}
	// #endregion
	addrStr := listener.Addr().String()
	// Always print startup message (default log level is warn, so Info is not shown)
	fmt.Fprintf(os.Stderr, "Server started successfully. Listening on %s\n", addrStr)
	s.Log.Info("server started successfully", "listening_on", addrStr)
	// #region agent log
	if f, err := os.OpenFile("/Users/yuanshentay/GitHub/fluency_api/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(f).Encode(map[string]interface{}{"timestamp": time.Now().UnixMilli(), "location": "server.go:Run", "message": "after Log.Info, calling Serve", "hypothesisId": "H2", "sessionId": "debug-session", "runId": "run1", "data": map[string]interface{}{}})
		f.Write([]byte("\n"))
		f.Close()
	}
	// #endregion
	return http.Serve(listener, s.Routes())
}
