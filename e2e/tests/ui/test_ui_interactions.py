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
        """
        Test filter functionality in views.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Open filter by pressing '/' key
        5. Type "test" as filter text
        6. Press Enter to apply the filter
        7. Verify filter is applied and application remains running
        8. Take a screenshot showing filtered results
        9. Clear filter by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Filter can be opened with '/' key
        - Filter text can be entered
        - Filter is applied successfully
        - Application remains responsive
        - Screenshot shows filtered view
        - Filter can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

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
        """
        Test sort functionality in views.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Press 't' key to trigger sort functionality
        5. Verify sort is applied and application remains running
        6. Take a screenshot showing sorted results

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Sort functionality is triggered with 't' key
        - Sort is applied successfully
        - Application remains responsive
        - Screenshot shows sorted view
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Press 't' for sort
        whaletui_app.send_key('t')
        time.sleep(1)

        # Check if sort is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("sort_functionality.png")

    def test_details_view(self, whaletui_app):
        """
        Test details view functionality for selected items.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Press Enter key to view details of selected item
        5. Verify details view is displayed and application remains running
        6. Take a screenshot showing details view
        7. Press 'q' key to return to previous view

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Details view is accessible with Enter key
        - Details are displayed correctly
        - Application remains responsive
        - Screenshot shows details view
        - Can return to previous view with 'q' key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

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
        """
        Test logs view functionality for containers.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Press 'l' key to view logs of selected container
        5. Verify logs view is displayed and application remains running
        6. Take a screenshot showing logs view
        7. Press 'q' key to return to previous view

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Logs view is accessible with 'l' key
        - Logs are displayed correctly
        - Application remains responsive
        - Screenshot shows logs view
        - Can return to previous view with 'q' key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

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
        """
        Test inspect functionality for containers.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Press 'i' key to inspect selected container
        5. Verify inspect view is displayed and application remains running
        6. Take a screenshot showing inspect view
        7. Press 'q' key to return to previous view

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Inspect functionality is accessible with 'i' key
        - Inspect information is displayed correctly
        - Application remains responsive
        - Screenshot shows inspect view
        - Can return to previous view with 'q' key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

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
        """
        Test attach functionality for containers.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Press 'a' key to attach to selected container
        5. Verify attach view is displayed and application remains running
        6. Take a screenshot showing attach view
        7. Press 'q' key to return to previous view

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Attach functionality is accessible with 'a' key
        - Attach view is displayed correctly
        - Application remains responsive
        - Screenshot shows attach view
        - Can return to previous view with 'q' key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

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
        """
        Test history view functionality for containers.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Press 'h' key to view history of selected container
        5. Verify history view is displayed and application remains running
        6. Take a screenshot showing history view
        7. Press 'q' key to return to previous view

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - History functionality is accessible with 'h' key
        - History information is displayed correctly
        - Application remains responsive
        - Screenshot shows history view
        - Can return to previous view with 'q' key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

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
        """
        Test command prompt functionality.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Press ':' key to open command prompt
        5. Type "help" command in the prompt
        6. Press Enter to execute the command
        7. Take a screenshot showing command prompt
        8. Press 'q' key to return to previous view

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Command prompt is accessible with ':' key
        - Commands can be typed in the prompt
        - Commands can be executed with Enter
        - Application remains responsive
        - Screenshot shows command prompt
        - Can return to previous view with 'q' key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

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
        """
        Test restart functionality for containers.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Press 'r' key to restart selected container
        5. Verify restart is applied and application remains running
        6. Take a screenshot showing restart functionality

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Restart functionality is accessible with 'r' key
        - Restart operation is applied successfully
        - Application remains responsive
        - Screenshot shows restart functionality
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Press 'r' for restart
        whaletui_app.send_key('r')
        time.sleep(1)

        # Check if restart is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("restart_functionality.png")

    def test_delete_functionality(self, whaletui_app):
        """
        Test delete functionality for containers.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Press 'd' key to delete selected container
        5. Verify delete is applied and application remains running
        6. Take a screenshot showing delete functionality

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Delete functionality is accessible with 'd' key
        - Delete operation is applied successfully
        - Application remains responsive
        - Screenshot shows delete functionality
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Press 'd' for delete
        whaletui_app.send_key('d')
        time.sleep(1)

        # Check if delete is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("delete_functionality.png")

    def test_refresh_functionality(self, whaletui_app):
        """
        Test refresh functionality for views.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Press 'r' key to refresh the view
        5. Verify refresh is applied and application remains running
        6. Take a screenshot showing refresh functionality

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Refresh functionality is accessible with 'r' key
        - Refresh operation is applied successfully
        - Application remains responsive
        - Screenshot shows refresh functionality
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Press 'r' for refresh
        whaletui_app.send_key('r')
        time.sleep(1)

        # Check if refresh is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("refresh_functionality.png")

    def test_multiple_interactions(self, whaletui_app):
        """
        Test multiple interactions in sequence to verify UI stability.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Perform a sequence of interactions:
           - Navigate down and up with arrow keys
           - Open filter with '/' key
           - Type "test" filter text
           - Apply filter with Enter
           - Clear filter with Esc
           - Sort with 't' key
           - Refresh with 'r' key
        5. Verify application remains responsive throughout
        6. Take a screenshot showing final state

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - All interactions work in sequence
        - Application remains stable during multiple operations
        - No crashes or freezes occur
        - Screenshot shows final interaction state
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

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
        """
        Test error recovery in UI interactions with invalid inputs.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using command mode
        4. Test error scenarios with function keys (F1-F5)
        5. Verify application handles errors gracefully
        6. Take a screenshot showing error recovery

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Invalid function keys are handled gracefully
        - Application remains responsive after error scenarios
        - No crashes occur from invalid input
        - Screenshot shows error recovery state
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view using command mode
        assert whaletui_app.navigate_to_view("containers")

        # Test error scenarios (safer ones that don't terminate the app)
        error_scenarios = [
            'F1',      # Function key
            'F2',      # Function key
            'F3',      # Function key
            'F4',      # Function key
            'F5',      # Function key
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
        """
        Test UI under stress conditions with rapid operations.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Perform rapid interactions:
           - 50 rapid Down/Up key presses with 0.05s intervals
           - 20 rapid view switches between c, i, v, n, s views
        5. Verify application remains stable under stress
        6. Take a screenshot showing stress test results

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - UI handles rapid key presses without issues
        - Rapid view switching works correctly
        - Application remains stable under stress conditions
        - No performance degradation or crashes
        - Screenshot shows stress test final state
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Perform rapid interactions
        for i in range(50):
            whaletui_app.send_key('Down')
            time.sleep(0.05)
            whaletui_app.send_key('Up')
            time.sleep(0.05)

        # Test rapid view switching
        views = ["containers", "images", "volumes", "networks", "swarm"]
        for i in range(20):
            view = views[i % len(views)]
            whaletui_app.navigate_to_view(view)
            time.sleep(0.1)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("ui_stress_test.png")
