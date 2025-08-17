# Development commands
.PHONY: help build test test-docker clean build-all build-packages checksums release lint format

# Detect OS and set appropriate commands
ifeq ($(OS),Windows_NT)
	# Windows commands
	RM := rmdir /s /q
	MKDIR := mkdir
	RMDIR := rmdir /s /q
	CP := copy
	MV := move
	TOUCH := type nul >
	LS := dir
	CAT := type
	BINARY_EXT := .exe
	PACKAGE_EXT := .zip
	PACKAGE_CMD := powershell -Command "& {Compress-Archive -Path '$1' -DestinationPath '$2' -Force}"
	CHECKSUM_CMD := powershell -Command "& {Get-FileHash -Path '$1' -Algorithm SHA256 | Select-Object -ExpandProperty Hash}"
	CHECKSUM_FILE := checksums.sha256
	CHECKSUM_WRITE := powershell -Command "& {Get-ChildItem '$1' | ForEach-Object {Get-FileHash -Path $$_.FullName -Algorithm SHA256 | Select-Object Hash,Path | Format-Table -AutoSize | Out-File -FilePath '$2' -Encoding UTF8}}"
else
	# Unix/Linux commands
	RM := rm -rf
	MKDIR := mkdir -p
	RMDIR := rm -rf
	CP := cp
	MV := mv
	TOUCH := touch
	LS := ls -la
	CAT := cat
	AWK := awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## .*$$/ {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	FILE := file
	SIZE := ls -lh
	BINARY_EXT :=
	PACKAGE_EXT := .tar.gz
	PACKAGE_CMD := tar -czf
	CHECKSUM_CMD := sha256sum
	CHECKSUM_FILE := checksums.sha256
	CHECKSUM_WRITE := for file in $1/*; do if [ -f "$$file" ]; then sha256sum "$$file"; fi; done > $1/../$2
endif

# Default target
help: ## Show this help message
	@echo "Available commands:"
ifeq ($(OS),Windows_NT)
	@echo "  build        - Build the application for current platform"
	@echo "  test         - Run all tests locally"
	@echo "  test-docker  - Run tests in Docker container"
	@echo "  lint         - Run golangci-lint to check for code quality issues"
	@echo "  format       - Format Go code using goimports"
	@echo "  clean        - Clean build artifacts"
	@echo "  build-all    - Build for all supported platforms"
	@echo "  build-packages - Build packages for all platforms"
	@echo "  checksums    - Generate checksums for all packages"
	@echo "  release      - Prepare complete release"
else
	@$(AWK) $(MAKEFILE_LIST)
endif

# Build variables
BINARY_NAME=d5r
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_DIR=dist
LDFLAGS=-ldflags="-s -w -X main.Version=$(VERSION)"

build: ## Build the application for current platform
	@echo "Building for current platform..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME)$(BINARY_EXT) main.go
ifeq ($(OS),Windows_NT)
	@if exist "bin\$(BINARY_NAME)$(BINARY_EXT)" (echo Binary created: bin\$(BINARY_NAME)$(BINARY_EXT)) else (echo ERROR: Binary not created && exit /b 1)
	@echo "✓ Binary created: bin/$(BINARY_NAME)$(BINARY_EXT)"
	@powershell -Command "& {Get-Item 'bin\$(BINARY_NAME)$(BINARY_EXT)' | Select-Object -ExpandProperty Length | ForEach-Object {Write-Host ('  Size: ' + [math]::Round($$_/1MB,2) + ' MB')}}"
else
	@if [ ! -f "bin/$(BINARY_NAME)$(BINARY_EXT)" ]; then echo "ERROR: Binary not created"; exit 1; fi
	@echo "✓ Binary created: bin/$(BINARY_NAME)$(BINARY_EXT)"
	@echo "  Size: $(shell ls -lh bin/$(BINARY_NAME)$(BINARY_EXT) | awk '{print $$5}')"
	@echo "  Type: $(shell file bin/$(BINARY_NAME)$(BINARY_EXT))"
endif

test: ## Run all tests locally
	go test ./... -v

test-docker: ## Run tests in Docker container
	docker run --rm -v "${CURDIR}:/app" -w /app golang:1.25.0-alpine sh -c "apk add --no-cache git && go test ./... -v"

lint: ## Run golangci-lint to check for code quality issues
ifeq ($(OS),Windows_NT)
	$(shell go env GOPATH)/bin/golangci-lint.exe run
else
	$(shell go env GOPATH)/bin/golangci-lint run
endif

format: ## Format Go code using goimports
ifeq ($(OS),Windows_NT)
	$(shell go env GOPATH)/bin/goimports.exe -w .
else
	$(shell go env GOPATH)/bin/goimports -w .
endif

clean: ## Clean build artifacts
ifeq ($(OS),Windows_NT)
	@if exist bin rmdir /s /q bin
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
else
	$(RM) bin/ $(BUILD_DIR)/
endif

# Cross-compilation targets
build-all: ## Build for all supported platforms (Linux: full cross-compilation, Windows: current platform only)
	@echo "Building for all platforms..."
	@echo "OS variable: '$(OS)'"
	@echo "OS detection: $(if $(filter Windows_NT,$(OS)),Windows,Linux/Unix)"
ifeq ($(OS),Windows_NT)
	@echo "Windows detected - building for current platform only"
	@echo "Cleaning previous build artifacts..."
	@if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
	@if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
	@echo "Building for Windows..."
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	@echo "✓ Windows binary built successfully"
	@echo "Build complete. Binary is in $(BUILD_DIR)/"
else
	@echo "Linux/Unix detected - building for all platforms"
	@echo "Cleaning previous build artifacts..."
	@rm -rf $(BUILD_DIR)/
	@mkdir -p $(BUILD_DIR)
	
	@echo "Starting Linux builds..."
	@echo "Building Linux amd64..."
ifeq ($(OS),Windows_NT)
	@set GOOS=linux && set GOARCH=amd64 && go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
else
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
endif
ifeq ($(OS),Windows_NT)
	@if exist "$(BUILD_DIR)\$(BINARY_NAME)-linux-amd64" (echo ✓ Linux amd64 built successfully) else (echo ERROR: Linux amd64 binary not created && exit /b 1)
else
	@if [ ! -f "$(BUILD_DIR)/$(BINARY_NAME)-linux-amd64" ]; then echo "ERROR: Linux amd64 binary not created"; exit 1; fi
	@echo "✓ Linux amd64 built successfully: $(shell ls -lh $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64)"
endif
	
	@echo "Building Linux arm64..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 main.go
ifeq ($(OS),Windows_NT)
	@if exist "$(BUILD_DIR)\$(BINARY_NAME)-linux-arm64" (echo ✓ Linux arm64 built successfully) else (echo ERROR: Linux arm64 binary not created && exit /b 1)
else
	@if [ ! -f "$(BUILD_DIR)/$(BINARY_NAME)-linux-arm64" ]; then echo "ERROR: Linux arm64 binary not created"; exit 1; fi
	@echo "✓ Linux arm64 built successfully: $(shell ls -lh $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64)"
endif
	
	@echo "Building Linux armv7..."
	@GOOS=linux GOARCH=arm GOARM=7 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-armv7 main.go
ifeq ($(OS),Windows_NT)
	@if exist "$(BUILD_DIR)\$(BINARY_NAME)-linux-armv7" (echo ✓ Linux armv7 built successfully) else (echo ERROR: Linux armv7 binary not created && exit /b 1)
else
	@if [ ! -f "$(BUILD_DIR)/$(BINARY_NAME)-linux-armv7" ]; then echo "ERROR: Linux armv7 binary not created"; exit 1; fi
	@echo "✓ Linux armv7 built successfully: $(shell ls -lh $(BUILD_DIR)/$(BINARY_NAME)-linux-armv7)"
endif
	
	@echo "Building Linux ppc64le..."
	@GOOS=linux GOARCH=ppc64le go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-ppc64le main.go
ifeq ($(OS),Windows_NT)
	@if exist "$(BUILD_DIR)\$(BINARY_NAME)-linux-ppc64le" (echo ✓ Linux ppc64le built successfully) else (echo ERROR: Linux ppc64le binary not created && exit /b 1)
else
	@if [ ! -f "$(BUILD_DIR)/$(BINARY_NAME)-linux-ppc64le" ]; then echo "ERROR: Linux ppc64le binary not created"; exit 1; fi
	@echo "✓ Linux ppc64le built successfully: $(shell ls -lh $(BUILD_DIR)/$(BINARY_NAME)-linux-ppc64le)"
endif
	
	@echo "Building Linux s390x..."
	@GOOS=linux GOARCH=s390x go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-s390x main.go
ifeq ($(OS),Windows_NT)
	@if exist "$(BUILD_DIR)\$(BINARY_NAME)-linux-s390x" (echo ✓ Linux s390x built successfully) else (echo ERROR: Linux s390x binary not created && exit /b 1)
else
	@if [ ! -f "$(BUILD_DIR)/$(BINARY_NAME)-linux-s390x" ]; then echo "ERROR: Linux s390x binary not created"; exit 1; fi
	@echo "✓ Linux s390x built successfully: $(shell ls -lh $(BUILD_DIR)/$(BINARY_NAME)-linux-s390x)"
endif
	
	@echo "Starting Darwin builds..."
	@echo "Building Darwin amd64..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go || (echo "Failed to build Darwin amd64"; exit 1)
	@echo "Building Darwin arm64..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go || (echo "Failed to build Darwin arm64"; exit 1)
	
	@echo "Starting FreeBSD builds..."
	@echo "Building FreeBSD amd64..."
	@GOOS=freebsd GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-amd64 main.go || (echo "Failed to build FreeBSD amd64"; exit 1)
	@echo "Building FreeBSD arm64..."
	@GOOS=freebsd GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-arm64 main.go || (echo "Failed to build FreeBSD arm64"; exit 1)
	
	@echo "Starting Windows builds..."
	@echo "Building Windows amd64..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go || (echo "Failed to build Windows amd64"; exit 1)
	@echo "Building Windows arm64..."
	@GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe main.go || (echo "Failed to build Windows arm64"; exit 1)
	
	@echo "Build complete. Binaries are in $(BUILD_DIR)/"
	
	# Final verification that all expected binaries exist
	@echo "Final verification:"
ifeq ($(OS),Windows_NT)
	@for %%a in (amd64 arm64 armv7 ppc64le s390x) do ( \
		if exist "$(BUILD_DIR)\$(BINARY_NAME)-linux-%%a" ( \
			echo ✓ $(BINARY_NAME)-linux-%%a: exists \
		) else ( \
			echo ✗ $(BINARY_NAME)-linux-%%a: MISSING && exit /b 1 \
		) \
	)
else
	@for arch in amd64 arm64 armv7 ppc64le s390x; do \
		if [ -f "$(BUILD_DIR)/$(BINARY_NAME)-linux-$$arch" ]; then \
			echo "✓ $(BINARY_NAME)-linux-$$arch: $(shell ls -lh $(BUILD_DIR)/$(BINARY_NAME)-linux-$$arch | awk '{print $$5}')"; \
		else \
			echo "✗ $(BINARY_NAME)-linux-$$arch: MISSING"; \
			exit 1; \
		fi; \
	done
endif
	@echo "All Linux binaries verified successfully!"
	@echo "Final directory listing:"
ifeq ($(OS),Windows_NT)
	@dir $(BUILD_DIR)
else
	@ls -la $(BUILD_DIR) 2>/dev/null || echo "Directory listing failed, but build was successful"
endif
endif

# Package creation targets
build-packages: build-all ## Build packages for all platforms
	@echo "Creating packages..."
ifeq ($(OS),Windows_NT)
	@if not exist "$(BUILD_DIR)\packages" mkdir "$(BUILD_DIR)\packages"
	@echo "Windows detected - creating Windows packages only"
	@echo "Creating Windows amd64 package..."
	@cd $(BUILD_DIR) && powershell -Command "& {Compress-Archive -Path '$(BINARY_NAME)-windows-amd64.exe' -DestinationPath 'packages\$(BINARY_NAME)-v$(VERSION)-windows-amd64$(PACKAGE_EXT)' -Force}" || (echo "Failed to create Windows amd64 package" && exit /b 1)
	@echo "Windows packages created successfully!"
else
	@mkdir -p $(BUILD_DIR)/packages
	@echo "Linux/Unix detected - creating packages for all platforms"
	@echo "Creating Linux packages..."
	@echo "Creating Linux amd64 package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-linux-amd64$(PACKAGE_EXT) $(BINARY_NAME)-linux-amd64 || (echo "Failed to create Linux amd64 package"; exit 1)
	@echo "Creating Linux arm64 package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-linux-arm64$(PACKAGE_EXT) $(BINARY_NAME)-linux-arm64 || (echo "Failed to create Linux arm64 package"; exit 1)
	@echo "Creating Linux armv7 package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-linux-armv7$(PACKAGE_EXT) $(BINARY_NAME)-linux-armv7 || (echo "Failed to create Linux armv7 package"; exit 1)
	@echo "Creating Linux ppc64le package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-linux-ppc64le$(PACKAGE_EXT) $(BINARY_NAME)-linux-ppc64le || (echo "Failed to create Linux ppc64le package"; exit 1)
	@echo "Creating Linux s390x package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-linux-s390x$(PACKAGE_EXT) $(BINARY_NAME)-linux-s390x || (echo "Failed to create Linux s390x package"; exit 1)
	
	@echo "Creating Darwin packages..."
	@echo "Creating Darwin amd64 package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-darwin-amd64$(PACKAGE_EXT) $(BINARY_NAME)-darwin-amd64 || (echo "Failed to create Darwin amd64 package"; exit 1)
	@echo "Creating Darwin arm64 package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-darwin-arm64$(PACKAGE_EXT) $(BINARY_NAME)-darwin-arm64 || (echo "Failed to create Darwin arm64 package"; exit 1)
	
	@echo "Creating FreeBSD packages..."
	@echo "Creating FreeBSD amd64 package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-freebsd-amd64$(PACKAGE_EXT) $(BINARY_NAME)-freebsd-amd64 || (echo "Failed to create FreeBSD amd64 package"; exit 1)
	@echo "Creating FreeBSD arm64 package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-freebsd-arm64$(PACKAGE_EXT) $(BINARY_NAME)-freebsd-arm64 || (echo "Failed to create FreeBSD arm64 package"; exit 1)
	
	@echo "Creating Windows packages..."
	@echo "Creating Windows amd64 package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-windows-amd64$(PACKAGE_EXT) $(BINARY_NAME)-windows-amd64.exe || (echo "Failed to create Windows amd64 package"; exit 1)
	@echo "Creating Windows arm64 package..."
	@cd $(BUILD_DIR) && tar -czf packages/$(BINARY_NAME)-v$(VERSION)-windows-arm64$(PACKAGE_EXT) $(BINARY_NAME)-windows-arm64.exe || (echo "Failed to create Windows arm64 package"; exit 1)
	
	@echo "Packages created in $(BUILD_DIR)/packages/"
	@echo "Verifying all packages:"
ifeq ($(OS),Windows_NT)
	@dir $(BUILD_DIR)\packages
else
	@ls -la $(BUILD_DIR)/packages/ 2>/dev/null || echo "Packages directory listing failed"
endif
endif

# Checksums
checksums: build-packages ## Generate checksums for all packages
	@echo "Generating checksums..."
ifeq ($(OS),Windows_NT)
	@$(CHECKSUM_WRITE) $(BUILD_DIR)\packages $(BUILD_DIR)\$(CHECKSUM_FILE)
else
	@cd $(BUILD_DIR) && sha256sum packages/* > checksums.sha256
endif
	@echo "Checksums generated: $(BUILD_DIR)/$(CHECKSUM_FILE)"

# Release preparation
release: checksums ## Prepare complete release
	@echo "Release prepared in $(BUILD_DIR)/"
ifeq ($(OS),Windows_NT)
	@dir $(BUILD_DIR)
	@dir $(BUILD_DIR)\packages
else
	@ls -la $(BUILD_DIR)/ 2>/dev/null || echo "Main directory listing failed"
	@ls -la $(BUILD_DIR)/packages/ 2>/dev/null || echo "Packages directory listing failed"
endif
