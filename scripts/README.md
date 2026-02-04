# Automation scripts

Scripts that interact with the `fluency` CLI for monitoring and automation.

## Prerequisites

- `fluency` binary on your PATH (or use `--fluency /path/to/fluency`).
- A `site_credentials.json` file with a `tokenMap` of site hostnames to API tokens.
- For `audit_db_status_all_sites.sh`: **jq** (used to parse the JSON config).

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
