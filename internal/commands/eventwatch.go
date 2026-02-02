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
	Use:   "search_summary",
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
		return printEventwatchHits(cmd, resp, "BehaviorSummary")
	},
}

var eventwatchTimelineSearchCmd = &cobra.Command{
	Use:   "search_timeline",
	Short: "Run timeline search (fsm_behavior_search)",
	RunE: func(cmd *cobra.Command, args []string) error {
		from, to := eventwatchFrom, eventwatchTo
		if from == 0 && to == 0 {
			now := time.Now().UnixMilli()
			to = now
			from = now - int64(time.Hour/time.Millisecond)
		}
		resp, err := AppAPI.TimelineSearch(eventwatchQuery, from, to)
		if err != nil {
			return err
		}
		return printEventwatchHits(cmd, resp, "BehaviorEvent")
	},
}

var eventwatchRuleSearchCmd = &cobra.Command{
	Use:   "search_rule",
	Short: "Run rule search (eventwatch_bucket_search)",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.RuleSearch(eventwatchQuery)
		if err != nil {
			return err
		}
		return printEventwatchHits(cmd, resp, "BehaviorRule")
	},
}

func printEventwatchHits(cmd *cobra.Command, resp *model.ElasticSearchResult, sourceType string) error {
	if resp.Hits == nil || len(resp.Hits.Hits) == 0 {
		cmd.PrintErrln("No hits found.")
		return nil
	}
	for _, hit := range resp.Hits.Hits {
		switch sourceType {
		case "BehaviorSummary":
			var src model.BehaviorSummary
			if err := json.Unmarshal(hit.Source, &src); err != nil {
				cmd.PrintErrf("skip hit %s: invalid _source: %v\n", hit.Id, err)
				continue
			}
			cmd.Printf("Key: %s, RiskScore: %d\n", src.Key, src.RiskScore)
		case "BehaviorEvent":
			var src model.BehaviorEvent
			if err := json.Unmarshal(hit.Source, &src); err != nil {
				cmd.PrintErrf("skip hit %s: invalid _source: %v\n", hit.Id, err)
				continue
			}
			cmd.Printf("Key: %s, RiskScore: %d\n", src.Key, src.RiskScore)
		case "BehaviorRule":
			var src model.EventWatchBucket
			if err := json.Unmarshal(hit.Source, &src); err != nil {
				cmd.PrintErrf("skip hit %s: invalid _source: %v\n", hit.Id, err)
				continue
			}
			cmd.Printf("Name: %s, Group: %s, Repository: %s\n", src.Name, src.Group, src.Repository)
		default:
			cmd.PrintErrf("skip hit %s: unknown source type %q\n", hit.Id, sourceType)
		}
	}
	return nil
}

func init() {
	RootCmd.AddCommand(eventwatchCmd)
	eventwatchCmd.AddCommand(eventwatchSummarySearchCmd, eventwatchTimelineSearchCmd, eventwatchRuleSearchCmd)

	eventwatchSummarySearchCmd.Flags().StringVar(&eventwatchQuery, "query", "", "Search query")
	eventwatchSummarySearchCmd.Flags().Int64Var(&eventwatchFrom, "from", 0, "Range start (Unix ms); default: 1 hour ago")
	eventwatchSummarySearchCmd.Flags().Int64Var(&eventwatchTo, "to", 0, "Range end (Unix ms); default: now")

	eventwatchTimelineSearchCmd.Flags().StringVar(&eventwatchQuery, "query", "", "Search query")
	eventwatchTimelineSearchCmd.Flags().Int64Var(&eventwatchFrom, "from", 0, "Range start (Unix ms); default: 1 hour ago")
	eventwatchTimelineSearchCmd.Flags().Int64Var(&eventwatchTo, "to", 0, "Range end (Unix ms); default: now")

	eventwatchRuleSearchCmd.Flags().StringVar(&eventwatchQuery, "query", "", "Search query")
}
