package api

import (
	"github.com/SecurityDo/fluency_api/client"
	"github.com/SecurityDo/fluency_api/model"
)

// EventWatchService provides helpers for calling eventwatch/overview endpoints.
type EventWatchService struct {
	client *client.FluencyClient
}

// NewEventWatchService constructs an EventWatchService instance backed by the provided client.
func NewEventWatchService(client *client.FluencyClient) *EventWatchService {
	return &EventWatchService{client: client}
}

func (s *EventWatchService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

// SummarySearch calls /api/ds/overview_summary_search with the given search string and time range.
func (s *EventWatchService) SummarySearch(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	req := &ElasticSearchRequest{
		Options: &SimpleSearchOption{
			SearchStr:   searchString,
			RangeFrom:   rangeFrom,
			RangeTo:     rangeTo,
			RangeField:  "from",
			FetchLimit:  100,
			FetchOffset: 0,
			SortField:   "to",
			SortOrder:   "desc",
			Facets: &FacetsOption{
				Facets:         []*FacetEntry{},
				MustFilters:    []*FilterEntry{},
				MustNotFilters: []*FilterEntry{},
			},
		},
	}
	var resp model.ElasticSearchResult
	if err := s.call("behavior_summary_search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TimelineSearch calls /api/ds/fsm_behavior_search with the given search string and time range.
func (s *EventWatchService) TimelineSearch(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	req := &ElasticSearchRequest{
		Options: &SimpleSearchOption{
			SearchStr:   searchString,
			RangeFrom:   rangeFrom,
			RangeTo:     rangeTo,
			RangeField:  "timestamp",
			FetchLimit:  100,
			FetchOffset: 0,
			SortField:   "timestamp",
			SortOrder:   "desc",
			Facets: &FacetsOption{
				Facets:         []*FacetEntry{},
				MustFilters:    []*FilterEntry{},
				MustNotFilters: []*FilterEntry{},
			},
		},
	}
	var resp model.ElasticSearchResult
	if err := s.call("fsm_behavior_search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RuleSearch calls /api/ds/eventwatch_bucket_search with the given search string (no time range).
func (s *EventWatchService) RuleSearch(searchString string) (*model.ElasticSearchResult, error) {
	req := &ElasticSearchRequest{
		Options: &SimpleSearchOption{
			SearchStr:   searchString,
			FetchLimit:  100,
			FetchOffset: 0,
			SortField:   "name",
			SortOrder:   "asc",
			Facets: &FacetsOption{
				Facets:         []*FacetEntry{},
				MustFilters:    []*FilterEntry{},
				MustNotFilters: []*FilterEntry{},
			},
		},
	}
	var resp model.ElasticSearchResult
	if err := s.call("eventwatch_bucket_search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type RuleTestRequest struct {
	Bucket model.EventWatchBucket `json:"bucket"`
	Input  map[string]interface{} `json:"input"`
}

type RuleTestResponse struct {
	Hit bool `json:"hit"`
}

// RuleTest calls /api/ds/eventwatch_rule_test to check if the input matches the given bucket.
func (s *EventWatchService) RuleTest(bucket model.EventWatchBucket, input map[string]interface{}) (*RuleTestResponse, error) {
	req := &RuleTestRequest{
		Bucket: bucket,
		Input:  input,
	}
	var resp RuleTestResponse
	if err := s.call("eventwatch_rule_test", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type ElasticSearchRequest struct {
	Options   *SimpleSearchOption `json:"options,omitempty"`
	Partition string              `json:"partition,omitempty"`
	DataType  string              `json:"dataType,omitempty"`
}
