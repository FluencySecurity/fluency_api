package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/SecurityDo/fluency_api/client"
	"github.com/SecurityDo/fluency_api/internal/config"
)

// FluencyAppAPI defines the contract for interacting with the backend
/*
type FluencyAppAPI interface {
	Init(cluster, namespace, kubeCtx string) error
	Call(functionName string, functionArgs json.RawMessage) error

	// Domain specific methods
	AddStreamSource() error
	AddStreamSink() error
	// ... add other methods here
}*/

// Client is the concrete implementation of FluencyAppAPI
type Client struct {
	Logger    *slog.Logger
	Cluster   string
	Namespace string
	// Add k8s clients, rest clients, etc here

	// Embed the K8s helper
	k8sClient    *K8sClusterClient
	fluencyClient *client.FluencyClient
}

// Option 1: Constructor injection (Recommended)
func NewClient(logger *slog.Logger) *Client {
	// Fallback: If caller passes nil, use a "No-Op" or Default logger
	// so the code doesn't panic.
	if logger == nil {
		// Using os.Stderr by default is safe for libraries
		logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}
	return &Client{
		Logger:    logger,
		k8sClient: NewK8sClient(), // Initialize the helper
	}
}

// Ensure Client implements the interface
//var _ FluencyAppAPI = (*Client)(nil)

var FluencyAppAPI *Client

func (c *Client) Init(cluster, namespace string, kubeContext string) error {
	c.Cluster = cluster
	c.Namespace = namespace

	c.Logger.Debug("initializing k8s client", "context", kubeContext)

	// 1. Connect to Kubernetes
	if err := c.k8sClient.Connect(kubeContext); err != nil {
		return err
	}
	c.Logger.Info("connected to kubernetes cluster", "context", kubeContext)

	// TODO: Perform actual login / connection logic here
	c.Logger.Debug("Connecting to cluster ...\n", "cluster", cluster, "namespace", namespace)

	token, err := c.k8sClient.GetAppSecret(namespace, "app-secret", "token")
	if err != nil {
		return fmt.Errorf("failed to get app secret token: %s", err)
	}

	configText, err := c.k8sClient.GetAppConfig(namespace, "ingext-community-config", "site_config.json")
	if err != nil {
		return fmt.Errorf("failed to get app config: %s", err)
	}
	var config struct {
		SiteURL string `json:"siteURL"`
	}

	if err := json.Unmarshal([]byte(configText), &config); err != nil {
		c.Logger.Error("failed to parse site config", "error", err, "config", configText)
		return fmt.Errorf("failed to parse site config:  %s", err)
	}

	fluencyClient := client.NewFluencyClient(config.SiteURL, token, false, c.Logger)

	c.fluencyClient = fluencyClient

	c.Logger.Info("initialized fluency client",
		"siteURL", config.SiteURL,
		//"token", token,
	)

	return nil
}

// SetDebug enables or disables HTTP request/response dump logging (e.g. when --log-level debug).
func (c *Client) SetDebug(debug bool) {
	if c.fluencyClient != nil {
		c.fluencyClient.SetDebug(debug)
	}
}

// InitFromSiteConfig initializes the client from site_credentials.json (no Kubernetes).
// siteConfigPath is the path to site_credentials.json; site is the hostname key (e.g. "demo.cloud.fluencysecurity.com").
// If site is empty, the first site in tokenMap is used.
func (c *Client) InitFromSiteConfig(siteConfigPath, site string) error {
	creds, err := config.LoadSiteCredentials(siteConfigPath)
	if err != nil {
		return err
	}
	baseURL, token, err := config.ResolveSite(creds, site)
	if err != nil {
		return err
	}
	c.fluencyClient = client.NewFluencyClient(baseURL, token, false, c.Logger)
	c.Logger.Info("initialized fluency client from site config", "siteURL", baseURL)
	return nil
}

func (c *Client) Call(functionName string, functionArgs json.RawMessage) error {
	// TODO: Implementation
	return nil
}
