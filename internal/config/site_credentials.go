package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/SecurityDo/fluency_api/model"
)

// LoadSiteCredentials reads site_credentials.json from path and returns the struct.
func LoadSiteCredentials(path string) (*model.SiteCredentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read site credentials file %q: %w", path, err)
	}
	var creds model.SiteCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse site credentials: %w", err)
	}
	if creds.TokenMap == nil || len(creds.TokenMap) == 0 {
		return nil, fmt.Errorf("site credentials file has no tokenMap or tokenMap is empty")
	}
	return &creds, nil
}

// ResolveSite returns base URL (https://hostname) and token for the given site.
// If site is empty, the first site in tokenMap (sorted by key) is used.
func ResolveSite(creds *model.SiteCredentials, site string) (baseURL, token string, err error) {
	if creds == nil || creds.TokenMap == nil {
		return "", "", fmt.Errorf("site credentials are nil or empty")
	}
	if site != "" {
		token, ok := creds.TokenMap[site]
		if !ok {
			return "", "", fmt.Errorf("site %q not found in tokenMap", site)
		}
		return "https://" + site, token, nil
	}
	// Pick first site (deterministic order)
	keys := make([]string, 0, len(creds.TokenMap))
	for k := range creds.TokenMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	host := keys[0]
	return "https://" + host, creds.TokenMap[host], nil
}
