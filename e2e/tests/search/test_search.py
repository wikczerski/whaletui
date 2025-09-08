"""
Search functionality e2e tests for WhaleTUI application.
"""
import pytest
import time
from e2e.whaletui_controller import WhaleTUIController


class TestWhaleTUISearch:
    """Search functionality tests for WhaleTUI."""

    @pytest.mark.docker
    def test_search_in_containers(self, whaletui_app, docker_test_environment):
        """Test search functionality in containers view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Wait for containers to load
        time.sleep(3)

        # Perform search
        assert whaletui_app.search("test")

        # Wait for search results
        time.sleep(2)

        # Take screenshot
        whaletui_app.take_screenshot("search_containers.png")

        # Should show search results or no results message
        output = whaletui_app.get_screen_content()
        assert "search" in output.lower() or "no results" in output.lower()

    @pytest.mark.docker
    def test_search_in_images(self, whaletui_app, docker_test_environment):
        """Test search functionality in images view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to images view
        assert whaletui_app.navigate_to_view("images")

        # Wait for images to load
        time.sleep(3)

        # Perform search
        assert whaletui_app.search("nginx")

        # Wait for search results
        time.sleep(2)

        # Take screenshot
        whaletui_app.take_screenshot("search_images.png")

        # Should show search results or no results message
        output = whaletui_app.get_screen_content()
        assert "search" in output.lower() or "no results" in output.lower()

    @pytest.mark.docker
    def test_search_in_volumes(self, whaletui_app, docker_test_environment):
        """Test search functionality in volumes view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to volumes view
        assert whaletui_app.navigate_to_view("volumes")

        # Wait for volumes to load
        time.sleep(3)

        # Perform search
        assert whaletui_app.search("test")

        # Wait for search results
        time.sleep(2)

        # Take screenshot
        whaletui_app.take_screenshot("search_volumes.png")

        # Should show search results or no results message
        output = whaletui_app.get_screen_content()
        assert "search" in output.lower() or "no results" in output.lower()

    @pytest.mark.docker
    def test_search_in_networks(self, whaletui_app, docker_test_environment):
        """Test search functionality in networks view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to networks view
        assert whaletui_app.navigate_to_view("networks")

        # Wait for networks to load
        time.sleep(3)

        # Perform search
        assert whaletui_app.search("bridge")

        # Wait for search results
        time.sleep(2)

        # Take screenshot
        whaletui_app.take_screenshot("search_networks.png")

        # Should show search results or no results message
        output = whaletui_app.get_screen_content()
        assert "search" in output.lower() or "no results" in output.lower()

    @pytest.mark.docker
    def test_search_in_swarm_services(self, whaletui_app, docker_test_environment):
        """Test search functionality in swarm services view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to services view
        assert whaletui_app.navigate_to_view("services")

        # Wait for services to load
        time.sleep(3)

        # Perform search
        assert whaletui_app.search("test")

        # Wait for search results
        time.sleep(2)

        # Take screenshot
        whaletui_app.take_screenshot("search_services.png")

        # Should show search results or no results message
        output = whaletui_app.get_screen_content()
        assert "search" in output.lower() or "no results" in output.lower()

    @pytest.mark.docker
    def test_search_in_swarm_nodes(self, whaletui_app, docker_test_environment):
        """Test search functionality in swarm nodes view."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to nodes view
        assert whaletui_app.navigate_to_view("nodes")

        # Wait for nodes to load
        time.sleep(3)

        # Perform search
        assert whaletui_app.search("manager")

        # Wait for search results
        time.sleep(2)

        # Take screenshot
        whaletui_app.take_screenshot("search_nodes.png")

        # Should show search results or no results message
        output = whaletui_app.get_screen_content()
        assert "search" in output.lower() or "no results" in output.lower()

    def test_search_clear(self, whaletui_app):
        """Test clearing search functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Wait for containers to load
        time.sleep(3)

        # Perform search
        assert whaletui_app.search("test")

        # Wait for search results
        time.sleep(2)

        # Clear search (press Esc)
        whaletui_app.send_key('Esc')
        time.sleep(1)

        # Take screenshot
        whaletui_app.take_screenshot("search_cleared.png")

        # Should show normal view without search
        output = whaletui_app.get_screen_content()
        assert "search" not in output.lower() or "no search" in output.lower()

    def test_search_empty_term(self, whaletui_app):
        """Test search with empty term."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Wait for containers to load
        time.sleep(3)

        # Perform empty search
        assert whaletui_app.search("")

        # Wait for search results
        time.sleep(2)

        # Take screenshot
        whaletui_app.take_screenshot("search_empty.png")

        # Should handle empty search gracefully
        assert whaletui_app.is_running()

    @pytest.mark.docker
    def test_search_special_characters(self, whaletui_app, docker_test_environment):
        """Test search with special characters."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Wait for containers to load
        time.sleep(3)

        # Perform search with special characters
        assert whaletui_app.search("test-123_456")

        # Wait for search results
        time.sleep(2)

        # Take screenshot
        whaletui_app.take_screenshot("search_special_chars.png")

        # Should handle special characters gracefully
        assert whaletui_app.is_running()
