package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	"github.com/SecurityDo/fluency_api/model"
)

// AuditSearch calls the audit_search API with the given search string and time range.
func (c *Client) AuditSearch(searchString string, rangeFrom, rangeTo int64) (*model.ElasticSearchResult, error) {
	svc := fluencyAPI.NewAuditService(c.fluencyClient)
	resp, err := svc.AuditSearch(searchString, rangeFrom, rangeTo)
	if err != nil {
		c.Logger.Error("failed to run audit search", "error", err)
		return nil, fmt.Errorf("failed to run audit search: %w", err)
	}
	return resp, nil
}

// DbStatus calls the db_status API and returns the DB status response.
func (c *Client) DbStatus() (*model.DbStatusResponse, error) {
	svc := fluencyAPI.NewAuditService(c.fluencyClient)
	resp, err := svc.DbStatus()
	if err != nil {
		c.Logger.Error("failed to get db status", "error", err)
		return nil, fmt.Errorf("failed to get db status: %w", err)
	}
	return resp, nil
}
