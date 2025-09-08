#!/usr/bin/env python3
"""
Isolated Docker-based test runner for WhaleTUI e2e tests.
Runs all tests inside Docker with NO host volume mounts - everything stays in the container.
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

class IsolatedDockerTestRunner:
    """Isolated Docker-based test runner - no host volume mounts."""

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
        """Build the test Docker image with all code and dependencies."""
        logger.info("Building isolated test Docker image...")

        try:
            result = subprocess.run([
                "docker", "build",
                "-f", str(self.docker_dir / "Dockerfile.test"),
                "-t", "whaletui-e2e-test-isolated",
                str(self.project_root),
                "--no-cache"
            ], capture_output=True, text=True, timeout=300)

            if result.returncode == 0:
                logger.info("✅ Test image built successfully!")
                return True
            else:
                logger.error(f"❌ Failed to build test image: {result.stderr}")
                return False
        except subprocess.TimeoutExpired:
            logger.error("❌ Docker build timed out")
            return False
        except Exception as e:
            logger.error(f"❌ Error building test image: {e}")
            return False

    def run_isolated_tests(self, test_category: str = None, verbose: bool = False, generate_data: bool = False) -> bool:
        """Run tests in completely isolated Docker-in-Docker container."""
        logger.info("Running tests in isolated Docker-in-Docker container...")

        # Build command - DinD setup with no host Docker socket access
        cmd = [
            "docker", "run", "--rm",
            "--privileged",
            "-e", "DOCKER_TLS_CERTDIR=",
            "-e", "PYTHONPATH=/app/e2e",
            "-e", "DOCKER_BUILDKIT=1",
            "whaletui-e2e-test-isolated"
        ]

        if generate_data:
            cmd.append("python3 /app/docker/test_data_generator.py")
        elif test_category:
            cmd.extend([
                "python3", "-m", "pytest",
                f"/app/e2e/tests/{test_category}/",
                "-v" if verbose else "",
                "--html=/app/reports/test_report.html",
                "--self-contained-html"
            ])
        else:
            cmd.extend([
                "python3", "-m", "pytest",
                "/app/e2e/tests/",
                "-v" if verbose else "",
                "--html=/app/reports/test_report.html",
                "--self-contained-html"
            ])

        # Remove empty strings
        cmd = [arg for arg in cmd if arg]

        logger.info(f"Running command: {' '.join(cmd)}")

        try:
            result = subprocess.run(cmd, timeout=600)  # 10 minute timeout
            return result.returncode == 0
        except subprocess.TimeoutExpired:
            logger.error("❌ Test execution timed out")
            return False
        except Exception as e:
            logger.error(f"❌ Error running tests: {e}")
            return False

    def generate_test_data(self) -> bool:
        """Generate test data inside isolated container."""
        logger.info("Generating test data in isolated container...")
        return self.run_isolated_tests(generate_data=True)

    def run_all_tests(self, verbose: bool = False) -> bool:
        """Run all tests in isolated container."""
        logger.info("Running all tests in isolated container...")
        return self.run_isolated_tests(verbose=verbose)

    def run_category_tests(self, category: str, verbose: bool = False) -> bool:
        """Run specific category tests in isolated container."""
        logger.info(f"Running {category} tests in isolated container...")
        return self.run_isolated_tests(test_category=category, verbose=verbose)

    def cleanup(self):
        """Clean up Docker resources."""
        logger.info("Cleaning up Docker resources...")

        # Remove test containers
        subprocess.run(["docker", "container", "prune", "-f"], capture_output=True)

        # Remove test images
        subprocess.run(["docker", "image", "rm", "whaletui-e2e-test-isolated"], capture_output=True)

        logger.info("✅ Cleanup completed")

def main():
    """Main function."""
    parser = argparse.ArgumentParser(description="Run WhaleTUI e2e tests in isolated Docker")
    parser.add_argument("--build", action="store_true", help="Build test image")
    parser.add_argument("--test", choices=["unit", "integration", "ui", "performance", "search", "error_handling"],
                       help="Run specific test category")
    parser.add_argument("--generate-data", action="store_true", help="Generate test data only")
    parser.add_argument("--verbose", "-v", action="store_true", help="Verbose output")
    parser.add_argument("--cleanup", action="store_true", help="Clean up Docker resources")

    args = parser.parse_args()

    runner = IsolatedDockerTestRunner()

    if not runner.check_docker_available():
        logger.error("❌ Docker is not available!")
        return 1

    if args.cleanup:
        runner.cleanup()
        return 0

    if args.build or not any([args.test, args.generate_data]):
        if not runner.build_test_image():
            logger.error("❌ Failed to build test image!")
            return 1

    if args.generate_data:
        if not runner.generate_test_data():
            logger.error("❌ Failed to generate test data!")
            return 1
        logger.info("✅ Test data generated successfully!")
        return 0

    if args.test:
        if not runner.run_category_tests(args.test, args.verbose):
            logger.error(f"❌ {args.test} tests failed!")
            return 1
        logger.info(f"✅ {args.test} tests passed!")
        return 0

    # Run all tests by default
    if not runner.run_all_tests(args.verbose):
        logger.error("❌ Tests failed!")
        return 1

    logger.info("✅ All tests passed!")
    return 0

if __name__ == "__main__":
    sys.exit(main())
