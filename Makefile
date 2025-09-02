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
