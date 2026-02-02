# fluency CLI

`fluency` is a command-line interface tool for managing Fluency resources on Kubernetes. It allows users to manage streams, processors, integrations, and authentication through a structured, AWS-CLI-style interface.

## Features

* **Standardized CLI:** Follows the intuitive `noun verb [flags]` pattern (e.g., `fluency stream add source`).
* **Kubernetes Native:** Connects directly to your clusters using your local `kubeconfig` context.
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

Before running commands, configure the target Kubernetes cluster and namespace. This saves settings to `~/.fluency/config.yaml`.

```bash
# Set your default target
fluency config set --cluster <k8s-cluster> --namespace <app-namespace> --context <kubectlContext>  --provider <eks|aks|gke>

# Example
fluency config set --cluster datalake  --namespace fluency --provider eks --context arn:aws:eks:$Region:$AWSAccount:cluster/datalake 

```

You can view your current configuration at any time:

```bash
fluency config view
```

List all configured clusters or delete a cluster profile:

```bash
fluency config list
fluency config delete --cluster <cluster-name>
```

**Environment Variables**
You can override defaults using `FLUENCY_` prefixed variables:

```bash
export FLUENCY_CLUSTER=prod-cluster
export FLUENCY_NAMESPACE=fluency

```

## Usage

### Global flags

| Flag | Shorthand | Default | Description |
| --- | --- | --- | --- |
| `--cluster` |  | _none_ | Target Kubernetes cluster (required unless set via config). |
| `--namespace` | `-n` | `fluency` | Namespace of the fluency app. |
| `--log-level` | `-l` | `warn` | Log level: `debug`, `info`, `warn`, or `error`. |
| `--version` | `-v` | `false` | Print CLI version (`1.1.0`) and exit. |

### Status (`status`)

Check the current namespace for running services and health checks for core fluency endpoints. Prints a table plus a summary of healthy/degraded/down services.

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
