package model

import (
	"encoding/json"
	"time"

	"golang.org/x/oauth2"
)

const (
	PLUGIN_SERVICE_ENABLE  = "enable"
	PLUGIN_SERVICE_DISABLE = "disable"
	PLUGIN_SERVICE_STOP    = "stop"
	PLUGIN_SERVICE_UPDATE  = "update"
	PLUGIN_SERVICE_INIT    = "init"

	PLUGIN_EVENT_SAVESTATE       = "saveState"
	PLUGIN_EVENT_EVENTDUMP       = "eventDump"
	PLUGIN_EVENT_TRANSLATIONDUMP = "translationDump"
	PLUGIN_EVENT_RESOURCEDUMP    = "resourceDump"
	PLUGIN_EVENT_ELASTICDUMP     = "elasticDump"
	PLUGIN_EVENT_NOTIFICATION    = "notification"

	INTEGRATION_CATO        = "Cato"
	INTEGRATION_RESTAPI     = "RESTAPI"
	INTEGRATION_SALESFORCE  = "Salesforce"
	INTEGRATION_MIMECAST    = "Mimecast"
	INTEGRATION_PROOFPOINT  = "Proofpoint"
	INTEGRATION_FALCON      = "Falcon"
	INTEGRATION_OKTA        = "Okta"
	INTEGRATION_SOPHOS      = "Sophos"
	INTEGRATION_Darktrace   = "Darktrace"
	INTEGRATION_SentinelOne = "SentinelOne"
	INTEGRATION_LDAP        = "LDAP"
	INTEGRATION_DataLake    = "DataLake"

	INTEGRATION_AWS_KINESIS_STREAM = "KinesisStream"
	INTEGRATION_AWS_S3BUCKET       = "S3Bucket"
	INTEGRATION_AWS_GUARDDUTY      = "AWSGuardDuty"
	INTEGRATION_AWS_LAMBDA         = "AWSLambda"
	INTEGRATION_AWS_FIREHOSE       = "FirehoseStream"
	INTEGRATION_AWS_EC2Audit       = "AWSEC2Audit"
	INTEGRATION_AWS_S3Notification = "S3Notification"

	INTEGRATION_AZURE_BLOBSTORAGE = "AzureBlobStorage"

	INTEGRATION_PROM_PUSH = "PROMPush"

	INTEGRATION_AWS_API  = "AWSAPI"
	INTEGRATION_AWS_User = "AWSUser"

// integration "FirehoseStream"
)

// S3Bucket,  AWSUser, PagerDuty, Slack, Office365Audit
type Integration struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Integration string    `json:"integration"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdOn"`
	UpdatedOn   time.Time `json:"updatedOn"`
	// DataSource  bool            `json:"dataSource"`
	// Disabled bool            `json:"disabled"`
	Secret json.RawMessage `json:"secret,omitempty"`
	Config json.RawMessage `json:"config"`

	InUse  bool `json:"inUse,omitempty"`
	Plugin bool `json:"plugin,omitempty"`

	Tags []*Tag `json:"tags,omitempty"`
}

func (r *Integration) IsS3Bucket() bool {
	return r.Integration == INTEGRATION_AWS_S3BUCKET
}

func (r *Integration) IsBlobStorage() bool {
	return r.Integration == INTEGRATION_AZURE_BLOBSTORAGE
}

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ClientSecret struct {
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
	PrivateKey   string `json:"privateKey" yaml:"privateKey"`
}

type Office365AuditConfig struct {
	TenantID   string `json:"tenantID" yaml:"tenantID"`
	ClientID   string `json:"clientId" yaml:"clientId"`
	Thumbprint string `json:"thumbprint" yaml:"thumbprint"`
}

type Office365AuditState struct {
	LastPoll  time.Time        `json:"lastPoll"`
	StreamMap map[string]int64 `json:"streamMap"`
}

// Office365ResourceWatch
type Office365ResourceWatchConfig struct {
	TenantID string `json:"tenantID" yaml:"tenantID"`
	// optional: only applicable for manual configuration
	ClientID string `json:"clientId" yaml:"clientId"`
}
type Office365ResourceWatchSecret struct {
	// optional: only applicable for manual configuration
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
}

// AzureAudit
type AzureAuditConfig struct {
	TenantID string `json:"tenantID" yaml:"tenantID"`
	// optional: only applicable for manual configuration
	ClientID string `json:"clientId" yaml:"clientId"`
}
type AzureAuditSecret struct {
	// optional: only applicable for manual configuration
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
}

type SlackConfig struct {
	Token string `json:"token" yaml:"token"`
}

