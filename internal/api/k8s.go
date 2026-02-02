package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K8sClusterClient wraps the standard kubernetes clientset
type K8sClusterClient struct {
	clientset *kubernetes.Clientset
	config    clientcmd.ClientConfig
}

type ServiceStatus struct {
	Name         string
	Type         string
	ExternalIP   string
	ReadyPods    int
	NotReadyPods int
	Status       string
}

type ServiceEndpoint struct {
	Name        string
	Namespace   string
	Port        int
	HealthCheck string
}

type EndpointHealthResult struct {
	Service   string
	Namespace string
	Port      int
	Path      string
	Success   bool
	Error     string
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// NewK8sClient creates a new instance (but doesn't connect yet)
func NewK8sClient() *K8sClusterClient {
	return &K8sClusterClient{}
}

// Connect loads the kubeconfig and initializes the clientset for a specific context
func (k *K8sClusterClient) Connect(kubeContext string) error {
	// 1. Resolve kubeconfig path (default to ~/.kube/config)
	// You can make this configurable via flags if needed
	home, _ := os.UserHomeDir()
	kubeconfigPath := filepath.Join(home, ".kube", "config")

	if envFile := os.Getenv("KUBECONFIG"); envFile != "" {
		kubeconfigPath = envFile
	}

	// 2. Define loading rules
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}

	// 3. Define overrides (this is how we select the specific context)
	configOverrides := &clientcmd.ConfigOverrides{}
	if kubeContext != "" {
		configOverrides.CurrentContext = kubeContext
	}

	// 4. Build the config
	// This creates a lazy-loader that reads the file + applies overrides
	k.config = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	clientConfig, err := k.config.ClientConfig()
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig (context: %s): %w", kubeContext, err)
	}

	// 5. Create the Clientset
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %w", err)
	}

	k.clientset = clientset
	return nil
}

// IsConnected returns true if the client has been connected to a cluster.
func (k *K8sClusterClient) IsConnected() bool {
	return k != nil && k.clientset != nil
}

// GetAppSecretToken fetches a secret from k8s and returns the "token" field
func (k *K8sClusterClient) GetAppSecret(namespace, secretName string, key string) (string, error) {
	if k.clientset == nil {
		return "", fmt.Errorf("k8s client not initialized")
	}

	// Fetch the Secret
	// We use context.TODO() here, but ideally, pass a context down from the command
	secret, err := k.clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get secret '%s' in namespace '%s': %w", secretName, namespace, err)
	}

	// Look for the 'token' key in the secret data
	// K8s secrets data is []byte, so we cast it to string
	if tokenBytes, ok := secret.Data[key]; ok {
		return string(tokenBytes), nil
	}

	// Fallback: Check 'data' key if 'token' doesn't exist (common pattern)
	if dataBytes, ok := secret.Data["data"]; ok {
		return string(dataBytes), nil
	}

	return "", fmt.Errorf("secret '%s' found, but contains no key %s", secretName, key)
}

func (k *K8sClusterClient) GetAppConfig(namespace, configName string, key string) (string, error) {
	if k.clientset == nil {
		return "", fmt.Errorf("k8s client not initialized")
	}

	// Fetch the Secret
	// We use context.TODO() here, but ideally, pass a context down from the command
	config, err := k.clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get configMap '%s' in namespace '%s': %w", configName, namespace, err)
	}

	// Look for the 'token' key in the secret data
	// K8s secrets data is []byte, so we cast it to string
	if b, ok := config.Data[key]; ok {
		return string(b), nil
	}

	return "", fmt.Errorf("configmap '%s' found, but contains no key %s", configName, key)
}

func (k *K8sClusterClient) CheckNamespaceServices(namespace string) ([]ServiceStatus, error) {
	if k.clientset == nil {
		return nil, fmt.Errorf("k8s client not initialized")
	}

	services, err := k.clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services in namespace %s: %w", namespace, err)
	}

	serviceStatuses := make([]ServiceStatus, 0, len(services.Items))
	for _, service := range services.Items {
		externalIP := ""
		if service.Spec.Type == "LoadBalancer" {
			switch {
			case len(service.Status.LoadBalancer.Ingress) > 0 && service.Status.LoadBalancer.Ingress[0].IP != "":
				externalIP = service.Status.LoadBalancer.Ingress[0].IP
			case len(service.Status.LoadBalancer.Ingress) > 0 && service.Status.LoadBalancer.Ingress[0].Hostname != "":
				externalIP = service.Status.LoadBalancer.Ingress[0].Hostname
			default:
				externalIP = "pending"
			}
		}

		endpointSlices, err := k.clientset.DiscoveryV1().EndpointSlices(namespace).List(
			context.TODO(),
			metav1.ListOptions{
				// Standard label to select slices for a service
				LabelSelector: fmt.Sprintf("kubernetes.io/service-name=%s", service.Name),
			},
		)
		if err != nil {
			serviceStatuses = append(serviceStatuses, ServiceStatus{
				Name:   service.Name,
				Type:   string(service.Spec.Type),
				Status: "Unknown",
			})
			continue
		}

		readyIPs := 0
		notReadyIPs := 0

		for _, slice := range endpointSlices.Items {
			for _, ep := range slice.Endpoints {
				if ep.Conditions.Ready != nil && *ep.Conditions.Ready {
					readyIPs++
				} else {
					notReadyIPs++
				}
			}
		}

		status := "Down"
		switch {
		case readyIPs == 0 && notReadyIPs == 0:
			status = "No pods"
		case readyIPs > 0 && notReadyIPs == 0:
			status = "Healthy"
		case readyIPs > 0 && notReadyIPs > 0:
			status = "Degraded"
		default:
			status = "Down"
		}

		serviceStatuses = append(serviceStatuses, ServiceStatus{
			Name:         service.Name,
			Type:         string(service.Spec.Type),
			ExternalIP:   externalIP,
			ReadyPods:    readyIPs,
			NotReadyPods: notReadyIPs,
			Status:       status,
		})
	}
	return serviceStatuses, nil
}

func (k *K8sClusterClient) CheckServiceEndpoints(defaultNamespace string, services []*ServiceEndpoint) ([]EndpointHealthResult, error) {
	if k.clientset == nil {
		return nil, fmt.Errorf("k8s client not initialized")
	}

	endpointResults := make([]EndpointHealthResult, 0, len(services))
	for _, service := range services {
		namespace := service.Namespace
		if namespace == "" {
			namespace = defaultNamespace
		}

		path := service.HealthCheck
		if path == "" {
			path = "/"
		}

		result := k.clientset.CoreV1().Services(namespace).ProxyGet("http", service.Name, fmt.Sprintf("%d", service.Port), path, nil)

		_, err := result.DoRaw(context.TODO())

		endpointResults = append(endpointResults, EndpointHealthResult{
			Service:   service.Name,
			Namespace: namespace,
			Port:      service.Port,
			Path:      path,
			Success:   err == nil,
			Error:     errorString(err),
		})
	}
	return endpointResults, nil
}
