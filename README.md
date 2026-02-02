# fluency CLI

`fluency` is a command-line interface tool for managing Fluency resources. It allows users to manage streams, processors, integrations, and authentication through a structured, AWS-CLI-style interface.

## Features

* **Standardized CLI:** Follows the intuitive `noun verb [flags]` pattern (e.g., `fluency stream add source`).
* **Site config:** Run commands using a `site_credentials.json` file (site hostname → API token); no Kubernetes required.
* **Kubernetes optional:** Can also connect via your local `kubeconfig` context when cluster/namespace are configured.
* **Pipe Friendly:** Designed for automation—strictly separates data output (STDOUT) from logs (STDERR) and supports reading files from STDIN.
* **Smart Config:** Hierarchical configuration (Flags > Env Vars > Config File).

## Installation

### From Source

Requirements: Go 1.21+

```bash
# Clone the repository
git clone https://github.com/your-org/fluency_api.git
cd fluency_api

# Build the binary
go build -o fluency cmd/fluency/main.go

# (Optional) Move to path
install -m 755 fluency /usr/local/bin/.

```

## Configuration

### Option 1: Site config (recommended, no Kubernetes)

Place a `site_credentials.json` file in your current directory (or pass `--site-config <path>`). The file maps site hostnames to API tokens:

```json
{
  "tokenMap": {
    "demo.cloud.fluencysecurity.com": "api-token",
    "expo.app.fluencyplatform.com": "api-token"
  }
}
```

Then run any command; the first site in `tokenMap` is used by default. To target a specific site, use `--site` or set a **default site** so you don't need `--site` every time:

**Default site (when `--site` is not set):**

- **Command:** `fluency config set-default-site demo.cloud.fluencysecurity.com` (saves to `~/.fluency/config.yaml`)
- **Config file:** or add to `~/.fluency/config.yaml`:
  ```yaml
  default-site: demo.cloud.fluencysecurity.com
  ```
- **Environment:** `export FLUENCY_DEFAULT_SITE=demo.cloud.fluencysecurity.com`

```bash
fluency --site demo.cloud.fluencysecurity.com stream list source
fluency --site-config /path/to/site_credentials.json auth list-user
```

### Option 2: Kubernetes cluster and namespace

Configure the target cluster and namespace (saved to `~/.fluency/config.yaml`):

```bash
fluency config set --cluster <k8s-cluster> --namespace <app-namespace> --context <kubectlContext> --provider <eks|aks|gke>
fluency config view              # show cluster + site settings (default site, site config path, available sites)
fluency config list              # list configured clusters
fluency config list-sites        # list sites from site_credentials.json (marks default with *)
fluency config set-default-site <hostname>   # set default site (e.g. demo.cloud.fluencysecurity.com)
fluency config delete --cluster <cluster-name>
```

**Environment variables**

Override with `FLUENCY_` prefixed variables:

```bash
export FLUENCY_SITE_CONFIG=/path/to/site_credentials.json
export FLUENCY_SITE=demo.cloud.fluencysecurity.com
export FLUENCY_DEFAULT_SITE=demo.cloud.fluencysecurity.com   # default when --site is not set
# Or for Kubernetes mode:
export FLUENCY_CLUSTER=prod-cluster
export FLUENCY_NAMESPACE=fluency
```

## Usage

### Global flags

| Flag | Shorthand | Default | Description |
| --- | --- | --- | --- |
| `--site-config` |  | `./site_credentials.json` | Path to site_credentials.json (site hostname → token). |
| `--site` |  | _first in tokenMap_ | Site hostname to use (e.g. `demo.cloud.fluencysecurity.com`). |
| `--cluster` |  | _none_ | Kubernetes cluster (required only when not using site config). |
| `--namespace` | `-n` | `fluency` | Namespace of the fluency app (Kubernetes mode). |
| `--log-level` | `-l` | `warn` | Log level: `debug`, `info`, `warn`, or `error`. |
| `--version` | `-v` | `false` | Print CLI version (`1.1.0`) and exit. |

### Status (`status`)

When connected via Kubernetes, check the current namespace for running services and health checks. When using site config, this command reports that status is only available in cluster mode.

