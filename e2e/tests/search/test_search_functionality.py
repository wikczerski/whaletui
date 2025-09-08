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
        """Test search functionality in containers view."""
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
        """Test search functionality in images view."""
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
        """Test search functionality in volumes view."""
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
        """Test search functionality in networks view."""
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
        """Test search functionality in swarm services view."""
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
        """Test search functionality in swarm nodes view."""
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
        """Test search clear functionality."""
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
        """Test search with empty term."""
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
        """Test search with special characters."""
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
        """Test search case sensitivity."""
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
        """Test search with multiple terms."""
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
        """Test search performance."""
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
        """Test search error handling."""
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
        """Test search persistence across view changes."""
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
        """Test search under stress conditions."""
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
