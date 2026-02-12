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

# Build base fluency command args
FLUENCY_ARGS=()
if [[ -n "$SITE_CONFIG" ]]; then
  FLUENCY_ARGS+=(--site-config "$SITE_CONFIG")
fi
if [[ -n "$SITE" ]]; then
  FLUENCY_ARGS+=(--site "$SITE")
fi

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
    echo "$collector: OK"
    if [[ "$VERBOSE" == true ]]; then
      echo "$output"
    fi
  fi
done

if [[ ${#failed_collectors[@]} -gt 0 ]]; then
  echo "" >&2
  echo "Failed collectors:" >&2
  for i in "${!failed_collectors[@]}"; do
    echo "  ${failed_collectors[i]}: ${failed_reasons[i]}" >&2
  done
  exit 1
fi
exit 0
