"""
UI navigation tests for WhaleTUI.
"""
import pytest
import time
from e2e.whaletui_controller import WhaleTUIController
from tests.utils.test_helpers import TestHelpers


class TestUINavigation:
    """UI navigation tests for WhaleTUI."""

    def test_main_screen_display(self, whaletui_app):
        """Test that the main screen displays correctly."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Check for key UI elements
        output = whaletui_app.get_screen_content()
        assert "Containers" in output
        assert "Enter" in output
        assert "Quit" in output

        # Take screenshot
        whaletui_app.take_screenshot("main_screen_display.png")

    def test_view_navigation_containers(self, whaletui_app):
        """Test navigation to containers view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Check if we're in containers view
        output = whaletui_app.get_screen_content()
        assert "containers" in output.lower() or "container" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("containers_view.png")

    def test_view_navigation_images(self, whaletui_app):
        """Test navigation to images view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to images view
        whaletui_app.send_text("i")
        time.sleep(2)

        # Check if we're in images view
        output = whaletui_app.get_screen_content()
        assert "images" in output.lower() or "image" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("images_view.png")

    def test_view_navigation_volumes(self, whaletui_app):
        """Test navigation to volumes view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to volumes view
        whaletui_app.send_text("v")
        time.sleep(2)

        # Check if we're in volumes view
        output = whaletui_app.get_screen_content()
        assert "volumes" in output.lower() or "volume" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("volumes_view.png")

    def test_view_navigation_networks(self, whaletui_app):
        """Test navigation to networks view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to networks view
        whaletui_app.send_text("n")
        time.sleep(2)

        # Check if we're in networks view
        output = whaletui_app.get_screen_content()
        assert "networks" in output.lower() or "network" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("networks_view.png")

    def test_view_navigation_swarm(self, whaletui_app):
        """Test navigation to swarm view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to swarm view
        whaletui_app.send_text("s")
        time.sleep(2)

        # Check if we're in swarm view
        output = whaletui_app.get_screen_content()
        assert "swarm" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("swarm_view.png")

    def test_view_navigation_nodes(self, whaletui_app):
        """Test navigation to nodes view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to nodes view
        whaletui_app.send_text("nodes")
        time.sleep(2)

        # Check if we're in nodes view
        output = whaletui_app.get_screen_content()
        assert "nodes" in output.lower() or "node" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("nodes_view.png")

    def test_view_navigation_services(self, whaletui_app):
        """Test navigation to services view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to services view
        whaletui_app.send_text("services")
        time.sleep(2)

        # Check if we're in services view
        output = whaletui_app.get_screen_content()
        assert "services" in output.lower() or "service" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("services_view.png")

    def test_help_screen_navigation(self, whaletui_app):
        """Test help screen navigation."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Press 'h' for help
        whaletui_app.send_key('h')
        time.sleep(1)

        # Check if help screen is displayed
        output = whaletui_app.get_screen_content()
        # Help screen should show some change in the output
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("help_screen.png")

        # Press 'q' to return to main screen
        whaletui_app.send_key('q')
        time.sleep(1)

    def test_quit_application(self, whaletui_app):
        """Test quitting the application."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Press 'q' to quit
        whaletui_app.send_key('q')
        time.sleep(1)

        # Application should be stopped
        assert not whaletui_app.is_running()

    def test_application_restart(self, whaletui_app):
        """Test application restart functionality."""
        # Start application
        whaletui_app.start()
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Quit application
        whaletui_app.send_key('q')
        time.sleep(1)
        assert not whaletui_app.is_running()

        # Restart application
        whaletui_app.start()
        assert whaletui_app.wait_for_screen("Details", timeout=10)
        assert whaletui_app.is_running()

    def test_keyboard_navigation(self, whaletui_app):
        """Test keyboard navigation within views."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test arrow key navigation
        whaletui_app.send_key('Down')
        time.sleep(0.5)
        whaletui_app.send_key('Up')
        time.sleep(0.5)
        whaletui_app.send_key('Right')
        time.sleep(0.5)
        whaletui_app.send_key('Left')
        time.sleep(0.5)

        # Test page navigation
        whaletui_app.send_key('PageDown')
        time.sleep(0.5)
        whaletui_app.send_key('PageUp')
        time.sleep(0.5)

        # Test home/end navigation
        whaletui_app.send_key('Home')
        time.sleep(0.5)
        whaletui_app.send_key('End')
        time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("keyboard_navigation.png")

    def test_view_refresh(self, whaletui_app):
        """Test view refresh functionality."""
        whaletui_app.start(['--refresh', '2'])  # 2 second refresh

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Wait for refresh cycle
        time.sleep(3)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("view_refresh.png")

    def test_ui_responsiveness(self, whaletui_app):
        """Test UI responsiveness to user input."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test rapid key presses
        for i in range(10):
            whaletui_app.send_key('Down')
            time.sleep(0.1)
            whaletui_app.send_key('Up')
            time.sleep(0.1)

        # Application should still be responsive
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("ui_responsiveness.png")

    def test_ui_error_handling(self, whaletui_app):
        """Test UI error handling."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test invalid key combinations
        whaletui_app.send_key('Ctrl+C')
        time.sleep(0.5)
        whaletui_app.send_key('Ctrl+Z')
        time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("ui_error_handling.png")

    @pytest.mark.slow
    def test_ui_stability_long_session(self, whaletui_app):
        """Test UI stability during long session."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate through different views multiple times
        views = ["c", "i", "v", "n", "s"]

        for i in range(3):  # 3 cycles
            for view in views:
                whaletui_app.send_text(view)
                time.sleep(1)

                # Test some interactions
                whaletui_app.send_key('Down')
                time.sleep(0.5)
                whaletui_app.send_key('Up')
                time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("ui_stability_long_session.png")
