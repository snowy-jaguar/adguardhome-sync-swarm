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

#!/bin/sh

API_PORT="${API_PORT:-8080}"
URL="http://localhost:$API_PORT/control"
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
