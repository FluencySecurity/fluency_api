package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	"github.com/SecurityDo/fluency_api/model"
)

func (c *Client) ListDatalakes() (entries []*model.Datalake, err error) {

	datalakeService := fluencyAPI.NewDatalakeService(c.fluencyClient)

	resp, err := datalakeService.ListDatalake()

	if err != nil {
		c.Logger.Error("failed to list datalakes", "error", err)
		return nil, fmt.Errorf("failed to list datalakes: %w", err)
	}
	return resp, nil
}

func (c *Client) AddDatalake(name string, managed bool, integrationID string) (err error) {

	datalakeService := fluencyAPI.NewDatalakeService(c.fluencyClient)

	err = datalakeService.AddDatalake(name, managed, integrationID)

	if err != nil {
		c.Logger.Error("failed to add datalakes", "error", err)
		return fmt.Errorf("failed to add datalakes: %w", err)
	}
	return nil
}

func (c *Client) AddDatalakeIndex(lake, index string, schema string) (err error) {

	datalakeService := fluencyAPI.NewDatalakeService(c.fluencyClient)

	err = datalakeService.AddDatalakeIndex(lake, index, schema)

	if err != nil {
		c.Logger.Error("failed to add datalake index", "error", err)
		return fmt.Errorf("failed to add datalake index: %w", err)
	}
	return nil
}

func (c *Client) ListDatalakeIndex(lake string) (entries []*model.DatalakeIndex, err error) {

	datalakeService := fluencyAPI.NewDatalakeService(c.fluencyClient)

	entries, err = datalakeService.ListDatalakeIndex(lake)

	if err != nil {
		c.Logger.Error("failed to list datalake index", "error", err)
		return nil, fmt.Errorf("failed to list datalake index: %w", err)
	}
	return entries, nil
}

func (c *Client) DeleteDatalakeIndex(lake, index string) (err error) {

	datalakeService := fluencyAPI.NewDatalakeService(c.fluencyClient)

	err = datalakeService.DeleteDatalakeIndex(lake, index)

	if err != nil {
		c.Logger.Error("failed to delete datalake index", "error", err)
		return fmt.Errorf("failed to delete datalake index: %w", err)
	}
	return nil
}
