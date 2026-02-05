package api

import (
	"fmt"

	"github.com/SecurityDo/fluency_api/api"
	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	model "github.com/SecurityDo/fluency_api/model"
)

func (c *Client) AddDataSource(source *model.DataSourceConfig) (resp *fluencyAPI.AddDataSourceResponse, err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	resp, err = platformService.AddDataSource(source)

	if err != nil {
		c.Logger.Error("failed to add data source", "error", err)
		return nil, fmt.Errorf("failed to add data source: %s", err.Error())
	}
	return resp, nil
}

func (c *Client) DeleteDataSource(id string) (err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	err = platformService.DeleteDataSource(id)

	if err != nil {
		c.Logger.Error("failed to delete data source", "error", err)
		return fmt.Errorf("failed to delete data source: %s", err.Error())
	}
	return nil
}

func (c *Client) ListDataSource() (entries []*model.DataSourceConfig, err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	entries, err = platformService.ListDataSource()

	if err != nil {
		c.Logger.Error("failed to list data source", "error", err)
		return nil, fmt.Errorf("failed to list data source: %s", err.Error())
	}
	return entries, nil
}

func (c *Client) AddDataSink(sink *model.DataSinkConfig) (resp *fluencyAPI.AddDataSinkResponse, err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	resp, err = platformService.AddDataSink(sink)

	if err != nil {
		c.Logger.Error("failed to add data sink", "error", err)
		return nil, fmt.Errorf("failed to add data sink: %s", err.Error())
	}
	return resp, nil
}

func (c *Client) DeleteDataSink(id string) (err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	err = platformService.DeleteDataSink(id)

	if err != nil {
		c.Logger.Error("failed to delete data sink", "error", err)
		return fmt.Errorf("failed to delete data sink: %s", err.Error())
	}
	return nil
}

func (c *Client) ListDataSink() (entries []*model.DataSinkConfig, err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	entries, err = platformService.ListDataSink()

	if err != nil {
		c.Logger.Error("failed to list data sink", "error", err)
		return nil, fmt.Errorf("failed to list data sink: %s", err.Error())
	}
	return entries, nil
}

func (c *Client) AddRouter(routerConfig *model.RouterConfig) (id string, err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	resp, err := platformService.AddRouter(routerConfig)

	if err != nil {
		c.Logger.Error("failed to add router", "error", err)
		return "", fmt.Errorf("failed to add router: %s", err.Error())
	}
	return resp.ID, nil
}

func (c *Client) AddSimpleRouter(processorName string, routerName string) (id string, err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	id, err = platformService.AddSimpleRouter(processorName, routerName)

	if err != nil {
		c.Logger.Error("failed to add router", "error", err)
		return "", fmt.Errorf("failed to add router: %s", err.Error())
	}
	return id, nil
}

func (c *Client) SetRouterSink(routerID, sinkID string) (err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	result, err := platformService.GetRouter(routerID)

	if err != nil {
		c.Logger.Error("failed to get router by ID", "error", err)
		return fmt.Errorf("failed get router by ID: %s", err.Error())
	}
	if len(result.Pipes) == 0 {
		c.Logger.Error("router has no pipes", "routerID", routerID)
		return fmt.Errorf("router has no pipes: %s", routerID)
	}
	pipeConfig := result.Pipes[len(result.Pipes)-1] // Get the last pipe

	for _, sid := range pipeConfig.SinkIDs {
		if sid == sinkID {
			c.Logger.Error("sink already exists in router", "routerID", routerID, "sinkID", sinkID)
			return fmt.Errorf("sink already exists in router: %s, sink: %s", routerID, sinkID)
		}
	}
	pipeConfig.SinkIDs = append(pipeConfig.SinkIDs, sinkID)

	req := &api.PipeUpdateReq{
		RouterID:   pipeConfig.RouterID,
		PipeConfig: pipeConfig,
	}

	err = platformService.UpdatePipe(req)
	if err != nil {
		c.Logger.Error("failed to update pipe with new sink", "error", err)
		return fmt.Errorf("failed to update pipe with new sink: %s", err.Error())
	}

	return nil
}

func (c *Client) SetSourceRouter(sourceID, routerID string) (err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	err = platformService.SetDataSourceRouter(&api.SourceSetRouterReq{
		DataSourceID: sourceID,
		RouterID:     routerID,
	})

	if err != nil {
		c.Logger.Error("failed to connect source to router", "error", err)
		return fmt.Errorf("failed to connect source to router: %s", err.Error())
	}
	return nil
}

// PlatformStatus calls the platform ListConfigs API and returns the metrics from the configuration snapshot.
func (c *Client) PlatformStatus() ([]*fluencyAPI.ComponentStat, error) {
	svc := fluencyAPI.NewPlatformService(c.fluencyClient)
	resp, err := svc.ListConfigs()
	if err != nil {
		c.Logger.Error("failed to get platform status", "error", err)
		return nil, fmt.Errorf("platform status: %w", err)
	}
	return resp.Metrics, nil
}

/*
func (c *Client) AddPipe(routerConfig *model.StreamPipeConfig) (id string, err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	resp, err := platformService.AddRouter(routerConfig)

	if err != nil {
		c.Logger.Error("failed to add router", "error", err)
		return "", fmt.Errorf("failed to add router: %s", err.Error())
	}
	return resp.ID, nil
}*/
