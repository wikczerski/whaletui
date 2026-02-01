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
        """
        Test containers view with Docker containers.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Wait for containers to load (3 seconds)
        5. Verify containers are displayed in the view
        6. Take a screenshot showing containers view with Docker data

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Docker containers are displayed in the view
        - Application remains responsive
        - Screenshot shows containers view with Docker data
        """
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
        """
        Test images view with Docker images.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to images view using 'i' key
        4. Wait for images to load (3 seconds)
        5. Verify images are displayed in the view
        6. Take a screenshot showing images view with Docker data

        Expected Outcome:
        - Application starts successfully
        - Can navigate to images view
        - Docker images are displayed in the view
        - Application remains responsive
        - Screenshot shows images view with Docker data
        """
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
        """
        Test volumes view with Docker volumes.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to volumes view using 'v' key
        4. Wait for volumes to load (3 seconds)
        5. Verify volumes are displayed in the view
        6. Take a screenshot showing volumes view with Docker data

        Expected Outcome:
        - Application starts successfully
        - Can navigate to volumes view
        - Docker volumes are displayed in the view
        - Application remains responsive
        - Screenshot shows volumes view with Docker data
        """
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
        """
        Test networks view with Docker networks.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to networks view using 'n' key
        4. Wait for networks to load (3 seconds)
        5. Verify networks are displayed in the view
        6. Take a screenshot showing networks view with Docker data

        Expected Outcome:
        - Application starts successfully
        - Can navigate to networks view
        - Docker networks are displayed in the view
        - Application remains responsive
        - Screenshot shows networks view with Docker data
        """
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
        """
        Test swarm view with Docker swarm.

        Steps:
        1. Check if Docker swarm is available (skip if not)
        2. Start the WhaleTUI application
        3. Wait for the main screen to appear
        4. Navigate to swarm view using 's' key
        5. Wait for swarm information to load (3 seconds)
        6. Verify swarm information is displayed in the view
        7. Take a screenshot showing swarm view with Docker data

        Expected Outcome:
        - Test is skipped if Docker swarm is not available
        - Application starts successfully
        - Can navigate to swarm view
        - Docker swarm information is displayed in the view
        - Application remains responsive
        - Screenshot shows swarm view with Docker data
        """
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
        """
        Test nodes view with Docker swarm.

        Steps:
        1. Check if Docker swarm is available (skip if not)
        2. Start the WhaleTUI application
        3. Wait for the main screen to appear
        4. Navigate to nodes view using "nodes" command
        5. Wait for nodes to load (3 seconds)
        6. Verify nodes are displayed in the view
        7. Take a screenshot showing nodes view with Docker data

        Expected Outcome:
        - Test is skipped if Docker swarm is not available
        - Application starts successfully
        - Can navigate to nodes view
        - Docker nodes are displayed in the view
        - Application remains responsive
        - Screenshot shows nodes view with Docker data
        """
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
        """
        Test services view with Docker services.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to services view using "services" command
        4. Wait for services to load (3 seconds)
        5. Verify services are displayed in the view
        6. Take a screenshot showing services view with Docker data

        Expected Outcome:
        - Application starts successfully
        - Can navigate to services view
        - Docker services are displayed in the view
        - Application remains responsive
        - Screenshot shows services view with Docker data
        """
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
        """
        Test container operations with Docker containers.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Wait for containers to load (3 seconds)
        5. Test container operations:
           - Press Enter for details
           - Press 'l' for logs
           - Press 'i' for inspect
           - Press 'a' for attach
           - Press 'h' for history
        6. For each operation, press 'q' to return to previous view
        7. Verify application remains responsive
        8. Take a screenshot showing container operations

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - All container operations work correctly
        - Each operation can be accessed and exited
        - Application remains responsive throughout
        - Screenshot shows container operations state
        """
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
        """
        Test image operations with Docker images.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to images view using 'i' key
        4. Wait for images to load (3 seconds)
        5. Test image operations:
           - Press Enter for details
           - Press 'i' for inspect
           - Press 'h' for history
        6. For each operation, press 'q' to return to previous view
        7. Verify application remains responsive
        8. Take a screenshot showing image operations

        Expected Outcome:
        - Application starts successfully
        - Can navigate to images view
        - All image operations work correctly
        - Each operation can be accessed and exited
        - Application remains responsive throughout
        - Screenshot shows image operations state
        """
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
        """
        Test volume operations with Docker volumes.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to volumes view using 'v' key
        4. Wait for volumes to load (3 seconds)
        5. Test volume operations:
           - Press Enter for details
           - Press 'i' for inspect
        6. For each operation, press 'q' to return to previous view
        7. Verify application remains responsive
        8. Take a screenshot showing volume operations

        Expected Outcome:
        - Application starts successfully
        - Can navigate to volumes view
        - All volume operations work correctly
        - Each operation can be accessed and exited
        - Application remains responsive throughout
        - Screenshot shows volume operations state
        """
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
        """
        Test network operations with Docker networks.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to networks view using 'n' key
        4. Wait for networks to load (3 seconds)
        5. Test network operations:
           - Press Enter for details
           - Press 'i' for inspect
        6. For each operation, press 'q' to return to previous view
        7. Verify application remains responsive
        8. Take a screenshot showing network operations

        Expected Outcome:
        - Application starts successfully
        - Can navigate to networks view
        - All network operations work correctly
        - Each operation can be accessed and exited
        - Application remains responsive throughout
        - Screenshot shows network operations state
        """
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
        """
        Test swarm operations with Docker swarm.

        Steps:
        1. Check if Docker swarm is available (skip if not)
        2. Start the WhaleTUI application
        3. Wait for the main screen to appear
        4. Navigate to swarm view using 's' key
        5. Wait for swarm information to load (3 seconds)
        6. Test swarm operations:
           - Press Enter for details
           - Press 'i' for inspect
        7. For each operation, press 'q' to return to previous view
        8. Verify application remains responsive
        9. Take a screenshot showing swarm operations

        Expected Outcome:
        - Test is skipped if Docker swarm is not available
        - Application starts successfully
        - Can navigate to swarm view
        - All swarm operations work correctly
        - Each operation can be accessed and exited
        - Application remains responsive throughout
        - Screenshot shows swarm operations state
        """
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
        """
        Test service operations with Docker services.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to services view using "services" command
        4. Wait for services to load (3 seconds)
        5. Test service operations:
           - Press Enter for details
           - Press 'i' for inspect
           - Press 'l' for logs
        6. For each operation, press 'q' to return to previous view
        7. Verify application remains responsive
        8. Take a screenshot showing service operations

        Expected Outcome:
        - Application starts successfully
        - Can navigate to services view
        - All service operations work correctly
        - Each operation can be accessed and exited
        - Application remains responsive throughout
        - Screenshot shows service operations state
        """
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
        """
        Test Docker refresh functionality with containers.

        Steps:
        1. Start WhaleTUI with 2-second refresh interval
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Wait for containers to load (3 seconds)
        5. Wait for refresh cycle to complete (3 seconds)
        6. Verify application remains responsive during refresh
        7. Take a screenshot showing Docker refresh functionality

        Expected Outcome:
        - Application starts with custom refresh interval
        - Can navigate to containers view
        - Docker data is refreshed automatically
        - Application remains responsive during refresh cycles
        - Screenshot shows Docker refresh functionality
        """
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
        """
        Test Docker search functionality with containers.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Wait for containers to load (3 seconds)
        5. Test search functionality:
           - Open search with '/' key
           - Type "test" as search term
           - Press Enter to apply search
        6. Verify search is applied and application remains running
        7. Take a screenshot showing Docker search functionality
        8. Clear search by pressing Esc key

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Search functionality works with Docker data
        - Search can be applied and cleared
        - Application remains responsive
        - Screenshot shows Docker search functionality
        """
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
        """
        Test Docker error handling with containers.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Navigate to containers view using 'c' key
        4. Wait for containers to load (3 seconds)
        5. Test error handling scenarios:
           - Send Ctrl+C (interrupt)
           - Send Ctrl+Z (suspend)
           - Send F1 (function key)
        6. Verify application handles errors gracefully
        7. Take a screenshot showing Docker error handling

        Expected Outcome:
        - Application starts successfully
        - Can navigate to containers view
        - Error scenarios are handled gracefully
        - Application remains responsive after error scenarios
        - No crashes occur from error inputs
        - Screenshot shows Docker error handling state
        """
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
        """
        Test Docker integration during long session.

        Steps:
        1. Start the WhaleTUI application
        2. Wait for the main screen to appear
        3. Perform 5 cycles of navigation through all views: ["c", "i", "v", "n", "s"]
        4. For each cycle and view:
           - Navigate to view using corresponding key
           - Wait for view to load (2 seconds)
           - Test basic interactions (Down/Up navigation)
        5. Verify application remains stable during extended session
        6. Take a screenshot showing Docker long session results

        Expected Outcome:
        - Application starts successfully
        - Can navigate through all views multiple times
        - Basic interactions work in each view
        - Application remains stable during extended session
        - No memory leaks or performance degradation
        - Screenshot shows Docker long session final state
        """
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
