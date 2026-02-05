#!/usr/bin/env bash
# Run billing_run_all_sites.sh and send the generated CSV to a Slack webhook as an attachment.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
SLACK_WEBHOOK_URL="${SLACK_WEBHOOK_URL:-}"
CSV_PATH=""   # default set after run
PASS_ARGS=()

# Parse arguments (collect pass-through args for the billing script)
while [[ $# -gt 0 ]]; do
  case "$1" in
    --slack-webhook-url)
      SLACK_WEBHOOK_URL="$2"
      shift 2
      ;;
    --csv)
      CSV_PATH="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [--slack-webhook-url URL] [--csv PATH] [billing script options...]"
      echo ""
      echo "Runs billing_run_all_sites.sh and posts the CSV as a Slack attachment."
      echo ""
      echo "Slack webhook: set --slack-webhook-url or SLACK_WEBHOOK_URL environment variable."
      echo "  --csv PATH   CSV file to attach (default: repo root billing_results.csv)"
      echo ""
      echo "Billing script options (passed through):"
      echo "  --site-config PATH   Path to site_credentials.json"
      echo "  --report-file PATH   Path to billing report JSON"
      echo "  --output PATH        CSV output path"
      echo "  --fluency PATH       Path to fluency binary"
      exit 0
      ;;
    *)
      PASS_ARGS+=("$1")
      shift
      ;;
  esac
done

if [[ -z "$SLACK_WEBHOOK_URL" ]]; then
  echo "Error: Slack webhook URL required. Set SLACK_WEBHOOK_URL or use --slack-webhook-url URL." >&2
  exit 2
fi

# Run billing script from repo root so relative paths (reports/billing.json, etc.) work
# set +u so empty PASS_ARGS does not trigger "unbound variable" when no extra args passed
set +e
set +u
output=$(cd "$REPO_ROOT" && "$SCRIPT_DIR/billing_run_all_sites.sh" "${PASS_ARGS[@]}" 2>&1)
exit_code=$?
set -u
set -e

# Always print output to console
echo "$output"

# Resolve CSV path: explicit --csv, or default at repo root (billing script default when run from REPO_ROOT)
if [[ -z "$CSV_PATH" ]]; then
  CSV_PATH="$REPO_ROOT/billing_results.csv"
fi

# Build Slack message with CSV as attachment
header="*Fluency Billing Report*
Run from: $(hostname) at $(date -u +"%Y-%m-%d %H:%M:%S UTC")"

if [[ -f "$CSV_PATH" ]]; then
  csv_content=$(cat "$CSV_PATH")
  csv_filename=$(basename "$CSV_PATH")
  # Wrap in code block so Slack shows monospace and title doesn't run into first line
  csv_block=$'```\n'"$csv_content"$'\n```'
  # Attachment: title = filename, text = CSV in code block. (Incoming Webhooks cannot send downloadable files; use Slack Web API files.upload with a bot token for that.)
  payload=$(jq -n \
    --arg text "$header" \
    --arg title "$csv_filename" \
    --arg fallback "billing_results.csv" \
    --arg body "$csv_block" \
    '{text: $text, attachments: [{fallback: $fallback, title: $title, text: $body}]}')
else
  payload=$(jq -n --arg text "${header}

(CSV file not found: ${CSV_PATH})" '{text: $text}')
fi

resp=$(curl -s -w "%{http_code}" -o /tmp/billing_slack_resp_$$.txt -X POST -H "Content-Type: application/json" -d "$payload" "$SLACK_WEBHOOK_URL")
rm -f /tmp/billing_slack_resp_$$.txt

if [[ "$resp" != "200" ]]; then
  echo "Warning: Slack webhook returned HTTP $resp" >&2
fi

exit "$exit_code"
