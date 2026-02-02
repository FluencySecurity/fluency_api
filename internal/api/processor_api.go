package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	model "github.com/SecurityDo/fluency_api/model"
)

func (c *Client) AddProcessor(name, content, processorType, description string) (err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	if processorType == "" {
		processorType = "fpl_processor" // Default to JavaScript if not specified
	}

	entry := &model.FPLScript{
		Name:        name,
		ScriptText:  content,
		Type:        processorType,
		Description: description,
	}

	err = platformService.AddProcessor(entry)

	if err != nil {
		c.Logger.Error("failed to add processor", "error", err)
		return fmt.Errorf("failed to add processor: %s", err.Error())
	}
	return nil
}

func (c *Client) DeleteProcessor(name string) (err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	err = platformService.DeleteProcessor(name)

	if err != nil {
		c.Logger.Error("failed to delete processor", "name", name, "error", err)
		return fmt.Errorf("failed to delete processor %s: %s", name, err.Error())
	}
	return nil
}

func (c *Client) ListProcessor() (entries []*model.FPLScript, err error) {

	platformService := fluencyAPI.NewPlatformService(c.fluencyClient)

	entries, err = platformService.ListProcessors()

	if err != nil {
		c.Logger.Error("failed to list processor", "error", err)
		return nil, fmt.Errorf("failed to list processor: %s", err.Error())
	}
	return entries, nil
}
