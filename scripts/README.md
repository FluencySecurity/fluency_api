# Automation scripts

Scripts that interact with the `fluency` CLI for monitoring and automation.

## Prerequisites

- `fluency` binary on your PATH (or use `--fluency /path/to/fluency`).
- A `site_credentials.json` file with a `tokenMap` of site hostnames to API tokens.
- **jq** (used to parse the JSON config and for Slack payloads).

To run monitoring in a Docker container (e.g. on a 24/7 device) and exec in to run these scripts, see the **Docker (24/7 monitoring)** section in the main [README](../README.md).

## audit_db_status_all_sites.sh

Runs `fluency audit db_status` for **every site** in `site_credentials.json` and reports which sites are up or down.

### Usage

```bash
# From repo root (uses ./site_credentials.json)
./scripts/audit_db_status_all_sites.sh

# Or via bash
bash scripts/audit_db_status_all_sites.sh

# Custom site config path
./scripts/audit_db_status_all_sites.sh --site-config /path/to/site_credentials.json

# Verbose: print db_status output for each site
./scripts/audit_db_status_all_sites.sh -v

# Custom fluency binary
./scripts/audit_db_status_all_sites.sh --fluency ./fluency
```

### Exit codes

- `0` – All sites OK.
- `1` – One or more sites down (see stderr for list).
- `2` – Script error (missing config, invalid JSON, jq missing, etc.).

### Customizing “down” criteria

By default, a site is **down** if `fluency audit db_status` exits non-zero (e.g. API error, timeout). You can add your own rules in the `is_site_down` function in the script, for example:

- Treat **etcd status: red** as down (grep stdout).
- Treat **master status** not OK as down.
- Treat **index queue length** above a threshold as down.

Edit the section marked “Add your custom criteria below” in `audit_db_status_all_sites.sh`.

---

## audit_db_status_slack_alert.sh

Runs `audit_db_status_all_sites.sh`; if any site is down (exit code non-zero), sends an alert to a **Slack Incoming Webhook** with the check output.

### Prerequisites

- Everything required for `audit_db_status_all_sites.sh` (fluency, jq, site_credentials.json).
- A Slack Incoming Webhook URL (create one in Slack: App → Incoming Webhooks → Add to Slack).

### Configuring the webhook

Use either:

- **Environment variable:** `export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/..."`
- **CLI flag:** `--slack-webhook-url "https://hooks.slack.com/services/..."`

### Usage

```bash
# Webhook from environment
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
./scripts/audit_db_status_slack_alert.sh

# Webhook from command line
./scripts/audit_db_status_slack_alert.sh --slack-webhook-url "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

# With audit script options (passed through)
./scripts/audit_db_status_slack_alert.sh --slack-webhook-url "https://..." --site-config /path/to/site_credentials.json -v
```

### Exit codes

- Same as `audit_db_status_all_sites.sh`: `0` = all OK, `1` = one or more sites down, `2` = script/config error (including missing webhook URL).
- When exit is `1`, the script posts to Slack and then exits 1.

---

## platform_status_all_sites.sh

Runs `fluency stream status` (platform status) for **every site** in `site_credentials.json` and reports which sites succeed or fail. A site is considered **OK** if it has input data (total input > 0), otherwise it's marked as **DOWN**.

### Usage

```bash
# From repo root (uses ./site_credentials.json)
./scripts/platform_status_all_sites.sh

# Custom site config or fluency binary
./scripts/platform_status_all_sites.sh --site-config /path/to/site_credentials.json --fluency ./fluency

# Verbose: print full error output for failed sites
./scripts/platform_status_all_sites.sh -v
```

### Exit codes

- `0` – All sites OK (have input data).
- `1` – One or more sites down or failed (see stderr for list).
- `2` – Script error (missing config, jq missing, etc.).

---

## platform_status_slack_alert.sh

