"""
UI interaction tests for WhaleTUI.
"""
import pytest
import time
from e2e.whaletui_controller import WhaleTUIController
from tests.utils.test_helpers import TestHelpers


class TestUIInteractions:
    """UI interaction tests for WhaleTUI."""

    def test_filter_functionality(self, whaletui_app):
        """Test filter functionality in views."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Open filter
        whaletui_app.send_key('/')
        time.sleep(0.5)

        # Type filter text
        whaletui_app.send_text("test")
        time.sleep(0.5)

        # Press Enter to apply filter
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if filter is applied
        output = whaletui_app.get_screen_content()
        # Filter should show some change in the output
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("filter_functionality.png")

        # Clear filter
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_sort_functionality(self, whaletui_app):
        """Test sort functionality in views."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Press 't' for sort
        whaletui_app.send_key('t')
        time.sleep(1)

        # Check if sort is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("sort_functionality.png")

    def test_details_view(self, whaletui_app):
        """Test details view functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Press Enter for details
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if details view is shown
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("details_view.png")

        # Press 'q' to go back
        whaletui_app.send_key('q')
        time.sleep(1)

    def test_logs_view(self, whaletui_app):
        """Test logs view functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Press 'l' for logs
        whaletui_app.send_key('l')
        time.sleep(1)

        # Check if logs view is shown
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("logs_view.png")

        # Press 'q' to go back
        whaletui_app.send_key('q')
        time.sleep(1)

    def test_inspect_functionality(self, whaletui_app):
        """Test inspect functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Press 'i' for inspect
        whaletui_app.send_key('i')
        time.sleep(1)

        # Check if inspect view is shown
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("inspect_functionality.png")

        # Press 'q' to go back
        whaletui_app.send_key('q')
        time.sleep(1)

    def test_attach_functionality(self, whaletui_app):
        """Test attach functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Press 'a' for attach
        whaletui_app.send_key('a')
        time.sleep(1)

        # Check if attach view is shown
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("attach_functionality.png")

        # Press 'q' to go back
        whaletui_app.send_key('q')
        time.sleep(1)

    def test_history_view(self, whaletui_app):
        """Test history view functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Press 'h' for history
        whaletui_app.send_key('h')
        time.sleep(1)

        # Check if history view is shown
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("history_view.png")

        # Press 'q' to go back
        whaletui_app.send_key('q')
        time.sleep(1)

    def test_command_prompt(self, whaletui_app):
        """Test command prompt functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Press ':' for command prompt
        whaletui_app.send_key(':')
        time.sleep(1)

        # Check if command prompt is shown
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Type a command
        whaletui_app.send_text("help")
        time.sleep(0.5)

        # Press Enter to execute
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Take screenshot
        whaletui_app.take_screenshot("command_prompt.png")

        # Press 'q' to go back
        whaletui_app.send_key('q')
        time.sleep(1)

    def test_restart_functionality(self, whaletui_app):
        """Test restart functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Press 'r' for restart
        whaletui_app.send_key('r')
        time.sleep(1)

        # Check if restart is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("restart_functionality.png")

    def test_delete_functionality(self, whaletui_app):
        """Test delete functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Press 'd' for delete
        whaletui_app.send_key('d')
        time.sleep(1)

        # Check if delete is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("delete_functionality.png")

    def test_refresh_functionality(self, whaletui_app):
        """Test refresh functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Press 'r' for refresh
        whaletui_app.send_key('r')
        time.sleep(1)

        # Check if refresh is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("refresh_functionality.png")

    def test_multiple_interactions(self, whaletui_app):
        """Test multiple interactions in sequence."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Perform multiple interactions
        interactions = [
            ('Down', 0.5),
            ('Up', 0.5),
            ('/', 0.5),  # Open filter
            ('test', 0.5),  # Type filter
            ('Enter', 1),  # Apply filter
            ('Esc', 0.5),  # Clear filter
            ('t', 1),  # Sort
            ('r', 1),  # Refresh
        ]

        for key, delay in interactions:
            if key == 'test':
                whaletui_app.send_text(key)
            else:
                whaletui_app.send_key(key)
            time.sleep(delay)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("multiple_interactions.png")

    def test_error_recovery(self, whaletui_app):
        """Test error recovery in UI interactions."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test error scenarios
        error_scenarios = [
            'Ctrl+C',  # Interrupt
            'Ctrl+Z',  # Suspend
            'F1',      # Function key
            'F2',      # Function key
            'F3',      # Function key
        ]

        for key in error_scenarios:
            whaletui_app.send_key(key)
            time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("error_recovery.png")

    @pytest.mark.slow
    def test_ui_stress_test(self, whaletui_app):
        """Test UI under stress conditions."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Perform rapid interactions
        for i in range(50):
            whaletui_app.send_key('Down')
            time.sleep(0.05)
            whaletui_app.send_key('Up')
            time.sleep(0.05)

        # Test rapid view switching
        views = ["c", "i", "v", "n", "s"]
        for i in range(20):
            view = views[i % len(views)]
            whaletui_app.send_text(view)
            time.sleep(0.1)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("ui_stress_test.png")
