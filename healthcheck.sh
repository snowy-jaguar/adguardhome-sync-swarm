#!/bin/sh

API_PORT="${API_PORT:-8080}"
URL="http://localhost:$API_PORT/control/status"
API_USERNAME="${API_USERNAME:-}"
API_PASSWORD="${API_PASSWORD:-}"

# Compose curl auth flag only if both username and password are set
AUTH=""
if [ -n "$API_USERNAME" ] && [ -n "$API_PASSWORD" ]; then
  AUTH="$API_USERNAME:$API_PASSWORD"
fi

# Make the request
RESPONSE=$(curl -s ${AUTH:+-u "$AUTH"} "$URL")
STATUS=$(echo "$RESPONSE" | grep -c '"running":true')

# Check status
if [ "$STATUS" -eq 1 ]; then
  exit 0
else
  echo "Healthcheck failed: service not running. Response: $RESPONSE"
  exit 1
fi