Runs `platform_status_all_sites.sh`; if any site is down (exit code non-zero), sends an alert to a **Slack Incoming Webhook** with the failed sites and their error messages.

### Prerequisites

- Everything required for `platform_status_all_sites.sh` (fluency, jq, site_credentials.json).
- A Slack Incoming Webhook URL (create one in Slack: App → Incoming Webhooks → Add to Slack).

### Configuring the webhook

Use either:

- **Environment variable:** `export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/..."`
- **CLI flag:** `--slack-webhook-url "https://hooks.slack.com/services/..."`

### Usage

```bash
# Webhook from environment
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
./scripts/platform_status_slack_alert.sh

# Webhook from command line
./scripts/platform_status_slack_alert.sh --slack-webhook-url "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

# With platform status script options (passed through)
./scripts/platform_status_slack_alert.sh --slack-webhook-url "https://..." --site-config /path/to/site_credentials.json -v
```

### Exit codes

- Same as `platform_status_all_sites.sh`: `0` = all OK, `1` = one or more sites down, `2` = script/config error (including missing webhook URL).
- When exit is `1`, the script posts to Slack and then exits 1.

---

## eventwatch_status_all_sites.sh

Runs `fluency eventwatch search_timeline` for the **last 2 hours** for every site in `site_credentials.json`. A site is **OK** if there is at least one timeline event in that window; otherwise it is reported with error **no timeline events**. API/connection failures are reported with a summarized error. Lists all sites that have errors.

### Usage

```bash
# From repo root (uses ./site_credentials.json)
./scripts/eventwatch_status_all_sites.sh

# Custom site config or fluency binary
./scripts/eventwatch_status_all_sites.sh --site-config /path/to/site_credentials.json --fluency ./fluency

# Verbose: print full error output for failed sites
./scripts/eventwatch_status_all_sites.sh -v
```

### Exit codes

- `0` – All sites OK (at least one timeline event in last 2 hours).
- `1` – One or more sites have errors (no timeline events or API failure); see stderr for list.
- `2` – Script error (missing config, jq missing, etc.).

---

## eventwatch_status_slack_alert.sh

Runs `eventwatch_status_all_sites.sh`; if any site has errors (exit code non-zero), sends an alert to a **Slack Incoming Webhook** with the list of failed sites and their errors.

### Prerequisites

- Everything required for `eventwatch_status_all_sites.sh` (fluency, jq, site_credentials.json).
- A Slack Incoming Webhook URL.

### Configuring the webhook

Use either:

- **Environment variable:** `export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/..."`
- **CLI flag:** `--slack-webhook-url "https://hooks.slack.com/services/..."`

### Usage

```bash
# Webhook from environment
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
./scripts/eventwatch_status_slack_alert.sh

# Webhook from command line
./scripts/eventwatch_status_slack_alert.sh --slack-webhook-url "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

# With eventwatch status script options (passed through)
./scripts/eventwatch_status_slack_alert.sh --slack-webhook-url "https://..." --site-config /path/to/site_credentials.json -v
```

### Exit codes

- Same as `eventwatch_status_all_sites.sh`: `0` = all OK, `1` = one or more sites have errors, `2` = script/config error (including missing webhook URL).
- When exit is `1`, the script posts to Slack and then exits 1.

---

## collector_status_all.sh

Runs `fluency collector list`, then gets status for each collector using `fluency collector status --collector <name>`. Reports which collectors succeed or fail.

### Usage

```bash
# From repo root (uses ./site_credentials.json)
./scripts/collector_status_all.sh

# With specific site
./scripts/collector_status_all.sh --site-config /path/to/site_credentials.json --site <site-name>

# Custom fluency binary
./scripts/collector_status_all.sh --fluency ./fluency

# Verbose: print full status output for each collector
./scripts/collector_status_all.sh -v
```

### Exit codes

