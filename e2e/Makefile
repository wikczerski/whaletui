# E2E Tests Makefile
# This Makefile provides easy commands for running e2e tests with Docker volumes

.PHONY: help build test test-unit test-integration test-ui test-performance test-all clean shell generate-data logs

# Default target
help:
	@echo "E2E Tests Makefile"
	@echo "=================="
	@echo ""
	@echo "Available commands:"
	@echo "  build          - Build the Docker image for e2e tests"
	@echo "  test           - Run all e2e tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-ui        - Run UI tests only"
	@echo "  test-performance - Run performance tests only"
	@echo "  test-all       - Run all test categories"
	@echo "  generate-data  - Generate test data in the DinD environment"
	@echo "  shell          - Open a shell in the test container"
	@echo "  logs           - Show logs from the test container"
	@echo "  clean          - Clean up Docker containers and volumes"
	@echo ""

# Build the Docker image
build:
	@echo "🔨 Building Docker image for e2e tests..."
	cd .. && docker build -f e2e/docker/Dockerfile.test -t whaletui-e2e-test-isolated . --no-cache
	@echo "✅ Docker image built successfully!"

# Run all e2e tests
test: build
	@echo "🧪 Running all e2e tests..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && python3 -m pytest /app/e2e/tests/ -v"
	@echo "✅ All tests completed!"

# Run unit tests only
test-unit: build
	@echo "🧪 Running unit tests..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && python3 -m pytest /app/e2e/tests/unit/ -v"
	@echo "✅ Unit tests completed!"

# Run integration tests only
test-integration: build
	@echo "🧪 Running integration tests..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && python3 -m pytest /app/e2e/tests/integration/ -v"
	@echo "✅ Integration tests completed!"

# Run UI tests only
test-ui: build
	@echo "🧪 Running UI tests..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && python3 -m pytest /app/e2e/tests/ui/ -v"
	@echo "✅ UI tests completed!"

# Run performance tests only
test-performance: build
	@echo "🧪 Running performance tests..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && python3 -m pytest /app/e2e/tests/performance/ -v"
	@echo "✅ Performance tests completed!"

# Run all test categories
test-all: test-unit test-integration test-ui test-performance
	@echo "🎉 All test categories completed!"

# Generate test data
generate-data: build
	@echo "📊 Generating test data..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test python3 docker/test_data_generator.py
	@echo "✅ Test data generated!"

# Open a shell in the test container
shell: build
	@echo "🐚 Opening shell in test container..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test /bin/bash

# Show logs from the test container
logs:
	@echo "📋 Showing logs from test container..."
	docker-compose -f docker/docker-compose.volume.yml logs whaletui-test

# Clean up Docker containers and volumes
clean:
	@echo "🧹 Cleaning up Docker containers and volumes..."
	docker-compose -f docker/docker-compose.volume.yml down -v
	docker system prune -f
	@echo "✅ Cleanup completed!"

# Quick test run (without rebuilding)
quick-test:
	@echo "⚡ Running quick test (no rebuild)..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && python3 -m pytest /app/e2e/tests/unit/ -v"

# Test specific file
test-file:
	@echo "🧪 Running specific test file..."
	@read -p "Enter test file path (e.g., tests/unit/test_basic.py): " file; \
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && python3 -m pytest /app/e2e/$$file -v"

# Test with coverage
test-coverage: build
	@echo "🧪 Running tests with coverage..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && python3 -m pytest /app/e2e/tests/ --cov=whaletui_controller --cov-report=html"
	@echo "✅ Coverage report generated!"

# Debug mode (with verbose output)
debug-test: build
	@echo "🐛 Running tests in debug mode..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && python3 -m pytest /app/e2e/tests/ -v -s --tb=long"
	@echo "✅ Debug test completed!"
