#!/bin/sh

# Copyright 2025 snowy-jaguar
# Contact: @snowyjaguar (Discord)
# Contact: contact@snowyjaguar.xyz (Email)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

API_PORT="${API_PORT:-8080}"
URL="http://localhost:$API_PORT/api/v1/status"
API_USERNAME="${API_USERNAME:-}"
API_PASSWORD="${API_PASSWORD:-}"

# Log function
log() {
  echo "[healthcheck] $*"
}

# Compose curl auth flag only if both username and password are set
AUTH=""
if [ -n "$API_USERNAME" ] && [ -n "$API_PASSWORD" ]; then
  AUTH="$API_USERNAME:$API_PASSWORD"
fi

# Check if the sync container is running and reachable
log "Checking sync container status..."

SYNC_RESPONSE=$(curl -s -w "%{http_code}" -o /dev/null ${AUTH:+-u "$AUTH"} "$URL")
if [ "$SYNC_RESPONSE" -ne 200 ]; then
  log "Healthcheck failed: Sync container authentication failed (HTTP $SYNC_RESPONSE)."
  exit 1
fi

log "Sync container authentication successful."

# Check if the sync container is running and reachable
SYNC_RESPONSE=$(curl -s ${AUTH:+-u "$AUTH"} "$URL")
if [ -z "$SYNC_RESPONSE" ]; then
  log "Healthcheck failed: Sync container is unreachable."
  exit 1
fi

log "Sync container is reachable."

# Check if the Origin instance is reachable
log "Checking origin instance ($ORIGIN_URL)..."
check_instance() {
  local host=$1
  local url="http://$host:3000"  # Assuming AdGuardHome is running on port 3000
  RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$url")
  if [ "$RESPONSE" -ne 200 ]; then
    log "AdGuardHome instance $host is not responsive (HTTP $RESPONSE)."
    return 1
  fi
  log "AdGuardHome instance $host is responsive."
  return 0
}

# Check the origin URL (static, defined as ORIGIN_URL)
check_instance "$ORIGIN_URL"

# Check dynamic replica URLs (REPLICA1_URL, REPLICA2_URL, etc.)
log "Checking replica instances..."
replica_counter=1
while true; do
  replica_url_var="REPLICA${replica_counter}_URL"
  replica_url="${!replica_url_var}"
  
  if [ -z "$replica_url" ]; then
    break  # No more replica URLs found, exit the loop
  fi

  check_instance "$replica_url" || true  # Allow non-fatal failures

  replica_counter=$((replica_counter + 1))
done

# If everything is fine, mark the healthcheck as successful
log "Healthcheck passed."
exit 0
