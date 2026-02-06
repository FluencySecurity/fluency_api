#!/usr/bin/env bash
# Run `fluency stream status` (platform status) for every site in site_credentials.json.
# Reports which sites succeed or fail. On success, prints the full API output.

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
      echo "  Run fluency stream status (platform status) for all sites in site_credentials.json."
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

# Summarize long error output for failed-sites report (e.g. extract HTTP status or truncate).
summarize_error() {
  local msg="$1"
  local short
  # Prefer HTTP status phrase (e.g. "503 Service Temporarily Unavailable")
  short=$(echo "$msg" | grep -oE '[0-9]{3} [A-Za-z ]+' 2>/dev/null | tail -1)
  if [[ -n "$short" ]]; then
    echo "$short"
    return
  fi
  # Else use quoted error= value if present (strip surrounding noise)
  short=$(echo "$msg" | sed -n 's/.*error="\([^"]*\)".*/\1/p' | sed 's/^.*: //' | tr '\n' ' ' | sed 's/  */ /g' | sed 's/^ *//;s/ *$//')
  if [[ -n "$short" ]]; then
    [[ ${#short} -gt 80 ]] && short="${short:0:77}..."
    echo "$short"
    return
  fi
  # Else truncate
  short=$(echo "$msg" | tr '\n' ' ' | sed 's/  */ /g' | sed 's/^ *//;s/ *$//')
  [[ ${#short} -gt 80 ]] && short="${short:0:77}..."
  echo "${short:-command failed}"
}

# Parse stream status output: sum only values from "ID, input" blocks and print total.
# Data lines: "timestamp - value B" (value may contain commas).
total_input_from_status() {
  awk '
    /, input[[:space:]]*$/ { in_input = 1; next }
    /, [^[:space:]]+[[:space:]]*$/ { in_input = 0; next }
    in_input && / - .* B[[:space:]]*$/ {
      line = $0
      sub(/.* - /, "", line)
      sub(/ B[[:space:]]*$/, "", line)
      gsub(/,/, "", line)
      total += line + 0
      next
    }
    END { printf "%.6g\n", total }
  '
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
  output=$( "$FLUENCY" --site-config "$SITE_CONFIG" --site "$site" stream status --timeslots 2 2>&1 )
  exit_code=$?
  set -e

  if [[ "$exit_code" -eq 0 ]]; then
    total_input=$(echo "$output" | total_input_from_status)
    if awk -v t="${total_input:-0}" 'BEGIN { exit (t+0 > 0) ? 0 : 1 }'; then
      echo "$site: OK"
    else
      echo "$site: NO INPUT DATA"
      failed_sites+=("$site")
      failed_reasons+=("no input data")
    fi
  else
    echo "$site: FAIL"
    failed_sites+=("$site")
    err_msg=$(echo "$output" | tr '\n' ' ' | sed 's/  */ /g' | sed 's/^ *//;s/ *$//')
    failed_reasons+=("$(summarize_error "${err_msg:-command failed}")")
    if [[ "$VERBOSE" == true ]]; then
      echo "$output"
    else
      echo "$output" | sed 's/^/  /' >&2
    fi
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
