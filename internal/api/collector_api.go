package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	"github.com/SecurityDo/fluency_api/model"
)

// CollectorList calls the collector_list API (kargs: {}) and returns collectors for web.
func (c *Client) CollectorList() ([]*model.CollectorForWeb, error) {
	svc := fluencyAPI.NewCollectorService(c.fluencyClient)
	entries, err := svc.CollectorList()
	if err != nil {
		c.Logger.Error("failed to list collectors", "error", err)
		return nil, fmt.Errorf("failed to list collectors: %w", err)
	}
	return entries, nil
}

// CollectorStatus calls the collector_status API with collector name and cargs; returns response as map[string]interface{}.
func (c *Client) CollectorStatus(collector string, cargs map[string]interface{}) (map[string]interface{}, error) {
	svc := fluencyAPI.NewCollectorService(c.fluencyClient)
	out, err := svc.CollectorStatus(collector, cargs)
	if err != nil {
		c.Logger.Error("failed to get collector status", "collector", collector, "error", err)
		return nil, fmt.Errorf("failed to get collector status: %w", err)
	}
	return out, nil
}
