#!/bin/bash

# Build script for d5r
# This script helps with local testing of the build process

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="d5r"
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_DIR="dist"

echo -e "${GREEN}Building $BINARY_NAME version $VERSION${NC}"

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -rf "$BUILD_DIR" bin/

# Create build directory
mkdir -p "$BUILD_DIR"

# Build variables
LDFLAGS="-ldflags=-s -w -X main.Version=$VERSION"

echo -e "${YELLOW}Building for all platforms...${NC}"

# Linux builds
echo "Building Linux binaries..."
GOOS=linux GOARCH=amd64 go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-linux-amd64" main.go
GOOS=linux GOARCH=arm64 go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-linux-arm64" main.go
GOOS=linux GOARCH=arm GOARM=7 go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-linux-armv7" main.go
GOOS=linux GOARCH=ppc64le go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-linux-ppc64le" main.go
GOOS=linux GOARCH=s390x go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-linux-s390x" main.go

# Darwin builds
echo "Building Darwin binaries..."
GOOS=darwin GOARCH=amd64 go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-darwin-amd64" main.go
GOOS=darwin GOARCH=arm64 go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-darwin-arm64" main.go

# FreeBSD builds
echo "Building FreeBSD binaries..."
GOOS=freebsd GOARCH=amd64 go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-freebsd-amd64" main.go
GOOS=freebsd GOARCH=arm64 go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-freebsd-arm64" main.go

# Windows builds
echo "Building Windows binaries..."
GOOS=windows GOARCH=amd64 go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-windows-amd64.exe" main.go
GOOS=windows GOARCH=arm64 go build $LDFLAGS -o "$BUILD_DIR/$BINARY_NAME-windows-arm64.exe" main.go

echo -e "${GREEN}Build complete!${NC}"
echo "Binaries are in $BUILD_DIR/"

# Create packages directory
mkdir -p "$BUILD_DIR/packages"

echo -e "${YELLOW}Creating packages...${NC}"

# Linux packages
echo "Creating Linux packages..."
cd "$BUILD_DIR"
tar -czf "packages/$BINARY_NAME-v$VERSION-linux-amd64.tar.gz" "$BINARY_NAME-linux-amd64"
tar -czf "packages/$BINARY_NAME-v$VERSION-linux-arm64.tar.gz" "$BINARY_NAME-linux-arm64"
tar -czf "packages/$BINARY_NAME-v$VERSION-linux-armv7.tar.gz" "$BINARY_NAME-linux-armv7"
tar -czf "packages/$BINARY_NAME-v$VERSION-linux-ppc64le.tar.gz" "$BINARY_NAME-linux-ppc64le"
tar -czf "packages/$BINARY_NAME-v$VERSION-linux-s390x.tar.gz" "$BINARY_NAME-linux-s390x"

# Darwin packages
echo "Creating Darwin packages..."
tar -czf "packages/$BINARY_NAME-v$VERSION-darwin-amd64.tar.gz" "$BINARY_NAME-darwin-amd64"
tar -czf "packages/$BINARY_NAME-v$VERSION-darwin-arm64.tar.gz" "$BINARY_NAME-darwin-arm64"

# FreeBSD packages
echo "Creating FreeBSD packages..."
tar -czf "packages/$BINARY_NAME-v$VERSION-freebsd-amd64.tar.gz" "$BINARY_NAME-freebsd-amd64"
tar -czf "packages/$BINARY_NAME-v$VERSION-freebsd-arm64.tar.gz" "$BINARY_NAME-freebsd-arm64"

# Windows packages
echo "Creating Windows packages..."
if command -v zip >/dev/null 2>&1; then
    zip "packages/$BINARY_NAME-v$VERSION-windows-amd64.zip" "$BINARY_NAME-windows-amd64.exe"
    zip "packages/$BINARY_NAME-v$VERSION-windows-arm64.zip" "$BINARY_NAME-windows-arm64.exe"
else
    echo -e "${YELLOW}Warning: zip command not found, skipping Windows packages${NC}"
fi

cd ..

echo -e "${GREEN}Packages created!${NC}"

# Generate checksums
echo -e "${YELLOW}Generating checksums...${NC}"
cd "$BUILD_DIR/packages"
if command -v sha256sum >/dev/null 2>&1; then
    sha256sum * > ../checksums.sha256
elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 * > ../checksums.sha256
else
    echo -e "${RED}Error: Neither sha256sum nor shasum found${NC}"
    exit 1
fi
cd ../..

echo -e "${GREEN}Checksums generated: $BUILD_DIR/checksums.sha256${NC}"

# Show results
echo -e "${GREEN}Build Summary:${NC}"
echo "Binaries: $(ls -1 "$BUILD_DIR"/* | grep -v packages | wc -l)"
echo "Packages: $(ls -1 "$BUILD_DIR/packages"/* | wc -l)"
echo ""
echo "Files created:"
ls -la "$BUILD_DIR/"
echo ""
echo "Packages created:"
ls -la "$BUILD_DIR/packages/"
echo ""
echo "Checksums:"
cat "$BUILD_DIR/checksums.sha256"

echo -e "${GREEN}Build process completed successfully!${NC}"
