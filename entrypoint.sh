#!/bin/sh
set -e

for secret in /run/secrets/*; do
  filename=$(basename "$secret")
  value=$(cat "$secret")
  
  # Remove prefix "adguardhomesync_" case-insensitive and convert to uppercase
  varname=$(echo "$filename" | sed -E 's/^adguardhomesync_//I' | tr '[:lower:]' '[:upper:]')

  export "$varname=$value"
done

exec "$@"
