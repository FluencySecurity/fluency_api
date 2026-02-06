#!/usr/bin/env bash
# Run eventwatch_status_all_sites.sh; if any site has errors (no timeline events or API failure),
# send an alert to a Slack webhook with the list of sites and their errors.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SLACK_WEBHOOK_URL="${SLACK_WEBHOOK_URL:-}"
PASS_ARGS=()

# Parse arguments (collect pass-through args for the eventwatch status script)
while [[ $# -gt 0 ]]; do
  case "$1" in
    --slack-webhook-url)
      SLACK_WEBHOOK_URL="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [--slack-webhook-url URL] [eventwatch status script options...]"
      echo ""
      echo "Runs eventwatch_status_all_sites.sh; if any site has errors, posts to Slack."
      echo ""
      echo "Slack webhook: set --slack-webhook-url or SLACK_WEBHOOK_URL environment variable."
      echo ""
      echo "Eventwatch status script options (passed through):"
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

# Run the eventwatch status check and capture output and exit code
set +e
set +u
output=$("$SCRIPT_DIR/eventwatch_status_all_sites.sh" "${PASS_ARGS[@]}" 2>&1)
exit_code=$?
set -u
set -e

if [[ "$exit_code" -eq 0 ]]; then
  echo "$output"
  exit 0
fi

# Some sites have errors: send Slack alert
echo "$output"

# Build Slack message with failed sites and reasons (from "Failed sites:" block)
failed_section=$(echo "$output" | sed -n '/^Failed sites:/,$p')
message="*Fluency EventWatch status check: one or more sites have errors*

\`\`\`
${failed_section}
\`\`\`
Run from: $(hostname) at $(date -u +"%Y-%m-%d %H:%M:%S UTC")"

payload=$(jq -n --arg text "$message" '{text: $text}')
resp=$(curl -s -w "%{http_code}" -o /tmp/eventwatch_status_slack_resp_$$.txt -X POST -H "Content-Type: application/json" -d "$payload" "$SLACK_WEBHOOK_URL")
rm -f /tmp/eventwatch_status_slack_resp_$$.txt

if [[ "$resp" != "200" ]]; then
  echo "Warning: Slack webhook returned HTTP $resp" >&2
fi

exit "$exit_code"
