"""
Test helper utilities for WhaleTUI e2e tests.
"""
import os
import time
import logging
from typing import List, Dict, Any, Optional
from pathlib import Path


class TestHelpers:
    """Helper utilities for e2e tests."""

    @staticmethod
    def wait_for_condition(condition_func, timeout: int = 10, interval: float = 0.5) -> bool:
        """
        Wait for a condition to be true.

        Steps:
        1. Record start time
        2. Enter polling loop:
           - Call condition function
           - If condition is met, return True
           - If timeout exceeded, return False
           - Sleep for specified interval
        3. Continue until condition is met or timeout

        Expected Outcome:
        - Returns True if condition is met within timeout
        - Returns False if timeout is exceeded
        - Polls condition function at specified intervals
        - Handles timeout gracefully

        Args:
            condition_func: Function that returns True when condition is met
            timeout: Maximum time to wait in seconds
            interval: Time between checks in seconds

        Returns:
            True if condition was met, False if timeout
        """
        start_time = time.time()
        while time.time() - start_time < timeout:
            if condition_func():
                return True
            time.sleep(interval)
        return False

    @staticmethod
    def create_test_data_file(filename: str, content: str) -> str:
        """
        Create a test data file.

        Steps:
        1. Create test data directory (/app/test-data/test_data) if it doesn't exist
        2. Construct full file path by joining directory and filename
        3. Open file for writing with UTF-8 encoding
        4. Write content to file
        5. Close file
        6. Return full path to created file

        Expected Outcome:
        - Test data directory is created if needed
        - File is created with specified content
        - File is written with UTF-8 encoding
        - Full path to created file is returned
        - No errors occur during file creation

        Args:
            filename: Name of the file to create
            content: Content to write to the file

        Returns:
            Path to the created file
        """
        test_data_dir = Path("/app/test-data/test_data")
        test_data_dir.mkdir(exist_ok=True)

        file_path = test_data_dir / filename
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)

        return str(file_path)

    @staticmethod
    def cleanup_test_data():
        """
        Clean up test data files.

        Steps:
        1. Check if test data directory exists (/app/test-data/test_data)
        2. If directory exists, remove it and all contents recursively
        3. If directory doesn't exist, do nothing

        Expected Outcome:
        - Test data directory is removed if it exists
        - All test data files are cleaned up
        - No errors occur if directory doesn't exist
        - Cleanup is performed safely
        """
        test_data_dir = Path("/app/test-data/test_data")
        if test_data_dir.exists():
            import shutil
            shutil.rmtree(test_data_dir)

    @staticmethod
    def get_docker_test_containers() -> List[Dict[str, Any]]:
        """
        Get list of test containers for Docker tests.

        Steps:
        1. Define test container configurations
        2. Return list of container dictionaries with:
           - name: Container name
           - image: Docker image to use
           - ports: Port mappings
           - environment: Environment variables

        Expected Outcome:
        - Returns list of test container configurations
        - Each container has required fields (name, image, ports, environment)
        - Containers are suitable for Docker integration tests
        - Configurations include common test images (nginx, redis, postgres)

        Returns:
            List of container information dictionaries
        """
        return [
            {
                "name": "test-nginx",
                "image": "nginx:alpine",
                "ports": ["80:80"],
                "environment": {}
            },
            {
                "name": "test-redis",
                "image": "redis:alpine",
                "ports": ["6379:6379"],
                "environment": {}
            },
            {
                "name": "test-postgres",
                "image": "postgres:13-alpine",
                "ports": ["5432:5432"],
                "environment": {
                    "POSTGRES_DB": "testdb",
                    "POSTGRES_USER": "testuser",
                    "POSTGRES_PASSWORD": "testpass"
                }
            }
        ]

    @staticmethod
    def get_docker_test_networks() -> List[str]:
        """
        Get list of test networks for Docker tests.

        Steps:
        1. Define test network names
        2. Return list of network names for testing

        Expected Outcome:
        - Returns list of test network names
        - Network names are suitable for Docker integration tests
        - Names follow consistent naming pattern
        - Networks can be used for testing network functionality

        Returns:
            List of network names
        """
        return [
            "test-network-1",
            "test-network-2",
            "test-network-3"
        ]

    @staticmethod
    def get_docker_test_volumes() -> List[str]:
        """
        Get list of test volumes for Docker tests.

        Steps:
        1. Define test volume names
        2. Return list of volume names for testing

        Expected Outcome:
        - Returns list of test volume names
        - Volume names are suitable for Docker integration tests
        - Names follow consistent naming pattern
        - Volumes can be used for testing volume functionality

        Returns:
            List of volume names
        """
        return [
            "test-volume-1",
            "test-volume-2",
            "test-volume-3"
        ]

    @staticmethod
    def get_expected_ui_elements() -> Dict[str, List[str]]:
        """
        Get expected UI elements for different views.

        Steps:
        1. Define expected UI elements for each view type
        2. Map view names to their expected column headers/elements
        3. Return dictionary with view-to-elements mapping

        Expected Outcome:
        - Returns dictionary mapping view names to expected elements
        - Each view has appropriate column headers/elements
        - Elements match actual UI structure
        - Can be used for UI validation in tests

        Returns:
            Dictionary mapping view names to expected elements
        """
        return {
            "containers": [
                "Id", "Name", "Image", "Status", "Ports", "Created"
            ],
            "images": [
                "Repository", "Tag", "Image ID", "Created", "Size"
            ],
            "volumes": [
                "Driver", "Name", "Created", "Size", "Mountpoint"
            ],
            "networks": [
                "Network ID", "Name", "Driver", "Scope", "Created"
            ],
            "swarm": [
                "Swarm", "Node", "Manager", "Worker"
            ],
            "nodes": [
                "ID", "Hostname", "Status", "Availability", "Manager Status"
            ],
            "services": [
                "ID", "Name", "Mode", "Replicas", "Image", "Ports"
            ]
        }

    @staticmethod
    def get_test_key_sequences() -> Dict[str, List[str]]:
        """
        Get test key sequences for different operations.

        Steps:
        1. Define key sequences for common operations
        2. Map operation names to their corresponding key sequences
        3. Return dictionary with operation-to-keys mapping

        Expected Outcome:
        - Returns dictionary mapping operation names to key sequences
        - Key sequences match actual application shortcuts
        - Covers all major operations (quit, help, refresh, etc.)
        - Can be used for automated UI testing

        Returns:
            Dictionary mapping operation names to key sequences
        """
        return {
            "quit": ["q"],
            "help": ["h"],
            "refresh": ["r"],
            "filter": ["/"],
            "clear_filter": ["Esc"],
            "details": ["Enter"],
            "delete": ["d"],
            "restart": ["r"],
            "logs": ["l"],
            "inspect": ["i"],
            "attach": ["a"],
            "history": ["h"],
            "sort": ["t"],
            "command": [":"],
            "navigate_up": ["Up"],
            "navigate_down": ["Down"],
            "navigate_left": ["Left"],
            "navigate_right": ["Right"],
            "page_up": ["PageUp"],
            "page_down": ["PageDown"],
            "home": ["Home"],
            "end": ["End"]
        }

    @staticmethod
    def get_test_search_terms() -> List[str]:
        """
        Get test search terms for different scenarios.

        Steps:
        1. Define common search terms for testing
        2. Include various types of terms (names, statuses, types)
        3. Return list of search terms

        Expected Outcome:
        - Returns list of test search terms
        - Terms cover various search scenarios
        - Includes common Docker-related terms
        - Can be used for search functionality testing

        Returns:
            List of search terms
        """
        return [
            "test",
            "nginx",
            "redis",
            "postgres",
            "alpine",
            "latest",
            "running",
            "exited",
            "bridge",
            "overlay"
        ]

    @staticmethod
    def get_test_error_scenarios() -> List[Dict[str, Any]]:
        """
        Get test error scenarios.

        Steps:
        1. Define various error scenarios for testing
        2. Include invalid configurations, connection errors, etc.
        3. Map each scenario to expected error type
        4. Return list of error scenario dictionaries

        Expected Outcome:
        - Returns list of error scenario dictionaries
        - Each scenario has name, args, and expected error
        - Covers common error conditions
        - Can be used for error handling testing

        Returns:
            List of error scenario dictionaries
        """
        return [
            {
                "name": "invalid_docker_host",
                "args": ["--host", "invalid-host:9999"],
                "expected_error": "connection"
            },
            {
                "name": "invalid_log_level",
                "args": ["--log-level", "INVALID"],
                "expected_error": "log level"
            },
            {
                "name": "invalid_refresh_interval",
                "args": ["--refresh", "-1"],
                "expected_error": "refresh"
            },
            {
                "name": "invalid_theme_file",
                "args": ["--theme", "nonexistent-theme.yaml"],
                "expected_error": "theme"
            }
        ]

    @staticmethod
    def get_performance_test_scenarios() -> List[Dict[str, Any]]:
        """
        Get performance test scenarios.

        Steps:
        1. Define various performance test scenarios
        2. Include startup time, view switching, search, refresh, memory usage
        3. Map each scenario to maximum acceptable time
        4. Return list of performance test scenario dictionaries

        Expected Outcome:
        - Returns list of performance test scenario dictionaries
        - Each scenario has name, description, and max time
        - Covers key performance metrics
        - Can be used for performance testing

        Returns:
            List of performance test scenario dictionaries
        """
        return [
            {
                "name": "startup_time",
                "description": "Test application startup time",
                "max_time": 5.0
            },
            {
                "name": "view_switching",
                "description": "Test view switching performance",
                "max_time": 3.0
            },
            {
                "name": "search_performance",
                "description": "Test search performance",
                "max_time": 2.0
            },
            {
                "name": "refresh_performance",
                "description": "Test refresh performance",
                "max_time": 1.0
            },
            {
                "name": "memory_usage",
                "description": "Test memory usage over time",
                "max_time": 10.0
            }
        ]

    @staticmethod
    def setup_test_environment():
        """
        Set up test environment.

        Steps:
        1. Create base test directory (/app/test-data)
        2. Create necessary subdirectories:
           - screenshots: For test screenshots
           - reports: For test reports
           - test_data: For test data files
           - logs: For test logs
        3. Ensure all directories exist

        Expected Outcome:
        - Base test directory is created
        - All required subdirectories are created
        - Test environment is ready for use
        - No errors occur during setup
        """
        # Create necessary directories in writable location
        base_dir = Path("/app/test-data")
        directories = [
            "screenshots",
            "reports",
            "test_data",
            "logs"
        ]

        for directory in directories:
            (base_dir / directory).mkdir(exist_ok=True)

    @staticmethod
    def cleanup_test_environment():
        """
        Clean up test environment.

        Steps:
        1. Clean up test data files
        2. Clean up old screenshots (older than 7 days)
        3. Remove old screenshot files from screenshots directory

        Expected Outcome:
        - Test data files are cleaned up
        - Old screenshots are removed
        - Test environment is cleaned up
        - No errors occur during cleanup
        """
        # Clean up test data
        TestHelpers.cleanup_test_data()

        # Clean up screenshots older than 7 days
        screenshots_dir = Path("/app/test-data/screenshots")
        if screenshots_dir.exists():
            current_time = time.time()
            for file in screenshots_dir.iterdir():
                if file.is_file() and current_time - file.stat().st_mtime > 7 * 24 * 3600:
                    file.unlink()

    @staticmethod
    def log_test_result(test_name: str, result: str, details: str = ""):
        """
        Log test result.

        Steps:
        1. Get logger instance
        2. Log test name and result
        3. Log additional details if provided

        Expected Outcome:
        - Test result is logged with appropriate level
        - Test name and result are recorded
        - Additional details are logged if provided
        - Logging is performed safely

        Args:
            test_name: Name of the test
            result: Test result (PASS, FAIL, SKIP)
            details: Additional details
        """
        logger = logging.getLogger(__name__)
        logger.info(f"Test: {test_name} - {result}")
        if details:
            logger.info(f"Details: {details}")

    @staticmethod
    def create_test_report(test_results: List[Dict[str, Any]]) -> str:
        """
        Create a test report.

        Steps:
        1. Create reports directory if it doesn't exist
        2. Generate unique report filename with timestamp
        3. Open report file for writing with UTF-8 encoding
        4. Write report header
        5. Write each test result with details
        6. Close file
        7. Return path to created report file

        Expected Outcome:
        - Reports directory is created if needed
        - Report file is created with unique timestamp
        - Report contains all test results with details
        - File is written with UTF-8 encoding
        - Path to report file is returned

        Args:
            test_results: List of test result dictionaries

        Returns:
            Path to the created report file
        """
        reports_dir = Path("/app/test-data/reports")
        reports_dir.mkdir(exist_ok=True)

        report_file = reports_dir / f"test_report_{int(time.time())}.txt"

        with open(report_file, 'w', encoding='utf-8') as f:
            f.write("WhaleTUI E2E Test Report\n")
            f.write("=" * 50 + "\n\n")

            for result in test_results:
                f.write(f"Test: {result.get('name', 'Unknown')}\n")
                f.write(f"Result: {result.get('result', 'Unknown')}\n")
                f.write(f"Duration: {result.get('duration', 0):.2f}s\n")
                if result.get('details'):
                    f.write(f"Details: {result['details']}\n")
                f.write("-" * 30 + "\n")

        return str(report_file)
