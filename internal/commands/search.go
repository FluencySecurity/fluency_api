package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	searchQuery   string
	searchFrom    int64
	searchTo      int64
	searchOutFile string
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search index zoom histogram (lv3)",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, to := searchFrom, searchTo
		if from == 0 && to == 0 {
			now := time.Now().UnixMilli()
			to = now
			from = now - int64(time.Hour/time.Millisecond)
		}

		resp, err := AppAPI.IndexZoomHistogramLv3(searchQuery, from, to)
		if err != nil {
			return err
		}

		hitCount := 0
		if resp != nil && resp.Hits != nil {
			hitCount = int(resp.Hits.TotalHits)
		}

		if hitCount == 0 {
			cmd.PrintErrln("No hits found.")
			return nil
		}

		outPath := searchOutFile
		if outPath == "" {
			outPath = fmt.Sprintf("search_results_%d.json", time.Now().UnixMilli())
		}

		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()

		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(resp.Hits.Hits); err != nil {
			return fmt.Errorf("failed to write results: %w", err)
		}

		cmd.Printf("%d result(s) written to %s\n", hitCount, outPath)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&searchQuery, "query", "", "Search query string")
	searchCmd.Flags().Int64Var(&searchFrom, "from", 0, "Range start (Unix ms); default: 1 hour ago")
	searchCmd.Flags().Int64Var(&searchTo, "to", 0, "Range end (Unix ms); default: now")
	searchCmd.Flags().StringVar(&searchOutFile, "out", "", "Output file path (default: search_results_<timestamp>.json)")
}
