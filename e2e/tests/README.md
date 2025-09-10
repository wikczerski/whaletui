# WhaleTUI Comprehensive E2E Test Suite

This directory contains a comprehensive end-to-end test suite for the WhaleTUI application, organized into logical categories for better maintainability and execution control.

## Directory Structure

```
tests/
├── __init__.py                 # Test package initialization
├── conftest.py                 # Pytest configuration and fixtures
├── unit/                       # Unit tests
│   ├── __init__.py
│   └── test_basic.py          # Basic functionality tests
├── integration/                # Integration tests
│   ├── __init__.py
│   └── test_docker_integration.py  # Docker integration tests
├── ui/                         # UI tests
│   ├── __init__.py
│   ├── test_ui_navigation.py  # UI navigation tests
│   └── test_ui_interactions.py # UI interaction tests
├── performance/                # Performance tests
│   ├── __init__.py
│   └── test_performance_metrics.py  # Performance metrics tests
├── docker/                     # Docker-specific tests
│   ├── __init__.py
│   └── test_docker_integration.py  # Docker integration tests
├── search/                     # Search functionality tests
│   ├── __init__.py
│   └── test_search_functionality.py  # Search tests
├── error_handling/             # Error handling tests
│   ├── __init__.py
│   └── test_error_scenarios.py  # Error scenario tests
├── fixtures/                   # Test fixtures
│   ├── __init__.py
│   └── docker_fixtures.py     # Docker test fixtures
└── utils/                      # Test utilities
    ├── __init__.py
    └── test_helpers.py         # Helper utilities
```

## Test Categories

### Unit Tests (`unit/`)
- **Purpose**: Test individual components and basic functionality
- **Scope**: Application startup, shutdown, basic commands
- **Execution Time**: Fast (< 1 minute)
- **Dependencies**: Minimal

### Integration Tests (`integration/`)
- **Purpose**: Test integration with external systems (Docker)
- **Scope**: Docker operations, container management, network operations
- **Execution Time**: Medium (1-5 minutes)
- **Dependencies**: Docker daemon

### UI Tests (`ui/`)
- **Purpose**: Test user interface interactions and navigation
- **Scope**: View switching, keyboard navigation, UI responsiveness
- **Execution Time**: Medium (2-5 minutes)
- **Dependencies**: Minimal

### Performance Tests (`performance/`)
- **Purpose**: Test application performance and resource usage
- **Scope**: Startup time, memory usage, CPU usage, response times
- **Execution Time**: Slow (5-15 minutes)
- **Dependencies**: Minimal

### Docker Tests (`docker/`)
- **Purpose**: Test Docker-specific functionality
- **Scope**: Container operations, image management, volume operations
- **Execution Time**: Medium (3-10 minutes)
- **Dependencies**: Docker daemon

### Search Tests (`search/`)
- **Purpose**: Test search functionality across all views
- **Scope**: Search operations, filtering, search performance
- **Execution Time**: Medium (2-5 minutes)
- **Dependencies**: Minimal

### Error Handling Tests (`error_handling/`)
- **Purpose**: Test error scenarios and recovery
- **Scope**: Invalid inputs, network errors, permission errors
- **Execution Time**: Medium (3-8 minutes)
- **Dependencies**: Minimal

## Test Markers

Tests are categorized using pytest markers:

- `@pytest.mark.slow`: Tests that take longer to execute
- `@pytest.mark.integration`: Integration tests
- `@pytest.mark.docker`: Tests that require Docker
- `@pytest.mark.ui`: UI-related tests
- `@pytest.mark.performance`: Performance tests
- `@pytest.mark.error_handling`: Error handling tests
- `@pytest.mark.search`: Search functionality tests

## Running Tests

### Using Makefile (Recommended)

```bash
# Install dependencies
make e2e-install

# Run all tests
make e2e-test

# Run specific test categories
make e2e-test-unit
make e2e-test-integration
make e2e-test-ui
make e2e-test-performance
make e2e-test-docker
make e2e-test-search
make e2e-test-error-handling

# Run tests with specific options
make e2e-test-verbose
make e2e-test-parallel
make e2e-test-no-slow
make e2e-test-docker-only
make e2e-test-ui-only
make e2e-test-performance-only

# Run quick tests (unit + UI)
make e2e-test-quick

# Run full test suite
make e2e-test-full

# Clean test artifacts
make e2e-clean
```

### Using Python Scripts

```bash
# Run comprehensive tests
cd e2e
python run_comprehensive_tests.py

# Run specific categories
python run_comprehensive_tests.py --category unit
python run_comprehensive_tests.py --category ui
python run_comprehensive_tests.py --category performance

# Run with specific markers
python run_comprehensive_tests.py --marker "docker"
python run_comprehensive_tests.py --marker "slow"
python run_comprehensive_tests.py --exclude-slow

# Run specific tests
python run_comprehensive_tests.py --test tests/unit/test_basic.py
python run_comprehensive_tests.py --test tests/ui/test_ui_navigation.py::TestUINavigation::test_main_screen_display

# Run with options
python run_comprehensive_tests.py --verbose --parallel
python run_comprehensive_tests.py --docker-only --verbose
```

