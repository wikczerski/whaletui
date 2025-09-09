"""
Pytest configuration for comprehensive WhaleTUI e2e tests.
"""
import os
import sys
import pytest
import subprocess
import time
from pathlib import Path

# Add the project root to the Python path
project_root = Path(__file__).parent.parent.parent
sys.path.insert(0, str(project_root))

# Debug: Print paths for troubleshooting
print(f"DEBUG: project_root = {project_root}")
print(f"DEBUG: Docker binary exists = {Path('/app/whaletui/whaletui').exists()}")
print(f"DEBUG: Current working directory = {Path.cwd()}")

from e2e.whaletui_controller import WhaleTUIController
from tests.utils.test_helpers import TestHelpers


@pytest.fixture(scope="session")
def whaletui_binary():
    """Build and return the path to the WhaleTUI binary."""
    # Check for binary in Docker container location first
    docker_binary_path = Path("/app/whaletui/whaletui")
    if docker_binary_path.exists():
        print(f"DEBUG: Using Docker binary at {docker_binary_path}")
        return str(docker_binary_path)

    # Fallback to local build
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

    print(f"DEBUG: Using local binary at {binary_path}")
    return str(binary_path)


@pytest.fixture
def whaletui_app(whaletui_binary):
    """Create a WhaleTUI application controller for testing."""
    controller = WhaleTUIController(whaletui_binary)
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


def pytest_sessionstart(session):
    """Set up test session."""
    TestHelpers.setup_test_environment()


def pytest_sessionfinish(session, exitstatus):
    """Clean up test session."""
    TestHelpers.cleanup_test_environment()
