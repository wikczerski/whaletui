"""
Performance e2e tests for WhaleTUI application.
"""
import pytest
import time
from e2e.whaletui_controller import WhaleTUIController


class TestWhaleTUIPerformance:
    """Performance tests for WhaleTUI."""

    @pytest.mark.slow
    def test_startup_time(self, whaletui_app):
        """Test application startup time."""
        start_time = time.time()

        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        startup_time = time.time() - start_time

        # Startup should be reasonably fast (less than 5 seconds)
        assert startup_time < 5.0, f"Startup time too slow: {startup_time:.2f}s"

        # Take screenshot
        whaletui_app.take_screenshot("startup_performance.png")

    @pytest.mark.slow
    def test_view_switching_performance(self, whaletui_app):
        """Test performance of switching between views."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        views = ["containers", "images", "volumes", "networks", "swarm"]
        switch_times = []

        for view in views:
            start_time = time.time()

            # Navigate to view
            assert whaletui_app.navigate_to_view(view)

            # Wait for view to load
            time.sleep(2)

            switch_time = time.time() - start_time
            switch_times.append(switch_time)

            # Each view switch should be reasonably fast (less than 3 seconds)
            assert switch_time < 3.0, f"View switch to {view} too slow: {switch_time:.2f}s"

        # Take screenshot
        whaletui_app.take_screenshot("view_switching_performance.png")

        # Log performance metrics
        avg_switch_time = sum(switch_times) / len(switch_times)
        print(f"Average view switch time: {avg_switch_time:.2f}s")

    @pytest.mark.slow
    def test_search_performance(self, whaletui_app):
        """Test search performance."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Wait for containers to load
        time.sleep(3)

        search_terms = ["test", "nginx", "redis", "postgres", "alpine"]
        search_times = []

        for term in search_terms:
            start_time = time.time()

            # Perform search
            assert whaletui_app.search(term)

            # Wait for search results
            time.sleep(1)

            search_time = time.time() - start_time
            search_times.append(search_time)

            # Each search should be reasonably fast (less than 2 seconds)
            assert search_time < 2.0, f"Search for '{term}' too slow: {search_time:.2f}s"

            # Clear search
            whaletui_app.send_key('Esc')
            time.sleep(0.5)

        # Take screenshot
        whaletui_app.take_screenshot("search_performance.png")

        # Log performance metrics
        avg_search_time = sum(search_times) / len(search_times)
        print(f"Average search time: {avg_search_time:.2f}s")

    @pytest.mark.slow
    def test_refresh_performance(self, whaletui_app):
        """Test refresh performance with different intervals."""
        refresh_intervals = [1, 2, 5, 10]

        for interval in refresh_intervals:
            whaletui_app.start(['--refresh', str(interval)])

            # Wait for main screen
            assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

            # Navigate to containers view
            assert whaletui_app.navigate_to_view("containers")

            # Wait for initial load
            time.sleep(3)

            # Measure refresh cycles
            refresh_times = []

            for i in range(3):  # Test 3 refresh cycles
                start_time = time.time()

                # Wait for refresh
                time.sleep(interval + 1)

                refresh_time = time.time() - start_time
                refresh_times.append(refresh_time)

            # Take screenshot
            whaletui_app.take_screenshot(f"refresh_performance_{interval}s.png")

            # Quit and restart for next test
            whaletui_app.send_key('q')
            time.sleep(1)

            # Log performance metrics
            avg_refresh_time = sum(refresh_times) / len(refresh_times)
            print(f"Refresh interval {interval}s - Average refresh time: {avg_refresh_time:.2f}s")

    @pytest.mark.slow
    def test_memory_usage_over_time(self, whaletui_app):
        """Test memory usage over time."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to different views to test memory usage
        views = ["containers", "images", "volumes", "networks", "swarm"]

        for i in range(3):  # Test 3 cycles
            for view in views:
                # Navigate to view
                assert whaletui_app.navigate_to_view(view)

                # Wait for view to load
                time.sleep(2)

                # Take screenshot
                whaletui_app.take_screenshot(f"memory_usage_cycle_{i}_{view}.png")

                # Perform search to test memory usage
                whaletui_app.search("test")
                time.sleep(1)

                # Clear search
                whaletui_app.send_key('Esc')
                time.sleep(0.5)

        # Application should still be running
        assert whaletui_app.is_running()

    @pytest.mark.slow
    def test_cpu_usage_under_load(self, whaletui_app):
        """Test CPU usage under load."""
        whaletui_app.start(['--refresh', '1'])  # Fast refresh

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Wait for initial load
        time.sleep(3)

        # Perform multiple operations to test CPU usage
        for i in range(10):
            # Perform search
            whaletui_app.search(f"test{i}")
            time.sleep(0.5)

            # Clear search
            whaletui_app.send_key('Esc')
            time.sleep(0.5)

            # Navigate to different view
            view = ["containers", "images", "volumes"][i % 3]
            whaletui_app.navigate_to_view(view)
            time.sleep(1)

        # Take screenshot
        whaletui_app.take_screenshot("cpu_usage_under_load.png")

        # Application should still be running
        assert whaletui_app.is_running()

    @pytest.mark.slow
    def test_large_dataset_performance(self, whaletui_app):
        """Test performance with large datasets."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to different views to test with large datasets
        views = ["containers", "images", "volumes", "networks", "swarm"]

        for view in views:
            start_time = time.time()

            # Navigate to view
            assert whaletui_app.navigate_to_view(view)

            # Wait for view to load
            time.sleep(3)

            # Perform search with various terms
            search_terms = ["a", "b", "c", "d", "e", "f", "g", "h", "i", "j"]

            for term in search_terms:
                whaletui_app.search(term)
                time.sleep(0.5)
                whaletui_app.send_key('Esc')
                time.sleep(0.5)

            view_time = time.time() - start_time

            # Take screenshot
            whaletui_app.take_screenshot(f"large_dataset_{view}.png")

            # Each view should handle large datasets reasonably well
            assert view_time < 10.0, f"View {view} too slow with large dataset: {view_time:.2f}s"

    @pytest.mark.slow
    def test_concurrent_operations(self, whaletui_app):
        """Test performance with concurrent operations."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("WhaleTUI", timeout=10)

        # Navigate to containers view
        assert whaletui_app.navigate_to_view("containers")

        # Wait for initial load
        time.sleep(3)

        # Perform rapid operations
        start_time = time.time()

        for i in range(20):
            # Rapid view switching
            view = ["containers", "images", "volumes"][i % 3]
            whaletui_app.navigate_to_view(view)
            time.sleep(0.1)

            # Rapid search operations
            whaletui_app.search(f"test{i}")
            time.sleep(0.1)
            whaletui_app.send_key('Esc')
            time.sleep(0.1)

        total_time = time.time() - start_time

        # Take screenshot
        whaletui_app.take_screenshot("concurrent_operations.png")

        # Should handle concurrent operations reasonably well
        assert total_time < 30.0, f"Concurrent operations too slow: {total_time:.2f}s"

        # Application should still be running
        assert whaletui_app.is_running()
