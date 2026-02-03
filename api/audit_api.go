package api

import (
	"github.com/SecurityDo/fluency_api/client"
	"github.com/SecurityDo/fluency_api/model"
)

// AuditService provides helpers for calling audit endpoints.
type AuditService struct {
	client *client.FluencyClient
}

// NewAuditService constructs an AuditService instance backed by the provided client.
func NewAuditService(client *client.FluencyClient) *AuditService {
	return &AuditService{client: client}
}

func (s *AuditService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

// AuditSearch calls /api/ds/audit_search with the given search string and time range.
func (s *AuditService) AuditSearch(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	req := &ElasticSearchRequest{
		Options: &SimpleSearchOption{
			SearchStr:  searchString,
			RangeFrom:  rangeFrom,
			RangeTo:    rangeTo,
			RangeField: "createdOn",
			FetchLimit: 100,
			FetchOffset: 0,
			SortField:   "createdOn",
			SortOrder:   "desc",
			Facets: &FacetsOption{
				Facets:         []*FacetEntry{},
				MustFilters:    []*FilterEntry{},
				MustNotFilters: []*FilterEntry{},
			},
		},
	}
	var resp model.ElasticSearchResult
	if err := s.call("audit_search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DbStatus calls /api/ds/db_status with an empty request and returns the DB status.
func (s *AuditService) DbStatus() (*model.DbStatusResponse, error) {
	var resp model.DbStatusResponse
	if err := s.call("db_status", map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
