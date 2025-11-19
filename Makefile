# For local development

# Build the executable for current platform with whaletui as the name
build:
	go build -o whaletui .

build-win:
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

mockery:
	mockery --config=.mockery.yml

# Docker testing environment setup
docker-test:
	@echo "Setting up Docker testing environment..."
	@echo "Initializing swarm..."
	docker swarm init --advertise-addr 127.0.0.1
	@echo "Creating swarm networks..."
	docker network create --driver overlay --attachable test-network-1
	docker network create --driver overlay --attachable test-network-2
	@echo "Creating volumes..."
	docker volume create test-volume-1
	docker volume create test-volume-2
	docker volume create test-volume-3
	@echo "Creating containers..."
	docker run -d --name test-container-1 --network test-network-1 -v test-volume-1:/data nginx:alpine
	docker run -d --name test-container-2 --network test-network-2 -v test-volume-2:/data redis:alpine
	docker run -d --name test-container-3 --network test-network-1 --network test-network-2 -v test-volume-3:/data postgres:13-alpine
	@echo "Pulling additional images..."
	docker pull busybox:latest
	docker pull alpine:latest
	@echo "Creating swarm services..."
	docker service create --name test-service-1 --replicas 2 --network test-network-1 nginx:alpine
	docker service create --name test-service-2 --replicas 1 --network test-network-2 redis:alpine
	@echo "Docker testing environment setup complete!"

# Docker testing environment cleanup
docker-cleanup:
	@echo "Cleaning up Docker testing environment..."
	@echo "Removing swarm services..."
	-docker service rm test-service-1 test-service-2
	@echo "Leaving swarm..."
	-docker swarm leave --force
	@echo "Removing containers..."
	-docker rm -f test-container-1 test-container-2 test-container-3
	@echo "Removing volumes..."
	-docker volume rm test-volume-1 test-volume-2 test-volume-3
	@echo "Removing networks..."
	-docker network rm test-network-1 test-network-2
	@echo "Docker testing environment cleanup complete!"

# E2E testing
e2e-setup: docker-test
	@echo "E2E test environment ready"

e2e-test: e2e-setup
	@echo "Running E2E tests..."
	go test -v -timeout 30m ./e2e/...

e2e-test-verbose: e2e-setup
	@echo "Running E2E tests with verbose output..."
	go test -v -timeout 30m -count=1 ./e2e/... -args -test.v

e2e-test-short: e2e-setup
	@echo "Running E2E tests (short mode)..."
	go test -v -timeout 10m -short ./e2e/...

e2e-cleanup: docker-cleanup
	@echo "E2E test environment cleaned up"

e2e-full: e2e-test e2e-cleanup
	@echo "E2E tests complete with cleanup"

# Run specific e2e test
e2e-test-container:
	@echo "Running container tests..."
	go test -v -timeout 10m ./e2e -run TestContainer

e2e-test-image:
	@echo "Running image tests..."
	go test -v -timeout 10m ./e2e -run TestImage

e2e-test-volume:
	@echo "Running volume tests..."
	go test -v -timeout 10m ./e2e -run TestVolume

e2e-test-network:
	@echo "Running network tests..."
	go test -v -timeout 10m ./e2e -run TestNetwork

e2e-test-swarm:
	@echo "Running swarm tests..."
	go test -v -timeout 10m ./e2e -run TestSwarm

e2e-test-workflow:
	@echo "Running workflow tests..."
	go test -v -timeout 10m ./e2e -run TestWorkflow

# TUI-specific tests
e2e-test-tui:
	@echo "Running TUI interaction tests..."
	go test -v -timeout 10m ./e2e -run TestTUI

e2e-test-tui-interaction:
	@echo "Running TUI interaction tests..."
	go test -v -timeout 10m ./e2e -run TestTUIInteraction

e2e-test-tui-workflow:
	@echo "Running TUI workflow tests..."
	go test -v -timeout 10m ./e2e -run TestTUIWorkflow

# E2E test coverage
e2e-coverage:
	@echo "Running E2E tests with coverage..."
	go test -v -timeout 30m -coverprofile=e2e_coverage.out ./e2e/...
	go tool cover -html=e2e_coverage.out -o e2e_coverage.html
	@echo "Coverage report generated: e2e_coverage.html"

