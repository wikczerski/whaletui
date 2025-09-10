# WhaleTUI E2E Tests

This directory contains end-to-end tests for the WhaleTUI application using Python and pexpect.

## Overview

The e2e tests use the following approach:
- **pexpect**: Python library for controlling interactive terminal applications
- **pytest**: Python testing framework
- **Screenshots**: Automatic screenshot capture for debugging

## Test Structure

- `conftest.py`: Pytest configuration and fixtures
- `whaletui_controller.py`: Controller class for interacting with WhaleTUI
- `test_basic.py`: Basic functionality tests
- `test_docker_integration.py`: Docker integration tests
- `test_search.py`: Search functionality tests
- `test_error_handling.py`: Error handling tests
- `test_performance.py`: Performance tests
- `run_tests.py`: Test runner script

## Prerequisites

1. **Python 3.7+** installed
2. **Docker** running (for integration tests)
3. **WhaleTUI binary** built (will be built automatically if not present)

## Installation

1. Install Python dependencies:
   ```bash
   cd e2e
   python -m pip install -r requirements.txt
   ```

2. Or use the test runner:
   ```bash
   python run_tests.py --install-deps
   ```

## Running Tests

### Basic Usage

```bash
# Run all tests
python run_tests.py

# Run specific test file
python run_tests.py --test test_basic.py

# Run specific test
python run_tests.py --test test_basic.py::TestWhaleTUIBasic::test_application_starts

# Run with verbose output
python run_tests.py --verbose

# Run tests in parallel
python run_tests.py --parallel
```

### Using pytest directly

```bash
# Run all tests
pytest

# Run specific test file
pytest test_basic.py

# Run tests with specific markers
pytest -m "not slow"
pytest -m "docker"

# Run with verbose output
pytest -v

# Run tests in parallel
pytest -n auto
```

### Using Makefile

The project Makefile includes e2e test commands:

```bash
# Run all e2e tests
make e2e-test

# Run basic e2e tests
make e2e-test-basic

# Run search e2e tests
make e2e-test-search

# Run performance e2e tests
make e2e-test-performance

# Run error handling e2e tests
make e2e-test-error-handling

# Run with verbose output
make e2e-test-verbose

# Run specific test
make e2e-test-specific TEST=TestWhaleTUIBasic::test_application_starts
```

## Test Categories

### Basic Tests (`test_basic.py`)
- Application startup
- Application shutdown
- Help screen
- Theme command
- Connect command
- Invalid commands
- Command line flags

### Docker Integration Tests (`test_docker_integration.py`)
- Containers view
- Images view
- Volumes view
- Networks view
- Swarm view
- Nodes view
- Services view
- View navigation
- Refresh functionality

### Search Tests (`test_search.py`)
- Search in different views
- Search clear functionality
- Empty search terms
- Special characters in search
- Search performance

### Error Handling Tests (`test_error_handling.py`)
- Invalid Docker host
- Invalid command line arguments
- Application interrupts
- Memory usage
- CPU usage

### Performance Tests (`test_performance.py`)
- Startup time
- View switching performance
- Search performance
- Refresh performance
- Memory usage over time
- CPU usage under load
- Large dataset performance
- Concurrent operations

## Test Markers

- `@pytest.mark.slow`: Slow tests (deselect with `-m "not slow"`)
- `@pytest.mark.integration`: Integration tests
- `@pytest.mark.docker`: Tests that require Docker

## Screenshots

Screenshots are automatically captured during tests and saved to the `screenshots/` directory. This helps with debugging and visual verification of the application state.

## Reports

Test reports are generated in the `reports/` directory:
- `report.html`: HTML test report
- `junit.xml`: JUnit XML report for CI/CD integration

## Configuration

### Environment Variables

- `LINES`: Terminal height (default: 24)
- `COLUMNS`: Terminal width (default: 80)

### Test Timeouts

- Default timeout: 30 seconds
- Can be overridden per test
- Global timeout: 60 seconds

## Debugging

### Enable Debug Logging

```bash
# Set log level to DEBUG
export WHALETUI_LOG_LEVEL=DEBUG

# Run tests with debug output
python run_tests.py --verbose
```

### View Screenshots

Screenshots are saved in the `screenshots/` directory and can be viewed to understand the application state during tests.

### Check Logs

- `whaletui_output.log`: WhaleTUI application output
- `reports/report.html`: Detailed test report

## CI/CD Integration

The tests are designed to work in CI/CD environments:

```yaml
# Example GitHub Actions workflow
- name: Run E2E Tests
  run: |
    cd e2e
    python run_tests.py --install-deps
    python run_tests.py --parallel
```

## Troubleshooting

### Common Issues

1. **Docker not running**: Some tests require Docker to be running
2. **Binary not found**: The WhaleTUI binary will be built automatically
3. **Permission issues**: Ensure proper permissions for the test directory
4. **Timeout issues**: Increase timeout values for slow environments

### Getting Help

- Check the test logs in `reports/`
- View screenshots in `screenshots/`
- Run tests with `--verbose` flag for detailed output
- Check the WhaleTUI application logs
