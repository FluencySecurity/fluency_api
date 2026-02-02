package api

import (
	"github.com/SecurityDo/fluency_api/client"
	"github.com/SecurityDo/fluency_api/model"
)

// PlatformService provides helpers for calling platform_* endpoints.
type ApplicationService struct {
	client *client.FluencyClient
}

// NewPlatformService constructs a PlatformService instance backed by the provided client.
func NewApplicationService(client *client.FluencyClient) *ApplicationService {
	return &ApplicationService{client: client}
}

func (s *ApplicationService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

func (s *ApplicationService) ListAppTemplates() (*ListAppTemplateResponse, error) {
	var resp ListAppTemplateResponse
	if err := s.call("platform_list_application_template", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type ListAppTemplateResponse struct {
	Entries []*model.ApplicationTemplateConfig `json:"entries"`
}

type ListAppInstanceResponse struct {
	Entries []*model.InstanceState `json:"entries"`
}

func (s *ApplicationService) ListAppInstances() (*ListAppInstanceResponse, error) {
	var resp ListAppInstanceResponse
	if err := s.call("platform_list_application_template", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type InstallAppInstanceRequest struct {
	Config *model.InstanceConfig `json:"config"`
}

func (s *ApplicationService) InstallAppTemplate(req *InstallAppInstanceRequest) error {
	//var resp ListAppTemplateResponse
	if err := s.call("platform_install_application_instance", req, nil); err != nil {
		return err
	}
	return nil
}

type UnInstallAppInstanceRequest struct {
	Application string `json:"application,omitempty" yaml:"application,omitempty"`
	Instance    string `json:"instance,omitempty" yaml:"instance,omitempty"`
}

func (s *ApplicationService) UnInstallAppTemplate(req *UnInstallAppInstanceRequest) error {
	//var resp ListAppTemplateResponse
	if err := s.call("platform_uninstall_application_instance", req, nil); err != nil {
		return err
	}
	return nil
}

func (s *ApplicationService) AddAppTemplate(content string) (id string, err error) {
	request := &GenericDAORequest[model.AppTemplateEntry]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.AppTemplateEntry]{
			//Id:    args.Id,
			Entry: &model.AppTemplateEntry{Content: content},
		},
	}
	var resp GenericDaoAddResponse
	if err := s.call("platform_application_template_dao", request, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (s *ApplicationService) DeleteAppTemplate(name string) (err error) {
	request := &GenericDAORequest[model.AppTemplateEntry]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.AppTemplateEntry]{
			Id: name,
		},
	}
	//var resp GenericDaoAddResponse
	if err := s.call("platform_application_template_dao", request, nil); err != nil {
		return err
	}
	return nil

}

func (s *ApplicationService) UpdateAppTemplate(name string, content string) (err error) {
	request := &GenericDAORequest[model.AppTemplateEntry]{
		Action: "update",
		Args: &GenericDAORequestArgs[model.AppTemplateEntry]{
			Id: name,
			Entry: &model.AppTemplateEntry{
				Content: content,
			},
		},
	}
	//var resp GenericDaoAddResponse
	if err := s.call("platform_application_template_dao", request, nil); err != nil {
		return err
	}
	return nil

}

type GetAppInstanceRequest struct {
	Application string `json:"application,omitempty" yaml:"application,omitempty"`
	Instance    string `json:"instance,omitempty" yaml:"instance,omitempty"`
}

type GetAppInstanceResponse struct {
	//State   *applicationAPI.InstanceState       `json:"state"`
	//Actions []*applicationAPI.AppInstanceAction `json:"actions"`
	Outputs []*model.InputParameter `json:"outputs,omitempty"`
}

func (s *ApplicationService) GetAppInstance(app string, instance string) (res *GetAppInstanceResponse, err error) {
	request := &GetAppInstanceRequest{
		Application: app,
		Instance:    instance,
	}
	var resp GetAppInstanceResponse
	if err := s.call("platform_get_application_instance", request, &resp); err != nil {
		return nil, err
	}
	return &resp, nil

}
