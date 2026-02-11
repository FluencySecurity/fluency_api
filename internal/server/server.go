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
	API          *api.Client
	Addr         string
	Log          *slog.Logger
	endpointList []EndpointEntry // set in Routes() for GET /api/endpoints
}

// routeDef defines one route: method, path, description, and handler.
type routeDef struct {
	Method      string
	Path        string
	Description string
	Handler     func(http.ResponseWriter, *http.Request)
}

// routeDefinitions returns all API routes with handlers bound to this server.
func (s *Server) routeDefinitions() []routeDef {
	return []routeDef{
		{"GET", "/health", "Health check", s.healthHandler},
		{"GET", "/api/endpoints", "List all API endpoints", s.endpointsHandler},
		{"GET", "/api/eventwatch/summary-search", "EventWatch summary search", s.eventwatchSummarySearch},
		{"GET", "/api/eventwatch/timeline-search", "EventWatch timeline search", s.eventwatchTimelineSearch},
		{"GET", "/api/eventwatch/rule-search", "EventWatch rule search", s.eventwatchRuleSearch},
		{"GET", "/api/collector/list", "List collectors (CollectorForWeb)", s.collectorList},
		{"POST", "/api/collector/status", "Collector status (body: collector, cargs)", s.collectorStatus},
		{"GET", "/api/audit/search", "Audit search", s.auditSearch},
		{"GET", "/api/audit/db_status", "DB status (audit)", s.auditDbStatus},
		{"GET", "/api/datalake", "List datalakes", s.datalakeList},
		{"POST", "/api/datalake", "Add datalake", s.datalakeAdd},
		{"GET", "/api/datalake/index", "List datalake indexes (query: lake)", s.datalakeListIndex},
		{"POST", "/api/datalake/index", "Add datalake index", s.datalakeAddIndex},
		{"DELETE", "/api/datalake/index", "Delete datalake index (query: lake, index)", s.datalakeDeleteIndex},
		{"POST", "/api/auth/users", "Add user", s.authAddUser},
		{"DELETE", "/api/auth/users", "Delete user (query: name)", s.authDeleteUser},
		{"GET", "/api/auth/users", "List users", s.authListUsers},
		{"POST", "/api/auth/tokens", "Add token", s.authAddToken},
		{"DELETE", "/api/auth/tokens", "Delete token (query: name)", s.authDeleteToken},
		{"GET", "/api/auth/tokens", "List tokens", s.authListTokens},
		{"GET", "/api/application/templates", "List application templates", s.applicationListTemplates},
		{"POST", "/api/application/install", "Install application instance", s.applicationInstall},
		{"GET", "/api/application/instance", "Get application instance (query: app, instance)", s.applicationGetInstance},
		{"DELETE", "/api/application/instance", "Uninstall instance (query: app, instance)", s.applicationUninstall},
		{"POST", "/api/application/templates", "Add template", s.applicationAddTemplate},
		{"DELETE", "/api/application/templates", "Delete template (query: app)", s.applicationDeleteTemplate},
		{"PUT", "/api/application/templates", "Update template (query: app)", s.applicationUpdateTemplate},
		{"GET", "/api/processor", "List processors", s.processorList},
		{"POST", "/api/processor", "Add processor", s.processorAdd},
		{"DELETE", "/api/processor", "Delete processor (query: name)", s.processorDelete},
		{"GET", "/api/integration", "List integrations", s.integrationList},
		{"POST", "/api/integration", "Add integration", s.integrationAdd},
		{"DELETE", "/api/integration", "Delete integration (query: id)", s.integrationDelete},
		{"POST", "/api/stream/connect-sink", "Connect router to sink", s.streamConnectSink},
		{"POST", "/api/stream/connect-router", "Connect source to router", s.streamConnectRouter},
		{"POST", "/api/stream/router", "Add stream router", s.streamAddRouter},
		{"POST", "/api/stream/source", "Add stream source", s.streamAddSource},
		{"DELETE", "/api/stream/source", "Delete stream source (query: id)", s.streamDeleteSource},
		{"GET", "/api/stream/source", "List stream sources", s.streamListSources},
		{"POST", "/api/stream/sink", "Add stream sink", s.streamAddSink},
		{"DELETE", "/api/stream/sink", "Delete stream sink (query: id)", s.streamDeleteSink},
		{"GET", "/api/stream/sink", "List stream sinks", s.streamListSinks},
		{"GET", "/api/stream/status", "Platform status", s.platformStatus},
		{"GET", "/api/eks/pod-role", "Get EKS pod role", s.eksGetPodRole},
		{"POST", "/api/eks/test-assumed-role", "Test assumed role", s.eksTestAssumedRole},
		{"POST", "/api/eks/assumed-role", "Add assumed role", s.eksAddAssumedRole},
		{"DELETE", "/api/eks/assumed-role", "Delete assumed role (query: id)", s.eksDeleteAssumedRole},
		{"GET", "/api/eks/assumed-role", "List assumed roles", s.eksListAssumedRoles},
		{"GET", "/api/status", "System status (Kubernetes only)", s.statusCheck},
		{"POST", "/api/import/processor", "Import processors from repo", s.importProcessor},
		{"POST", "/api/import/application", "Import application templates from repo", s.importApplication},
		{"POST", "/api/fpl/run", "Run FPL v2 report (run_report)", s.fplRunReport},
		{"GET", "/api/fpl/task", "Get FPL task by ID (query: id)", s.fplGetTask},
		{"GET", "/api/fpl/results", "Get FPL results by task ID (query: id)", s.fplGetResults},
	}
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
	routes := s.routeDefinitions()
	s.endpointList = make([]EndpointEntry, 0, len(routes))
	for _, r := range routes {
		mux.HandleFunc(r.Method+" "+r.Path, r.Handler)
		s.endpointList = append(s.endpointList, EndpointEntry{Method: r.Method, Path: r.Path, Description: r.Description})
	}
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
