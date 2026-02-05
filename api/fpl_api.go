package api

import (
	"github.com/SecurityDo/fluency_api/client"
	"github.com/SecurityDo/fluency_api/model"
)

// FPLService provides helpers for calling FPL report/task/result endpoints.
type FPLService struct {
	client *client.FluencyClient
}

// NewFPLService constructs an FPLService instance backed by the provided client.
func NewFPLService(client *client.FluencyClient) *FPLService {
	return &FPLService{client: client}
}

func (s *FPLService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

// FPLIDRequest is the kargs payload for get_fpl_task and get_fpl_result.
type FPLIDRequest struct {
	ID uint `json:"id"`
	FPL bool `json:"fpl"`
}

// RunFPLV2ReportResponse is the response from run_fplv2_report (task id).
type RunFPLV2ReportResponse struct {
	ID uint `json:"taskID"`
}

// RunReport calls /api/ds/run_fplv2_report with the given request and returns the task id.
func (s *FPLService) RunReport(req *model.RunFPLV2Report) (uint, error) {
	var resp RunFPLV2ReportResponse
	if err := s.call("run_fplv2_report", req, &resp); err != nil {
		return 0, err
	}
	return resp.ID, nil
}

// GetTaskByID calls /api/ds/get_fpl_task with kargs {id: id} and returns the task response.
func (s *FPLService) GetTaskByID(id uint) (*model.GetTaskResponse, error) {
	req := &FPLIDRequest{ID: id, FPL: true}
	var resp model.GetTaskResponse
	if err := s.call("get_fplv2_task", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetResultsByID calls /api/ds/get_fpl_result with kargs {id: id} and returns the metric result.
func (s *FPLService) GetResultsByID(id uint) (*model.GetMetricFPLResultResponse, error) {
	req := &FPLIDRequest{ID: id, FPL: true}
	var resp model.GetMetricFPLResultResponse
	if err := s.call("get_fplv2_result", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
