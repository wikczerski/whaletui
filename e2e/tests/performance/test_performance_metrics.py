"""
Performance tests for WhaleTUI.
"""
import pytest
import time
from e2e.whaletui_controller import WhaleTUIController
from tests.utils.test_helpers import TestHelpers


class TestPerformanceMetrics:
    """Performance tests for WhaleTUI."""

    @pytest.mark.slow
    def test_startup_time(self, whaletui_app):
        """Test application startup time."""
        start_time = time.time()

        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

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
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        views = ["c", "i", "v", "n", "s"]
        switch_times = []

        for view in views:
            start_time = time.time()

            # Navigate to view
            whaletui_app.send_text(view)
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
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        search_terms = ["test", "nginx", "redis", "postgres", "alpine"]
        search_times = []

        for term in search_terms:
            start_time = time.time()

            # Perform search
            whaletui_app.send_key('/')
            time.sleep(0.1)
            whaletui_app.send_text(term)
            time.sleep(0.1)
            whaletui_app.send_key('Enter')
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
            assert whaletui_app.wait_for_screen("Details", timeout=10)

            # Navigate to containers view
            whaletui_app.send_text("c")
            time.sleep(2)

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
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to different views to test memory usage
        views = ["c", "i", "v", "n", "s"]

        for i in range(3):  # Test 3 cycles
            for view in views:
                # Navigate to view
                whaletui_app.send_text(view)
                time.sleep(2)

                # Take screenshot
                whaletui_app.take_screenshot(f"memory_usage_cycle_{i}_{view}.png")

                # Perform search to test memory usage
                whaletui_app.send_key('/')
                time.sleep(0.1)
                whaletui_app.send_text("test")
                time.sleep(0.1)
                whaletui_app.send_key('Enter')
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
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Perform multiple operations to test CPU usage
        for i in range(10):
            # Perform search
            whaletui_app.send_key('/')
            time.sleep(0.1)
            whaletui_app.send_text(f"test{i}")
            time.sleep(0.1)
            whaletui_app.send_key('Enter')
            time.sleep(0.5)

            # Clear search
            whaletui_app.send_key('Esc')
            time.sleep(0.1)

            # Navigate to different view
            view = ["c", "i", "v"][i % 3]
            whaletui_app.send_text(view)
            time.sleep(0.5)

        # Take screenshot
        whaletui_app.take_screenshot("cpu_usage_under_load.png")

        # Application should still be running
        assert whaletui_app.is_running()

    @pytest.mark.slow
    def test_large_dataset_performance(self, whaletui_app):
        """Test performance with large datasets."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to different views to test with large datasets
        views = ["c", "i", "v", "n", "s"]

        for view in views:
            start_time = time.time()

            # Navigate to view
            whaletui_app.send_text(view)
            time.sleep(3)

            # Perform search with various terms
            search_terms = ["a", "b", "c", "d", "e", "f", "g", "h", "i", "j"]

            for term in search_terms:
                whaletui_app.send_key('/')
                time.sleep(0.1)
                whaletui_app.send_text(term)
                time.sleep(0.1)
                whaletui_app.send_key('Enter')
                time.sleep(0.5)
                whaletui_app.send_key('Esc')
                time.sleep(0.1)

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
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Perform rapid operations
        start_time = time.time()

        for i in range(20):
            # Rapid view switching
            view = ["c", "i", "v"][i % 3]
            whaletui_app.send_text(view)
            time.sleep(0.1)

            # Rapid search operations
            whaletui_app.send_key('/')
            time.sleep(0.05)
            whaletui_app.send_text(f"test{i}")
            time.sleep(0.05)
            whaletui_app.send_key('Enter')
            time.sleep(0.1)
            whaletui_app.send_key('Esc')
            time.sleep(0.05)

        total_time = time.time() - start_time

        # Take screenshot
        whaletui_app.take_screenshot("concurrent_operations.png")

        # Should handle concurrent operations reasonably well
        assert total_time < 30.0, f"Concurrent operations too slow: {total_time:.2f}s"

        # Application should still be running
        assert whaletui_app.is_running()

    @pytest.mark.slow
    def test_ui_responsiveness(self, whaletui_app):
        """Test UI responsiveness under load."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test UI responsiveness
        start_time = time.time()

        for i in range(50):
            whaletui_app.send_key('Down')
            time.sleep(0.05)
            whaletui_app.send_key('Up')
            time.sleep(0.05)

        total_time = time.time() - start_time

        # UI should be responsive
        assert total_time < 10.0, f"UI responsiveness too slow: {total_time:.2f}s"

        # Take screenshot
        whaletui_app.take_screenshot("ui_responsiveness.png")

        # Application should still be running
        assert whaletui_app.is_running()

    @pytest.mark.slow
    def test_keyboard_input_performance(self, whaletui_app):
        """Test keyboard input performance."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test keyboard input performance
        start_time = time.time()

        # Test rapid key presses
        for i in range(100):
            whaletui_app.send_key('Down')
            time.sleep(0.01)
            whaletui_app.send_key('Up')
            time.sleep(0.01)

        total_time = time.time() - start_time

        # Keyboard input should be responsive
        assert total_time < 15.0, f"Keyboard input too slow: {total_time:.2f}s"

        # Take screenshot
        whaletui_app.take_screenshot("keyboard_input_performance.png")

        # Application should still be running
        assert whaletui_app.is_running()

    @pytest.mark.slow
    def test_screen_rendering_performance(self, whaletui_app):
        """Test screen rendering performance."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test screen rendering performance
        start_time = time.time()

        # Perform operations that trigger screen updates
        for i in range(20):
            whaletui_app.send_key('Down')
            time.sleep(0.1)
            whaletui_app.send_key('Up')
            time.sleep(0.1)

            # Test search to trigger screen updates
            whaletui_app.send_key('/')
            time.sleep(0.1)
            whaletui_app.send_text(f"test{i}")
            time.sleep(0.1)
            whaletui_app.send_key('Enter')
            time.sleep(0.5)
            whaletui_app.send_key('Esc')
            time.sleep(0.1)

        total_time = time.time() - start_time

        # Screen rendering should be reasonably fast
        assert total_time < 20.0, f"Screen rendering too slow: {total_time:.2f}s"

        # Take screenshot
        whaletui_app.take_screenshot("screen_rendering_performance.png")

        # Application should still be running
        assert whaletui_app.is_running()

    @pytest.mark.slow
    def test_application_stability_performance(self, whaletui_app):
        """Test application stability under performance load."""
        whaletui_app.start(['--refresh', '1'])  # Fast refresh

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test application stability under load
        for i in range(100):
            # Perform various operations
            whaletui_app.send_key('Down')
            time.sleep(0.01)
            whaletui_app.send_key('Up')
            time.sleep(0.01)

            # Test search
            if i % 10 == 0:
                whaletui_app.send_key('/')
                time.sleep(0.01)
                whaletui_app.send_text(f"test{i}")
                time.sleep(0.01)
                whaletui_app.send_key('Enter')
                time.sleep(0.1)
                whaletui_app.send_key('Esc')
                time.sleep(0.01)

            # Test view switching
            if i % 20 == 0:
                view = ["c", "i", "v"][i % 3]
                whaletui_app.send_text(view)
                time.sleep(0.1)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("application_stability_performance.png")

    @pytest.mark.slow
    def test_memory_leak_detection(self, whaletui_app):
        """Test for memory leaks during extended use."""
        whaletui_app.start()

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Perform operations that might cause memory leaks
        for i in range(50):
            # Navigate through views
            view = ["c", "i", "v", "n", "s"][i % 5]
            whaletui_app.send_text(view)
            time.sleep(0.5)

            # Perform search operations
            whaletui_app.send_key('/')
            time.sleep(0.1)
            whaletui_app.send_text(f"test{i}")
            time.sleep(0.1)
            whaletui_app.send_key('Enter')
            time.sleep(0.5)
            whaletui_app.send_key('Esc')
            time.sleep(0.1)

            # Test keyboard navigation
            for j in range(5):
                whaletui_app.send_key('Down')
                time.sleep(0.01)
                whaletui_app.send_key('Up')
                time.sleep(0.01)

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("memory_leak_detection.png")

    @pytest.mark.slow
    def test_performance_under_stress(self, whaletui_app):
        """Test performance under stress conditions."""
        whaletui_app.start(['--refresh', '1'])  # Fast refresh

        # Wait for main screen
        assert whaletui_app.wait_for_screen("Details", timeout=10)

        # Navigate to containers view
        whaletui_app.send_text("c")
        time.sleep(2)

        # Test performance under stress
        start_time = time.time()

        # Perform rapid operations
        for i in range(200):
            # Rapid key presses
            whaletui_app.send_key('Down')
            time.sleep(0.005)
            whaletui_app.send_key('Up')
            time.sleep(0.005)

            # Rapid search operations
            if i % 10 == 0:
                whaletui_app.send_key('/')
                time.sleep(0.005)
                whaletui_app.send_text(f"test{i}")
                time.sleep(0.005)
                whaletui_app.send_key('Enter')
                time.sleep(0.05)
                whaletui_app.send_key('Esc')
                time.sleep(0.005)

            # Rapid view switching
            if i % 20 == 0:
                view = ["c", "i", "v"][i % 3]
                whaletui_app.send_text(view)
                time.sleep(0.01)

        total_time = time.time() - start_time

        # Should handle stress reasonably well
        assert total_time < 60.0, f"Performance under stress too slow: {total_time:.2f}s"

        # Application should still be running
        assert whaletui_app.is_running()

        # Take screenshot
        whaletui_app.take_screenshot("performance_under_stress.png")
