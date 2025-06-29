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

# Export all secrets as environment variables
for secret in /run/secrets/*; do
  varname=$(basename "$secret" | tr '[:lower:]' '[:upper:]' | tr '-' '_')
  export "${varname}=$(cat "$secret")"
done

# Apply shared username/password to listed targets
if [ -n "$ADGUARDHOMESYNC_SHARED_TARGETS" ]; then
  IFS=',' read -r -a instances <<< "$ADGUARDHOMESYNC_SHARED_TARGETS"

  for instance in "${instances[@]}"; do
    upper_instance=$(echo "$instance" | tr '[:lower:]' '[:upper:]')

    [ -n "$ADGUARDHOME_SHARED_USERNAME" ] && export "${upper_instance}_USERNAME=${ADGUARDHOME_SHARED_USERNAME}"
    [ -n "$ADGUARDHOME_SHARED_PASSWORD" ] && export "${upper_instance}_PASSWORD=${ADGUARDHOME_SHARED_PASSWORD}"
  done
fi

# Run the main command
exec "$@"