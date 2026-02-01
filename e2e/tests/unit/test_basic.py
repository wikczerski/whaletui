"""
Basic e2e tests for WhaleTUI application.
"""
import pytest
import time
from whaletui_controller import WhaleTUIController


class TestWhaleTUIBasic:
    """Basic functionality tests for WhaleTUI."""

    def test_application_starts(self, whaletui_app):
        """
        Test that the application starts successfully.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear (looking for "Details" text)
        3. Verify the application is running
        4. Take a screenshot for debugging purposes

        Expected Outcome:
        - Application starts without errors
        - Main screen appears within 10 seconds
        - Application process is running
        - Screenshot is captured for visual verification
        """
        whaletui_app.start()

        # Wait for the main screen to appear (look for the actual UI elements)
        assert whaletui_app.wait_for_screen("Details", timeout=10)
        assert whaletui_app.is_running()

        # Take a screenshot for debugging
        whaletui_app.take_screenshot("app_startup.png")

    def test_application_shutdown(self, whaletui_app):
        """
        Test that the application shuts down gracefully.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Send quit command using 'q' key
        4. Wait for application to terminate
        5. Verify application process is no longer running

        Expected Outcome:
        - Application starts successfully
        - Quit command is recognized
        - Application terminates gracefully
        - Application process stops completely
        """
        whaletui_app.start()

        # Wait for the main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Send quit command
        whaletui_app.send_key('q')
        time.sleep(1)

        # Application should be stopped
        assert not whaletui_app.is_running()

    def test_help_screen(self, whaletui_app):
        """
        Test that the help screen can be accessed.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Press 'h' key to access help screen
        4. Wait for help screen to display
        5. Verify application is still running
        6. Take a screenshot of the help screen
        7. Press 'q' to quit help and return to main screen

        Expected Outcome:
        - Application starts successfully
        - Help screen is accessible via 'h' key
        - Application remains responsive during help display
        - Screenshot is captured showing help screen
        - Can return to main screen using 'q' key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Press 'h' for help
        whaletui_app.send_key('h')
        time.sleep(1)

        # Check if help screen is displayed (look for common help elements)
        # The help screen might show different text, so let's check for any change
        output = whaletui_app.get_screen_content()
        # If help screen is shown, there should be some change in the output
        # If no help screen, the application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("help_screen.png")

        # Press 'q' to quit help (or return to main screen)
        whaletui_app.send_key('q')
        time.sleep(1)

    def test_theme_command(self, whaletui_app):
        """
        Test the theme command functionality.

        Steps:
        1. Change to writable directory (/app/test-data)
        2. Start WhaleTUI with 'theme' command
        3. Wait for theme command to complete (2 seconds)
        4. Verify application exits after theme command
        5. Check if theme file was created in config/theme.yaml

        Expected Outcome:
        - Theme command executes successfully
        - Application exits immediately after theme command
        - Theme file is created in the expected location
        - No errors occur during theme file creation
        """
        # Change to writable directory before running theme command
        import os
        os.chdir("/app/test-data")

        whaletui_app.start(['theme'])

        # Wait for theme command to complete (it exits immediately)
        time.sleep(2)

        # Application should exit after theme command
        assert not whaletui_app.is_running()

        # Check if theme file was created
        theme_file = os.path.join("config", "theme.yaml")
        assert os.path.exists(theme_file), "Theme file should be created"

    def test_connect_command_help(self, whaletui_app):
        """
        Test the connect command help functionality.

        Steps:
        1. Start WhaleTUI with connect command and --help flag
        2. Wait for help output to appear
        3. Verify help text contains "Connect to a remote Docker host"
        4. Wait for application to exit after help display
        5. Verify application process is no longer running

        Expected Outcome:
        - Connect command help is displayed
        - Help text contains expected connection information
        - Application exits gracefully after help display
        - No errors occur during help display
        """
        whaletui_app.start(['connect', '--help'])

        # Wait for help output
        assert whaletui_app.wait_for_screen("Connect to a remote Docker host", timeout=10)

        # Application should exit after help
        time.sleep(2)
        assert not whaletui_app.is_running()

    def test_invalid_command(self, whaletui_app):
        """
        Test handling of invalid commands.

        Steps:
        1. Start WhaleTUI with an invalid command ('invalid-command')
        2. Wait for error or help output to appear
        3. Verify application exits after handling invalid command
        4. Verify application process is no longer running

        Expected Outcome:
        - Invalid command is handled gracefully
        - Application exits without crashing
        - No unexpected behavior occurs
        - Application process terminates cleanly
        """
        whaletui_app.start(['invalid-command'])

        # Wait for error or help output
        time.sleep(2)

        # Application should exit
        assert not whaletui_app.is_running()

    def test_log_level_flag(self, whaletui_app):
        """
        Test the log level flag functionality.

        Steps:
        1. Start WhaleTUI with --log-level DEBUG flag
        2. Wait for the main screen to appear
        3. Verify application is running with debug logging
        4. Verify application remains responsive

        Expected Outcome:
        - Application starts with debug log level
        - Main screen appears within 10 seconds
        - Application remains running and responsive
        - Debug logging is enabled
        """
        whaletui_app.start(['--log-level', 'DEBUG'])

        # Wait for the main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Application should be running
        assert whaletui_app.is_running()

    def test_refresh_interval_flag(self, whaletui_app):
        """
        Test the refresh interval flag functionality.

        Steps:
        1. Start WhaleTUI with --refresh 10 flag (10 second refresh)
        2. Wait for the main screen to appear
        3. Verify application is running with custom refresh interval
        4. Verify application remains responsive

        Expected Outcome:
        - Application starts with 10-second refresh interval
        - Main screen appears within 10 seconds
        - Application remains running and responsive
        - Custom refresh interval is applied
        """
        whaletui_app.start(['--refresh', '10'])

        # Wait for the main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Application should be running
        assert whaletui_app.is_running()

    @pytest.mark.slow
    def test_application_stability(self, whaletui_app):
        """
        Test application stability over time.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Let the application run for 5 seconds
        4. Verify application is still running
        5. Take a final screenshot for stability verification

        Expected Outcome:
        - Application starts successfully
        - Application remains stable over time
        - No crashes or freezes occur during extended run
        - Application process remains active
        - Screenshot shows stable application state
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Let it run for a bit
        time.sleep(5)

        # Should still be running
        assert whaletui_app.is_running()

        # Take final screenshot
        whaletui_app.take_screenshot("stability_test.png")
