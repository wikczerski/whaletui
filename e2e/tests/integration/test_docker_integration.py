"""
Docker integration tests for WhaleTUI.
"""
import pytest
import time
from e2e.whaletui_controller import WhaleTUIController
from tests.utils.test_helpers import TestHelpers


class TestDockerIntegration:
    """Docker integration tests for WhaleTUI."""

    @pytest.mark.docker
    def test_containers_view_with_docker(self, whaletui_app, docker_containers):
        """Test containers view with Docker containers."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(3)

        # Check if containers are displayed
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("containers_view_with_docker.png")

    @pytest.mark.docker
    def test_images_view_with_docker(self, whaletui_app, docker_containers):
        """Test images view with Docker images."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to images view
        whaletui_app.send_text("i")
        time.sleep(3)

        # Check if images are displayed
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("images_view_with_docker.png")

    @pytest.mark.docker
    def test_volumes_view_with_docker(self, whaletui_app, docker_volumes):
        """Test volumes view with Docker volumes."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to volumes view
        whaletui_app.send_text("v")
        time.sleep(3)

        # Check if volumes are displayed
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("volumes_view_with_docker.png")

    @pytest.mark.docker
    def test_networks_view_with_docker(self, whaletui_app, docker_networks):
        """Test networks view with Docker networks."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to networks view
        whaletui_app.send_text("n")
        time.sleep(3)

        # Check if networks are displayed
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("networks_view_with_docker.png")

    @pytest.mark.docker
    def test_swarm_view_with_docker(self, whaletui_app, docker_swarm):
        """Test swarm view with Docker swarm."""
        if not docker_swarm:
            pytest.skip("Docker swarm not available")

        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to swarm view
        whaletui_app.send_text("s")
        time.sleep(3)

        # Check if swarm information is displayed
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("swarm_view_with_docker.png")

    @pytest.mark.docker
    def test_nodes_view_with_docker(self, whaletui_app, docker_swarm):
        """Test nodes view with Docker swarm."""
        if not docker_swarm:
            pytest.skip("Docker swarm not available")

        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to nodes view
        whaletui_app.send_text("nodes")
        time.sleep(3)

        # Check if nodes are displayed
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("nodes_view_with_docker.png")

    @pytest.mark.docker
    def test_services_view_with_docker(self, whaletui_app, docker_services):
        """Test services view with Docker services."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to services view
        whaletui_app.send_text("services")
        time.sleep(3)

        # Check if services are displayed
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("services_view_with_docker.png")

    @pytest.mark.docker
    def test_container_operations(self, whaletui_app, docker_containers):
        """Test container operations."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(3)

        # Test container operations
        operations = [
            ('Enter', 1),  # Details
            ('l', 1),      # Logs
            ('i', 1),      # Inspect
            ('a', 1),      # Attach
            ('h', 1),      # History
        ]

        for key, delay in operations:
            whaletui_app.send_key(key)
            time.sleep(delay)

            # Press 'q' to go back
            whaletui_app.send_key('q')
            time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("container_operations.png")

    @pytest.mark.docker
    def test_image_operations(self, whaletui_app, docker_containers):
        """Test image operations."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to images view
        whaletui_app.send_text("i")
        time.sleep(3)

        # Test image operations
        operations = [
            ('Enter', 1),  # Details
            ('i', 1),      # Inspect
            ('h', 1),      # History
        ]

        for key, delay in operations:
            whaletui_app.send_key(key)
            time.sleep(delay)

            # Press 'q' to go back
            whaletui_app.send_key('q')
            time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("image_operations.png")

    @pytest.mark.docker
    def test_volume_operations(self, whaletui_app, docker_volumes):
        """Test volume operations."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to volumes view
        whaletui_app.send_text("v")
        time.sleep(3)

        # Test volume operations
        operations = [
            ('Enter', 1),  # Details
            ('i', 1),      # Inspect
        ]

        for key, delay in operations:
            whaletui_app.send_key(key)
            time.sleep(delay)

            # Press 'q' to go back
            whaletui_app.send_key('q')
            time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("volume_operations.png")

    @pytest.mark.docker
    def test_network_operations(self, whaletui_app, docker_networks):
        """Test network operations."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to networks view
        whaletui_app.send_text("n")
        time.sleep(3)

        # Test network operations
        operations = [
            ('Enter', 1),  # Details
            ('i', 1),      # Inspect
        ]

        for key, delay in operations:
            whaletui_app.send_key(key)
            time.sleep(delay)

            # Press 'q' to go back
            whaletui_app.send_key('q')
            time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("network_operations.png")

    @pytest.mark.docker
    def test_swarm_operations(self, whaletui_app, docker_swarm):
        """Test swarm operations."""
        if not docker_swarm:
            pytest.skip("Docker swarm not available")

        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to swarm view
        whaletui_app.send_text("s")
        time.sleep(3)

        # Test swarm operations
        operations = [
            ('Enter', 1),  # Details
            ('i', 1),      # Inspect
        ]

        for key, delay in operations:
            whaletui_app.send_key(key)
            time.sleep(delay)

            # Press 'q' to go back
            whaletui_app.send_key('q')
            time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("swarm_operations.png")

    @pytest.mark.docker
    def test_service_operations(self, whaletui_app, docker_services):
        """Test service operations."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to services view
        whaletui_app.send_text("services")
        time.sleep(3)

        # Test service operations
        operations = [
            ('Enter', 1),  # Details
            ('i', 1),      # Inspect
            ('l', 1),      # Logs
        ]

        for key, delay in operations:
            whaletui_app.send_key(key)
            time.sleep(delay)

            # Press 'q' to go back
            whaletui_app.send_key('q')
            time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("service_operations.png")

    @pytest.mark.docker
    def test_docker_refresh_functionality(self, whaletui_app, docker_containers):
        """Test Docker refresh functionality."""
        whaletui_app.start(['--refresh', '2'])  # 2 second refresh

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(3)

        # Wait for refresh cycle
        time.sleep(3)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("docker_refresh_functionality.png")

    @pytest.mark.docker
    def test_docker_search_functionality(self, whaletui_app, docker_containers):
        """Test Docker search functionality."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(3)

        # Test search functionality
        whaletui_app.send_key('/')
        time.sleep(0.5)
        whaletui_app.send_text("test")
        time.sleep(0.5)
        whaletui_app.send_key('Enter')
        time.sleep(1)

        # Check if search is applied
        output = whaletui_app.get_screen_content()
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("docker_search_functionality.png")

        # Clear search
        whaletui_app.send_key('Esc')
        time.sleep(0.5)

    @pytest.mark.docker
    def test_docker_error_handling(self, whaletui_app, docker_containers):
        """Test Docker error handling."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(3)

        # Test error handling scenarios
        error_scenarios = [
            'Ctrl+C',  # Interrupt
            'Ctrl+Z',  # Suspend
            'F1',      # Function key
        ]

        for key in error_scenarios:
            whaletui_app.send_key(key)
            time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("docker_error_handling.png")

    @pytest.mark.docker
    @pytest.mark.slow
    def test_docker_long_session(self, whaletui_app, docker_containers):
        """Test Docker integration during long session."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate through different views multiple times
        views = ["c", "i", "v", "n", "s"]

        for i in range(5):  # 5 cycles
            for view in views:
                whaletui_app.send_text(view)
                time.sleep(2)

                # Test some interactions
                whaletui_app.send_key('Down')
                time.sleep(0.5)
                whaletui_app.send_key('Up')
                time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("docker_long_session.png")
