package api

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/SecurityDo/fluency_api/client"
	"github.com/SecurityDo/fluency_api/model"
)

// PlatformService provides helpers for calling platform_* endpoints.
type PlatformService struct {
	client *client.FluencyClient
}

// NewPlatformService constructs a PlatformService instance backed by the provided client.
func NewPlatformService(client *client.FluencyClient) *PlatformService {
	return &PlatformService{client: client}
}

func ApiCall(client *client.FluencyClient, function string, payload interface{}, out interface{}) error {
	res, err := client.GenericCall("api/ds", function, payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling %s: %v\n", function, err.Error())
		return err
	}
	if out == nil {
		return nil
	}
	if res == nil {
		err = fmt.Errorf("empty response from %s", function)
		fmt.Fprintln(os.Stderr, err.Error())
		return err
	}
	if err := json.Unmarshal(res.GetBytes(), out); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s response: %v\n", function, err.Error())
		return err
	}
	return nil
}

func (s *PlatformService) call(function string, payload interface{}, out interface{}) error {
	return ApiCall(s.client, function, payload, out)
}

// Data source and sink configuration structures.

type SourceSetRouterReq struct {
	RouterID     string `json:"routerID"`
	DataSourceID string `json:"dataSourceID"`
}

type RouterAddSourceReq struct {
	RouterID     string `json:"routerID"`
	DataSourceID string `json:"dataSourceID"`
}

type RouterDeleteSourceReq struct {
	RouterID     string `json:"routerID"`
	DataSourceID string `json:"dataSourceID"`
}

type RouterAddPipeReq struct {
	RouterID   string                  `json:"routerID"`
	PipeConfig *model.StreamPipeConfig `json:"pipeConfig"`
}

type RouterDeletePipeReq struct {
	RouterID string `json:"routerID"`
	PipeID   string `json:"pipeID"`
}

type RouterUpdatePipesReq struct {
	RouterID string   `json:"routerID"`
	PipeIDs  []string `json:"pipeIDs"`
}

type PipeUpdateReq struct {
	RouterID   string                  `json:"routerID"`
	PipeConfig *model.StreamPipeConfig `json:"pipeConfig"`
}

type FPLProcessorValidateRequest struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

type FPLProcessorValidateResult struct {
	Console string `json:"console"`
	Error   string `json:"error"`
	OK      bool   `json:"ok"`
}

type FPLProcessorTestRequest struct {
	Name   string `json:"name"`
	Script string `json:"script"`
	Source string `json:"source"`
}

type FPLProcessorTestResult struct {
	Console    string `json:"console"`
	Error      string `json:"error"`
	NewContent string `json:"newContent"`
	Status     string `json:"status"`
}

