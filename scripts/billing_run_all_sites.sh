#!/usr/bin/env bash
# Run `fluency fpl run --file ./reports/billing.json` for every site in site_credentials.json.
# On success: poll task until state is completed, then print "site state". If state is aborted, exit with error. On run failure: print site and error message.

set -euo pipefail

SITE_CONFIG="site_credentials.json"
REPORT_FILE="./reports/billing.json"
FLUENCY="./fluency"
CSV_OUTPUT="./billing_results.csv"

# Parse arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --site-config)
      SITE_CONFIG="$2"
      shift 2
      ;;
    --report-file)
      REPORT_FILE="$2"
      shift 2
      ;;
    --fluency)
      FLUENCY="$2"
      shift 2
      ;;
    --output)
      CSV_OUTPUT="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [--site-config PATH] [--report-file PATH] [--fluency PATH] [--output CSV_PATH]"
      echo "  Run fluency fpl run for all sites in site_credentials.json."
      echo "  Default report file: ./reports/billing.json. Writes results to CSV (default: ./billing_results.csv)."
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

if [[ ! -f "$REPORT_FILE" ]]; then
  echo "Error: report file not found: $REPORT_FILE" >&2
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

echo "siteURL,eventCount,totalDataIngress,dailyIngressRate,tickets,error" > "$CSV_OUTPUT"

# Append a CSV row with error message (one line, escapes " for CSV)
csv_error_row() {
  local site="$1"
  local err="$2"
  local one_line=$(echo "$err" | tr '\n' ' ' | sed 's/  */ /g' | sed 's/^ *//;s/ *$//')
  local escaped="${one_line//\"/\"\"}"
  echo "$site,,,,,\"$escaped\"" >> "$CSV_OUTPUT"
}

for site in "${sites[@]}"; do
  set +e
  output=$( "$FLUENCY" --site-config "$SITE_CONFIG" --site "$site" fpl run --file "$REPORT_FILE" 2>&1 )
  exit_code=$?
  set -e

  if [[ "$exit_code" -eq 0 ]]; then
    # Success: get report id, then poll fpl get until state is completed (or aborted)
    report_id=$(echo "$output" | head -n1 | tr -d '[:space:]')
    while true; do
      set +e
      status_output=$( "$FLUENCY" --site-config "$SITE_CONFIG" --site "$site" fpl get --id "$report_id" 2>&1 )
      status_code=$?
      set -e
      if [[ "$status_code" -ne 0 ]]; then
        echo "$site $status_output"
        csv_error_row "$site" "$status_output"
        break
      fi
      state=$(echo "$status_output" | sed -n 's/.*state: //p' | tr -d '[:space:]')
      if [[ "$state" == "aborted" ]]; then
        echo "Error: $site report $report_id aborted" >&2
        exit 1
      fi
      if [[ "$state" == "completed" ]]; then
        results_output=$( "$FLUENCY" --site-config "$SITE_CONFIG" --site "$site" fpl results --id "$report_id" 2>&1 )
        count=$(echo "$results_output" | jq -r '.objects[0].table.rows[0].count' 2>/dev/null)
        total=$(echo "$results_output" | jq -r '.objects[0].table.rows[0].total' 2>/dev/null)
        daily=$(echo "$results_output" | jq -r '.objects[0].table.rows[0].daily' 2>/dev/null)
        tickets=$(echo "$results_output" | jq -r '.objects[0].table.rows[0].tickets' 2>/dev/null)
        if [[ -n "$count" && "$count" != "null" && -n "$total" && "$total" != "null" ]]; then
          echo "$site count:$count daily:$daily tickets:$tickets total:$total"
          echo "$site,$count,$total,$daily,$tickets," >> "$CSV_OUTPUT"
        else
          echo "$site invalid output error"
          csv_error_row "$site" "invalid output error"
        fi
        break
      fi
      # registered, running, or other: wait 5s and poll again
      sleep 5
    done
  else
    # Failure: print site and error message, write to CSV
    echo "$site $output"
    csv_error_row "$site" "$output"
  fi
done
