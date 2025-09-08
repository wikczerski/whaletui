"""
Docker integration e2e tests for WhaleTUI application.
"""
import pytest
import time
from e2e.whaletui_controller import WhaleTUIController


class TestWhaleTUIDockerIntegration:
    """Docker integration tests for WhaleTUI."""

    @pytest.mark.docker
    def test_containers_view(self, whaletui_app, docker_test_environment):
        """Test the containers view functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Wait for containers to load
        time.sleep(3)

        # Take screenshot
        whaletui_app.take_screenshot("containers_view.png")

        # Should show some containers or empty state
        output = whaletui_app.get_screen_content()
        assert "containers" in output.lower() or "no containers" in output.lower()

    @pytest.mark.docker
    def test_images_view(self, whaletui_app, docker_test_environment):
        """Test the images view functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to images view
        assert whaletui_app.navigate_to_view("images")

        # Wait for images to load
        time.sleep(3)

        # Take screenshot
        whaletui_app.take_screenshot("images_view.png")

        # Should show some images or empty state
        output = whaletui_app.get_screen_content()
        assert "images" in output.lower() or "no images" in output.lower()

    @pytest.mark.docker
    def test_volumes_view(self, whaletui_app, docker_test_environment):
        """Test the volumes view functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to volumes view
        assert whaletui_app.navigate_to_view("volumes")

        # Wait for volumes to load
        time.sleep(3)

        # Take screenshot
        whaletui_app.take_screenshot("volumes_view.png")

        # Should show some volumes or empty state
        output = whaletui_app.get_screen_content()
        assert "volumes" in output.lower() or "no volumes" in output.lower()

    @pytest.mark.docker
    def test_networks_view(self, whaletui_app, docker_test_environment):
        """Test the networks view functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to networks view
        assert whaletui_app.navigate_to_view("networks")

        # Wait for networks to load
        time.sleep(3)

        # Take screenshot
        whaletui_app.take_screenshot("networks_view.png")

        # Should show some networks or empty state
        output = whaletui_app.get_screen_content()
        assert "networks" in output.lower() or "no networks" in output.lower()

    @pytest.mark.docker
    def test_swarm_view(self, whaletui_app, docker_test_environment):
        """Test the swarm view functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to swarm view
        assert whaletui_app.navigate_to_view("swarm")

        # Wait for swarm to load
        time.sleep(3)

        # Take screenshot
        whaletui_app.take_screenshot("swarm_view.png")

        # Should show swarm information or empty state
        output = whaletui_app.get_screen_content()
        assert "swarm" in output.lower() or "no swarm" in output.lower()

    @pytest.mark.docker
    def test_nodes_view(self, whaletui_app, docker_test_environment):
        """Test the nodes view functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to nodes view
        assert whaletui_app.navigate_to_view("nodes")

        # Wait for nodes to load
        time.sleep(3)

        # Take screenshot
        whaletui_app.take_screenshot("nodes_view.png")

        # Should show nodes information or empty state
        output = whaletui_app.get_screen_content()
        assert "nodes" in output.lower() or "no nodes" in output.lower()

    @pytest.mark.docker
    def test_services_view(self, whaletui_app, docker_test_environment):
        """Test the services view functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to services view
        assert whaletui_app.navigate_to_view("services")

        # Wait for services to load
        time.sleep(3)

        # Take screenshot
        whaletui_app.take_screenshot("services_view.png")

        # Should show services information or empty state
        output = whaletui_app.get_screen_content()
        assert "services" in output.lower() or "no services" in output.lower()

    @pytest.mark.docker
    def test_view_navigation(self, whaletui_app, docker_test_environment):
        """Test navigation between different views."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Test navigation to different views
        views = ["containers", "images", "volumes", "networks", "swarm"]

        for view in views:
            assert whaletui_app.navigate_to_view(view)
            time.sleep(2)  # Wait for view to load

            # Take screenshot for each view
            whaletui_app.take_screenshot(f"navigation_{view}.png")

            # Verify we're in the correct view
            output = whaletui_app.get_screen_content()
            assert view in output.lower() or f"no {view}" in output.lower()

    @pytest.mark.docker
    @pytest.mark.slow
    def test_refresh_functionality(self, whaletui_app, docker_test_environment):
        """Test the refresh functionality."""
        whaletui_app.start(['--refresh', '2'])  # 2 second refresh

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Wait for initial load
        time.sleep(3)

        # Take initial screenshot
        whaletui_app.take_screenshot("refresh_initial.png")

        # Wait for refresh cycle
        time.sleep(3)

        # Take screenshot after refresh
        whaletui_app.take_screenshot("refresh_after.png")

        # Application should still be running
        assert whaletui_app.is_running()