type EventTailReq struct {
	ID     string `json:"id"`
	Status string `json:"status,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type EventTailResponse struct {
	Entries []string `json:"entries"`
}

type ProcessorTailReq struct {
	PipeID        string `json:"pipeID"`
	ProcessorName string `json:"processorName,omitempty"`
	WorkerIndex   int    `json:"workerIndex,omitempty"`
	Limit         int    `json:"limit,omitempty"`
}

type ProcessorTailResponse struct {
	Entries []*model.LogEvent `json:"entries"`
}

type ProcessorPipesReq struct {
	ProcessorName  string   `json:"processorName,omitempty"`
	ProcessorNames []string `json:"processorNames,omitempty"`
}

type PipeInfo struct {
	PipeID      string `json:"pipeID"`
	PipeName    string `json:"pipeName"`
	Router      string `json:"router"`
	WorkerCount int    `json:"workerCount"`
}

type ProcessorPipesResponse struct {
	Pipes   []*PipeInfo       `json:"pipes,omitempty"`
	Entries []*ProcessorPipes `json:"entries,omitempty"`
}

type ProcessorPipes struct {
	ProcessorName string      `json:"processorName"`
	Pipes         []*PipeInfo `json:"pipes"`
}

type ComponentMetricReq struct {
	Component string `json:"component"`
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Interval  string `json:"interval"`
}

type ProcessorMetricReq struct {
	Channel   string `json:"channel,omitempty"`
	PipeID    string `json:"pipeID,omitempty"`
	Processor string `json:"processor"`
	From      string `json:"from"`
	To        string `json:"to"`
	Interval  string `json:"interval"`
}

type GetComponentStateReq struct {
	ID string `json:"id"`
}

type GetComponentInfoReq struct {
	ID string `json:"id"`
}

type GetComponentStateResponse struct {
	State json.RawMessage `json:"state"`
}

type SetComponentTagsReq struct {
	ID   string      `json:"id"`
	Tags []model.Tag `json:"tags"`
}

type PluginTailReq struct {
	ID    string `json:"id"`
	Limit int    `json:"limit,omitempty"`
}

type PluginTailResponse struct {
	Lines []string `json:"lines"`
}

type ListComponentErrorReq struct {
	ID string `json:"id"`
}

type ListComponentErrorResponse struct {
	Errors []*model.PluginNotification `json:"errors"`
}

type SourceReloadReq struct {
	DataSourceID string `json:"dataSourceID"`
}

type PlatformMetricReq struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Interval string `json:"interval"`
}

type PlatformMetric struct {
	Tenant    string    `json:"tenant,omitempty"`
	Name      string    `json:"name,omitempty"`
	ID        string    `json:"id,omitempty"`
	Pipe      string    `json:"pipe,omitempty"`
	Processor string    `json:"processor,omitempty"`
	Unit      string    `json:"unit,omitempty"`
	Slots     []int64   `json:"slots"`
	Values    []float64 `json:"values"`
}

type PlatformMetricsResponse struct {
	Metrics []*PlatformMetric `json:"metrics"`
}

type SimpleSearchOption struct {
	DataType       string        `json:"dataType,omitempty"`
	SearchStr      string        `json:"searchStr"`
	DateFacetField string        `json:"dateFacetField,omitempty"`
	FetchOffset    int           `json:"fetchOffset"`
	FetchLimit     int           `json:"fetchLimit"`
	SortField      string        `json:"sortField,omitempty"`
	SortOrder      string        `json:"sortOrder,omitempty"`
	Facets         *FacetsOption `json:"facets,omitempty"`
	RangeFrom      int64         `json:"range_from,omitempty"`
	RangeTo        int64         `json:"range_to,omitempty"`
	RangeField     string        `json:"range_field,omitempty"`
}

type FacetEntry struct {
	Field string `json:"field"`
	Order string `json:"order"`
	Size  int    `json:"size"`
}

type FilterEntry struct {
	Field string        `json:"field"`
	Terms []interface{} `json:"terms"`
}

type FacetsOption struct {
	Facets         []*FacetEntry  `json:"facets,omitempty"`
	MustFilters    []*FilterEntry `json:"mustFilters,omitempty"`
	MustNotFilters []*FilterEntry `json:"mustNotFilters,omitempty"`
}

type DeviceInfo struct {
	Name     string `json:"name,omitempty"`
	Category string `json:"category,omitempty"`
}

type DeviceEntry struct {
	Name        string      `json:"name,omitempty"`
	IPs         []string    `json:"ips,omitempty"`
	Hostname    string      `json:"hostname,omitempty"`
	Group       string      `json:"group,omitempty"`
	Device      *DeviceInfo `json:"device,omitempty"`
	Description string      `json:"description,omitempty"`
	Exclude     bool        `json:"exclude,omitempty"`
	CreatedOn   time.Time   `json:"createdOn"`
	UpdatedOn   time.Time   `json:"updatedOn"`
	Timestamp   int64       `json:"timestamp"`
}

type DeviceStatesRequest struct {
	Names []string `json:"names,omitempty"`
}

type DeviceState struct {
	Name     string `json:"name,omitempty"`
	LastSeen int64  `json:"lastSeen"`
}

type DeviceStatesResponse struct {
	States []*DeviceState `json:"states,omitempty"`
}

type DeviceMetricsRequest struct {
	Names []string `json:"names,omitempty"`
}

type ImportStatStat struct {
	ImportSource string    `json:"importSource"`
	Slots        []int64   `json:"slots"`
	Values       []float64 `json:"values"`
}

type DeviceMetricsResponse struct {
	Metrics []*ImportStatStat `json:"metrics,omitempty"`
}

/*
type Integration struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Integration string          `json:"integration"`
	Description string          `json:"description"`
	CreatedOn   time.Time       `json:"createdOn"`
	UpdatedOn   time.Time       `json:"updatedOn"`
	Secret      json.RawMessage `json:"secret"`
	Config      json.RawMessage `json:"config"`
}

type GenericDAORequestArgs[T any] struct {
	Id    string `json:"id,omitempty"`
	Entry *T     `json:"entry,omitempty"`
	Flag  bool   `json:"flag,omitempty"`
}

type GenericDAORequest[T any] struct {
	Action string                    `json:"action"`
	Args   *GenericDAORequestArgs[T] `json:"args,omitempty"`
}

type GenericDaoAddResponse struct {
	ID string                    `json:"id"`
}*/

type ClientSecret struct {
	ClientSecret string `json:"clientSecret"`
	PrivateKey   string `json:"privateKey"`
}

type Office365AuditConfig struct {
	TenantID   string `json:"tenantID"`
	ClientID   string `json:"clientId"`
	Thumbprint string `json:"thumbprint"`
}

type SlackConfig struct {
	Token string `json:"token"`
}

type PagerDutyConfig struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

type PagerDutySecret struct {
	IntegrationKey string `json:"integrationKey"`
}

type S3BucketConfig struct {
	Prefix    string `json:"prefix"`
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
	Mode      string `json:"mode"`
	AWSUser   string `json:"awsUser"`
	AccessKey string `json:"accessKey"`
}

type AWSUserSecret struct {
	SecretKey string `json:"secretKey"`
}

type AWSUserConfig struct {
	AccessKey string `json:"accessKey"`
}

type DuoConfig struct {
	IntegrationKey string `json:"integrationKey"`
	APIHostname    string `json:"apiHostname"`
}

type DuoSecret struct {
	SecretKey string `json:"secretKey"`
}

type Endpoint struct {
	AuthURL  string `json:"authURL"`
	TokenURL string `json:"tokenURL"`
}

type OAuth2Config struct {
	ClientID    string   `json:"clientID"`
	Endpoint    Endpoint `json:"endpoint"`
	RedirectURL string   `json:"redirectURL"`
	Scopes      []string `json:"scopes"`
}

type GSuiteConfig struct {
	Config *OAuth2Config `json:"config"`
}

type FalconConfig struct {
	ClientID string `json:"clientID"`
	BaseURL  string `json:"baseURL"`
}

type FalconSecret struct {
	ClientSecret string `json:"clientSecret"`
}

type MimecastConfig struct {
	ApplicationID  string `json:"applicationID"`
	ApplicationKey string `json:"applicationKey"`
	AccessKey      string `json:"accessKey"`
	BaseURL        string `json:"baseURL"`
}

type MimecastSecret struct {
	SecretKey string `json:"secretKey"`
}

type ListPluginsResponse struct {
	Plugins []string `json:"plugins"`
}

type ListConfigsResponse struct {
	Sources      []*model.DataSourceConfig    `json:"sources"`
	Sinks        []*model.DataSinkConfig      `json:"sinks"`
	Routers      []*model.RouterConfig        `json:"routers"`
	Pipes        []*model.StreamPipeConfig    `json:"pipes"`
	Metrics      []*ComponentStat             `json:"metrics"`
	Channels     []*model.ChannelConfig       `json:"channels"`
	Connections  []*model.RouterInput         `json:"connections"`
	Integrations []*model.Integration         `json:"integrations"`
	Errors       []*model.PluginNotification  `json:"errors"`
	ErrorStates  []*model.ComponentErrorState `json:"errorStates"`
}

type ComponentStat struct {
	Tenant    string `json:"tenant,omitempty"`
	Action    string `json:"action"`
	Namespace string `json:"namespace"`
	Component string `json:"component"`
	ID        string `json:"id"`
	Unit      string `json:"unit"`

	Application string `json:"application,omitempty"`
	AppInstance string `json:"appInstance,omitempty"`

	Slots  []int64   `json:"slots"`
	Values []float64 `json:"values"`
	// Values    []model.SamplePair `json:"values"`
}

type DataSourceEntryResponse struct {
	Entry *model.DataSourceConfig `json:"entry"`
}

type AddDataSourceResponse struct {
	ID     string          `json:"id"`
	Secret json.RawMessage `json:"secret,omitempty"`
	URL    string          `json:"url,omitempty"`
}
type ListDataSourceResponse struct {
	Entries []*model.DataSourceConfig `json:"entries"`
}

type DataSinkEntryResponse struct {
	Entry *model.DataSinkConfig `json:"entry"`
}

type AddDataSinkResponse struct {
	ID string `json:"id"`
}

type AddIntegrationResponse struct {
	ID string `json:"id"`
}

type ListDataSinkResponse struct {
	Entries []*model.DataSinkConfig `json:"entries"`
}

type RouterEntryResponse struct {
	Entry *model.RouterConfig       `json:"entry"`
	Pipes []*model.StreamPipeConfig `json:"pipes"`
}

type AddRouterResponse struct {
	ID string `json:"id"`
}

type ChannelEntryResponse struct {
	Entry *model.ChannelConfig `json:"entry"`
}

type AddChannelResponse struct {
	ID string `json:"id"`
}

type ProcessorEntryResponse struct {
	Entry *model.FPLScript `json:"entry"`
}

type ListProcessorsResponse struct {
	Entries []*model.FPLScript `json:"entries"`
}

type GenericOKResponse struct {
}

type ImportDeviceSearchRequest struct {
	Options *SimpleSearchOption `json:"options,omitempty"`
}

type ImportDeviceSearchResponse struct {
	Results json.RawMessage `json:"results"`
}

type DeviceEntryResponse struct {
	Entry *DeviceEntry `json:"entry"`
}

type AddDeviceResponse struct {
	ID string `json:"id"`
}

type ImportDeviceDAORequestArgs struct {
	Id    string       `json:"id,omitempty"`
	Entry *DeviceEntry `json:"entry,omitempty"`
}

type ImportDeviceDAORequest struct {
	Action string                      `json:"action"`
	Args   *ImportDeviceDAORequestArgs `json:"args,omitempty"`
}

type IntegrationDAORequestArgs struct {
	Id    string             `json:"id,omitempty"`
	Entry *model.Integration `json:"entry,omitempty"`
}

type IntegrationDAORequest struct {
	Action string                     `json:"action"`
	Args   *IntegrationDAORequestArgs `json:"args,omitempty"`
}

type IntegrationDAOResponse struct {
	Entry   *model.Integration   `json:"entry,omitempty"`
	Entries []*model.Integration `json:"entries,omitempty"`
}

type RouterAddPipeResponse struct {
	ID string `json:"id"`
}

type GenericDaoRequestArgs[T any] struct {
	Id    string `json:"id,omitempty"`
	Entry T      `json:"entry,omitempty"`
	Flag  bool   `json:"flag"`
}

type GenericDaoRequest[T any] struct {
	Action string                    `json:"action"`
	Args   *GenericDaoRequestArgs[T] `json:"args"`
}

type GenericDAORequestArgs[T any] struct {
	Id    string `json:"id,omitempty"`
	Entry *T     `json:"entry,omitempty"`
	Flag  bool   `json:"flag,omitempty"`
}

type GenericDAORequest[T any] struct {
	Action string                    `json:"action"`
	Args   *GenericDAORequestArgs[T] `json:"args,omitempty"`
}

type GenericDaoAddResponse struct {
	ID string `json:"id"`
}

type GenericDaoListResponse[T any] struct {
	Entries []*T `json:"entries"`
}

func (s *PlatformService) GetPodRole() (role, arn string, err error) {
	var resp struct {
		Role string `json:"role"`
		ARN  string `json:"arn"`
	}
	if err := s.call("get_pod_role", nil, &resp); err != nil {
		return "", "", err
	}
	return resp.Role, resp.ARN, nil

}

func (s *PlatformService) TestAssumedRole(roleARN string, externalID string) (err error) {
	request := &model.InstanceRoleTestRequest{
		RoleARN:    roleARN,
		ExternalID: externalID,
	}

	if err := s.call("platform_instancerole_test", request, nil); err != nil {
		return err
	}
	return nil

}

func (s *PlatformService) AddAssumedRole(name string, roleARN string, externalID string) (id string, err error) {
	request := &GenericDAORequest[model.InstanceRole]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.InstanceRole]{
			//Id:    args.Id,
			Entry: &model.InstanceRole{
				DisplayName: name,
				RoleARN:     roleARN,
				ExternalID:  externalID,
			},
		},
	}
	var resp GenericDaoAddResponse
	if err := s.call("platform_instancerole_dao", request, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil

}

func (s *PlatformService) DeleteAssumedRole(roleID string) (err error) {
	request := &GenericDAORequest[model.InstanceRole]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.InstanceRole]{
			Id: roleID,
		},
	}
	//var resp GenericDaoAddResponse
	if err := s.call("platform_instancerole_dao", request, nil); err != nil {
		return err
	}
	return nil

}

func (s *PlatformService) ListAssumedRole() (roles []*model.InstanceRole, err error) {
	request := &GenericDAORequest[model.InstanceRole]{
		Action: "list",
	}
	var resp GenericDaoListResponse[model.InstanceRole]
	if err := s.call("platform_instancerole_dao", request, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil

}

// ListPlugins fetches the available platform plugins.
func (s *PlatformService) ListPlugins() ([]string, error) {
	var resp ListPluginsResponse
	if err := s.call("platform_list_plugins", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Plugins, nil
}

// ListConfigs returns the current platform configuration snapshot.
func (s *PlatformService) ListConfigs() (*ListConfigsResponse, error) {
	var resp ListConfigsResponse
	if err := s.call("platform_list_configs", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDataSource retrieves a data source by id.
func (s *PlatformService) GetDataSource(id string) (*model.DataSourceConfig, error) {
	req := &GenericDAORequest[model.DataSourceConfig]{
		Action: "get",
		Args: &GenericDAORequestArgs[model.DataSourceConfig]{
			Id: id,
		},
	}
	var resp DataSourceEntryResponse
	if err := s.call("platform_datasource_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// AddDataSource creates a new data source entry.
func (s *PlatformService) AddDataSource(entry *model.DataSourceConfig) (*AddDataSourceResponse, error) {
	req := &GenericDAORequest[model.DataSourceConfig]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.DataSourceConfig]{
			Entry: entry,
		},
	}
	var resp AddDataSourceResponse
	if err := s.call("platform_datasource_dao", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDataSource removes a data source by id.
func (s *PlatformService) DeleteDataSource(id string) error {
	req := &GenericDAORequest[model.DataSourceConfig]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.DataSourceConfig]{
			Id: id,
		},
	}
	return s.call("platform_datasource_dao", req, nil)
}

func (s *PlatformService) ListDataSource() (entries []*model.DataSourceConfig, err error) {
	req := &GenericDAORequest[model.DataSourceConfig]{
		Action: "list",
	}

	var resp ListDataSourceResponse
	if err := s.call("platform_datasource_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// GetDataSink retrieves a data sink by id.
func (s *PlatformService) GetDataSink(id string) (*model.DataSinkConfig, error) {
	req := &GenericDAORequest[model.DataSinkConfig]{
		Action: "get",
		Args: &GenericDAORequestArgs[model.DataSinkConfig]{
			Id: id,
		},
	}
	var resp DataSinkEntryResponse
	if err := s.call("platform_datasink_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// AddDataSink creates a new data sink entry.
func (s *PlatformService) AddDataSink(entry *model.DataSinkConfig) (*AddDataSinkResponse, error) {
	req := &GenericDAORequest[model.DataSinkConfig]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.DataSinkConfig]{
			Entry: entry,
		},
	}
	var resp AddDataSinkResponse
	if err := s.call("platform_datasink_dao", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDataSink removes a data sink by id.
func (s *PlatformService) DeleteDataSink(id string) error {
	req := &GenericDAORequest[model.DataSinkConfig]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.DataSinkConfig]{
			Id: id,
		},
	}
	return s.call("platform_datasink_dao", req, nil)
}

// ListDataSinks lists configured data sinks.
func (s *PlatformService) ListDataSink() ([]*model.DataSinkConfig, error) {
	req := &GenericDAORequest[model.DataSinkConfig]{
		Action: "list",
	}
	var resp ListDataSinkResponse
	if err := s.call("platform_datasink_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// GetRouter loads a router and its pipes by id.
func (s *PlatformService) GetRouter(id string) (*RouterEntryResponse, error) {
	req := &GenericDAORequest[model.RouterConfig]{
		Action: "get",
		Args: &GenericDAORequestArgs[model.RouterConfig]{
			Id: id,
		},
	}
	var resp RouterEntryResponse
	if err := s.call("platform_router_dao", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type AddSimpleRouterReq struct {
	Processor string `json:"processor,omitempty"`
	Router    string `json:"router"`
	Pipe      string `json:"pipe,omitempty"`
	// multiple pipes
	// Pipes      []string `json:"pipes,omitempty"`
	// Processors []string `json:"processors,omitempty"`
}

func (s *PlatformService) AddSimpleRouter(processorName, routerName string) (string, error) {
	if routerName == "" {
		routerName = processorName
	}
	req := &AddSimpleRouterReq{
		Processor: processorName,
		Router:    routerName,
		Pipe:      processorName,
	}
	var resp AddRouterResponse
	if err := s.call("platform_add_simple_router", req, &resp); err != nil {
		return "", err
	}
	return resp.ID, nil
}

// AddRouter registers a new router.
func (s *PlatformService) AddRouter(entry *model.RouterConfig) (*AddRouterResponse, error) {
	req := &GenericDAORequest[model.RouterConfig]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.RouterConfig]{
			Entry: entry,
		},
	}
	var resp AddRouterResponse
	if err := s.call("platform_router_dao", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteRouter removes a router by id.
func (s *PlatformService) DeleteRouter(id string) error {
	req := &GenericDAORequest[model.RouterConfig]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.RouterConfig]{
			Id: id,
		},
	}
	return s.call("platform_router_dao", req, nil)
}

// GetChannel returns a channel configuration by id.
func (s *PlatformService) GetChannel(id string) (*model.ChannelConfig, error) {
	req := &GenericDAORequest[model.ChannelConfig]{
		Action: "get",
		Args: &GenericDAORequestArgs[model.ChannelConfig]{
			Id: id,
		},
	}
	var resp ChannelEntryResponse
	if err := s.call("platform_channel_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// AddChannel creates a new channel entry.
func (s *PlatformService) AddChannel(entry *model.ChannelConfig) (*AddChannelResponse, error) {
	req := &GenericDAORequest[model.ChannelConfig]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.ChannelConfig]{
			Entry: entry,
		},
	}
	var resp AddChannelResponse
	if err := s.call("platform_channel_dao", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteChannel removes the specified channel configuration.
func (s *PlatformService) DeleteChannel(id string) error {
	req := &GenericDAORequest[model.ChannelConfig]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.ChannelConfig]{
			Id: id,
		},
	}
	return s.call("platform_channel_dao", req, nil)
}

// ListProcessors returns all configured processors.
func (s *PlatformService) ListProcessors() ([]*model.FPLScript, error) {
	req := &GenericDAORequest[model.FPLScript]{
		Action: "list",
	}
	var resp ListProcessorsResponse
	if err := s.call("platform_processor_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// GetProcessor returns a processor definition by name.
func (s *PlatformService) GetProcessor(name string) (*model.FPLScript, error) {
	req := &GenericDAORequest[*model.FPLScript]{
		Action: "get",
		Args: &GenericDAORequestArgs[*model.FPLScript]{
			Id: name,
		},
	}
	var resp ProcessorEntryResponse
	if err := s.call("platform_processor_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// AddProcessor registers a new processor definition.
func (s *PlatformService) AddProcessor(entry *model.FPLScript) error {
	req := &GenericDAORequest[model.FPLScript]{
		Action: "add",
		Args: &GenericDAORequestArgs[model.FPLScript]{
			Entry: entry,
		},
	}
	return s.call("platform_processor_dao", req, nil)
}

// UpdateProcessor modifies an existing processor definition.
func (s *PlatformService) UpdateProcessor(entry *model.FPLScript) error {
	req := &GenericDAORequest[model.FPLScript]{
		Action: "update",
		Args: &GenericDAORequestArgs[model.FPLScript]{
			Entry: entry,
		},
	}
	return s.call("platform_processor_dao", req, nil)
}

// DeleteProcessor removes a processor by name.
func (s *PlatformService) DeleteProcessor(name string) error {
	req := &GenericDAORequest[model.FPLScript]{
		Action: "delete",
		Args: &GenericDAORequestArgs[model.FPLScript]{
			Id: name,
		},
	}
	return s.call("platform_processor_dao", req, nil)
}

// ValidateProcessor compiles and validates a processor script.
func (s *PlatformService) ValidateProcessor(req *FPLProcessorValidateRequest) (*FPLProcessorValidateResult, error) {
	var resp FPLProcessorValidateResult
	if err := s.call("platform_processor_validate", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TestProcessor executes a processor script against sample data.
func (s *PlatformService) TestProcessor(req *FPLProcessorTestRequest) (*FPLProcessorTestResult, error) {
	var resp FPLProcessorTestResult
	if err := s.call("platform_processor_test", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetDataSourceRouter assigns (or removes) a router from a data source.
func (s *PlatformService) SetDataSourceRouter(req *SourceSetRouterReq) error {
	return s.call("platform_source_set_router", req, nil)
}

/*
// AddRouterSource attaches a data source to a router. Deprecated in the backend but kept for compatibility.
func (s *PlatformService) AddRouterSource(req *RouterAddSourceReq) error {
	return s.call("platform_router_add_source", req, nil)
}

// DeleteRouterSource detaches a data source from a router. Deprecated in the backend but kept for compatibility.
func (s *PlatformService) DeleteRouterSource(req *RouterDeleteSourceReq) error {
	return s.call("platform_router_delete_source", req, nil)
}*/

// AddRouterPipe creates a new pipe under a router.
func (s *PlatformService) AddRouterPipe(req *RouterAddPipeReq) (*RouterAddPipeResponse, error) {
	var resp RouterAddPipeResponse
	if err := s.call("platform_router_add_pipe", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteRouterPipe removes a pipe from a router.
func (s *PlatformService) DeleteRouterPipe(req *RouterDeletePipeReq) error {
	return s.call("platform_router_delete_pipe", req, nil)
}

// UpdateRouterPipes replaces the set of pipes attached to a router.
func (s *PlatformService) UpdateRouterPipes(req *RouterUpdatePipesReq) error {
	return s.call("platform_router_update_pipes", req, nil)
}

// UpdatePipe updates a pipe configuration on a router.
func (s *PlatformService) UpdatePipe(req *PipeUpdateReq) error {
	return s.call("platform_pipe_update", req, nil)
}

// EventTail retrieves recent events for a component.
func (s *PlatformService) EventTail(req *EventTailReq) (*EventTailResponse, error) {
	var resp EventTailResponse
	if err := s.call("platform_event_tail", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ProcessorTail fetches processor trace logs.
func (s *PlatformService) ProcessorTail(req *ProcessorTailReq) (*ProcessorTailResponse, error) {
	var resp ProcessorTailResponse
	if err := s.call("platform_processor_tail", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ProcessorPipes lists pipes attached to a processor.
func (s *PlatformService) ProcessorPipes(req *ProcessorPipesReq) (*ProcessorPipesResponse, error) {
	var resp ProcessorPipesResponse
	if err := s.call("platform_processor_pipes", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ComponentMetrics retrieves metrics for a given component.
func (s *PlatformService) ComponentMetrics(req *ComponentMetricReq) (*PlatformMetricsResponse, error) {
	var resp PlatformMetricsResponse
	if err := s.call("platform_component_metrics", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ProcessorMetrics fetches metrics scoped to a processor.
func (s *PlatformService) ProcessorMetrics(req *ProcessorMetricReq) (*PlatformMetricsResponse, error) {
	var resp PlatformMetricsResponse
	if err := s.call("platform_processor_metrics", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetComponentState returns the state payload for a component.
func (s *PlatformService) GetComponentState(id string) (*GetComponentStateResponse, error) {
	req := &GetComponentStateReq{ID: id}
	var resp GetComponentStateResponse
	if err := s.call("platform_get_component_state", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetComponentTags updates component tags.
func (s *PlatformService) SetComponentTags(req *SetComponentTagsReq) error {
	return s.call("platform_set_component_tags", req, nil)
}

// PluginTail retrieves tail logs for a plugin data source.
func (s *PlatformService) PluginTail(req *PluginTailReq) (*PluginTailResponse, error) {
	var resp PluginTailResponse
	if err := s.call("platform_plugin_tail", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListComponentErrors returns alerts for a component id.
func (s *PlatformService) ListComponentErrors(id string) (*ListComponentErrorResponse, error) {
	req := &ListComponentErrorReq{ID: id}
	var resp ListComponentErrorResponse
	if err := s.call("platform_list_component_errors", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ClearComponentError acknowledges component errors for the given id.
func (s *PlatformService) ClearComponentError(id string) error {
	req := &ListComponentErrorReq{ID: id}
	return s.call("platform_clear_component_error", req, nil)
}

// GetComponentInfo loads detailed component information including logs and alerts.
func (s *PlatformService) GetComponentInfo(id string) (*model.ComponentInfo, error) {
	req := &GetComponentInfoReq{ID: id}
	var resp model.ComponentInfo
	if err := s.call("platform_get_component_info", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SourceReload triggers a data source reload operation.
func (s *PlatformService) SourceReload(id string) error {
	req := &SourceReloadReq{DataSourceID: id}
	return s.call("platform_source_reload", req, nil)
}

/*
// ProfileTotal retrieves aggregated platform CPU and memory metrics.
func (s *PlatformService) ProfileTotal(req *PlatformMetricReq) (*PlatformMetricsResponse, error) {
	var resp PlatformMetricsResponse
	if err := s.call("platform_profile_total", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ProfileComponent gathers metrics for specific components or processors.
func (s *PlatformService) ProfileComponent(req *PlatformMetricReq) (*PlatformMetricsResponse, error) {
	var resp PlatformMetricsResponse
	if err := s.call("platform_profile_component", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ImportDeviceSearch performs a facet search over imported devices.
func (s *PlatformService) ImportDeviceSearch(req *ImportDeviceSearchRequest) (json.RawMessage, error) {
	var raw json.RawMessage
	if err := s.call("platform_import_device_search", req, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// GetImportDevice fetches a device entry by id.
func (s *PlatformService) GetImportDevice(id string) (*DeviceEntry, error) {
	req := &GenericDAORequest[DeviceEntry]{
		Action: "get",
		Args: &GenericDAORequestArgs[DeviceEntry]{
			Id: id,
		},
	}
	var resp DeviceEntryResponse
	if err := s.call("platform_import_device_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// AddImportDevice registers a new device.
func (s *PlatformService) AddImportDevice(entry *DeviceEntry) (*AddDeviceResponse, error) {
	req := &GenericDAORequest[DeviceEntry]{
		Action: "add",
		Args: &GenericDAORequestArgs[DeviceEntry]{
			Entry: entry,
		},
	}
	var resp AddDeviceResponse
	if err := s.call("platform_import_device_dao", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateImportDevice updates an existing device entry.
func (s *PlatformService) UpdateImportDevice(entry *DeviceEntry) error {
	req := &GenericDAORequest[DeviceEntry]{
		Action: "update",
		Args: &GenericDAORequestArgs[DeviceEntry]{
			Entry: entry,
		},
	}
	return s.call("platform_import_device_dao", req, nil)
}

// DeleteImportDevice removes a device by id.
func (s *PlatformService) DeleteImportDevice(id string) error {
	req := &GenericDAORequest[DeviceEntry]{
		Action: "delete",
		Args: &GenericDAORequestArgs[DeviceEntry]{
			Id: id,
		},
	}
	return s.call("platform_import_device_dao", req, nil)
}

// GetImportDeviceStates loads last-seen state for devices.
func (s *PlatformService) GetImportDeviceStates(req *DeviceStatesRequest) (*DeviceStatesResponse, error) {
	var resp DeviceStatesResponse
	if err := s.call("platform_import_device_states", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetImportDeviceMetrics fetches metrics for import devices.
func (s *PlatformService) GetImportDeviceMetrics(req *DeviceMetricsRequest) (*DeviceMetricsResponse, error) {
	var resp DeviceMetricsResponse
	if err := s.call("platform_import_device_metrics", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}*/

// ListIntegrations lists all platform integrations.
func (s *PlatformService) ListIntegrations() ([]*model.Integration, error) {
	req := &IntegrationDAORequest{Action: "list"}
	var resp IntegrationDAOResponse
	if err := s.call("platform_integration_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// GetIntegration fetches an integration by id.
func (s *PlatformService) GetIntegration(id string) (*model.Integration, error) {
	req := &IntegrationDAORequest{
		Action: "get",
		Args:   &IntegrationDAORequestArgs{Id: id},
	}
	var resp IntegrationDAOResponse
	if err := s.call("platform_integration_dao", req, &resp); err != nil {
		return nil, err
	}
	return resp.Entry, nil
}

// AddIntegration registers a new integration entry.
func (s *PlatformService) AddIntegration(entry *model.Integration) (id string, err error) {
	req := &IntegrationDAORequest{
		Action: "add",
		Args:   &IntegrationDAORequestArgs{Entry: entry},
	}
	var resp AddIntegrationResponse
	err = s.call("platform_integration_dao", req, &resp)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// UpdateIntegration updates an existing integration.
func (s *PlatformService) UpdateIntegration(entry *model.Integration) error {
	req := &IntegrationDAORequest{
		Action: "update",
		Args:   &IntegrationDAORequestArgs{Entry: entry},
	}
	return s.call("platform_integration_dao", req, nil)
}

// DeleteIntegration removes an integration by id.
func (s *PlatformService) DeleteIntegration(id string) error {
	req := &IntegrationDAORequest{
		Action: "delete",
		Args:   &IntegrationDAORequestArgs{Id: id},
	}
	return s.call("platform_integration_dao", req, nil)
}
