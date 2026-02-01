# E2E Tests Makefile
# This Makefile provides easy commands for running e2e tests with Docker volumes

.PHONY: help build test test-complete test-parallel test-unit test-unit-parallel test-integration test-ui test-performance test-all clean shell generate-data logs

# Default number of parallel workers (override with `make ... WORKERS=4`)
WORKERS ?= auto

# Default target
help:
	@echo "E2E Tests Makefile"
	@echo "=================="
	@echo ""
	@echo "Available commands:"
	@echo "  test-complete  - Complete e2e test run (build once, generate data once, run all tests)"
	@echo "  test-parallel  - Run all tests in parallel with multiple workers (faster)"
	@echo "  build          - Build the Docker image for e2e tests"
	@echo "  test           - Run all e2e tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-unit-parallel - Run unit tests in parallel (faster)"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-ui        - Run UI tests only"
	@echo "  test-performance - Run performance tests only"
	@echo "  test-all       - Run all test categories"
	@echo "  generate-data  - Generate test data in the DinD environment"
	@echo "  shell          - Open a shell in the test container"
	@echo "  logs           - Show logs from the test container"
	@echo "  clean          - Clean up Docker containers and volumes"
	@echo ""
	@echo "Parallel execution options:"
	@echo "  Use 'make test-parallel WORKERS=4' to specify number of workers (default: auto)"
	@echo "  Use 'make test-unit-parallel WORKERS=2' to run unit tests in parallel"
	@echo ""

# Build the Docker image
build:
	@echo "🔨 Building Docker image for e2e tests..."
	cd .. && docker build -f e2e/docker/Dockerfile.test -t whaletui-e2e-test-isolated . --no-cache
	@echo "✅ Docker image built successfully!"

# Run all e2e tests
test: build
	@echo "🧪 Running all e2e tests..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/whaletui/e2e/tests/ -v"
	@echo "✅ All tests completed!"

# Complete e2e test run (build once, generate data once, run all tests)
test-complete:
	@echo "🚀 Starting complete e2e test run..."
	@echo "🔨 Building Docker image..."
	cd .. && docker build -f e2e/docker/Dockerfile.test -t whaletui-e2e-test-isolated . --no-cache
	@echo "📊 Generating test data..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test python3 /app/docker/test_data_generator.py
	@echo "🧪 Running all tests..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/whaletui/e2e/tests/ -v"
	@echo "✅ Complete e2e test run finished!"

# Run all tests in parallel with multiple workers (faster)
test-parallel: build
	@echo "🚀 Running all tests in parallel..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/whaletui/e2e/tests/ -v"
	@echo "✅ Parallel tests completed!"

# Run unit tests only
test-unit: build
	@echo "🧪 Running unit tests..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/whaletui/e2e/tests/unit/ -v"
	@echo "✅ Unit tests completed!"

# Run unit tests in parallel
test-unit-parallel: build
	@echo "🧪 Running unit tests in parallel..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/whaletui/e2e/tests/unit/ -v"
	@echo "✅ Unit tests completed!"

# Run integration tests only
test-integration: build
	@echo "🧪 Running integration tests..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/whaletui/e2e/tests/integration/ -v"
	@echo "✅ Integration tests completed!"

# Run UI tests only
test-ui: build
	@echo "🧪 Running UI tests..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/whaletui/e2e/tests/ui/ -v"
	@echo "✅ UI tests completed!"

# Run all test categories
test-all: test-unit test-integration test-ui
	@echo "🎉 All test categories completed!"

# Generate test data
generate-data: build
	@echo "📊 Generating test data..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test python3 /app/docker/test_data_generator.py
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
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/whaletui/e2e/tests/unit/ -v"

# Test specific file
test-file:
	@echo "🧪 Running specific test file..."
	@read -p "Enter test file path (e.g., tests/unit/test_basic.py): " file; \
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/e2e/$$file -v"

# Test with coverage
test-coverage: build
	@echo "🧪 Running tests with coverage..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/whaletui/e2e/tests/ --cov=whaletui_controller --cov-report=html"
	@echo "✅ Coverage report generated!"

# Debug mode (with verbose output)
debug-test: build
	@echo "🐛 Running tests in debug mode..."
	docker-compose -f docker/docker-compose.volume.yml run --rm whaletui-test bash -c "cd /app/test-data && PYTHONPATH=/app/whaletui/e2e python3 -m pytest -n $(WORKERS) /app/whaletui/e2e/tests/ -v -s --tb=long"
	@echo "✅ Debug test completed!"
