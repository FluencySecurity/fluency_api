#!/usr/bin/env bash
# Run the db_status check for all sites; if any site is down, send an alert to a Slack webhook.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SLACK_WEBHOOK_URL="${SLACK_WEBHOOK_URL:-}"
PASS_ARGS=()

# Parse arguments (collect pass-through args for the audit script)
while [[ $# -gt 0 ]]; do
  case "$1" in
    --slack-webhook-url)
      SLACK_WEBHOOK_URL="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [--slack-webhook-url URL] [audit script options...]"
      echo ""
      echo "Runs audit_db_status_all_sites.sh; if any site is down, posts to Slack."
      echo ""
      echo "Slack webhook: set --slack-webhook-url or SLACK_WEBHOOK_URL environment variable."
      echo ""
      echo "Audit script options (passed through):"
      echo "  --site-config PATH   Path to site_credentials.json"
      echo "  --fluency PATH       Path to fluency binary"
      echo "  -v, --verbose        Verbose output"
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

# Run the db_status check and capture output and exit code
set +e
set +u
output=$("$SCRIPT_DIR/audit_db_status_all_sites.sh" "${PASS_ARGS[@]}" 2>&1)
exit_code=$?
set -u
set -e

if [[ "$exit_code" -eq 0 ]]; then
  echo "$output"
  exit 0
fi

# Something is down: send Slack alert
echo "$output"

# Build Slack message with down sites and reasons (from "Down sites (with reason):" block)
down_section=$(echo "$output" | sed -n '/^Down sites (with reason):/,$p')
message="*Fluency DB status check: one or more sites down*

\`\`\`
${down_section}
\`\`\`
Run from: $(hostname) at $(date -u +"%Y-%m-%d %H:%M:%S UTC")"

payload=$(jq -n --arg text "$message" '{text: $text}')
resp=$(curl -s -w "%{http_code}" -o /tmp/slack_alert_resp_$$.txt -X POST -H "Content-Type: application/json" -d "$payload" "$SLACK_WEBHOOK_URL")
rm -f /tmp/slack_alert_resp_$$.txt

if [[ "$resp" != "200" ]]; then
  echo "Warning: Slack webhook returned HTTP $resp" >&2
fi

exit "$exit_code"
