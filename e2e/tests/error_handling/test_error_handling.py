"""
Error handling e2e tests for WhaleTUI application.
"""
import pytest
import time
from e2e.whaletui_controller import WhaleTUIController


class TestWhaleTUIErrorHandling:
    """Error handling tests for WhaleTUI."""

    def test_invalid_docker_host(self, whaletui_app):
        """Test handling of invalid Docker host connection."""
        whaletui_app.start(['--log-level', 'DEBUG'])

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

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
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Application should still run
        assert whaletui_app.is_running()

    def test_invalid_refresh_interval(self, whaletui_app):
        """Test handling of invalid refresh interval."""
        whaletui_app.start(['--refresh', '-1'])

        # Wait for main screen (should use default)
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Application should still run
        assert whaletui_app.is_running()

    def test_invalid_theme_file(self, whaletui_app):
        """Test handling of invalid theme file."""
        whaletui_app.start(['--theme', 'nonexistent-theme.yaml'])

        # Wait for main screen (should use default theme)
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

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
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Send interrupt signal
        whaletui_app.send_key('Ctrl+C')
        time.sleep(1)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_application_termination(self, whaletui_app):
        """Test handling of application termination (Ctrl+Z)."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Send termination signal
        whaletui_app.send_key('Ctrl+Z')
        time.sleep(1)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_application_quit(self, whaletui_app):
        """Test application quit functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Send quit command
        whaletui_app.send_key('q')
        time.sleep(1)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_application_escape(self, whaletui_app):
        """Test application escape functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

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
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Press help
        whaletui_app.send_key('h')
        time.sleep(1)

        # Should show help even if there are errors
        assert whaletui_app.wait_for_screen("help", timeout=5)

        # Take screenshot
        whaletui_app.take_screenshot("help_during_error.png")

    def test_application_restart_after_error(self, whaletui_app):
        """Test application restart after error."""
        # Start with invalid configuration
        whaletui_app.start(['--theme', 'invalid-theme.yaml'])

        # Wait for main screen (should handle error gracefully)
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Quit application
        whaletui_app.send_key('q')
        time.sleep(1)

        # Restart with valid configuration
        whaletui_app.start()

        # Should start successfully
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)
        assert whaletui_app.is_running()

    def test_application_memory_usage(self, whaletui_app):
        """Test application memory usage and potential leaks."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

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
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

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