```bash
fluency status
```

### Authentication (`auth`)

Manage users and access tokens.

```bash
# Users
fluency auth add-user --name foo@gmail.com --role admin --displayName "Foo Bar" --org fluency
fluency auth del-user --name foo@gmail.com
fluency auth list-user

# API tokens
fluency auth add-token --name ci-bot --role analyst --displayName "CI Bot"
fluency auth del-token --name ci-bot
fluency auth list-token
```

### Streams (`stream`)

Manage data pipelines (sources, sinks, routers).

```bash
# Sources
fluency stream add-source --name clickstream-v1 --source-type plugin --integration-id <integration-id>
fluency stream add-source --name hec-ingest --source-type hec
fluency stream list-source
fluency stream del-source --id <source-id>

# Sinks
fluency stream add-sink --name datalake-out --sink-type datalake --datalake managed --index <index-name>
fluency stream add-sink --name hec-out --sink-type hec --url https://hec.example --token <token>
fluency stream add-sink --name webhook-out --sink-type webhook --url https://example.com/hook
fluency stream list-sink
fluency stream del-sink --id <sink-id>

# Routers and wiring
fluency stream add-router --processor my-processor --router-name main-router
fluency stream connect-router --source-id <source-id> --router-id <router-id>
fluency stream connect-sink --router-id <router-id> --sink-id <sink-id>
```

### Processors (`processor`)

Deploy data processors. Supports piping input via `-` and file loading via `@path`.

```bash
fluency processor add --name filter-logic --content @./scripts/filter.js [--type fpl_processor] [--desc "Filter logic"]
cat ./scripts/transform.js | fluency processor add --name transform-logic --content -
fluency processor list
fluency processor del --name filter-logic
```

### Integrations (`integration`)

Manage third-party connections.

```bash
fluency integration add --integration slack --name alert-bot --description "Send alerts to Slack" \
  --config key1=value1 --config-bool enabled=true --config-int retries=3 \
  --config-json 'tags=["a","b"]' --secret api_key=xxx --add-source
fluency integration list
fluency integration del --id <integration-id>
```

### Data Lake (`datalake`)

Manage datalakes and their indexes.

```bash
fluency datalake add --datalake my-datalake --managed --integration <integration-id>
fluency datalake list
fluency datalake add-index --datalake my-datalake --index events --schema "fluency default"
fluency datalake list-index --datalake my-datalake
fluency datalake del-index --datalake my-datalake --index events
```

### EKS Pod Identity Roles (`eks`)

```bash
fluency eks add-assumed-role --name ingest-role --roleArn <role-arn> [--externalId <external-id>]
fluency eks list-assumed-role
fluency eks del-assumed-role --id <role-id>
fluency eks get-pod-role
fluency eks test-assumed-role --roleArn <role-arn> [--externalId <external-id>]
```

### Applications (`application`)

```bash
# List and manage templates
fluency application list
fluency application add --content @./template.yaml
fluency application update --app <template> --content @./template.yaml
fluency application del --app <template>

# Install and manage instances
fluency application install --app <template> --instance <instance> --displayName "My App" --config key=value --secret secretKey=value
fluency application uninstall --app <template> --instance <instance>
fluency application get-instance --app <template> --instance <instance>
```

### Import (`import`)

Import resources from a GitHub repository.

```bash
fluency import processor --type fpl_processor
fluency import application
```

## Development

### Project Structure

The project follows the Standard Go Project Layout:

| Path | Description |
| --- | --- |
| `cmd/fluency/` | Application entry point (`main.go`). |
| `internal/commands/` | Cobra command definitions and flag parsing. |
| `internal/api/` | Business logic and Kubernetes client (`client-go`). |
| `internal/config/` | Configuration loading (Viper). |

### Kubernetes Dependency Note

This project uses `client-go` v0.35.0. If you change versions, ensure all k8s libraries match exactly to avoid build errors:

```bash
go get k8s.io/client-go@v0.35.0 k8s.io/api@v0.35.0 k8s.io/apimachinery@v0.35.0
go get "github.com/google/go-github/v64/github"
go mod tidy

```

## License

[MIT](https://www.google.com/search?q=LICENSE)
