"""
Pytest configuration and fixtures for WhaleTUI e2e tests.
"""
import os
import sys
import pytest
import subprocess
import time
from pathlib import Path

# Add the project root to the Python path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root))

from whaletui_controller import WhaleTUIController


@pytest.fixture(scope="session")
def whaletui_binary():
    """Build and return the path to the WhaleTUI binary."""
    binary_path = project_root / "whaletui.exe"

    if not binary_path.exists():
        # Build the binary if it doesn't exist
        print("Building WhaleTUI binary...")
        result = subprocess.run(
            ["go", "build", "-o", "whaletui.exe", "."],
            cwd=project_root,
            capture_output=True,
            text=True
        )
        if result.returncode != 0:
            pytest.fail(f"Failed to build WhaleTUI: {result.stderr}")

    return str(binary_path)


@pytest.fixture
def whaletui_app(whaletui_binary):
    """Create a WhaleTUI application controller for testing."""
    controller = WhaleTUIController(whaletui_binary)
    # Set the theme file path for the test environment
    controller.theme_path = "/app/e2e/config/theme.yaml"
    print(f"DEBUG: Set theme_path to: {controller.theme_path}")
    yield controller
    controller.cleanup()


@pytest.fixture(scope="session")
def docker_test_environment():
    """Set up Docker test environment if needed."""
    # Check if Docker is running
    try:
        result = subprocess.run(
            ["docker", "version", "--format", "{{.Server.Version}}"],
            capture_output=True,
            text=True,
            timeout=5
        )
        if result.returncode != 0:
            pytest.skip("Docker is not running or not accessible")
    except (subprocess.TimeoutExpired, FileNotFoundError):
        pytest.skip("Docker is not installed or not accessible")


@pytest.fixture(autouse=True)
def test_timeout():
    """Set default timeout for all tests."""
    return 30


def pytest_configure(config):
    """Configure pytest with custom markers."""
    config.addinivalue_line(
        "markers", "slow: marks tests as slow (deselect with '-m \"not slow\"')"
    )
    config.addinivalue_line(
        "markers", "integration: marks tests as integration tests"
    )
    config.addinivalue_line(
        "markers", "docker: marks tests that require Docker"
    )
    config.addinivalue_line(
        "markers", "ui: marks tests as UI tests"
    )
    config.addinivalue_line(
        "markers", "performance: marks tests as performance tests"
    )
    config.addinivalue_line(
        "markers", "error_handling: marks tests as error handling tests"
    )
    config.addinivalue_line(
        "markers", "search: marks tests as search tests"
    )
