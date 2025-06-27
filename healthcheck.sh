#!/bin/sh

API_PORT="${API_PORT:-8080}"
URL="http://localhost:$API_PORT/control/status"
API_USERNAME="${API_USERNAME:-}"
API_PASSWORD="${API_PASSWORD:-}"

if [ -n "$API_USERNAME" ] && [ -n "$API_PASSWORD" ]; then
  STATUS=$(curl -s -u "$API_USERNAME:$API_PASSWORD" "$URL" | grep -c '"running":true')
else
  STATUS=$(curl -s "$URL" | grep -c '"running":true')
fi

if [ "$STATUS" -eq 1 ]; then
  exit 0
else
  echo "Healthcheck failed: service not running"
  exit 1
fi
