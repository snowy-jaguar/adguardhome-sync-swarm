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

#!/bin/bash
set -euo pipefail

# Set version info dynamically
VERSION=$(date +%Y.%m)
BUILD=$(date -u +"%Y%m%dT%H%M%SZ")

# Image metadata
IMAGE_NAME="ghcr.io/snowy-jaguar/adguardhome-sync-swarm"

# Platforms to build for
PLATFORMS="linux/amd64,linux/arm64"

# Build the image
echo "🔧 Building Docker image..."
docker buildx build \
  --platform "$PLATFORMS" \
  --build-arg VERSION="$VERSION" \
  --build-arg BUILD="$BUILD" \
  --tag "$IMAGE_NAME:$VERSION" \
  --tag "$IMAGE_NAME:$VERSION-$BUILD" \
  --tag "$IMAGE_NAME:latest" \
  --file Dockerfile \
  . \
  --push

echo "✅ Build complete: $IMAGE_NAME:$VERSION-$BUILD"
