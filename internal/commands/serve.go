package commands

import (
	"encoding/json"
	"os"
	"time"

	"github.com/SecurityDo/fluency_api/internal/server"
	"github.com/spf13/cobra"
)

var serveAddress string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the HTTP API server",
	Long:  "Starts an HTTP server that exposes the Fluency API (e.g. EventWatch search). Uses the same config as the CLI (site-config or cluster/namespace).",
	RunE: func(cmd *cobra.Command, args []string) error {
		// #region agent log
		if f, err := os.OpenFile("/Users/yuanshentay/GitHub/fluency_api/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			json.NewEncoder(f).Encode(map[string]interface{}{"timestamp": time.Now().UnixMilli(), "location": "serve.go:RunE", "message": "serve RunE entered", "hypothesisId": "H1", "sessionId": "debug-session", "runId": "run1", "data": map[string]interface{}{"addr": serveAddress, "appAPI_nil": AppAPI == nil}})
			f.Write([]byte("\n"))
			f.Close()
		}
		// #endregion
		logger := AppAPI.Logger
		srv := server.New(AppAPI, serveAddress, logger)
		return srv.Run()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&serveAddress, "address", "a", ":3000", "Listen address (e.g. :3000 or 127.0.0.1:3000)")
}
