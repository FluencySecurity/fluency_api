package api

import (
	"fmt"

	"github.com/SecurityDo/fluency_api/api"
	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	"github.com/SecurityDo/fluency_api/model"
)

func (c *Client) ListAppTemplates() (templates []*model.ApplicationTemplateConfig, err error) {

	applicationService := fluencyAPI.NewApplicationService(c.fluencyClient)

	resp, err := applicationService.ListAppTemplates()

	if err != nil {
		c.Logger.Error("failed to list application templates", "error", err)
		return nil, fmt.Errorf("failed to list application templates: %w", err)
	}
	return resp.Entries, nil
}

func (c *Client) InstallAppInstance(application, instance string, displayName string, parameters []*model.InputParameter) (err error) {

	applicationService := fluencyAPI.NewApplicationService(c.fluencyClient)

	req := &api.InstallAppInstanceRequest{
		Config: &model.InstanceConfig{
			Application:     application,
			Instance:        instance,
			DisplayName:     displayName,
			InputParameters: parameters,
		},
	}

	err = applicationService.InstallAppTemplate(req)

	if err != nil {
		c.Logger.Error("failed to install application instance", "error", err)
		return fmt.Errorf("failed to install application instance: %w", err)
	}

	return nil
}

func (c *Client) GetAppInstance(application, instance string) (res *api.GetAppInstanceResponse, err error) {

	applicationService := fluencyAPI.NewApplicationService(c.fluencyClient)
	res, err = applicationService.GetAppInstance(application, instance)

	if err != nil {
		c.Logger.Error("failed to get application instance", "error", err)
		return nil, fmt.Errorf("failed to get application instance: %w", err)
	}

	return res, nil
}

func (c *Client) UnInstallAppInstance(application, instance string) (err error) {

	applicationService := fluencyAPI.NewApplicationService(c.fluencyClient)

	req := &api.UnInstallAppInstanceRequest{
		Application: application,
		Instance:    instance,
	}

	err = applicationService.UnInstallAppTemplate(req)

	if err != nil {
		c.Logger.Error("failed to uninstall application instance", "error", err)
		return fmt.Errorf("failed to uninstall application instance: %w", err)
	}
	return nil
}

func (c *Client) AddTemplate(content string) (id string, err error) {

	applicationService := fluencyAPI.NewApplicationService(c.fluencyClient)

	id, err = applicationService.AddAppTemplate(content)

	if err != nil {
		c.Logger.Error("failed to add application template", "error", err)
		return "", fmt.Errorf("failed to add application template: %w", err)
	}
	return id, nil
}

func (c *Client) DeleteTemplate(name string) (err error) {

	applicationService := fluencyAPI.NewApplicationService(c.fluencyClient)

	err = applicationService.DeleteAppTemplate(name)

	if err != nil {
		c.Logger.Error("failed to delete application template", "error", err)
		return fmt.Errorf("failed to delete application template %s: %w", name, err)
	}
	return nil
}

func (c *Client) UpdateTemplate(name string, content string) (err error) {

	applicationService := fluencyAPI.NewApplicationService(c.fluencyClient)

	err = applicationService.UpdateAppTemplate(name, content)

	if err != nil {
		c.Logger.Error("failed to update application template", "error", err)
		return fmt.Errorf("failed to update application template %s: %w", name, err)
	}
	return nil
}
