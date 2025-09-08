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

        Args:
            filename: Name of the file to create
            content: Content to write to the file

        Returns:
            Path to the created file
        """
        test_data_dir = Path(__file__).parent.parent.parent / "test_data"
        test_data_dir.mkdir(exist_ok=True)

        file_path = test_data_dir / filename
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)

        return str(file_path)

    @staticmethod
    def cleanup_test_data():
        """Clean up test data files."""
        test_data_dir = Path(__file__).parent.parent.parent / "test_data"
        if test_data_dir.exists():
            import shutil
            shutil.rmtree(test_data_dir)

    @staticmethod
    def get_docker_test_containers() -> List[Dict[str, Any]]:
        """
        Get list of test containers for Docker tests.

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
        """Set up test environment."""
        # Create necessary directories
        directories = [
            "screenshots",
            "reports",
            "test_data",
            "logs"
        ]

        for directory in directories:
            Path(directory).mkdir(exist_ok=True)

    @staticmethod
    def cleanup_test_environment():
        """Clean up test environment."""
        # Clean up test data
        TestHelpers.cleanup_test_data()

        # Clean up screenshots older than 7 days
        screenshots_dir = Path("screenshots")
        if screenshots_dir.exists():
            current_time = time.time()
            for file in screenshots_dir.iterdir():
                if file.is_file() and current_time - file.stat().st_mtime > 7 * 24 * 3600:
                    file.unlink()

    @staticmethod
    def log_test_result(test_name: str, result: str, details: str = ""):
        """
        Log test result.

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

        Args:
            test_results: List of test result dictionaries

        Returns:
            Path to the created report file
        """
        reports_dir = Path("reports")
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
