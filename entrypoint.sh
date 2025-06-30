# Copyright 2025 snowy-jaguar
# Contact: @snowyjaguar (Discord)
# Contact: contact@snowyjaguar.xyz (Email)

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#!/bin/sh
set -e

log() {
  echo "[entrypoint] $*"
}

# Export all secrets as environment variables
log "Importing secrets from /run/secrets..."
for secret in /run/secrets/*; do
  varname=$(basename "$secret" | tr '[:lower:]' '[:upper:]' | tr '-' '_')
  export "${varname}=$(cat "$secret")"
done

# Apply shared username/password to listed targets
if [ -n "$ADGUARDHOMESYNC_SHARED_TARGETS" ]; then
  IFS=',' read -r -a instances <<< "$ADGUARDHOMESYNC_SHARED_TARGETS"

  for instance in "${instances[@]}"; do
    upper_instance=$(echo "$instance" | tr '[:lower:]' '[:upper:]')

    if [ -n "$ADGUARDHOME_SHARED_USERNAME" ]; then
      export "${upper_instance}_USERNAME=${ADGUARDHOME_SHARED_USERNAME}"
      log "Set ${upper_instance}_USERNAME from shared username"
    fi

    if [ -n "$ADGUARDHOME_SHARED_PASSWORD" ]; then
      export "${upper_instance}_PASSWORD=${ADGUARDHOME_SHARED_PASSWORD}"
      log "Set ${upper_instance}_PASSWORD from shared password"
    fi

    if [ -n "$ADGUARDHOME_SHARED_CREDENTIALS" ]; then
      if echo "$ADGUARDHOME_SHARED_CREDENTIALS" | grep -q ':::'; then
        username=$(echo "$ADGUARDHOME_SHARED_CREDENTIALS" | cut -d ':' -f1)
        password=$(echo "$ADGUARDHOME_SHARED_CREDENTIALS" | cut -d ':' -f4)
      else
        IFS=$'\n' read -r username password <<EOF
$ADGUARDHOME_SHARED_CREDENTIALS
EOF
      fi
      export "${upper_instance}_USERNAME=$username"
      export "${upper_instance}_PASSWORD=$password"
      log "Set ${upper_instance}_USERNAME and _PASSWORD from shared credentials"
    fi
  done
fi

# Handle ORIGIN_CREDENTIALS if ORIGIN_USERNAME / ORIGIN_PASSWORD not set
if [ -n "$ORIGIN_CREDENTIALS" ] && [ -z "$ORIGIN_USERNAME" ] && [ -z "$ORIGIN_PASSWORD" ]; then
  if echo "$ORIGIN_CREDENTIALS" | grep -q ':::'; then
    ORIGIN_USERNAME=$(echo "$ORIGIN_CREDENTIALS" | cut -d ':' -f1)
    ORIGIN_PASSWORD=$(echo "$ORIGIN_CREDENTIALS" | cut -d ':' -f4)
  else
    IFS=$'\n' read -r ORIGIN_USERNAME ORIGIN_PASSWORD <<EOF
$ORIGIN_CREDENTIALS
EOF
  fi
  export ORIGIN_USERNAME ORIGIN_PASSWORD
  log "Set ORIGIN_USERNAME and ORIGIN_PASSWORD from ORIGIN_CREDENTIALS"
else
  log "Skipped ORIGIN_CREDENTIALS (already set or not provided)"
fi

# Handle REPLICA#_CREDENTIALS dynamically
log "Processing dynamic REPLICA# credentials..."
env | grep -E '^REPLICA[0-9]+_CREDENTIALS=' | while IFS='=' read -r full_var value; do
  prefix=$(echo "$full_var" | cut -d'_' -f1)
  user_var="${prefix}_USERNAME"
  pass_var="${prefix}_PASSWORD"

  eval existing_user=\$$user_var
  eval existing_pass=\$$pass_var

  if [ -z "$existing_user" ] && [ -z "$existing_pass" ]; then
    if echo "$value" | grep -q ':::'; then
      username=$(echo "$value" | cut -d ':' -f1)
      password=$(echo "$value" | cut -d ':' -f4)
    else
      IFS=$'\n' read -r username password <<EOF
$value
EOF
    fi
    export "${user_var}=$username"
    export "${pass_var}=$password"
    log "Set ${user_var} and ${pass_var} from ${full_var}"
  else
    log "Skipped ${full_var} (explicit credentials already set)"
  fi
done

log "Starting service..."
exec "$@"