### Using pytest directly

```bash
# Run all tests
cd e2e
pytest tests/

# Run specific categories
pytest tests/unit/
pytest tests/ui/
pytest tests/performance/

# Run with markers
pytest tests/ -m "docker"
pytest tests/ -m "not slow"
pytest tests/ -m "ui and not slow"

# Run specific tests
pytest tests/unit/test_basic.py
pytest tests/ui/test_ui_navigation.py::TestUINavigation::test_main_screen_display

# Run with options
pytest tests/ -v --parallel
pytest tests/ -v --html=reports/report.html
```

## Test Configuration

### Environment Variables

- `WHALETUI_LOG_LEVEL`: Set log level for tests (default: INFO)
- `WHALETUI_TIMEOUT`: Set default timeout for tests (default: 30)
- `WHALETUI_SCREENSHOTS`: Enable/disable screenshots (default: true)

### Test Data

Test data is stored in the `test_data/` directory and is automatically cleaned up after tests.

### Screenshots

Screenshots are automatically captured during tests and stored in the `screenshots/` directory. They are useful for debugging and visual verification.

### Reports

Test reports are generated in the `reports/` directory in HTML and JUnit XML formats.

## Writing New Tests

### Test Structure

```python
import pytest
from e2e.whaletui_controller import WhaleTUIController
from tests.utils.test_helpers import TestHelpers

class TestNewFeature:
    """Tests for new feature."""

    def test_feature_basic(self, whaletui_app):
        """Test basic feature functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test feature
        # ... test code ...

        # Take screenshot
        whaletui_app.take_screenshot("feature_basic.png")

    @pytest.mark.slow
    def test_feature_performance(self, whaletui_app):
        """Test feature performance."""
        # ... test code ...
```

### Test Guidelines

1. **Use descriptive test names**: Test names should clearly describe what is being tested
2. **Add docstrings**: Each test should have a docstring explaining its purpose
3. **Use appropriate markers**: Mark tests with appropriate pytest markers
4. **Take screenshots**: Use `whaletui_app.take_screenshot()` for visual verification
5. **Clean up**: Tests should clean up after themselves
6. **Handle errors gracefully**: Tests should handle expected errors appropriately
7. **Use timeouts**: Set appropriate timeouts for operations
8. **Test edge cases**: Include tests for edge cases and error conditions

### Test Categories

When writing new tests, place them in the appropriate category:

- **Unit tests**: Basic functionality, commands, configuration
- **Integration tests**: External system integration (Docker, SSH)
- **UI tests**: User interface, navigation, interactions
- **Performance tests**: Performance metrics, resource usage
- **Docker tests**: Docker-specific functionality
- **Search tests**: Search and filtering functionality
- **Error handling tests**: Error scenarios and recovery

## Debugging Tests

### Common Issues

1. **Test timeouts**: Increase timeout values or check for blocking operations
2. **Screenshot failures**: Ensure screenshot directory exists and is writable
3. **Docker tests failing**: Check if Docker daemon is running
4. **UI tests failing**: Check if application is responding to input

### Debugging Tools

1. **Screenshots**: Check screenshots in the `screenshots/` directory
2. **Logs**: Check logs in the `logs/` directory
3. **Verbose output**: Use `--verbose` flag for detailed output
4. **Single test execution**: Run individual tests for debugging

### Test Isolation

Tests should be isolated and not depend on each other. Use fixtures for setup and teardown.

## Continuous Integration

The test suite is designed to work with CI/CD pipelines:

1. **Fast feedback**: Unit and UI tests run quickly
2. **Comprehensive coverage**: Full test suite covers all functionality
3. **Parallel execution**: Tests can run in parallel for faster execution
4. **Artifact collection**: Screenshots and reports are collected for analysis
5. **Docker support**: Docker tests can be run in CI environments

## Maintenance

### Regular Tasks

1. **Update test data**: Keep test data current with application changes
2. **Review screenshots**: Check screenshots for UI changes
3. **Update dependencies**: Keep Python dependencies updated
4. **Clean artifacts**: Regularly clean up test artifacts
5. **Review test coverage**: Ensure adequate test coverage

### Adding New Features

When adding new features to WhaleTUI:

1. **Add unit tests**: Test basic functionality
2. **Add UI tests**: Test user interface interactions
3. **Add integration tests**: Test external system integration
4. **Add performance tests**: Test performance impact
5. **Add error handling tests**: Test error scenarios
6. **Update documentation**: Update this README if needed

## Troubleshooting

### Common Problems

1. **Import errors**: Check Python path and dependencies
2. **Permission errors**: Check file permissions for screenshots and reports
3. **Docker errors**: Check Docker daemon status and permissions
4. **Timeout errors**: Check application responsiveness and timeout values
5. **Memory errors**: Check for memory leaks in long-running tests

### Getting Help

1. **Check logs**: Look at log files for error details
2. **Run with verbose output**: Use `--verbose` flag for detailed information
3. **Check screenshots**: Look at captured screenshots for visual issues
4. **Run individual tests**: Isolate problematic tests
5. **Check dependencies**: Ensure all dependencies are installed correctly
