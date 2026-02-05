package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	"github.com/SecurityDo/fluency_api/model"
)

// RunReport calls the run_fplv2_report API with the given request and returns the task id.
func (c *Client) RunReport(req *model.RunFPLV2Report) (uint, error) {
	svc := fluencyAPI.NewFPLService(c.fluencyClient)
	id, err := svc.RunReport(req)
	if err != nil {
		c.Logger.Error("failed to run FPL report", "error", err)
		return 0, fmt.Errorf("failed to run FPL report: %w", err)
	}
	return id, nil
}

// GetTaskByID calls the get_fpl_task API with the given id and returns the task response.
func (c *Client) GetTaskByID(id uint) (*model.GetTaskResponse, error) {
	svc := fluencyAPI.NewFPLService(c.fluencyClient)
	resp, err := svc.GetTaskByID(id)
	if err != nil {
		c.Logger.Error("failed to get FPL task", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get FPL task: %w", err)
	}
	return resp, nil
}

// GetResultsByID calls the get_fpl_result API with the given id and returns the metric result.
func (c *Client) GetResultsByID(id uint) (*model.GetMetricFPLResultResponse, error) {
	svc := fluencyAPI.NewFPLService(c.fluencyClient)
	resp, err := svc.GetResultsByID(id)
	if err != nil {
		c.Logger.Error("failed to get FPL result", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get FPL result: %w", err)
	}
	return resp, nil
}
