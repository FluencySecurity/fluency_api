#!/usr/bin/env bash
# Run fluency collector list, then get status for each collector.
# Reports which collectors succeed or fail.

set -euo pipefail

SITE_CONFIG="site_credentials.json"
SITE=""
FLUENCY="./fluency"
VERBOSE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --site-config)
      SITE_CONFIG="$2"
      shift 2
      ;;
    --site)
      SITE="$2"
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
      echo "Usage: $0 [--site-config PATH] [--site NAME] [--fluency PATH] [-v|--verbose]"
      echo "  Run fluency collector list, then get status for each collector."
      echo "  Collectors with API errors are reported as failed."
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
  echo "Error: jq is required for collector status validation. Install jq and try again." >&2
  exit 2
fi

# Build base fluency command args
FLUENCY_ARGS=()
if [[ -n "$SITE_CONFIG" ]]; then
  FLUENCY_ARGS+=(--site-config "$SITE_CONFIG")
fi
if [[ -n "$SITE" ]]; then
  FLUENCY_ARGS+=(--site "$SITE")
fi

# Validate collector status JSON: output list of errors or HEALTHY
validate_collector_json() {
  local output="$1"
  local errors=()
  local line

  if ! jq -e . <<< "$output" &>/dev/null; then
    echo "invalid JSON response"
    return 0
  fi

  # 1. Queues with value > 1000
  while IFS= read -r queue_name; do
    [[ -n "$queue_name" ]] && errors+=("queue $queue_name is too big")
  done < <(jq -r '.queues // {} | to_entries[] | select((.value | tonumber? // 0) > 1000) | .key' <<< "$output" 2>/dev/null)

  # 2. Procinfo: required services must appear in the string
  local required_services=(collectorsync_service ddos_filter_service ddos_summary_service device_mgr fsb_client mgmt_service)
  local procinfo
  procinfo=$(jq -r '.procinfo // ""' <<< "$output" 2>/dev/null)
  for svc in "${required_services[@]}"; do
    if [[ "$procinfo" != *"$svc"* ]]; then
      errors+=("service not running: $svc")
    fi
  done

  # 3. Files: value starting with "ls" (e.g. ls: cannot access ...)
  while IFS= read -r line; do
    [[ -n "$line" ]] && errors+=("$line")
  done < <(jq -r '.files // {} | to_entries[] | select(.value | startswith("ls")) | "\(.key): `\(.value | gsub("\n"; " ") | gsub("`"; "'"'"'") )`"' <<< "$output" 2>/dev/null)

  # 4. Files: value not starting with "ls" (ls -l output) – check if file date is older than 1 day
  local now_epoch file_epoch date_part fkey fval
  now_epoch=$(date +%s 2>/dev/null) || now_epoch=0
  while IFS= read -r line; do
    [[ -z "$line" ]] && continue
    fkey="${line%%$'\t'*}"
    fval="${line#*$'\t'}"
    # Extract "Month DD HH:MM" or "Month DD  YYYY" from ls -l style output
    date_part=$(echo "$fval" | grep -oE '[A-Za-z]{3} +[0-9]{1,2} +[0-9]{2}:[0-9]{2}' | head -1)
    if [[ -z "$date_part" ]]; then
      date_part=$(echo "$fval" | grep -oE '[A-Za-z]{3} +[0-9]{1,2} +[0-9]{4}' | head -1)
    fi
    if [[ -n "$date_part" && "$now_epoch" -gt 0 ]]; then
      file_epoch=$(date -d "$date_part" "+%s" 2>/dev/null) || file_epoch=$(date -j -f "%b %d %H:%M" "$date_part" "+%s" 2>/dev/null)
      if [[ -z "$file_epoch" ]]; then
        file_epoch=$(date -j -f "%b %d %Y" "$date_part" "+%s" 2>/dev/null)
      fi
      if [[ -n "$file_epoch" && $((now_epoch - file_epoch)) -gt 86400 ]]; then
        errors+=("$fkey: file not accessed in the past day")
      fi
    fi
  done < <(jq -r '.files // {} | to_entries[] | select(.value | startswith("ls") | not) | [.key, (.value | gsub("\n"; " "))] | @tsv' <<< "$output" 2>/dev/null)

  if [[ ${#errors[@]} -eq 0 ]]; then
    echo "HEALTHY"
  else
    for err in "${errors[@]}"; do
      echo "$err"
    done
  fi
}

# Summarize long error output for failed-collectors report
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

# Get collector list
set +e
list_output=$("$FLUENCY" "${FLUENCY_ARGS[@]}" collector list 2>&1)
list_exit_code=$?
set -e

if [[ "$list_exit_code" -ne 0 ]]; then
  echo "Error: failed to list collectors" >&2
  echo "$list_output" >&2
  exit 2
fi

# Parse collector names and lastPoll from output (format: "Name: <name>, LastPoll: <timestamp>")
collectors=()
lastpolls=()
while IFS= read -r line; do
  if [[ "$line" =~ ^Name:\ ([^,]+),\ LastPoll:\ ([0-9-]+) ]]; then
    collectors+=("${BASH_REMATCH[1]}")
    lastpolls+=("${BASH_REMATCH[2]}")
  fi
done <<< "$list_output"

if [[ ${#collectors[@]} -eq 0 ]]; then
  if echo "$list_output" | grep -q "No collectors found"; then
    echo "No collectors found."
    exit 0
  fi
  echo "Warning: no collectors parsed from output" >&2
  if [[ "$VERBOSE" == true ]]; then
    echo "$list_output"
  fi
  exit 0
fi

failed_collectors=()
failed_reasons=()
for i in "${!collectors[@]}"; do
  collector="${collectors[$i]}"
  lastpoll="${lastpolls[$i]}"
  
  # Check if lastPoll < 0 (collector offline)
  if [[ "$lastpoll" -lt 0 ]]; then
    echo "$collector: OFFLINE"
    failed_collectors+=("$collector")
    failed_reasons+=("collector offline")
    continue
  fi

  # Skip status call for "local" collector
  if [[ "$collector" == "local" ]]; then
    echo "local: OK"
    continue
  fi
  
  # Call status API for collectors with lastPoll >= 0
  set +e
  output=$("$FLUENCY" "${FLUENCY_ARGS[@]}" collector status --collector "$collector" 2>&1)
  exit_code=$?
  set -e

  if [[ "$exit_code" -ne 0 ]]; then
    echo "$collector: FAIL"
    failed_collectors+=("$collector")
    err_msg=$(echo "$output" | tr '\n' ' ' | sed 's/  */ /g' | sed 's/^ *//;s/ *$//')
    failed_reasons+=("$(summarize_error "${err_msg:-command failed}")")
    if [[ "$VERBOSE" == true ]]; then
      echo "$output"
    else
      echo "$output" | sed 's/^/  /' >&2
    fi
  else
    validation_result=$(validate_collector_json "$output")
    if [[ "$validation_result" == "HEALTHY" ]]; then
      echo "$collector: HEALTHY"
      if [[ "$VERBOSE" == true ]]; then
        echo "$output"
      fi
    else
      echo "$collector:"
      echo "$validation_result" | sed 's/^/  - /'
      failed_collectors+=("$collector")
      failed_reasons+=("$validation_result")
    fi
  fi
done

if [[ ${#failed_collectors[@]} -gt 0 ]]; then
  echo "" >&2
  echo "Failed collectors:" >&2
  for i in "${!failed_collectors[@]}"; do
    echo "  ${failed_collectors[i]}:"
    echo "${failed_reasons[i]}" | sed 's/^/    - /' >&2
  done
  exit 1
fi
exit 0
