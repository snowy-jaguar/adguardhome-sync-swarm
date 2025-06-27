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