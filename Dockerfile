# Build stage: compile fluency CLI
FROM golang:1.25.0-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o fluency ./cmd/fluency/main.go

# Runtime stage: monitoring scripts + cron
FROM alpine:3.19
RUN apk add --no-cache bash jq curl dcron
COPY --from=builder /build/fluency /app/fluency
COPY scripts/ /app/scripts/
WORKDIR /app

# Wrapper so cron job runs with container env (SLACK_WEBHOOK_URL etc.)
RUN echo '#!/bin/sh' > /app/scripts/run_monitoring.sh \
    && echo '. /app/.env.monitoring 2>/dev/null' >> /app/scripts/run_monitoring.sh \
    && echo 'cd /app && /app/scripts/start_monitoring.sh --site-config /app/site_credentials.json --fluency /app/fluency' >> /app/scripts/run_monitoring.sh \
    && chmod +x /app/scripts/run_monitoring.sh

# Cron: run monitoring every 2 hours (at :00)
RUN echo '0 */2 * * * /app/scripts/run_monitoring.sh' | crontab -

# At startup: write env for cron jobs, then run crond in foreground
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
# Run via sh so the script is never executed as a binary (avoids exec format error from CRLF or shebang)
ENTRYPOINT ["/bin/sh", "/docker-entrypoint.sh"]
