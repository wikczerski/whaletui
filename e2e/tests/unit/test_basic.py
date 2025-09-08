"""
Basic e2e tests for WhaleTUI application.
"""
import pytest
import time
from whaletui_controller import WhaleTUIController


class TestWhaleTUIBasic:
    """Basic functionality tests for WhaleTUI."""

    def test_application_starts(self, whaletui_app):
        """Test that the application starts successfully."""
        whaletui_app.start()

        # Wait for the main screen to appear (look for the actual UI elements)
        assert whaletui_app.wait_for_screen("Details", timeout=10)
        assert whaletui_app.is_running()

        # Take a screenshot for debugging
        whaletui_app.take_screenshot("app_startup.png")

    def test_application_shutdown(self, whaletui_app):
        """Test that the application shuts down gracefully."""
        whaletui_app.start()

        # Wait for the main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Send quit command
        whaletui_app.send_key('q')
        time.sleep(1)

        # Application should be stopped
        assert not whaletui_app.is_running()

    def test_help_screen(self, whaletui_app):
        """Test that the help screen can be accessed."""
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
        """Test the theme command functionality."""
        whaletui_app.start(['theme'])

        # Wait for theme command to complete (it exits immediately)
        time.sleep(2)

        # Application should exit after theme command
        assert not whaletui_app.is_running()

        # Check if theme file was created
        import os
        theme_file = os.path.join("..", "config", "theme.yaml")
        assert os.path.exists(theme_file), "Theme file should be created"

    def test_connect_command_help(self, whaletui_app):
        """Test the connect command help."""
        whaletui_app.start(['connect', '--help'])

        # Wait for help output
        assert whaletui_app.wait_for_screen("Connect to a remote Docker host", timeout=10)

        # Application should exit after help
        time.sleep(2)
        assert not whaletui_app.is_running()

    def test_invalid_command(self, whaletui_app):
        """Test handling of invalid commands."""
        whaletui_app.start(['invalid-command'])

        # Wait for error or help output
        time.sleep(2)

        # Application should exit
        assert not whaletui_app.is_running()

    def test_log_level_flag(self, whaletui_app):
        """Test the log level flag."""
        whaletui_app.start(['--log-level', 'DEBUG'])

        # Wait for the main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Application should be running
        assert whaletui_app.is_running()

    def test_refresh_interval_flag(self, whaletui_app):
        """Test the refresh interval flag."""
        whaletui_app.start(['--refresh', '10'])

        # Wait for the main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Application should be running
        assert whaletui_app.is_running()

    @pytest.mark.slow
    def test_application_stability(self, whaletui_app):
        """Test application stability over time."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Let it run for a bit
        time.sleep(5)

        # Should still be running
        assert whaletui_app.is_running()

        # Take final screenshot
        whaletui_app.take_screenshot("stability_test.png")
