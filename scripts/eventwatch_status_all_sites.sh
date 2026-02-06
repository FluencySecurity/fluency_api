#!/usr/bin/env bash
# Run fluency eventwatch search_timeline for the last 2 hours for every site in site_credentials.json.
# A site is OK if there is at least one timeline event in the last 2 hours; otherwise error "no timeline events".
# API/connection failures are reported with a summarized error. Lists all sites that have errors.

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
      echo "  Run fluency eventwatch search_timeline (last 2 hours) for all sites."
      echo "  Sites with no timeline events or API errors are reported as failed."
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

# Last 2 hours in Unix milliseconds
now_sec=$(date +%s)
TO_MS=$((now_sec * 1000))
FROM_MS=$((TO_MS - 2 * 3600 * 1000))

# Summarize long error output for failed-sites report
summarize_error() {
  local msg="$1"
  local short
  short=$(echo "$msg" | grep -oE '[0-9]{3} [A-Za-z ]+' 2>/dev/null | tail -1)
  if [[ -n "$short" ]]; then
    echo "$short"
    return
  fi
  short=$(echo "$msg" | sed -n 's/.*error="\([^"]*\)".*/\1/p' | sed 's/^.*: //' | tr '\n' ' ' | sed 's/  */ /g' | sed 's/^ *//;s/ *$//')
  if [[ -n "$short" ]]; then
    [[ ${#short} -gt 80 ]] && short="${short:0:77}..."
    echo "$short"
    return
  fi
  short=$(echo "$msg" | tr '\n' ' ' | sed 's/  */ /g' | sed 's/^ *//;s/ *$//')
  [[ ${#short} -gt 80 ]] && short="${short:0:77}..."
  echo "${short:-command failed}"
}

# Get site list (tokenMap keys, exclude those starting with _)
sites=()
while IFS= read -r line; do
  [[ -n "$line" ]] && sites+=("$line")
done < <(jq -r '.tokenMap | keys[] | select(startswith("_") | not)' "$SITE_CONFIG" | sort)

if [[ ${#sites[@]} -eq 0 ]]; then
  echo "No sites found in tokenMap (or all keys start with _)." >&2
  exit 0
fi

failed_sites=()
failed_reasons=()
for site in "${sites[@]}"; do
  set +e
  output=$("$FLUENCY" --site-config "$SITE_CONFIG" --site "$site" eventwatch search_timeline --from "$FROM_MS" --to "$TO_MS" 2>&1)
  exit_code=$?
  set -e

  if [[ "$exit_code" -ne 0 ]]; then
    echo "$site: FAIL"
    failed_sites+=("$site")
    err_msg=$(echo "$output" | tr '\n' ' ' | sed 's/  */ /g' | sed 's/^ *//;s/ *$//')
    failed_reasons+=("$(summarize_error "${err_msg:-command failed}")")
    if [[ "$VERBOSE" == true ]]; then
      echo "$output"
    else
      echo "$output" | sed 's/^/  /' >&2
    fi
  elif echo "$output" | grep -q "No hits found"; then
    echo "$site: NO TIMELINE EVENTS"
    failed_sites+=("$site")
    failed_reasons+=("no timeline events")
  else
    echo "$site: OK"
  fi
done

if [[ ${#failed_sites[@]} -gt 0 ]]; then
  echo "" >&2
  echo "Failed sites:" >&2
  for i in "${!failed_sites[@]}"; do
    echo "  ${failed_sites[i]}: ${failed_reasons[i]}" >&2
  done
  exit 1
fi
exit 0
