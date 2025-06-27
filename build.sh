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
