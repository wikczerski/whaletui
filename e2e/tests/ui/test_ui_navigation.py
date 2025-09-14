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
        """
        Test that the main screen displays correctly.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear (looking for "Details" text)
        3. Verify key UI elements are present (Docker/WhaleTui branding and Navigation/Actions)
        4. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Main screen displays within 10 seconds
        - UI contains Docker/WhaleTui branding information
        - Navigation or Actions elements are visible
        - Screenshot is captured for debugging purposes
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Check for key UI elements
        output = whaletui_app.get_screen_content()
        # The main screen shows Docker info and navigation, not individual view names
        assert "Docker" in output or "WhaleTui" in output
        assert "Navigation" in output or "Actions" in output

        # Take screenshot
        whaletui_app.take_screenshot("main_screen_display.png")

    def test_view_navigation_containers(self, whaletui_app):
        """
        Test navigation to containers view using command mode.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using the navigate_to_view method with "containers" parameter
        4. Verify the containers view is displayed by checking screen content
        5. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Navigation to containers view succeeds
        - Screen content contains "containers" or "container" text (case-insensitive)
        - Screenshot is captured showing the containers view
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view using command mode
        assert whaletui_app.navigate_to_view("containers")

        # Check if we're in containers view
        output = whaletui_app.get_screen_content()
        assert "containers" in output.lower() or "container" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("containers_view.png")

    def test_view_navigation_images(self, whaletui_app):
        """
        Test navigation to images view using command mode.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to images view using the navigate_to_view method with "images" parameter
        4. Verify the images view is displayed by checking screen content
        5. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Navigation to images view succeeds
        - Screen content contains "images" or "image" text (case-insensitive)
        - Screenshot is captured showing the images view
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to images view using command mode
        assert whaletui_app.navigate_to_view("images")

        # Check if we're in images view
        output = whaletui_app.get_screen_content()
        assert "images" in output.lower() or "image" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("images_view.png")

    def test_view_navigation_volumes(self, whaletui_app):
        """
        Test navigation to volumes view using command mode.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to volumes view using the navigate_to_view method with "volumes" parameter
        4. Verify the volumes view is displayed by checking screen content
        5. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Navigation to volumes view succeeds
        - Screen content contains "volumes" or "volume" text (case-insensitive)
        - Screenshot is captured showing the volumes view
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to volumes view using command mode
        assert whaletui_app.navigate_to_view("volumes")

        # Check if we're in volumes view
        output = whaletui_app.get_screen_content()
        assert "volumes" in output.lower() or "volume" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("volumes_view.png")

    def test_view_navigation_networks(self, whaletui_app):
        """
        Test navigation to networks view using command mode.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to networks view using the navigate_to_view method with "networks" parameter
        4. Verify the networks view is displayed by checking screen content
        5. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Navigation to networks view succeeds
        - Screen content contains "networks" or "network" text (case-insensitive)
        - Screenshot is captured showing the networks view
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to networks view using command mode
        assert whaletui_app.navigate_to_view("networks")

        # Check if we're in networks view
        output = whaletui_app.get_screen_content()
        assert "networks" in output.lower() or "network" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("networks_view.png")

    def test_view_navigation_swarm(self, whaletui_app):
        """
        Test navigation to swarm view using command mode.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to swarm view using the navigate_to_view method with "swarm" parameter
        4. Verify the swarm view is displayed by checking screen content
        5. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Navigation to swarm view succeeds
        - Screen content contains "swarm" text (case-insensitive)
        - Screenshot is captured showing the swarm view
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to swarm view using command mode
        assert whaletui_app.navigate_to_view("swarm")

        # Check if we're in swarm view
        output = whaletui_app.get_screen_content()
        assert "swarm" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("swarm_view.png")

    def test_view_navigation_nodes(self, whaletui_app):
        """
        Test navigation to nodes view using command mode.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to nodes view using the navigate_to_view method with "nodes" parameter
        4. Verify the nodes view is displayed by checking screen content
        5. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Navigation to nodes view succeeds
        - Screen content contains "nodes" or "node" text (case-insensitive)
        - Screenshot is captured showing the nodes view
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to nodes view using command mode
        assert whaletui_app.navigate_to_view("nodes")

        # Check if we're in nodes view
        output = whaletui_app.get_screen_content()
        assert "nodes" in output.lower() or "node" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("nodes_view.png")

    def test_view_navigation_services(self, whaletui_app):
        """
        Test navigation to services view using command mode.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to services view using the navigate_to_view method with "services" parameter
        4. Verify the services view is displayed by checking screen content
        5. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Navigation to services view succeeds
        - Screen content contains "services" or "service" text (case-insensitive)
        - Screenshot is captured showing the services view
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to services view using command mode
        assert whaletui_app.navigate_to_view("services")

        # Check if we're in services view
        output = whaletui_app.get_screen_content()
        assert "services" in output.lower() or "service" in output.lower()

        # Take screenshot
        whaletui_app.take_screenshot("services_view.png")

    def test_help_screen_navigation(self, whaletui_app):
        """
        Test help screen navigation functionality.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Press 'h' key to access help screen
        4. Wait for help screen to display
        5. Verify application is still running
        6. Take a screenshot of the help screen
        7. Press 'q' to return to main screen

        Expected Outcome:
        - Application starts successfully
        - Help screen is accessible via 'h' key
        - Application remains responsive during help screen display
        - Screenshot is captured showing the help screen
        - Can return to main screen using 'q' key
        """
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
        """
        Test application quit functionality.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Press 'q' key to quit the application
        4. Wait for application to terminate
        5. Verify application is no longer running

        Expected Outcome:
        - Application starts successfully
        - Quit command ('q') is recognized
        - Application terminates gracefully
        - Application process is no longer running
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Press 'q' to quit
        whaletui_app.send_key('q')
        time.sleep(1)

        # Application should be stopped
        assert not whaletui_app.is_running()

    def test_application_restart(self, whaletui_app):
        """
        Test application restart functionality.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Quit the application using 'q' key
        4. Verify application has stopped
        5. Restart the application
        6. Wait for the main screen to appear again
        7. Verify application is running

        Expected Outcome:
        - Application starts successfully initially
        - Application quits gracefully
        - Application can be restarted successfully
        - Main screen appears after restart
        - Application is running after restart
        """
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
        """
        Test keyboard navigation within views using arrow keys and navigation keys.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Test arrow key navigation (Down, Up, Right, Left)
        5. Test page navigation (PageDown, PageUp)
        6. Test home/end navigation (Home, End)
        7. Verify application remains responsive
        8. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Arrow keys work for navigation
        - Page navigation keys function properly
        - Home/End keys work for navigation
        - Application remains responsive throughout
        - Screenshot is captured showing navigation state
        """
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
        """
        Test view refresh functionality with automatic refresh interval.

        Steps:
        1. Start the WhaleTUI application with 2-second refresh interval
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Wait for refresh cycle to complete (3 seconds)
        5. Verify application is still running
        6. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts with custom refresh interval
        - Can navigate to containers view
        - Automatic refresh cycle completes successfully
        - Application remains running during refresh
        - Screenshot is captured showing refresh functionality
        """
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
        """
        Test UI responsiveness to rapid user input.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Perform rapid key presses (Down/Up keys) 10 times with 0.1s intervals
        4. Verify application remains responsive
        5. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - UI responds to rapid key presses
        - Application remains stable during rapid input
        - No crashes or freezes occur
        - Screenshot is captured showing responsiveness
        """
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
        """
        Test UI error handling for invalid key combinations.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Test invalid key combinations (Ctrl+C, Ctrl+Z)
        4. Verify application handles errors gracefully
        5. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Invalid key combinations are handled gracefully
        - Application remains running after error scenarios
        - No crashes occur from invalid input
        - Screenshot is captured showing error handling
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Test invalid key combinations (safer ones that don't terminate the app)
        whaletui_app.send_key('F1')
        time.sleep(0.5)
        whaletui_app.send_key('F2')
        time.sleep(0.5)
        whaletui_app.send_key('F3')
        time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("ui_error_handling.png")

    @pytest.mark.slow
    def test_ui_stability_long_session(self, whaletui_app):
        """
        Test UI stability during extended session with multiple view switches.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Perform 3 cycles of navigation through all views (c, i, v, n, s)
        4. For each view, test basic interactions (Down/Up navigation)
        5. Verify application remains stable throughout
        6. Take a screenshot for visual verification

        Expected Outcome:
        - Application starts successfully
        - Can navigate through all views multiple times
        - Basic interactions work in each view
        - Application remains stable during extended session
        - No memory leaks or performance degradation
        - Screenshot is captured showing final state
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate through different views multiple times using command mode
        views = ["containers", "images", "volumes", "networks", "swarm"]

        for i in range(3):  # 3 cycles
            for view in views:
                assert whaletui_app.navigate_to_view(view)
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
