package api

import (
	"fmt"
)

func (c *Client) CheckStatus() error {
	if c.k8sClient == nil || !c.k8sClient.IsConnected() {
		return fmt.Errorf("status is only available when connected via Kubernetes (--cluster and --namespace); when using site config (site_credentials.json), cluster health is not checked")
	}

	serviceStatuses, err := c.k8sClient.CheckNamespaceServices(c.Namespace)
	if err != nil {
		c.Logger.Error("failed to check namespace services", "error", err, "namespace", c.Namespace)
		return fmt.Errorf("failed to check namespace services: %w", err)
	}

	endpoints := []*ServiceEndpoint{
		{Name: "api", Port: 8002, HealthCheck: "/health-check"},
		{Name: "datalake-api", Port: 19010, HealthCheck: "/health-check"},
		{Name: "platform-service", Port: 28180, HealthCheck: "/health-check"},
		{Name: "search-service", Port: 19100, HealthCheck: "/health-check"},
	}

	endpointResults, err := c.k8sClient.CheckServiceEndpoints(c.Namespace, endpoints)
	if err != nil {
		c.Logger.Error("failed to check service endpoints", "error", err, "namespace", c.Namespace)
		return fmt.Errorf("failed to check service endpoints: %w", err)
	}

	fmt.Printf("\nServices in namespace %q\n", c.Namespace)
	fmt.Printf("%-24s %-12s %-18s %-18s %-10s\n", "Service", "Type", "External IP/Host", "Pods (ready/total)", "Status")
	fmt.Printf("%s\n", "--------------------------------------------------------------------------------")

	var healthyCount, degradedCount, downCount int
	for _, s := range serviceStatuses {
		external := s.ExternalIP
		if external == "" {
			external = "-"
		}
		totalPods := s.ReadyPods + s.NotReadyPods
		switch s.Status {
		case "Healthy":
			healthyCount++
		case "Degraded":
			degradedCount++
		default:
			downCount++
		}
		fmt.Printf("%-24s %-12s %-18s %5d/%-12d %-10s\n", s.Name, s.Type, external, s.ReadyPods, totalPods, s.Status)
	}

	fmt.Printf("\nHealth checks\n")
	fmt.Printf("%-24s %-12s %-6s %-18s %-8s %s\n", "Service", "Namespace", "Port", "Path", "Result", "Details")
	fmt.Printf("%s\n", "--------------------------------------------------------------------------------")

	passedChecks := 0
	for _, res := range endpointResults {
		resultText := "FAIL"
		if res.Success {
			resultText = "PASS"
			passedChecks++
		}
		fmt.Printf("%-24s %-12s %-6d %-18s %-8s %s\n", res.Service, res.Namespace, res.Port, res.Path, resultText, res.Error)
	}

	totalServices := len(serviceStatuses)
	totalChecks := len(endpointResults)
	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Services: %d total | %d healthy | %d degraded | %d down/no pods/unknown\n", totalServices, healthyCount, degradedCount, downCount)
	fmt.Printf("- Health checks: %d/%d passed\n\n", passedChecks, totalChecks)

	return nil
}
