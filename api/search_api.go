package api

import (
	"encoding/json"

	"github.com/SecurityDo/fluency_api/client"
	"github.com/SecurityDo/fluency_api/model"
)

// SearchService provides helpers for calling search endpoints.
type SearchService struct {
	client *client.FluencyClient
}

// NewSearchService constructs a SearchService instance backed by the provided client.
func NewSearchService(client *client.FluencyClient) *SearchService {
	return &SearchService{client: client}
}

// IndexZoomHistogramLv3 calls /api/ds/get_index_zoom_histogram_lv3 with the given search string and time range.
func (s *SearchService) IndexZoomHistogramLv3(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	req := &ElasticSearchRequest{
		Partition: "default",
		DataType:  "event",
		Options: &SimpleSearchOption{
			DataType:       "event",
			SearchStr:      searchString,
			DateFacetField: "@timestamp",
			RangeFrom:      rangeFrom,
			RangeTo:        rangeTo,
			FetchLimit:     100,
			FetchOffset:    0,
			SortField:      "@timestamp",
			SortOrder:      "desc",
			Facets: &FacetsOption{
				Facets:         []*FacetEntry{},
				MustFilters:    []*FilterEntry{},
				MustNotFilters: []*FilterEntry{},
			},
		},
	}
	res, err := s.client.GenericCall("api/ds", "get_index_zoom_histogram_lv3", req)
	if err != nil {
		return nil, err
	}
	// A nil response means the server returned "response": null — treat as no results.
	if res == nil {
		return &model.ElasticSearchResult{}, nil
	}
	var resp model.ElasticSearchResult
	if err := json.Unmarshal(res.GetBytes(), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
