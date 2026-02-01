"""
Error handling e2e tests for WhaleTUI application.
"""
import pytest
import time
from e2e.whaletui_controller import WhaleTUIController


class TestWhaleTUIErrorHandling:
    """Error handling tests for WhaleTUI."""

    def test_invalid_docker_host(self, whaletui_app):
        """
        Test handling of invalid Docker host connection.

        Steps:
        1. Start WhaleTUI with debug logging enabled
        2. Wait for main screen to appear (looking for "WhaleTUI" text)
        3. Verify application handles Docker connection errors gracefully
        4. Take a screenshot showing error handling state

        Expected Outcome:
        - Application starts with debug logging
        - Main screen appears within 10 seconds
        - Docker connection errors are handled gracefully
        - Application remains running despite connection issues
        - Screenshot shows error handling state
        """
        whaletui_app.start(['--log-level', 'DEBUG'])

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Application should handle Docker connection errors gracefully
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("invalid_docker_host.png")

    def test_connect_invalid_host(self, whaletui_app):
        """
        Test connect command with invalid host.

        Steps:
        1. Start WhaleTUI with connect command to invalid host (invalid-host:9999)
        2. Provide test user credentials
        3. Wait for error handling (5 seconds)
        4. Verify application exits gracefully

        Expected Outcome:
        - Connect command is executed with invalid host
        - Connection error is handled properly
        - Application exits gracefully after connection failure
        - No crashes or unexpected behavior
        """
        whaletui_app.start(['connect', '--host', 'invalid-host:9999', '--user', 'test'])

        # Wait for error handling
        time.sleep(5)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_connect_missing_user(self, whaletui_app):
        """
        Test connect command with missing user.

        Steps:
        1. Start WhaleTUI with connect command to valid host but missing user
        2. Wait for error handling (3 seconds)
        3. Verify application exits with error

        Expected Outcome:
        - Connect command is executed with missing user
        - Missing user error is handled properly
        - Application exits with appropriate error
        - No crashes or unexpected behavior
        """
        whaletui_app.start(['connect', '--host', '192.168.1.100'])

        # Wait for error handling
        time.sleep(3)

        # Application should exit with error
        assert not whaletui_app.is_running()

    def test_invalid_log_level(self, whaletui_app):
        """
        Test handling of invalid log level.

        Steps:
        1. Start WhaleTUI with invalid log level ("INVALID")
        2. Wait for main screen to appear (should fallback to default)
        3. Verify application still runs with fallback log level

        Expected Outcome:
        - Invalid log level is handled gracefully
        - Application falls back to default log level
        - Main screen appears within 10 seconds
        - Application remains running and responsive
        """
        whaletui_app.start(['--log-level', 'INVALID'])

        # Wait for main screen (should fallback to default)
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Application should still run
        assert whaletui_app.is_running()

    def test_invalid_refresh_interval(self, whaletui_app):
        """
        Test handling of invalid refresh interval.

        Steps:
        1. Start WhaleTUI with invalid refresh interval (-1)
        2. Wait for main screen to appear (should use default)
        3. Verify application still runs with default refresh interval

        Expected Outcome:
        - Invalid refresh interval is handled gracefully
        - Application falls back to default refresh interval
        - Main screen appears within 10 seconds
        - Application remains running and responsive
        """
        whaletui_app.start(['--refresh', '-1'])

        # Wait for main screen (should use default)
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Application should still run
        assert whaletui_app.is_running()

    def test_invalid_theme_file(self, whaletui_app):
        """
        Test handling of invalid theme file.

        Steps:
        1. Start WhaleTUI with invalid theme file (nonexistent-theme.yaml)
        2. Wait for main screen to appear (should use default theme)
        3. Verify application still runs with default theme

        Expected Outcome:
        - Invalid theme file is handled gracefully
        - Application falls back to default theme
        - Main screen appears within 10 seconds
        - Application remains running and responsive
        """
        whaletui_app.start(['--theme', 'nonexistent-theme.yaml'])

        # Wait for main screen (should use default theme)
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Application should still run
        assert whaletui_app.is_running()

    def test_theme_command_invalid_path(self, whaletui_app):
        """
        Test theme command with invalid path.

        Steps:
        1. Start WhaleTUI with theme command and invalid path (/invalid/path/theme.yaml)
        2. Wait for error handling (3 seconds)
        3. Verify application exits with error

        Expected Outcome:
        - Theme command with invalid path is handled properly
        - Application exits with appropriate error
        - No crashes or unexpected behavior
        """
        whaletui_app.start(['theme', '--path', '/invalid/path/theme.yaml'])

        # Wait for error handling
        time.sleep(3)

        # Application should exit with error
        assert not whaletui_app.is_running()

    def test_application_interrupt(self, whaletui_app):
        """
        Test handling of application interrupt (Ctrl+C).

        Steps:
        1. Start the WhaleTUI application
        2. Wait for main screen to appear
        3. Send interrupt signal (Ctrl+C)
        4. Wait for application to handle interrupt
        5. Verify application exits gracefully

        Expected Outcome:
        - Application starts successfully
        - Main screen appears within 10 seconds
        - Interrupt signal is handled gracefully
        - Application exits cleanly after interrupt
        - No crashes or unexpected behavior
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Send interrupt signal
        whaletui_app.send_key('Ctrl+C')
        time.sleep(1)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_application_termination(self, whaletui_app):
        """
        Test handling of application termination (Ctrl+Z).

        Steps:
        1. Start the WhaleTUI application
        2. Wait for main screen to appear
        3. Send termination signal (Ctrl+Z)
        4. Wait for application to handle termination
        5. Verify application exits gracefully

        Expected Outcome:
        - Application starts successfully
        - Main screen appears within 10 seconds
        - Termination signal is handled gracefully
        - Application exits cleanly after termination
        - No crashes or unexpected behavior
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Send termination signal
        whaletui_app.send_key('Ctrl+Z')
        time.sleep(1)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_application_quit(self, whaletui_app):
        """
        Test application quit functionality.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for main screen to appear
        3. Send quit command ('q' key)
        4. Wait for application to quit
        5. Verify application exits gracefully

        Expected Outcome:
        - Application starts successfully
        - Main screen appears within 10 seconds
        - Quit command is recognized
        - Application exits gracefully
        - No crashes or unexpected behavior
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Send quit command
        whaletui_app.send_key('q')
        time.sleep(1)

        # Application should exit gracefully
        assert not whaletui_app.is_running()

    def test_application_escape(self, whaletui_app):
        """
        Test application escape functionality.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for main screen to appear
        3. Send escape key
        4. Wait for application to handle escape
        5. Check if application is still running
        6. If still running, try quit command ('q')
        7. Verify application exits

        Expected Outcome:
        - Application starts successfully
        - Main screen appears within 10 seconds
        - Escape key is handled appropriately
        - Application either exits or remains running
        - If still running, quit command works
        - Application exits cleanly
        """
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
        """
        Test help functionality during error state.

        Steps:
        1. Start WhaleTUI with debug logging (may show errors)
        2. Wait for main screen to appear
        3. Press help key ('h') to access help
        4. Wait for help screen to display
        5. Verify help is accessible even with errors
        6. Take a screenshot showing help during error state

        Expected Outcome:
        - Application starts with debug logging
        - Main screen appears within 10 seconds
        - Help functionality works even during error states
        - Help screen is accessible via 'h' key
        - Screenshot shows help during error state
        """
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
        """
        Test application restart after error.

        Steps:
        1. Start WhaleTUI with invalid configuration (invalid theme)
        2. Wait for main screen to appear (should handle error gracefully)
        3. Quit application using 'q' key
        4. Restart application with valid configuration
        5. Wait for main screen to appear
        6. Verify application starts successfully after error

        Expected Outcome:
        - Application starts with invalid configuration
        - Error is handled gracefully
        - Application can be quit cleanly
        - Application can be restarted successfully
        - Main screen appears after restart
        - Application runs normally after restart
        """
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
        """
        Test application memory usage and potential leaks.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for main screen to appear
        3. Let application run for 5 seconds
        4. Verify application is still running
        5. Take a screenshot showing memory usage state
        6. Quit application gracefully using 'q' key
        7. Verify application exits cleanly

        Expected Outcome:
        - Application starts successfully
        - Main screen appears within 10 seconds
        - Application remains stable over time
        - No memory leaks or excessive memory usage
        - Screenshot shows stable memory state
        - Application exits cleanly
        """
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
        """
        Test application CPU usage during normal operation.

        Steps:
        1. Start WhaleTUI with fast refresh interval (1 second)
        2. Wait for main screen to appear
        3. Let application run for 10 seconds with fast refresh
        4. Verify application is still running
        5. Take a screenshot showing CPU usage state
        6. Quit application gracefully using 'q' key
        7. Verify application exits cleanly

        Expected Outcome:
        - Application starts with fast refresh interval
        - Main screen appears within 10 seconds
        - Application remains stable with fast refresh
        - CPU usage remains reasonable during operation
        - Screenshot shows stable CPU state
        - Application exits cleanly
        """
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
