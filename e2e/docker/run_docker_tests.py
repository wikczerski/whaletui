#!/usr/bin/env python3
"""
Docker-based test runner for WhaleTUI e2e tests.
Runs all tests inside Docker with Docker-in-Docker (DinD) and generates 200+ volumes for stress testing.
"""
import os
import sys
import subprocess
import argparse
import time
import logging
from pathlib import Path

# Set up logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class DockerTestRunner:
    """Docker-based test runner for WhaleTUI e2e tests."""

    def __init__(self, project_root: str = None):
        self.project_root = project_root or Path(__file__).parent.parent.parent
        self.docker_dir = Path(__file__).parent
        self.e2e_dir = self.docker_dir.parent

    def check_docker_available(self) -> bool:
        """Check if Docker is available."""
        try:
            result = subprocess.run(
                ["docker", "version", "--format", "{{.Server.Version}}"],
                capture_output=True,
                text=True,
                timeout=10
            )
            return result.returncode == 0
        except (subprocess.TimeoutExpired, FileNotFoundError):
            return False

    def build_test_image(self) -> bool:
        """Build the test Docker image."""
        logger.info("Building test Docker image...")

        try:
            result = subprocess.run([
                "docker", "build",
                "-f", str(self.docker_dir / "Dockerfile.test"),
                "-t", "whaletui-e2e-test",
                str(self.project_root)
            ], capture_output=True, text=True)

            if result.returncode == 0:
                logger.info("✅ Test image built successfully")
                return True
            else:
                logger.error(f"❌ Failed to build test image: {result.stderr}")
                return False

        except Exception as e:
            logger.error(f"❌ Error building test image: {e}")
            return False

    def run_tests_in_docker(self, test_category: str = None, verbose: bool = False) -> bool:
        """Run tests inside Docker container."""
        logger.info("Running tests inside Docker container...")

        # Build command
        cmd = [
            "docker", "run", "--rm",
            "--privileged",
            "-v", f"{self.project_root}:/app/whaletui",
            "-v", f"{self.e2e_dir}:/app/e2e",
            "-v", "whaletui-test-data:/app/test-data",
            "-v", "whaletui-screenshots:/app/screenshots",
            "-v", "whaletui-reports:/app/reports",
            "-e", "DOCKER_HOST=unix:///var/run/docker.sock",
            "-e", "PYTHONPATH=/app/e2e",
            "-e", "DOCKER_BUILDKIT=1",
            "whaletui-e2e-test",
            "/app/docker/run_tests_in_docker.sh"
        ]

        if test_category:
            cmd.extend(["--category", test_category])

        if verbose:
            cmd.extend(["--verbose"])

        try:
            logger.info(f"Running command: {' '.join(cmd)}")
            result = subprocess.run(cmd, text=True)

            if result.returncode == 0:
                logger.info("✅ Tests completed successfully")
                return True
            else:
                logger.error(f"❌ Tests failed with return code: {result.returncode}")
                return False

        except Exception as e:
            logger.error(f"❌ Error running tests: {e}")
            return False

    def run_tests_with_compose(self, test_category: str = None, verbose: bool = False) -> bool:
        """Run tests using Docker Compose."""
        logger.info("Running tests with Docker Compose...")

        # Change to docker directory
        os.chdir(self.docker_dir)

        # Build and run with compose
        try:
            # Build
            logger.info("Building with Docker Compose...")
            build_result = subprocess.run([
                "docker-compose", "-f", "docker-compose.test.yml", "build"
            ], capture_output=True, text=True)

            if build_result.returncode != 0:
                logger.error(f"❌ Failed to build with compose: {build_result.stderr}")
                return False

            # Run
            logger.info("Running tests with Docker Compose...")
            run_result = subprocess.run([
                "docker-compose", "-f", "docker-compose.test.yml", "up", "--abort-on-container-exit"
            ], text=True)

            if run_result.returncode == 0:
                logger.info("✅ Tests completed successfully")
                return True
            else:
                logger.error(f"❌ Tests failed with return code: {run_result.returncode}")
                return False

        except Exception as e:
            logger.error(f"❌ Error running tests with compose: {e}")
            return False
        finally:
            # Cleanup
            logger.info("Cleaning up Docker Compose...")
            subprocess.run([
                "docker-compose", "-f", "docker-compose.test.yml", "down", "-v"
            ], capture_output=True)

    def generate_test_data_only(self) -> bool:
        """Generate test data only."""
        logger.info("Generating test data only...")

        try:
            result = subprocess.run([
                "docker", "run", "--rm",
                "--privileged",
                "-v", f"{self.project_root}:/app/whaletui",
                "-v", f"{self.e2e_dir}:/app/e2e",
                "-v", "whaletui-test-data:/app/test-data",
                "-e", "DOCKER_HOST=unix:///var/run/docker.sock",
                "-e", "PYTHONPATH=/app/e2e",
                "whaletui-e2e-test",
                "python3", "/app/docker/test_data_generator.py"
            ], text=True)

            if result.returncode == 0:
                logger.info("✅ Test data generated successfully")
                return True
            else:
                logger.error(f"❌ Failed to generate test data: {result.returncode}")
                return False

        except Exception as e:
            logger.error(f"❌ Error generating test data: {e}")
            return False

    def cleanup_docker_resources(self):
        """Clean up Docker resources."""
        logger.info("Cleaning up Docker resources...")

        try:
            # Remove test containers
            subprocess.run([
                "docker", "container", "prune", "-f"
            ], capture_output=True)

            # Remove test images
            subprocess.run([
                "docker", "image", "prune", "-f"
            ], capture_output=True)

            # Remove test volumes
            subprocess.run([
                "docker", "volume", "prune", "-f"
            ], capture_output=True)

            # Remove test networks
            subprocess.run([
                "docker", "network", "prune", "-f"
            ], capture_output=True)

            logger.info("✅ Docker resources cleaned up")

        except Exception as e:
            logger.error(f"❌ Error cleaning up Docker resources: {e}")

    def show_test_reports(self):
        """Show test reports."""
        logger.info("Test reports available in Docker volumes:")
        logger.info("  - whaletui-reports: Contains HTML test reports")
        logger.info("  - whaletui-screenshots: Contains test screenshots")
        logger.info("  - whaletui-test-data: Contains generated test data")

        # Try to show reports from volume
        try:
            result = subprocess.run([
                "docker", "run", "--rm",
                "-v", "whaletui-reports:/app/reports",
                "alpine:latest",
                "ls", "-la", "/app/reports"
            ], capture_output=True, text=True)

            if result.returncode == 0:
                logger.info("Available reports:")
                logger.info(result.stdout)

        except Exception as e:
            logger.warning(f"Could not list reports: {e}")

