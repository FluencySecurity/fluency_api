#!/usr/bin/env bash
# Run all monitoring Slack alerts: platform status, eventwatch status, and audit DB status.
# Passes --slack-webhook-url, --site-config, --fluency, and -v to each script.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SLACK_WEBHOOK_URL="${SLACK_WEBHOOK_URL:-}"
PASS_ARGS=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --slack-webhook-url)
      SLACK_WEBHOOK_URL="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [--slack-webhook-url URL] [--site-config PATH] [--fluency PATH] [-v|--verbose]"
      echo ""
      echo "Runs all monitoring Slack alerts in sequence:"
      echo "  1. platform_status_slack_alert.sh"
      echo "  2. eventwatch_status_slack_alert.sh"
      echo "  3. audit_db_status_slack_alert.sh"
      echo ""
      echo "Slack webhook: set --slack-webhook-url or SLACK_WEBHOOK_URL environment variable."
      echo ""
      echo "Options (passed to each alert script):"
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

# Build args for child scripts: webhook + any pass-through
ARGS=("--slack-webhook-url" "$SLACK_WEBHOOK_URL" "${PASS_ARGS[@]}")

exit_code=0

echo "=== Platform status ==="
if ! "$SCRIPT_DIR/platform_status_slack_alert.sh" "${ARGS[@]}"; then
  exit_code=1
fi

echo ""
echo "=== EventWatch status ==="
if ! "$SCRIPT_DIR/eventwatch_status_slack_alert.sh" "${ARGS[@]}"; then
  exit_code=1
fi

echo ""
echo "=== Audit DB status ==="
if ! "$SCRIPT_DIR/audit_db_status_slack_alert.sh" "${ARGS[@]}"; then
  exit_code=1
fi

exit "$exit_code"
