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

for secret in /run/secrets/*; do
  filename=$(basename "$secret")
  value=$(cat "$secret")

  # Normalize the secret name to uppercase
  varname=$(echo "$filename" | tr '[:lower:]' '[:upper:]')

  # If prefixed with ADGUARDHOMESYNC_, strip it
  clean_name=$(echo "$varname" | sed -E 's/^ADGUARDHOMESYNC_//I')

  # Export only the clean, normalized name
  export "$clean_name=$value"
done

exec "$@"