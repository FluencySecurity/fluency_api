# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build -o fluency cmd/fluency/main.go

# Run tests
go test ./...

# Run a single test
go test ./api/... -run TestName

# Install binary
install -m 755 fluency /usr/local/bin/

# Docker (monitoring container)
docker compose up -d   # requires .env with SLACK_WEBHOOK_URL
```

## Architecture

**fluency** is a CLI tool (AWS CLI style: `fluency <noun> <verb> [flags]`) for managing Fluency data pipeline resources. It supports two connection modes: site credentials (recommended) or Kubernetes cluster+namespace.

### Request Flow

```
cmd/fluency/main.go
  → internal/commands/root.go     # cobra root; loads config, initializes API client
  → internal/commands/*.go        # subcommands (stream, auth, processor, etc.)
  → internal/api/client.go        # Client struct — bridges config to HTTP calls
  → client/client.go              # FluencyClient.GenericCall() — POST to backend RPC
  → https://{site}/api/{prefix}/{function}
```

### Key Packages

| Package | Role |
|---------|------|
| `internal/commands/` | Cobra subcommands; `root.go` handles global flags and `PersistentPreRunE` setup |
| `internal/api/` | Business logic for each resource type (`stream_api.go`, `auth_api.go`, etc.); `client.go` initializes the client from either site config or K8s secrets |
| `internal/config/` | Viper config (`~/.fluency/config.yaml`) + `site_credentials.json` parsing |
| `client/` | Low-level HTTP wrapper; sends Bearer-token-authenticated POSTs; 600s timeout |
| `api/` | Older/parallel API layer (e.g., `platform_api.go`); some overlap with `internal/api/` |
| `fsb/` | Custom flexible JSON type (`JNode`) used in all RPC payloads |
| `model/` | Shared data structs for streams, processors, integrations, etc. |
| `internal/server/` | Optional HTTP REST server (`fluency serve`) that wraps CLI commands as endpoints |

### Connection Modes

**Site config (no K8s needed):**
```json
// site_credentials.json
{ "tokenMap": { "demo.cloud.fluencysecurity.com": "api-token" } }
```

**Kubernetes mode:** reads token from Secret `app-secret` (key: `token`) and siteURL from ConfigMap `ingext-community-config` (key: `site_config.json`).

### Configuration Hierarchy (highest → lowest)
1. CLI flags
2. Environment variables (`FLUENCY_CLUSTER`, `FLUENCY_NAMESPACE`, `FLUENCY_SITE`, `FLUENCY_SITE_CONFIG`, `FLUENCY_DEFAULT_SITE`)
3. `~/.fluency/config.yaml`
4. Defaults

### RPC Protocol

All API calls use HTTP POST to `https://{site}/api/{prefix}/{function}` with a `fsb.CallRequest` body. Responses are `fsb.CallResponse` with verdict `OK`, `ERROR`, or `EXCEPTION`. Prefixes: `api/auth` (users/roles), `api/ds` (data services/streams).

### Docker / Monitoring

The Docker container runs monitoring scripts on a 1-hour cron. Scripts in `scripts/` check database health, platform stream I/O, collector status, and eventwatch across all sites, then alert via Slack webhook or SMTP email. Secrets (`site_credentials.json`, `.env`) are mounted at runtime and excluded from git.
