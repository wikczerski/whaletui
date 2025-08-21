#!/bin/bash

# Build script for whaletui with automatic version injection

set -e

# Get version information from git
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')

echo "Building whaletui..."
echo "Version: $VERSION"
echo "Commit: $COMMIT_SHA"
echo "Build Date: $BUILD_DATE"

# Build with version information
go build -ldflags "-X github.com/wikczerski/whaletui/cmd.Version=$VERSION -X github.com/wikczerski/whaletui/cmd.CommitSHA=$COMMIT_SHA -X github.com/wikczerski/whaletui/cmd.BuildDate=$BUILD_DATE" -o whaletui .

echo "Build complete: whaletui"
