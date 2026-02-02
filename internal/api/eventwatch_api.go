package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	"github.com/SecurityDo/fluency_api/model"
)

// SummarySearch calls the overview summary search API with the given search string and time range.
func (c *Client) SummarySearch(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	svc := fluencyAPI.NewEventWatchService(c.fluencyClient)
	resp, err := svc.SummarySearch(searchString, rangeFrom, rangeTo)
	if err != nil {
		c.Logger.Error("failed to run summary search", "error", err)
		return nil, fmt.Errorf("failed to run summary search: %w", err)
	}
	return resp, nil
}

// TimelineSearch calls the fsm_behavior_search API with the given search string and time range.
func (c *Client) TimelineSearch(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	svc := fluencyAPI.NewEventWatchService(c.fluencyClient)
	resp, err := svc.TimelineSearch(searchString, rangeFrom, rangeTo)
	if err != nil {
		c.Logger.Error("failed to run timeline search", "error", err)
		return nil, fmt.Errorf("failed to run timeline search: %w", err)
	}
	return resp, nil
}
