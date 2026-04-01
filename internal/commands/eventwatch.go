package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/SecurityDo/fluency_api/model"
	"github.com/spf13/cobra"
)

var (
	eventwatchQuery     string
	eventwatchFrom      int64
	eventwatchTo        int64
	ruleTestBucketFile  string
	ruleTestInputFile   string
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

var eventwatchRuleTestCmd = &cobra.Command{
	Use:   "rule_test",
	Short: "Test whether an input event matches a bucket rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		bucketData, err := os.ReadFile(ruleTestBucketFile)
		if err != nil {
			return fmt.Errorf("failed to read bucket file: %w", err)
		}
		if !json.Valid(bucketData) {
			return fmt.Errorf("bucket file %q is not valid JSON", ruleTestBucketFile)
		}
		var bucket model.EventWatchBucket
		if err := json.Unmarshal(bucketData, &bucket); err != nil {
			return fmt.Errorf("bucket file %q is not a valid JSON object: %w", ruleTestBucketFile, err)
		}

		inputData, err := os.ReadFile(ruleTestInputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
		if !json.Valid(inputData) {
			return fmt.Errorf("input file %q is not valid JSON", ruleTestInputFile)
		}
		var input map[string]interface{}
		if err := json.Unmarshal(inputData, &input); err != nil {
			return fmt.Errorf("input file %q is not a valid JSON object: %w", ruleTestInputFile, err)
		}

		hit, err := AppAPI.RuleTest(bucket, input)
		if err != nil {
			return err
		}
		cmd.Printf("hit: %v\n", hit)
		return nil
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
	eventwatchCmd.AddCommand(eventwatchSummarySearchCmd, eventwatchTimelineSearchCmd, eventwatchRuleSearchCmd, eventwatchRuleTestCmd)

	eventwatchSummarySearchCmd.Flags().StringVar(&eventwatchQuery, "query", "", "Search query")
	eventwatchSummarySearchCmd.Flags().Int64Var(&eventwatchFrom, "from", 0, "Range start (Unix ms); default: 1 hour ago")
	eventwatchSummarySearchCmd.Flags().Int64Var(&eventwatchTo, "to", 0, "Range end (Unix ms); default: now")

	eventwatchTimelineSearchCmd.Flags().StringVar(&eventwatchQuery, "query", "", "Search query")
	eventwatchTimelineSearchCmd.Flags().Int64Var(&eventwatchFrom, "from", 0, "Range start (Unix ms); default: 1 hour ago")
	eventwatchTimelineSearchCmd.Flags().Int64Var(&eventwatchTo, "to", 0, "Range end (Unix ms); default: now")

	eventwatchRuleSearchCmd.Flags().StringVar(&eventwatchQuery, "query", "", "Search query")

	eventwatchRuleTestCmd.Flags().StringVar(&ruleTestBucketFile, "bucket", "", "Path to JSON file containing the EventWatchBucket")
	eventwatchRuleTestCmd.Flags().StringVar(&ruleTestInputFile, "input", "", "Path to JSON file containing the input event")
	eventwatchRuleTestCmd.MarkFlagRequired("bucket")
	eventwatchRuleTestCmd.MarkFlagRequired("input")
}
