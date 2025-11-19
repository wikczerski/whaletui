# WhaleTUI E2E Testing Framework

## Overview

This directory contains end-to-end tests for the WhaleTUI application. The tests simulate real user interactions with the TUI and verify Docker operations across all domains.

## Architecture

```
e2e/
├── framework/           # Core test framework
│   ├── framework.go     # Test harness and lifecycle (Docker integration tests)
│   ├── docker_helper.go # Docker operation helpers
│   ├── tui_helper.go    # TUI interaction helpers (for integration tests)
│   ├── tui_test_framework.go # TUI testing with tcell simulation (for UI tests)
│   └── fixtures.go      # Test fixtures and data
├── testdata/            # Test data files
│   ├── themes/          # Test theme files
│   └── configs/         # Test configuration files
├── *_test.go            # Integration tests (Docker operations)
├── tui_*_test.go        # TUI interaction tests (UI testing)
└── README.md            # This file
```

## Test Types

### 1. Integration Tests (Docker Backend)
Tests in `container_test.go`, `image_test.go`, `volume_test.go`, etc. focus on:
- Docker API operations
- Business logic
- Data workflows
- **Do NOT** launch the TUI application
- Fast and reliable for backend testing

### 2. TUI Interaction Tests (UI Frontend)
Tests in `tui_interaction_test.go`, `tui_workflow_test.go` focus on:
- Actual UI rendering with tcell simulation
- Keyboard navigation and input
- Screen content verification
- User workflows through the interface
- **DO** launch tview components with simulated screen
- Test what users actually see and interact with

## Prerequisites

- Docker daemon running locally
- Go 1.25.0+
- Docker Swarm initialized (for swarm tests)

## Quick Start

### Setup Test Environment

```bash
# Create Docker test fixtures (containers, images, volumes, networks, swarm)
make e2e-setup
```

This creates:
- 3 test containers (nginx, redis, postgres)
- 2 test images (busybox, alpine)
- 3 test volumes
- 2 test networks
- 2 swarm services

### Run All Tests

```bash
# Run all e2e tests
make e2e-test

# Run with verbose output
make e2e-test-verbose

# Run and cleanup
make e2e-full
```

### Run Specific Tests

```bash
# Container tests only
go test -v ./e2e -run TestContainer

# Swarm tests only
go test -v ./e2e -run TestSwarm

# Navigation tests only
go test -v ./e2e -run TestNavigation

# TUI interaction tests only
go test -v ./e2e -run TestTUI

# Single test
go test -v ./e2e -run TestContainerList
```

### Cleanup

```bash
# Remove all test fixtures
make e2e-cleanup
```

## Test Framework

### Core Components

#### Framework (`framework/framework.go`)

The main test harness providing:
- `NewTestFramework()` - Initialize test environment
- `SetupDocker()` - Configure Docker client
- `TearDown()` - Cleanup resources
- `SimulateKey()` - Simulate keyboard input
- `WaitForCondition()` - Wait for async operations

#### Docker Helper (`framework/docker_helper.go`)

Docker operation utilities:
- `CreateTestContainer()` - Create container for testing
- `WaitForContainerState()` - Wait for container state change
- `CleanupContainer()` - Remove test container
- Similar helpers for images, volumes, networks, swarm

#### TUI Helper (`framework/tui_helper.go`)

TUI interaction utilities:
- `NavigateToView()` - Switch to specific view
- `SelectTableRow()` - Select row in table
- `PressKey()` - Simulate key press
- `VerifyTableContent()` - Assert table contents
- `VerifyModalVisible()` - Assert modal state

#### Fixtures (`framework/fixtures.go`)

Predefined test data:
- Container configurations
- Image specifications
- Volume definitions
- Network configurations
- Swarm service specs

## Writing Tests

### Test Structure

```go
func TestFeatureName(t *testing.T) {
    // Setup
    fw := framework.NewTestFramework(t)
    defer fw.TearDown()

    // Create test fixtures
    containerID := fw.CreateTestContainer("test-container", "nginx:alpine")

    // Execute test
    fw.NavigateToView("containers")
    fw.SelectTableRow(containerID)
    fw.PressKey('s') // Start container

    // Verify
    fw.WaitForContainerState(containerID, "running")
    fw.VerifyTableRowColor(containerID, "green")
}
```

### Best Practices

1. **Independence**: Each test should be independent and not rely on other tests
2. **Cleanup**: Always use `defer fw.TearDown()` to cleanup resources
3. **Descriptive Names**: Use clear, descriptive test names (e.g., `TestContainerStartStoppedContainer`)
4. **Wait Conditions**: Use `WaitForCondition()` for async operations
5. **Clear Assertions**: Use testify assertions with clear error messages
6. **Test Data**: Use fixtures from `framework/fixtures.go` for consistency

### Test Categories

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions
- **E2E Tests**: Test complete user workflows
- **Error Tests**: Test error handling and edge cases

## Test Coverage

The e2e tests cover:

- ✅ Container management (list, start, stop, restart, delete, inspect, logs, shell, exec)
- ✅ Image management (list, delete, inspect)
- ✅ Volume management (list, delete, inspect)
- ✅ Network management (list, delete, inspect)
- ✅ Swarm services (list, scale, remove, inspect, logs)
- ✅ Swarm nodes (list, drain, activate, remove, inspect)
- ✅ Navigation (view switching, keyboard shortcuts)
- ✅ Search and filtering
- ✅ Command mode
- ✅ Error handling and edge cases
- ✅ Theme and configuration

## Troubleshooting

### Docker Daemon Not Running

```
Error: Cannot connect to the Docker daemon
```

**Solution**: Start Docker daemon
```bash
sudo systemctl start docker  # Linux
# or start Docker Desktop on macOS/Windows
```

### Swarm Not Initialized

```
Error: This node is not a swarm manager
```

**Solution**: Initialize swarm
```bash
docker swarm init --advertise-addr 127.0.0.1
```

### Port Conflicts

```
Error: Bind for 0.0.0.0:8080 failed: port is already allocated
```

**Solution**: Stop conflicting containers or use different ports in test fixtures

### Test Timeout

```
Error: test timed out after 30m
```

**Solution**: Increase timeout or optimize slow tests
```bash
go test -v -timeout 60m ./e2e/...
```

### Cleanup Issues

If tests leave resources behind:
```bash
# Manual cleanup
make e2e-cleanup

# Nuclear option (removes ALL Docker resources)
docker system prune -a --volumes -f
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25.0'
      - name: Setup Docker
        run: |
          docker swarm init --advertise-addr 127.0.0.1
      - name: Run E2E Tests
        run: |
          make e2e-setup
          make e2e-test
      - name: Cleanup
        if: always()
        run: make e2e-cleanup
```

## Performance

- Average test execution time: ~2-5 seconds per test
- Full suite execution: ~5-10 minutes
- Parallel execution: Supported with `-p` flag

## Contributing

When adding new tests:

1. Follow existing test patterns
2. Add test to appropriate file (or create new file for new domain)
3. Update this README if adding new test categories
4. Ensure tests pass locally before submitting PR
5. Add test documentation in code comments

## Support

For issues or questions:
- Check existing tests for examples
- Review test framework code in `framework/`
- Consult main test plan: `e2e_test_plan.md`
