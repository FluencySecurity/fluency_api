#!/bin/sh
# Write env vars so cron jobs (which run with minimal env) can use them.
if [ -n "$SLACK_WEBHOOK_URL" ]; then
  printf 'export SLACK_WEBHOOK_URL="%s"\n' "$SLACK_WEBHOOK_URL" > /app/.env.monitoring
fi
# Run crond in background to avoid setpgid (needs extra caps in Docker). Keep container alive.
crond
exec sleep infinity
