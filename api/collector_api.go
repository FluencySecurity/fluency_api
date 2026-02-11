package api

import (
	"github.com/SecurityDo/fluency_api/client"
	"github.com/SecurityDo/fluency_api/model"
)

// CollectorService provides helpers for calling collector endpoints.
type CollectorService struct {
	client *client.FluencyClient
}

// NewCollectorService constructs a CollectorService instance backed by the provided client.
func NewCollectorService(c *client.FluencyClient) *CollectorService {
	return &CollectorService{client: c}
}

func (s *CollectorService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

type CollectorListResponse struct {
	Entries []*model.CollectorForWeb `json:"entries"`
}

// CollectorList calls the collector_list API with kargs {} and returns an array of CollectorForWeb.
func (s *CollectorService) CollectorList() ([]*model.CollectorForWeb, error) {
	var out CollectorListResponse
	if err := s.call("collector_list", map[string]interface{}{}, &out); err != nil {
		return nil, err
	}
	return out.Entries, nil
}

// CollectorStatusKargs is the payload for collector_status (kargs: collector + cargs).
type CollectorStatusKargs struct {
	Collector string                 `json:"collector"`
	Cargs     map[string]interface{} `json:"cargs,omitempty"`
}

// CollectorStatus calls the collector_status API with kargs {"collector":"name","cargs":{...}} and returns the response as map[string]interface{}.
func (s *CollectorService) CollectorStatus(collector string, cargs map[string]interface{}) (map[string]interface{}, error) {
	payload := CollectorStatusKargs{Collector: collector, Cargs: cargs}
	var out map[string]interface{}
	if err := s.call("get_system_status", payload, &out); err != nil {
		return nil, err
	}
	return out, nil
}
