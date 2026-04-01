package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	"github.com/SecurityDo/fluency_api/model"
)

// IndexZoomHistogramLv3 calls /api/ds/get_index_zoom_histogram_lv3 with the given search string and time range.
func (c *Client) IndexZoomHistogramLv3(searchString string, rangeFrom, rangeTo int64, mustFilters, mustNotFilters []*fluencyAPI.FilterEntry, facets []*fluencyAPI.FacetEntry) (*model.ElasticSearchResult, error) {
	svc := fluencyAPI.NewSearchService(c.fluencyClient)
	resp, err := svc.IndexZoomHistogramLv3(searchString, rangeFrom, rangeTo, mustFilters, mustNotFilters, facets)
	if err != nil {
		c.Logger.Error("failed to run index zoom histogram lv3 search", "error", err)
		return nil, fmt.Errorf("failed to run index zoom histogram lv3 search: %w", err)
	}
	return resp, nil
}
