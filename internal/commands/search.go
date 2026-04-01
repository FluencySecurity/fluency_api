package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	"github.com/spf13/cobra"
)

var (
	searchQuery      string
	searchFrom       int64
	searchTo         int64
	searchOutFile    string
	searchMust       []string
	searchMustNot    []string
	searchFacetFields []string
)

func parseFilterEntries(filters []string) []*fluencyAPI.FilterEntry {
	entries := make([]*fluencyAPI.FilterEntry, 0, len(filters))
	for _, f := range filters {
		field, termStr, ok := strings.Cut(f, "=")
		if !ok || field == "" {
			continue
		}
		rawTerms := strings.Split(termStr, ",")
		terms := make([]interface{}, 0, len(rawTerms))
		for _, t := range rawTerms {
			if t != "" {
				terms = append(terms, t)
			}
		}
		entries = append(entries, &fluencyAPI.FilterEntry{Field: field, Terms: terms})
	}
	return entries
}

func parseFacetEntries(fields []string) []*fluencyAPI.FacetEntry {
	entries := make([]*fluencyAPI.FacetEntry, 0, len(fields))
	for _, f := range fields {
		if f != "" {
			entries = append(entries, &fluencyAPI.FacetEntry{Field: f, Order: "desc", Size: 20})
		}
	}
	return entries
}

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

		mustFilters := parseFilterEntries(searchMust)
		mustNotFilters := parseFilterEntries(searchMustNot)
		facets := parseFacetEntries(searchFacetFields)

		resp, err := AppAPI.IndexZoomHistogramLv3(searchQuery, from, to, mustFilters, mustNotFilters, facets)
		if err != nil {
			return err
		}

		hitCount := 0
		if resp != nil && resp.Hits != nil {
			hitCount = int(resp.Hits.TotalHits)
		}

		hasAggs := resp != nil && len(resp.Aggregations) > 0

		if hitCount == 0 && !hasAggs {
			cmd.PrintErrln("No hits found.")
			return nil
		}

		outPath := searchOutFile
		if outPath == "" {
			outPath = fmt.Sprintf("search_results_%d.json", time.Now().UnixMilli())
		} else if !strings.HasSuffix(outPath, ".json") {
			outPath += ".json"
		}

		if hitCount > 0 {
			f, err := os.Create(outPath)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			defer f.Close()

			sources := make([]json.RawMessage, 0, len(resp.Hits.Hits))
			for _, hit := range resp.Hits.Hits {
				sources = append(sources, hit.Source)
			}
			enc := json.NewEncoder(f)
			enc.SetIndent("", "  ")
			if err := enc.Encode(sources); err != nil {
				return fmt.Errorf("failed to write results: %w", err)
			}
			cmd.Printf("%d result(s) written to %s\n", hitCount, outPath)
		}

		if hasAggs {
			aggPath := strings.TrimSuffix(outPath, ".json") + "_aggregations.json"
			af, err := os.Create(aggPath)
			if err != nil {
				return fmt.Errorf("failed to create aggregations file: %w", err)
			}
			defer af.Close()

			enc := json.NewEncoder(af)
			enc.SetIndent("", "  ")
			if err := enc.Encode(resp.Aggregations); err != nil {
				return fmt.Errorf("failed to write aggregations: %w", err)
			}
			cmd.Printf("Aggregations written to %s\n", aggPath)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&searchQuery, "query", "", "Search query string")
	searchCmd.Flags().Int64Var(&searchFrom, "from", 0, "Range start (Unix ms); default: 1 hour ago")
	searchCmd.Flags().Int64Var(&searchTo, "to", 0, "Range end (Unix ms); default: now")
	searchCmd.Flags().StringVar(&searchOutFile, "out", "", "Output file path (default: search_results_<timestamp>.json)")
	searchCmd.Flags().StringArrayVar(&searchMust, "must-filter", nil, "Must filter as field=val1,val2 (repeatable)")
	searchCmd.Flags().StringArrayVar(&searchMustNot, "must-not-filter", nil, "Must-not filter as field=val1,val2 (repeatable)")
	searchCmd.Flags().StringArrayVar(&searchFacetFields, "facet", nil, "Facet field name (repeatable; order=desc, size=20)")
}