- `0` – All collectors OK (status retrieved successfully).
- `1` – One or more collectors failed (API error); see stderr for list.
- `2` – Script error (missing config, failed to list collectors, etc.).

---

## collector_status_email_alert.sh

Runs `collector_status_all.sh`; if any collector has errors (exit code non-zero), sends an **SMTP email alert** with the list of failed collectors and their errors.

### Prerequisites

- Everything required for `collector_status_all.sh` (fluency, site_credentials.json).
- **python3** (with smtplib module - standard library).
- SMTP server configuration (host, port, from/to addresses, optionally username/password).

### Configuring SMTP

Use either environment variables or CLI flags:

**Environment variables:**
- `SMTP_HOST` – SMTP server hostname (required)
- `SMTP_PORT` – SMTP server port (default: 587)
- `SMTP_USER` – SMTP username (optional, for authenticated SMTP)
- `SMTP_PASS` – SMTP password (optional, for authenticated SMTP)
- `SMTP_FROM` – From email address (required)
- `SMTP_TO` – To email address (required)
- `SMTP_FROM_NAME` – From name (default: "Fluency Collector Monitor")
- `SMTP_USE_STARTTLS` – Set to `0` to disable STARTTLS (if your server does not support it)

**CLI flags:**
- `--smtp-host HOST`
- `--smtp-port PORT`
- `--smtp-user USER`
- `--smtp-pass PASS`
- `--smtp-from EMAIL`
- `--smtp-to EMAIL`
- `--smtp-from-name NAME`
- `--smtp-no-starttls` – Disable STARTTLS (if server does not support it)

### Usage

```bash
# Using environment variables
export SMTP_HOST="smtp.example.com"
export SMTP_PORT="587"
export SMTP_USER="user@example.com"
export SMTP_PASS="password"
export SMTP_FROM="alerts@example.com"
export SMTP_TO="admin@example.com"
./scripts/collector_status_email_alert.sh

# Using CLI flags
./scripts/collector_status_email_alert.sh \
  --smtp-host smtp.example.com \
  --smtp-port 587 \
  --smtp-user user@example.com \
  --smtp-pass password \
  --smtp-from alerts@example.com \
  --smtp-to admin@example.com

# With collector status script options (passed through)
./scripts/collector_status_email_alert.sh \
  --smtp-host smtp.example.com \
  --smtp-from alerts@example.com \
  --smtp-to admin@example.com \
  --site-config /path/to/site_credentials.json \
  --site <site-name> \
  -v
```

### Exit codes

- Same as `collector_status_all.sh`: `0` = all OK, `1` = one or more collectors failed, `2` = script/config error (including missing SMTP config).
- When exit is `1`, the script sends an email and then exits 1.

---

## billing_run_all_sites.sh

Runs the billing FPL report for every site in `site_credentials.json`, polls until completed (or aborted), then writes results to a CSV. See script header for details.

---

## billing_slack_alert.sh

Runs `billing_run_all_sites.sh` and sends the generated **CSV as a Slack attachment** (plus a short header). Console output is still printed locally.

### Prerequisites

- Everything required for `billing_run_all_sites.sh` (fluency, jq, site_credentials.json, reports/billing.json).
- A Slack Incoming Webhook URL.

### Configuring the webhook

- **Environment variable:** `export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/..."`
- **CLI flag:** `--slack-webhook-url "https://hooks.slack.com/services/..."`

### Usage

```bash
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
./scripts/billing_slack_alert.sh

./scripts/billing_slack_alert.sh --slack-webhook-url "https://..." --site-config /path/to/site_credentials.json

# Attach a different CSV (e.g. if you passed --output to the billing script)
./scripts/billing_slack_alert.sh --slack-webhook-url "https://..." --csv /path/to/billing_results.csv
```

By default the script attaches the CSV at `./billing_results.csv` (repo root). Use `--csv PATH` if the billing script wrote to a different path (e.g. via `--output`). Exits with the same exit code as the billing script.