def main():
    """Main function."""
    parser = argparse.ArgumentParser(description="Run WhaleTUI e2e tests in Docker")
    parser.add_argument(
        "--category", "-c",
        choices=["unit", "integration", "ui", "performance", "docker", "search", "error_handling"],
        help="Run tests for a specific category"
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Verbose output"
    )
    parser.add_argument(
        "--compose",
        action="store_true",
        help="Use Docker Compose instead of direct Docker run"
    )
    parser.add_argument(
        "--generate-data-only",
        action="store_true",
        help="Generate test data only (200+ volumes)"
    )
    parser.add_argument(
        "--cleanup",
        action="store_true",
        help="Clean up Docker resources"
    )
    parser.add_argument(
        "--show-reports",
        action="store_true",
        help="Show available test reports"
    )

    args = parser.parse_args()

    runner = DockerTestRunner()

    # Check Docker availability
    if not runner.check_docker_available():
        logger.error("❌ Docker is not available or not running")
        sys.exit(1)

    # Handle different commands
    if args.cleanup:
        runner.cleanup_docker_resources()
        return

    if args.show_reports:
        runner.show_test_reports()
        return

    if args.generate_data_only:
        success = runner.generate_test_data_only()
        sys.exit(0 if success else 1)

    # Build test image
    if not runner.build_test_image():
        sys.exit(1)

    # Run tests
    if args.compose:
        success = runner.run_tests_with_compose(args.category, args.verbose)
    else:
        success = runner.run_tests_in_docker(args.category, args.verbose)

    if success:
        logger.info("🎉 All tests completed successfully!")
        runner.show_test_reports()
    else:
        logger.error("❌ Some tests failed!")
        sys.exit(1)

if __name__ == "__main__":
    main()
