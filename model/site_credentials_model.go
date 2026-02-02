package model

// SiteCredentials is the structure of site_credentials.json.
// TokenMap maps site hostname (e.g. "demo.cloud.fluencysecurity.com") to API token.
type SiteCredentials struct {
	TokenMap map[string]string `json:"tokenMap"`
}
