#!/usr/bin/env bash
# Run `fluency audit db_status` for every site in site_credentials.json and report
# which sites are up or down. Down-criteria can be customized (see is_site_down).

set -euo pipefail

SITE_CONFIG="site_credentials.json"
FLUENCY="./fluency"
VERBOSE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --site-config)
      SITE_CONFIG="$2"
      shift 2
      ;;
    --fluency)
      FLUENCY="$2"
      shift 2
      ;;
    -v|--verbose)
      VERBOSE=true
      shift
      ;;
    -h|--help)
      echo "Usage: $0 [--site-config PATH] [--fluency PATH] [-v|--verbose]"
      echo "  Run fluency audit db_status for all sites in site_credentials.json."
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      exit 2
      ;;
  esac
done

if [[ ! -f "$SITE_CONFIG" ]]; then
  echo "Error: site config not found: $SITE_CONFIG" >&2
  exit 2
fi

if ! command -v jq &>/dev/null; then
  echo "Error: jq is required to parse site_credentials.json. Install jq and try again." >&2
  exit 2
fi

# Get site list (tokenMap keys, exclude those starting with _)
sites=()
while IFS= read -r line; do
  [[ -n "$line" ]] && sites+=("$line")
done < <(jq -r '.tokenMap | keys[] | select(startswith("_") | not)' "$SITE_CONFIG" | sort)

if [[ ${#sites[@]} -eq 0 ]]; then
  echo "No sites found in tokenMap (or all keys start with _)." >&2
  exit 0
fi

# Optional: add custom criteria by parsing output. Return 0 = down, 1 = up.
# Default: site is down if fluency command failed (exit code != 0).
is_site_down() {
  local site="$1"
  local exit_code="$2"
  local output="$3"

  # Default: command failure means down
  if [[ "$exit_code" -ne 0 ]]; then
    return 0
  fi

  # etcd status not green => down
  if ! echo "$output" | grep -q "etcd status: green"; then
    return 0
  fi

  # master status not green => down
  if ! echo "$output" | grep -q "master status: green"; then
    return 0
  fi

  # any index length > 1000 => down
  while IFS= read -r len; do
    [[ -n "$len" && "$len" -gt 1000 ]] && return 0
  done < <(echo "$output" | grep -oE "length: [0-9]+" | sed 's/length: //')

  return 1
}

down_sites=()
for site in "${sites[@]}"; do
  set +e
  output=$( "$FLUENCY" --site-config "$SITE_CONFIG" --site "$site" audit db_status 2>&1 )
  exit_code=$?
  set -e

  if is_site_down "$site" "$exit_code" "$output"; then
    down_sites+=("$site")
    status="DOWN"
  else
    status="OK"
  fi

  echo "$site: $status"
  if [[ "$VERBOSE" == true ]]; then
    echo "$output"
  fi
done

if [[ ${#down_sites[@]} -gt 0 ]]; then
  echo "" >&2
  echo "Down sites:" >&2
  for s in "${down_sites[@]}"; do echo "$s" >&2; done
  exit 1
fi
exit 0
