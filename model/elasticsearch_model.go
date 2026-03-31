package model

import (
	"net/http"

	elastic "github.com/olivere/elastic/v7"
)

type ElasticSearchResult struct {
	Header          http.Header            `json:"-"`
	TookInMillis    int64                  `json:"took,omitempty"`             // search time in milliseconds
	TerminatedEarly bool                   `json:"terminated_early,omitempty"` // request terminated early
	NumReducePhases int                    `json:"num_reduce_phases,omitempty"`
	Clusters        *SearchResultCluster   `json:"_clusters,omitempty"`    // 6.1.0+
	ScrollId        string                 `json:"_scroll_id,omitempty"`   // only used with Scroll and Scan operations
	Hits            *SearchHits            `json:"hits,omitempty"`         // the actual search hits
	Aggregations    elastic.Aggregations   `json:"aggregations,omitempty"` // results from aggregations
	TimedOut        bool                   `json:"timed_out,omitempty"`    // true if the search timed out
	Error           *elastic.ErrorDetails  `json:"error,omitempty"`        // only used in MultiGet
	Profile         *elastic.SearchProfile `json:"profile,omitempty"`      // profiling results, if optional Profile API was active for this search
	Status          int                    `json:"status,omitempty"`       // used in MultiSearch
	PitId           string                 `json:"pit_id,omitempty"`       // Point In Time ID
}

type SearchResultCluster struct {
	Successful int `json:"successful,omitempty"`
	Total      int `json:"total,omitempty"`
	Skipped    int `json:"skipped,omitempty"`
}

type SearchHits struct {
	TotalHits int64 `json:"total"` // total number of hits found
	// TotalHits *TotalHits   `json:"total,omitempty"`     // total number of hits found
	MaxScore *float64             `json:"max_score,omitempty"` // maximum score of all hits
	Hits     []*elastic.SearchHit `json:"hits,omitempty"`      // the actual hits returned
}

type NestedHit struct {
	Field  string     `json:"field"`
	Offset int        `json:"offset,omitempty"`
	Child  *NestedHit `json:"_nested,omitempty"`
}

type TotalHits struct {
	Value    int64  `json:"value"`    // value of the total hit count
	Relation string `json:"relation"` // how the value should be interpreted: accurate ("eq") or a lower bound ("gte")
}
