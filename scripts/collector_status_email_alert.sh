#!/usr/bin/env bash
# Run collector_status_all.sh; if any collector has errors, send an SMTP email alert.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SMTP_HOST="${SMTP_HOST:-}"
SMTP_PORT="${SMTP_PORT:-587}"
SMTP_USER="${SMTP_USER:-}"
SMTP_PASS="${SMTP_PASS:-}"
SMTP_FROM="${SMTP_FROM:-}"
SMTP_TO="${SMTP_TO:-}"
SMTP_FROM_NAME="${SMTP_FROM_NAME:-Fluency Collector Monitor}"
# Set to 0 or use --smtp-no-starttls if the server does not support STARTTLS
SMTP_USE_STARTTLS="${SMTP_USE_STARTTLS:-1}"
BODY_FILE=""
PASS_ARGS=()

# Parse arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --body-file)
      BODY_FILE="$2"
      shift 2
      ;;
    --smtp-no-starttls)
      SMTP_USE_STARTTLS="0"
      shift
      ;;
    --smtp-host)
      SMTP_HOST="$2"
      shift 2
      ;;
    --smtp-port)
      SMTP_PORT="$2"
      shift 2
      ;;
    --smtp-user)
      SMTP_USER="$2"
      shift 2
      ;;
    --smtp-pass)
      SMTP_PASS="$2"
      shift 2
      ;;
    --smtp-from)
      SMTP_FROM="$2"
      shift 2
      ;;
    --smtp-to)
      SMTP_TO="$2"
      shift 2
      ;;
    --smtp-from-name)
      SMTP_FROM_NAME="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: $0 [SMTP options...] [collector status script options...]"
      echo ""
      echo "Runs collector_status_all.sh; if any collector has errors, sends SMTP email."
      echo ""
      echo "Requires: python3 (with smtplib)"
      echo ""
      echo "SMTP options (can also be set via environment variables):"
      echo "  --smtp-host HOST       SMTP server hostname (SMTP_HOST)"
      echo "  --smtp-port PORT       SMTP server port (default: 587, SMTP_PORT)"
      echo "  --smtp-user USER       SMTP username (SMTP_USER)"
      echo "  --smtp-pass PASS       SMTP password (SMTP_PASS)"
      echo "  --smtp-from EMAIL      From email address (SMTP_FROM)"
      echo "  --smtp-to EMAIL        To email address (SMTP_TO)"
      echo "  --smtp-from-name NAME  From name (default: 'Fluency Collector Monitor', SMTP_FROM_NAME)"
      echo "  --smtp-no-starttls     Disable STARTTLS (use if server does not support it; or SMTP_USE_STARTTLS=0)"
      echo "  --body-file PATH       Write email body to this file (default: /tmp/collector_status_<timestamp>.txt)"
      echo ""
      echo "Collector status script options (passed through):"
      echo "  --site-config PATH     Path to site_credentials.json"
      echo "  --site NAME            Site name"
      echo "  --fluency PATH         Path to fluency binary"
      echo "  -v, --verbose          Verbose output"
      exit 0
      ;;
    *)
      PASS_ARGS+=("$1")
      shift
      ;;
  esac
done

# Validate required SMTP settings
if [[ -z "$SMTP_HOST" ]] || [[ -z "$SMTP_FROM" ]] || [[ -z "$SMTP_TO" ]]; then
  echo "Error: SMTP configuration required." >&2
  echo "Set SMTP_HOST, SMTP_FROM, SMTP_TO (and optionally SMTP_PORT, SMTP_USER, SMTP_PASS)" >&2
  echo "or use --smtp-host, --smtp-from, --smtp-to flags." >&2
  exit 2
fi

# Run the collector status check and capture output and exit code
set +e
set +u
output=$("$SCRIPT_DIR/collector_status_all.sh" "${PASS_ARGS[@]}" 2>&1)
exit_code=$?
set -u
set -e

if [[ "$exit_code" -eq 0 ]]; then
  echo "$output"
  exit 0
fi

# Some collectors have errors: send email alert
echo "$output"

# Build email message with failed collectors and reasons (from "Failed collectors:" block)
failed_section=$(echo "$output" | sed -n '/^Failed collectors:/,$p')
subject="Fluency Collector Status Alert: Failed Collectors"
body="Fluency Collector Status Check Alert

One or more collectors have errors:

${failed_section}

Run from: $(hostname) at $(date -u +"%Y-%m-%d %H:%M:%S UTC")

---
Full output:
${output}"

# Write body to file (default: /tmp/collector_status_<timestamp>.txt)
body_file="${BODY_FILE:-/tmp/collector_status_$(date +%Y%m%d-%H%M%S).txt}"
printf '%s' "$body" > "$body_file"
echo "Email body written to: $body_file" >&2

# Send email using Python smtplib
if ! command -v python3 &>/dev/null; then
  echo "Error: python3 is required to send email" >&2
  exit 2
fi

SMTP_USER_VAL="${SMTP_USER:-}"
SMTP_PASS_VAL="${SMTP_PASS:-}"
SMTP_USE_STARTTLS_VAL="${SMTP_USE_STARTTLS:-1}"
python3 <<EOF
import sys
import smtplib
from email.mime.text import MIMEText
from email.header import Header

smtp_host = "${SMTP_HOST}"
smtp_port = ${SMTP_PORT}
smtp_user = "${SMTP_USER_VAL}"
smtp_pass = "${SMTP_PASS_VAL}"
smtp_from = "${SMTP_FROM}"
smtp_to = "${SMTP_TO}"
smtp_from_name = "${SMTP_FROM_NAME}"
use_starttls = "${SMTP_USE_STARTTLS_VAL}" == "1"
subject = "${subject}"
body = """${body}"""

msg = MIMEText(body, 'plain', 'utf-8')
msg['From'] = Header(f"{smtp_from_name} <{smtp_from}>", 'utf-8')
msg['To'] = Header(smtp_to, 'utf-8')
msg['Subject'] = Header(subject, 'utf-8')

try:
    server = smtplib.SMTP(smtp_host, smtp_port)
    if use_starttls:
        server.starttls()
    if smtp_user:
        server.login(smtp_user, smtp_pass)
    server.sendmail(smtp_from, [smtp_to], msg.as_string())
    server.quit()
except Exception as e:
    print(f"Failed to send email: {e}", file=sys.stderr)
    sys.exit(1)
EOF

if [[ $? -ne 0 ]]; then
  echo "Error: Failed to send email" >&2
  exit 2
fi

exit "$exit_code"
