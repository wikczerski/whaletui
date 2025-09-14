#!/usr/bin/env python3
"""
Volume-based test runner for WhaleTUI e2e tests.
This script runs tests using Docker volumes instead of rebuilding the image.
"""
import subprocess
import sys
import logging
import time
import os
from pathlib import Path

# Set up logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class VolumeTestRunner:
    """Test runner that uses Docker volumes for live code changes."""

    def __init__(self):
        self.project_root = Path(__file__).parent.parent.parent
        self.e2e_root = Path(__file__).parent.parent
        self.docker_dir = Path(__file__).parent

    def build_base_image(self):
        """Build the base Docker image (only once)."""
        logger.info("🔨 Building base Docker image...")
        cmd = [
            "docker", "build",
            "-f", str(self.docker_dir / "Dockerfile.test"),
            "-t", "whaletui-e2e-test-isolated",
            str(self.project_root),
            "--no-cache"
        ]

        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode != 0:
            logger.error(f"❌ Failed to build Docker image: {result.stderr}")
            return False

        logger.info("✅ Base Docker image built successfully")
        return True

    def run_tests_with_volumes(self, test_category=None, verbose=False, generate_data=False):
        """Run tests using Docker volumes for live code changes."""
        logger.info("🚀 Running tests with volume mounts...")

        # Build the base image first
        if not self.build_base_image():
            return False

        # Prepare docker-compose command
        cmd = [
            "docker-compose",
            "-f", str(self.docker_dir / "docker-compose.volume.yml"),
            "run", "--rm", "whaletui-test"
        ]

        # Add test command
        test_cmd = ["python3", "-m", "pytest"]

        if verbose:
            test_cmd.append("-v")

        if test_category:
            test_cmd.extend([f"/app/e2e/tests/{test_category}/", "-v"])
        else:
            test_cmd.extend([
                "/app/e2e/tests/ui/",
                "/app/e2e/tests/unit/",
                "/app/e2e/tests/error_handling/",
                "/app/e2e/tests/performance/",
                "/app/e2e/tests/search/",
                "/app/e2e/tests/fixtures/",
                "/app/e2e/tests/utils/",
                "-v"
            ])

        test_cmd.extend([
            "--html=/app/reports/test_report.html",
            "--self-contained-html"
        ])

        cmd.extend(test_cmd)

        logger.info(f"Running command: {' '.join(cmd)}")

        # Run the tests
        result = subprocess.run(cmd, cwd=str(self.docker_dir))
        return result.returncode == 0

    def run_single_test(self, test_path, verbose=False):
        """Run a single test file."""
        logger.info(f"🧪 Running single test: {test_path}")

        # Build the base image first
        if not self.build_base_image():
            return False

        cmd = [
            "docker-compose",
            "-f", str(self.docker_dir / "docker-compose.volume.yml"),
            "run", "--rm", "whaletui-test",
            "python3", "-m", "pytest", f"/app/e2e/tests/{test_path}"
        ]

        if verbose:
            cmd.append("-v")

        logger.info(f"Running command: {' '.join(cmd)}")

        result = subprocess.run(cmd, cwd=str(self.docker_dir))
        return result.returncode == 0

    def run_interactive_shell(self):
        """Run an interactive shell in the test container."""
        logger.info("🐚 Starting interactive shell...")

        # Build the base image first
        if not self.build_base_image():
            return False

        cmd = [
            "docker-compose",
            "-f", str(self.docker_dir / "docker-compose.volume.yml"),
            "run", "--rm", "whaletui-test",
            "bash"
        ]

        logger.info(f"Running command: {' '.join(cmd)}")
        subprocess.run(cmd, cwd=str(self.docker_dir))

    def generate_test_data(self):
        """Generate test data in the isolated environment."""
        logger.info("📊 Generating test data...")

        # Build the base image first
        if not self.build_base_image():
            return False

        cmd = [
            "docker-compose",
            "-f", str(self.docker_dir / "docker-compose.volume.yml"),
            "run", "--rm", "whaletui-test",
            "python3", "/app/docker/test_data_generator.py"
        ]

        logger.info(f"Running command: {' '.join(cmd)}")
        result = subprocess.run(cmd, cwd=str(self.docker_dir))
        return result.returncode == 0

    def cleanup(self):
        """Clean up Docker resources."""
        logger.info("🧹 Cleaning up Docker resources...")

        # Stop and remove containers
        subprocess.run([
            "docker-compose",
            "-f", str(self.docker_dir / "docker-compose.volume.yml"),
            "down"
        ], cwd=str(self.docker_dir))

        # Remove unused images
        subprocess.run(["docker", "image", "prune", "-f"])

        logger.info("✅ Cleanup completed")

def main():
    """Main function."""
    import argparse

    parser = argparse.ArgumentParser(description="Volume-based WhaleTUI e2e test runner")
    parser.add_argument("--category", help="Test category to run (ui, unit, etc.)")
    parser.add_argument("--test", help="Single test file to run")
    parser.add_argument("--verbose", "-v", action="store_true", help="Verbose output")
    parser.add_argument("--generate-data", action="store_true", help="Generate test data")
    parser.add_argument("--shell", action="store_true", help="Run interactive shell")
    parser.add_argument("--cleanup", action="store_true", help="Clean up Docker resources")

    args = parser.parse_args()

    runner = VolumeTestRunner()

    try:
        if args.cleanup:
            runner.cleanup()
            return 0

        if args.shell:
            runner.run_interactive_shell()
            return 0

        if args.generate_data:
            success = runner.generate_test_data()
            return 0 if success else 1

        if args.test:
            success = runner.run_single_test(args.test, args.verbose)
        else:
            success = runner.run_tests_with_volumes(args.category, args.verbose, args.generate_data)

        if success:
            logger.info("🎉 All tests passed!")
            return 0
        else:
            logger.error("❌ Some tests failed!")
            return 1

    except KeyboardInterrupt:
        logger.info("⏹️  Test run interrupted by user")
        return 1
    except Exception as e:
        logger.error(f"💥 Unexpected error: {e}")
        return 1

if __name__ == "__main__":
    sys.exit(main())
