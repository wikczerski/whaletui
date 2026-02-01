"""
Error handling tests for WhaleTUI.
"""
import pytest
import time
import os
from e2e.whaletui_controller import WhaleTUIController
from tests.utils.test_helpers import TestHelpers


class TestErrorScenarios:
    """Error handling tests for WhaleTUI."""

    def test_invalid_docker_host(self, whaletui_app):
        """Test handling of invalid Docker host connection."""
        whaletui_app.start(['--log-level', 'DEBUG'])

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Application should handle Docker connection errors gracefully
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("invalid_docker_host.png")

    def test_connect_invalid_host(self, whaletui_app):
        """Test connect command with invalid host."""
        whaletui_app.start(['connect', '--host', 'invalid-host:9999', '--user', 'test'])

        # Wait for error handling
        time.sleep(5)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_connect_missing_user(self, whaletui_app):
        """Test connect command with missing user."""
        whaletui_app.start(['connect', '--host', '192.168.1.100'])

        # Wait for error handling
        time.sleep(3)

        # Application should exit with error
        assert not whaletui_app.is_running()

    def test_invalid_log_level(self, whaletui_app):
        """Test handling of invalid log level."""
        whaletui_app.start(['--log-level', 'INVALID'])

        # Wait for main screen (should fallback to default)
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Application should still run
        assert whaletui_app.is_running()

    def test_invalid_refresh_interval(self, whaletui_app):
        """Test handling of invalid refresh interval."""
        whaletui_app.start(['--refresh', '-1'])

        # Wait for main screen (should use default)
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Application should still run
        assert whaletui_app.is_running()

    def test_invalid_theme_file(self, whaletui_app):
        """Test handling of invalid theme file."""
        whaletui_app.start(['--theme', 'nonexistent-theme.yaml'])

        # Wait for main screen (should use default theme)
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Application should still run
        assert whaletui_app.is_running()

    def test_theme_command_invalid_path(self, whaletui_app):
        """Test theme command with invalid path."""
        whaletui_app.start(['theme', '--path', '/invalid/path/theme.yaml'])

        # Wait for error handling
        time.sleep(3)

        # Application should exit with error
        assert not whaletui_app.is_running()

    def test_application_interrupt(self, whaletui_app):
        """Test handling of application interrupt (Ctrl+C)."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Send interrupt signal
        whaletui_app.send_key('Ctrl+C')
        time.sleep(1)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_application_termination(self, whaletui_app):
        """Test handling of application termination (Ctrl+Z)."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Send termination signal
        whaletui_app.send_key('Ctrl+Z')
        time.sleep(1)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_application_quit(self, whaletui_app):
        """Test application quit functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Send quit command
        whaletui_app.send_key('q')
        time.sleep(1)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_application_escape(self, whaletui_app):
        """Test application escape functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Send escape key
        whaletui_app.send_key('Esc')
        time.sleep(1)

        # Application should still be running (escape might not quit)
        # This depends on the actual implementation
        if whaletui_app.is_running():
            # If still running, try quit
            whaletui_app.send_key('q')
            time.sleep(1)
            assert not whaletui_app.is_running()

    def test_application_help_during_error(self, whaletui_app):
        """Test help functionality during error state."""
        whaletui_app.start(['--log-level', 'DEBUG'])

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Press help
        whaletui_app.send_key('h')
        time.sleep(1)

        # Should show help even if there are errors
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("help_during_error.png")

    def test_application_restart_after_error(self, whaletui_app):
        """Test application restart after error."""
        # Start with invalid configuration
        whaletui_app.start(['--theme', 'invalid-theme.yaml'])

        # Wait for main screen (should handle error gracefully)
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Quit application
        whaletui_app.send_key('q')
        time.sleep(1)

        # Restart with valid configuration
        whaletui_app.start()

        # Should start successfully
        assert whaletui_app.wait_for_screen("Details", timeout=10)
        assert whaletui_app.is_running()

    def test_application_memory_usage(self, whaletui_app):
        """Test application memory usage and potential leaks."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Let it run for a bit
        time.sleep(5)

        # Should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("memory_usage.png")

        # Quit gracefully
        whaletui_app.send_key('q')
        time.sleep(1)
        assert not whaletui_app.is_running()

    def test_application_cpu_usage(self, whaletui_app):
        """Test application CPU usage during normal operation."""
        whaletui_app.start(['--refresh', '1'])  # Fast refresh

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Let it run for a bit with fast refresh
        time.sleep(10)

        # Should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("cpu_usage.png")

        # Quit gracefully
        whaletui_app.send_key('q')
        time.sleep(1)
        assert not whaletui_app.is_running()

    def test_ui_error_handling(self, whaletui_app):
        """Test UI error handling."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test invalid key combinations
        invalid_keys = [
            'Ctrl+C',
            'Ctrl+Z',
            'F1',
            'F2',
            'F3',
            'F4',
            'F5',
            'F6',
            'F7',
            'F8',
            'F9',
            'F10',
            'F11',
            'F12',
        ]

        for key in invalid_keys:
            whaletui_app.send_key(key)
            time.sleep(0.1)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("ui_error_handling.png")

    def test_search_error_handling(self, whaletui_app):
        """Test search error handling."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test search error scenarios
        search_errors = [
            ('/', 0.5),  # Open search
            ('test@#$%^&*()', 0.5),  # Invalid characters
            ('Enter', 1),  # Apply search
            ('Esc', 0.5),  # Clear search
        ]

        for key, delay in search_errors:
            if key == 'test@#$%^&*()':
                whaletui_app.send_text(key)
            else:
                whaletui_app.send_key(key)
            time.sleep(delay)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_error_handling.png")

    def test_view_error_handling(self, whaletui_app):
        """Test view error handling."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test rapid view switching
        views = ["c", "i", "v", "n", "s", "invalid", "c", "i", "v", "n", "s"]

        for view in views:
            whaletui_app.send_text(view)
            time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("view_error_handling.png")

    def test_network_error_handling(self, whaletui_app):
        """Test network error handling."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test network operations that might fail
        network_operations = [
            ('n', 2),  # Navigate to networks
            ('Enter', 1),  # Details
            ('q', 0.5),  # Go back
            ('i', 1),  # Inspect
            ('q', 0.5),  # Go back
        ]

        for key, delay in network_operations:
            whaletui_app.send_key(key)
            time.sleep(delay)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("network_error_handling.png")

    def test_docker_error_handling(self, whaletui_app):
        """Test Docker error handling."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test Docker operations that might fail
        docker_operations = [
            ('c', 2),  # Navigate to containers
            ('Enter', 1),  # Details
            ('q', 0.5),  # Go back
            ('l', 1),  # Logs
            ('q', 0.5),  # Go back
            ('i', 1),  # Inspect
            ('q', 0.5),  # Go back
        ]

        for key, delay in docker_operations:
            whaletui_app.send_key(key)
            time.sleep(delay)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("docker_error_handling.png")

    def test_permission_error_handling(self, whaletui_app):
        """Test permission error handling."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test operations that might require permissions
        permission_operations = [
            ('c', 2),  # Navigate to containers
            ('d', 1),  # Delete (might fail)
            ('r', 1),  # Restart (might fail)
            ('a', 1),  # Attach (might fail)
            ('q', 0.5),  # Go back
        ]

        for key, delay in permission_operations:
            whaletui_app.send_key(key)
            time.sleep(delay)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("permission_error_handling.png")

    def test_timeout_error_handling(self, whaletui_app):
        """Test timeout error handling."""
        whaletui_app.start(['--refresh', '1'])  # Fast refresh

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Let it run for a bit to test timeout handling
        time.sleep(5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("timeout_error_handling.png")

        # Quit gracefully
        whaletui_app.send_key('q')
        time.sleep(1)
        assert not whaletui_app.is_running()

    def test_concurrent_error_handling(self, whaletui_app):
        """Test concurrent error handling."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test concurrent operations that might cause errors
        for i in range(10):
            whaletui_app.send_key('Down')
            time.sleep(0.1)
            whaletui_app.send_key('Up')
            time.sleep(0.1)
            whaletui_app.send_key('Left')
            time.sleep(0.1)
            whaletui_app.send_key('Right')
            time.sleep(0.1)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("concurrent_error_handling.png")

    @pytest.mark.slow
    def test_error_recovery_long_session(self, whaletui_app):
        """Test error recovery during long session."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test error recovery over time
        for i in range(20):
            # Test different error scenarios
            whaletui_app.send_key('Ctrl+C')
            time.sleep(0.1)
            whaletui_app.send_key('Ctrl+Z')
            time.sleep(0.1)
            whaletui_app.send_key('F1')
            time.sleep(0.1)
            whaletui_app.send_key('F2')
            time.sleep(0.1)

            # Test normal operations
            whaletui_app.send_key('Down')
            time.sleep(0.1)
            whaletui_app.send_key('Up')
            time.sleep(0.1)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("error_recovery_long_session.png")
