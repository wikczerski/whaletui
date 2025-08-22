# For local development

# Build the executable for current platform with whaletui as the name
build:
	go build -o whaletui.exe .

# Build with automatic version injection (Unix/Linux/macOS)
build-versioned:
	@chmod +x scripts/build.sh
	./scripts/build.sh

# Build with automatic version injection (Windows)
build-versioned-windows:
	scripts\build.bat

# Build for release with version injection
release: build-versioned

# Build for Windows release
release-windows: build-versioned-windows

# Goreleaser commands
goreleaser-build:
	goreleaser build --snapshot --clean

goreleaser-release:
	goreleaser release --snapshot --clean

goreleaser-full:
	goreleaser release --clean

# Install Goreleaser (if not already installed)
install-goreleaser:
	go install github.com/goreleaser/goreleaser@latest

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	gofumpt -w .

format-check:
	gofumpt -d .

format-fix: fmt imports

imports:
	goimports -w .

vet:
	go vet ./...

test-all: vet imports fmt lint test