type PagerDutyConfig struct {
	//IntegrationKey string `json:"integrationKey"`

	URL string `json:"url" yaml:"url"`
}

type PagerDutySecret struct {
	IntegrationKey string `json:"integrationKey" yaml:"integrationKey"`
	Token          string `json:"token" yaml:"token"`
}

type CatoConfig struct {
	RoleARN string `json:"roleARN"`
	Prefix  string `json:"prefix"`
	Begin   string `json:"begin"`
	End     string `json:"end"`
	Bucket  string `json:"bucket"`
	Region  string `json:"region"`

	// read, write, read/write
	Mode string `json:"mode"`

	//AWSUser   string `json:"awsUser"`
	AccessKey string `json:"accessKey"`

	Secret *AWSUserSecret `json:"-"`

	Role    string        `json:"role,omitempty"`
	RoleObj *InstanceRole `json:"-"`
}

type AzureBlobStorageConfig struct {
	StorageAccount string `json:"storageAccount"`
	Container      string `json:"container"`
	// Location         string `json:"location"`
	AuthenticationMethod string `json:"authenticationMethod"` // "connectionString", "servicePrincipal"
	TenantID             string `json:"tenantID,omitempty"`
	ClientID             string `json:"clientID,omitempty"`
}

type AzureBlobStorageSecret struct {
	ConnectionString string `json:"connectionString,omitempty"`
	ClientSecret     string `json:"clientSecret,omitempty"`
}

type S3BucketConfig struct {
	Prefix string `json:"prefix" yaml:"prefix"`
	Begin  string `json:"begin" yaml:"begin"`
	End    string `json:"end" yaml:"end"`
	Bucket string `json:"bucket" yaml:"bucket"`
	Region string `json:"region" yaml:"region"`

	// read, write, read/write
	Mode string `json:"mode" yaml:"mode"`

	//AWSUser   string `json:"awsUser"`
	AccessKey string `json:"accessKey" yaml:"accessKey"`
	User      string `json:"user,omitempty" yaml:"user"`
	//Secret *AWSUserSecret `json:"-"`

	Role string `json:"role,omitempty" yaml:"role"`
	//RoleObj *InstanceRole `json:"-"`

	AuthInfo *AWSAuthInfo `json:"-"`
}

