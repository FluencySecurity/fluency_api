package commands

import (
	"encoding/json"
	"time"

	"github.com/SecurityDo/fluency_api/model"
	"github.com/spf13/cobra"
)

var (
	eventwatchQuery string
	eventwatchFrom  int64
	eventwatchTo    int64
)

var eventwatchCmd = &cobra.Command{
	Use:   "eventwatch",
	Short: "EventWatch service",
}

var eventwatchSummarySearchCmd = &cobra.Command{
	Use:   "summary_search",
	Short: "Run summary search",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, to := eventwatchFrom, eventwatchTo
		if from == 0 && to == 0 {
			now := time.Now().UnixMilli()
			to = now
			from = now - int64(time.Hour/time.Millisecond)
		}
		resp, err := AppAPI.SummarySearch(eventwatchQuery, from, to)
		if err != nil {
			return err
		}
		if resp.Hits == nil || len(resp.Hits.Hits) == 0 {
			cmd.PrintErrln("No hits found.")
			return nil
		}
		for _, hit := range resp.Hits.Hits {
			var src model.BehaviorEvent
			if err := json.Unmarshal(hit.Source, &src); err != nil {
				cmd.PrintErrf("skip hit %s: invalid _source: %v\n", hit.Id, err)
				continue
			}
			cmd.Printf("Key: %s, RiskScore: %d\n", src.Key, src.RiskScore)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(eventwatchCmd)
	eventwatchCmd.AddCommand(eventwatchSummarySearchCmd)

	eventwatchSummarySearchCmd.Flags().StringVar(&eventwatchQuery, "query", "", "Search query")
	eventwatchSummarySearchCmd.Flags().Int64Var(&eventwatchFrom, "from", 0, "Range start (Unix ms); default: 1 hour ago")
	eventwatchSummarySearchCmd.Flags().Int64Var(&eventwatchTo, "to", 0, "Range end (Unix ms); default: now")
}
