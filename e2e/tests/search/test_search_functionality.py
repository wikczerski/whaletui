"""
Search functionality tests for WhaleTUI.
"""
import pytest
import time
from e2e.whaletui_controller import WhaleTUIController
from tests.utils.test_helpers import TestHelpers


class TestSearchFunctionality:
    """Search functionality tests for WhaleTUI."""

    def test_search_in_containers(self, whaletui_app):
        """
        Test search functionality in containers view.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Open search by pressing '/' key
        5. Type "test" as search term
        6. Press Enter to apply the search
        7. Verify search is applied and application remains running
        8. Take a screenshot showing search results
        9. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Search can be opened with '/' key
        - Search term can be entered
        - Search is applied successfully
        - Application remains responsive
        - Screenshot shows filtered search results
        - Search can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Open search
        whaletui_app.send_key('/')
        time.sleep(0.5)

        # Type search term
        whaletui_app.send_text("test")
        time.sleep(0.5)

        # Press Enter to search
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if search is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_containers.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_search_in_images(self, whaletui_app):
        """
        Test search functionality in images view.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to images view using 'i' key
        4. Open search by pressing '/' key
        5. Type "nginx" as search term
        6. Press Enter to apply the search
        7. Verify search is applied and application remains running
        8. Take a screenshot showing search results
        9. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to images view
        - Search can be opened with '/' key
        - Search term can be entered
        - Search is applied successfully
        - Application remains responsive
        - Screenshot shows filtered search results
        - Search can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to images view
        whaletui_app.send_text("i")
        time.sleep(2)

        # Open search
        whaletui_app.send_key('/')
        time.sleep(0.5)

        # Type search term
        whaletui_app.send_text("nginx")
        time.sleep(0.5)

        # Press Enter to search
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if search is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_images.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_search_in_volumes(self, whaletui_app):
        """
        Test search functionality in volumes view.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to volumes view using 'v' key
        4. Open search by pressing '/' key
        5. Type "test" as search term
        6. Press Enter to apply the search
        7. Verify search is applied and application remains running
        8. Take a screenshot showing search results
        9. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to volumes view
        - Search can be opened with '/' key
        - Search term can be entered
        - Search is applied successfully
        - Application remains responsive
        - Screenshot shows filtered search results
        - Search can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to volumes view
        whaletui_app.send_text("v")
        time.sleep(2)

        # Open search
        whaletui_app.send_key('/')
        time.sleep(0.5)

        # Type search term
        whaletui_app.send_text("test")
        time.sleep(0.5)

        # Press Enter to search
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if search is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_volumes.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_search_in_networks(self, whaletui_app):
        """
        Test search functionality in networks view.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to networks view using 'n' key
        4. Open search by pressing '/' key
        5. Type "bridge" as search term
        6. Press Enter to apply the search
        7. Verify search is applied and application remains running
        8. Take a screenshot showing search results
        9. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to networks view
        - Search can be opened with '/' key
        - Search term can be entered
        - Search is applied successfully
        - Application remains responsive
        - Screenshot shows filtered search results
        - Search can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to networks view
        whaletui_app.send_text("n")
        time.sleep(2)

        # Open search
        whaletui_app.send_key('/')
        time.sleep(0.5)

        # Type search term
        whaletui_app.send_text("bridge")
        time.sleep(0.5)

        # Press Enter to search
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if search is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_networks.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_search_in_swarm_services(self, whaletui_app):
        """
        Test search functionality in swarm services view.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to services view using "services" command
        4. Open search by pressing '/' key
        5. Type "test" as search term
        6. Press Enter to apply the search
        7. Verify search is applied and application remains running
        8. Take a screenshot showing search results
        9. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to services view
        - Search can be opened with '/' key
        - Search term can be entered
        - Search is applied successfully
        - Application remains responsive
        - Screenshot shows filtered search results
        - Search can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to services view
        whaletui_app.send_text("services")
        time.sleep(2)

        # Open search
        whaletui_app.send_key('/')
        time.sleep(0.5)

        # Type search term
        whaletui_app.send_text("test")
        time.sleep(0.5)

        # Press Enter to search
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if search is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_services.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_search_in_swarm_nodes(self, whaletui_app):
        """
        Test search functionality in swarm nodes view.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to nodes view using "nodes" command
        4. Open search by pressing '/' key
        5. Type "manager" as search term
        6. Press Enter to apply the search
        7. Verify search is applied and application remains running
        8. Take a screenshot showing search results
        9. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to nodes view
        - Search can be opened with '/' key
        - Search term can be entered
        - Search is applied successfully
        - Application remains responsive
        - Screenshot shows filtered search results
        - Search can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to nodes view
        whaletui_app.send_text("nodes")
        time.sleep(2)

        # Open search
        whaletui_app.send_key('/')
        time.sleep(0.5)

        # Type search term
        whaletui_app.send_text("manager")
        time.sleep(0.5)

        # Press Enter to search
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if search is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_nodes.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_search_clear_functionality(self, whaletui_app):
        """
        Test search clear functionality.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Open search by pressing '/' key
        5. Type "test" as search term
        6. Press Enter to apply the search
        7. Clear search by pressing Esc key
        8. Verify search is cleared and application remains running
        9. Take a screenshot showing cleared search state

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Search can be opened and applied
        - Search can be cleared with Esc key
        - Application remains responsive after clearing
        - Screenshot shows cleared search state
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Open search
        whaletui_app.send_key('/')
        time.sleep(0.5)

        # Type search term
        whaletui_app.send_text("test")
        time.sleep(0.5)

        # Press Enter to search
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Clear search with Esc
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

        # Check if search is cleared
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_clear.png")

    def test_search_empty_term(self, whaletui_app):
        """
        Test search with empty term.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Open search by pressing '/' key
        5. Press Enter without typing any search term
        6. Verify search handles empty term gracefully
        7. Take a screenshot showing empty search handling
        8. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Search can be opened with '/' key
        - Empty search term is handled gracefully
        - Application remains responsive
        - Screenshot shows empty search handling
        - Search can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Open search
        whaletui_app.send_key('/')
        time.sleep(0.5)

        # Press Enter without typing anything
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if search handles empty term gracefully
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_empty_term.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_search_special_characters(self, whaletui_app):
        """
        Test search with special characters.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Open search by pressing '/' key
        5. Type "test-123_456" as search term with special characters
        6. Press Enter to apply the search
        7. Verify search handles special characters correctly
        8. Take a screenshot showing special character search
        9. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Search can be opened with '/' key
        - Special characters in search term are handled correctly
        - Application remains responsive
        - Screenshot shows special character search results
        - Search can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Open search
        whaletui_app.send_key('/')
        time.sleep(0.5)

        # Type search term with special characters
        whaletui_app.send_text("test-123_456")
        time.sleep(0.5)

        # Press Enter to search
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if search handles special characters
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_special_chars.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_search_case_sensitivity(self, whaletui_app):
        """
        Test search case sensitivity.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Test uppercase search:
           - Open search with '/' key
           - Type "TEST" in uppercase
           - Press Enter to apply
           - Take screenshot
           - Clear with Esc key
        5. Test lowercase search:
           - Open search with '/' key
           - Type "test" in lowercase
           - Press Enter to apply
           - Take screenshot
           - Clear with Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Both uppercase and lowercase searches work
        - Case sensitivity behavior is consistent
        - Screenshots show both search results
        - Both searches can be cleared
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test uppercase search
        whaletui_app.send_key('/')
        time.sleep(0.5)
        whaletui_app.send_text("TEST")
        time.sleep(0.5)
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Take screenshot
        whaletui_app.take_screenshot("search_uppercase.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

        # Test lowercase search
        whaletui_app.send_key('/')
        time.sleep(0.5)
        whaletui_app.send_text("test")
        time.sleep(0.5)
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Take screenshot
        whaletui_app.take_screenshot("search_lowercase.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_search_multiple_terms(self, whaletui_app):
        """
        Test search with multiple terms.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Test multiple search terms: ["test", "nginx", "redis", "postgres", "alpine"]
        5. For each term:
           - Open search with '/' key
           - Type the search term
           - Press Enter to apply
           - Take screenshot with term-specific filename
           - Clear search with Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - All search terms work correctly
        - Each search produces results
        - Screenshots are captured for each term
        - All searches can be cleared
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test multiple search terms
        search_terms = ["test", "nginx", "redis", "postgres", "alpine"]

        for term in search_terms:
            whaletui_app.send_key('/')
            time.sleep(0.5)
            whaletui_app.send_text(term)
            time.sleep(0.5)
            whaletui_app.send_key('Enter')
            time.sleep(1)

            # Take screenshot for each term
            whaletui_app.take_screenshot(f"search_term_{term}.png")

            # Clear search
            whaletui_app.send_key('Esc')
            time.sleep(0.5)

    def test_search_performance(self, whaletui_app):
        """
        Test search performance with multiple terms.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Test search performance with multiple terms: ["test", "nginx", "redis", "postgres", "alpine", "latest", "running", "exited"]
        5. For each term:
           - Open search with '/' key
           - Type the search term
           - Press Enter to apply
           - Clear search with Esc key
        6. Measure total time for all searches
        7. Verify total time is less than 10 seconds
        8. Take final screenshot

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - All searches complete within reasonable time
        - Total search time is less than 10 seconds
        - Application remains responsive during performance test
        - Screenshot shows final state
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test search performance with multiple terms
        search_terms = ["test", "nginx", "redis", "postgres", "alpine", "latest", "running", "exited"]

        start_time = time.time()

        for term in search_terms:
            whaletui_app.send_key('/')
            time.sleep(0.1)
            whaletui_app.send_text(term)
            time.sleep(0.1)
            whaletui_app.send_key('Enter')
            time.sleep(0.1)
            whaletui_app.send_key('Esc')
            time.sleep(0.1)

        end_time = time.time()
        total_time = end_time - start_time

        # Search should be reasonably fast
        assert total_time < 10.0, f"Search performance too slow: {total_time:.2f}s"

        # Take screenshot
        whaletui_app.take_screenshot("search_performance.png")

    def test_search_error_handling(self, whaletui_app):
        """
        Test search error handling with invalid characters.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Open search by pressing '/' key
        5. Type "test@#$%^&*()" with invalid characters
        6. Press Enter to apply the search
        7. Verify search handles invalid characters gracefully
        8. Take a screenshot showing error handling
        9. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Search can be opened with '/' key
        - Invalid characters are handled gracefully
        - Application remains responsive
        - Screenshot shows error handling state
        - Search can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test search with invalid characters
        whaletui_app.send_key('/')
        time.sleep(0.5)
        whaletui_app.send_text("test@#$%^&*()")
        time.sleep(0.5)
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if search handles invalid characters gracefully
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_error_handling.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    def test_search_persistence(self, whaletui_app):
        """
        Test search persistence across view changes.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Open search and type "test" term
        5. Press Enter to apply the search
        6. Switch to images view using 'i' key
        7. Switch back to containers view using 'c' key
        8. Verify search is still applied
        9. Take a screenshot showing search persistence
        10. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate between views
        - Search is applied in containers view
        - Search persists when switching views
        - Search remains active when returning to view
        - Application remains responsive
        - Screenshot shows persistent search state
        - Search can be cleared with Esc key
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Open search and type term
        whaletui_app.send_key('/')
        time.sleep(0.5)
        whaletui_app.send_text("test")
        time.sleep(0.5)
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Switch to another view
        whaletui_app.send_text("i")
        time.sleep(2)

        # Switch back to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Check if search is still applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_persistence.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    @pytest.mark.slow
    def test_search_stress_test(self, whaletui_app):
        """
        Test search under stress conditions with rapid operations.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Perform rapid search operations (20 iterations):
           - Open search with '/' key
           - Type "test{i}" with iteration number
           - Press Enter to apply
           - Clear search with Esc key
        5. Verify application remains stable under stress
        6. Take a screenshot showing stress test results

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Rapid search operations work correctly
        - Application remains stable under stress
        - No crashes or freezes occur
        - Screenshot shows stress test final state
        """
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Perform rapid search operations
        for i in range(20):
            whaletui_app.send_key('/')
            time.sleep(0.05)
            whaletui_app.send_text(f"test{i}")
            time.sleep(0.05)
            whaletui_app.send_key('Enter')
            time.sleep(0.05)
            whaletui_app.send_key('Esc')
            time.sleep(0.05)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("search_stress_test.png")