type AWSUserSecret struct {
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

type AWSAuthInfo struct {
	SecretKey string        `json:"-"`
	AccessKey string        `json:"accessKey"`
	RoleObj   *InstanceRole `json:"-"`
}

func (r *AWSAuthInfo) GetCopy() *AWSAuthInfoCopy {
	return &AWSAuthInfoCopy{
		SecretKey: r.SecretKey,
		AccessKey: r.AccessKey,
		RoleObj:   r.RoleObj,
	}
}

type AWSAuthInfoCopy struct {
	SecretKey string        `json:"secretKey,omitempty"`
	AccessKey string        `json:"accessKey,omitempty"`
	RoleObj   *InstanceRole `json:"roleObj,omitempty"`
}
type GSuiteSecret struct {
	ClientSecret string        `json:"clientSecret"`
	Token        *oauth2.Token `json:"token"`
}

type GSuiteConfig struct {
	Config *oauth2.Config `json:"config"`
}

type GSuiteState struct {
	LastPoll time.Time `json:"lastPoll"`
	// StreamMap map[string]int64 `json:"streamMap"`
}

type DuoConfig struct {
	//Customer       string `json:"customer"`
	IntegrationKey string `json:"integrationKey" yaml:"integrationKey"`
	//SecretKey      string `json:"secretKey"`
	APIHostname string `json:"apiHostname" yaml:"apiHostname"`
}

type DuoSecret struct {
	// Customer       string `json:"customer"`
	//IntegrationKey string `json:"integrationKey"`
	SecretKey string `json:"secretKey" yaml:"secretKey"`
	//APIHostname    string `json:"apiHostname"`
}
type DuoState struct {
	StateMap map[string]int64 `json:"stateMap"`
}

type MSDefenderATPConfig struct {
	ClientId string `json:"clientId" yaml:"clientId"`
	TenantId string `json:"tenantId" yaml:"tenantId"`
}

type MSDefenderATPSecret struct {
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
}

type OktaConfig struct {
	Domain string `json:"domain" yaml:"domain"`
}

type OktaSecret struct {
	Token string `json:"token" yaml:"token"`
}

type SophosEDRConfig struct {
	ClientID string `json:"clientID" yaml:"clientID"`
}

type SophosEDRSecret struct {
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
}

// integration Falcon
type FalconConfig struct {
	// ClientSecret string           `json:"clientSecret"`
	ClientID string `json:"clientID" yaml:"clientID"`
	BaseURL  string `json:"baseURL" yaml:"baseURL"`
	//StreamMap    map[string]int64 `json:"streamMap"`
	//Customer     string           `json:"customer"`
}

type FalconSecret struct {
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`
}

type FalconState struct {
	StreamMap map[string]int64 `json:"streamMap"`
}

type MimecastConfig struct {
	ApplicationID  string `json:"applicationID" yaml:"applicationID"`
	ApplicationKey string `json:"applicationKey" yaml:"applicationKey"`
	AccessKey      string `json:"accessKey" yaml:"accessKey"`
	// SecretKey      string                `json:"secretKey"`
	BaseURL string `json:"baseURL" yaml:"baseURL"`
}

type MimecastSecret struct {
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

type MimecastState struct {
	// Token          string                `json:"token"`
	LastPoll    time.Time             `json:"lastPoll"`
	APITokenMap map[string]string     `json:"apiTokenMap"`
	APIPollMap  map[string]*time.Time `json:"apiPollMap"`
}

// integration "Proofpoint"
type ProofPointState struct {
	LastPoll time.Time `json:"lastPoll"`
}

type ProofPointConfig struct {
	Principal string `json:"principal" yaml:"principal"`
	// Secret    string `json:"secret"`
}

type ProofPointSecret struct {
	// Principal string `json:"principal"`
	Secret string `json:"secret" yaml:"secret"`
}

type EventhubConfig struct {
	Eventhubs []*AzureEventhubConfig `json:"eventhubs" yaml:"eventhubs"`
}
type AzureEventhubConfig struct {
	// Eventhub    string `json:"eventhub"`
	Endpoint    string `json:"endpoint" yaml:"endpoint"`
	Description string `json:"description" yaml:"description"`
}

type EventhubState struct {
	LastPoll time.Time `json:"lastPoll"`
}

// integration "AWSAudit"
type AWSAuditConfig struct {
	// required
	AccessKey string `json:"accessKey" yaml:"accessKey"`
	// optional
	CloudTrails []*AWSCloudTrail `json:"cloudTrails" yaml:"cloudTrails"`
	// optional
	CloudWatches []*AWSCloudWatch `json:"cloudWatches" yaml:"cloudWatches"`
}

type AWSAuditSecret struct {
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

type AWSCloudTrail struct {
	Name     string `json:"name"`
	Region   string `json:"region"`
	QueueURL string `json:"queueURL,omitempty"`
	Bucket   string `json:"bucket,omitempty"`
	// User     string `json:"user"`
	// Disabled bool   `json:"disabled"`
	// UserObj *AWSUser `json:"-"`
}

type AWSCloudWatch struct {
	// Name   string `json:"name"`
	Region string `json:"region"`
	// User      string   `json:"user"`
	LogGroups []string `json:"logGroups"`
	// Disabled  bool     `json:"disabled"`

	//UserObj *AWSUser `json:"-"`
}

// integration "HEC"

type HECConfig struct {
	Ack bool `json:"ack"`
	// "/services/collector"
	Endpoint string          `json:"endpoint"`
	Props    json.RawMessage `json:"props"`
}
type HECSecret struct {
	Token string `json:"token"`
}

type BitdefenderAPIConfig struct {
	// default value: "https://cloud.gravityzone.bitdefender.com/api"
	URL       string `json:"url"`
	CompanyID string `json:"companyID"`
}
type BitdefenderAPISecret struct {
	APIKey string `json:"apiKey"`
}

type SalesforceConfig struct {
	URL       string `json:"baseURL" yaml:"baseURL"`
	Client_id string `json:"clientID" yaml:"clientID"`
}

type SalesforceSecret struct {
	Client_secret string `json:"clientSecret" yaml:"clientSecret"`
}

type SQSConfig struct {
	URL    string `json:"url" yaml:"url"`
	Region string `json:"region" yaml:"region"`

	// read, write, read/write
	Mode string `json:"mode" yaml:"mode"`

	//AWSUser   string `json:"awsUser"`
	AccessKey string `json:"accessKey" yaml:"accessKey"`

	Secret *AWSUserSecret `json:"-"`

	Role    string        `json:"role,omitempty" yaml:"role,omitempty"`
	RoleObj *InstanceRole `json:"-"`
}

type SQSSecret struct {
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

type S3NotificationConfig struct {
	QueueARN  string `json:"queueARN" yaml:"queueARN"`
	QueueURL  string `json:"queueURL" yaml:"queueURL"`
	Region    string `json:"region" yaml:"region"`
	AccessKey string `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	User      string `json:"user,omitempty" yaml:"user,omitempty"`
	Role      string `json:"role,omitempty" yaml:"role,omitempty"`
}

type S3NotificationSecret struct {
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

type ServiceNowConfig struct {
	URL           string `json:"url" yaml:"url"`
	Username      string `json:"username" yaml:"username"`
	Scope         string `json:"scope" yaml:"scope"`
	IncidentTable string `json:"incidentTable" yaml:"incidentTable"`
}

type ServiceNowSecret struct {
	Password string `json:"password" yaml:"password"`
}

type AWSCloudWatchConfig struct {
	// required
	AccessKey string `json:"accessKey" yaml:"accessKey"`
	// optional
	//CloudTrails []*AWSCloudTrail `json:"cloudTrails"`
	// optional
	CloudWatches []*AWSCloudWatch `json:"cloudWatches" yaml:"cloudWatches"`
}

type AWSCloudWatchSecret struct {
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

type AWSCloudTrailConfig struct {
	// required
	AccessKey string `json:"accessKey"`

	// optional
	CloudTrails []*AWSCloudTrail `json:"cloudTrails"`
	// optional
	//CloudWatches []*AWSCloudWatch `json:"cloudWatches"`
}

type AWSCloudTrailSecret struct {
	SecretKey string `json:"secretKey"`
}

type AWSGuardDutyConfig struct {
	// requird
	Regions []string `json:"regions" yaml:"regions"`
	// required
	AccessKey string `json:"accessKey,omitempty" yaml:"accessKey"`
	User      string `json:"user,omitempty" yaml:"user"`
	Role      string `json:"role,omitempty" yaml:"role"`
}

type AWSGuardDutySecret struct {
	// required
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

type AWSKinesisStreamConfig struct {
	// required
	Region string `json:"region" yaml:"region"`
	// requird
	Name string `json:"name" yaml:"name"`
	// requird
	Arn string `json:"arn" yaml:"arn"`

	// read, write, read/write
	Mode string `json:"mode" yaml:"mode"`

	// required
	AccessKey string `json:"accessKey,omitempty" yaml:"accessKey"`
	User      string `json:"user,omitempty" yaml:"user"`
	Role      string `json:"role,omitempty" yaml:"role"`
}

type AWSKinesisStreamSecret struct {
	// required
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

type AWSLambdaConfig struct {
	// requird
	Region       string `json:"region" yaml:"region"`
	FunctionName string `json:"functionName" yaml:"functionName"`
	// required
	AccessKey string `json:"accessKey,omitempty" yaml:"accessKey"`
	Role      string `json:"role,omitempty" yaml:"role"`
	User      string `json:"user,omitempty" yaml:"user"`
}

type AWSLambdaSecret struct {
	// required
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

// integration "FirehoseStream"
type AWSFirehoseConfig struct {
	// requird
	Region string `json:"region" yaml:"region"`
	Name   string `json:"name" yaml:"name"`
	// required
	AccessKey string `json:"accessKey,omitempty" yaml:"accessKey"`
	Role      string `json:"role,omitempty" yaml:"role"`
	User      string `json:"user,omitempty" yaml:"user"`
}

type AWSFirehoseSecret struct {
	// required
	SecretKey string `json:"secretKey" yaml:"secretKey"`
}

func GetConfigSecret[C any, S any](r *Integration) (*C, *S, error) {
	var config C
	var secret S

	err := json.Unmarshal(r.Config, &config)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(r.Secret, &secret)
	if err != nil {
		return nil, nil, err
	}
	return &config, &secret, nil
}

func GetConfig[C any](r *Integration) (*C, error) {
	var config C

	err := json.Unmarshal(r.Config, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetSecret[S any](r *Integration) (*S, error) {
	var secret S

	err := json.Unmarshal(r.Secret, &secret)
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

type DataLakeConfig struct {
	// integration type  INTEGRATION_AWS_S3BUCKET
	// StorageType string `json:"storageType" yaml:"storageType"`

	// if this flag is set, use cluster default cloud storage
	Managed bool `json:"managed" yaml:"managed"`
	// integration id of cloud storage (s3)
	StorageIntegration string `json:"storageIntegration" yaml:"storageIntegration"`
	Index              string `json:"index" yaml:"index"`
	// "fluency default"  or  "elastic/opensearch"  or "k8s logs"
	SchemaName string `json:"schemaName" yaml:"schemaName"`
	Schema     string `json:"schema" yaml:"schema"`
	// default is 200MB
	//BufferSizeInMB int64 `json:"bufferSizeInMB" yaml:"bufferSizeInMB"`
}
