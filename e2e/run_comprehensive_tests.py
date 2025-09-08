#!/usr/bin/env python3
"""
Comprehensive test runner script for WhaleTUI e2e tests.
"""
import os
import sys
import argparse
import subprocess
import time
from pathlib import Path


def setup_environment():
    """Set up the test environment."""
    # Create necessary directories
    os.makedirs("screenshots", exist_ok=True)
    os.makedirs("reports", exist_ok=True)
    os.makedirs("test_data", exist_ok=True)
    os.makedirs("logs", exist_ok=True)

    # Set up logging
    import logging
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )


def install_dependencies():
    """Install Python dependencies."""
    print("Installing Python dependencies...")
    result = subprocess.run([
        sys.executable, "-m", "pip", "install", "-r", "requirements.txt"
    ], capture_output=True, text=True)

    if result.returncode != 0:
        print(f"Failed to install dependencies: {result.stderr}")
        return False

    print("Dependencies installed successfully")
    return True


def run_test_category(category, verbose=False, parallel=False):
    """Run tests for a specific category."""
    cmd = [sys.executable, "-m", "pytest", f"tests/{category}/"]

    if verbose:
        cmd.append("-v")

    if parallel:
        cmd.extend(["-n", "auto"])

    print(f"Running {category} tests with command: {' '.join(cmd)}")

    result = subprocess.run(cmd, capture_output=False)
    return result.returncode == 0


def run_all_tests(verbose=False, parallel=False):
    """Run all tests."""
    cmd = [sys.executable, "-m", "pytest", "tests/"]

    if verbose:
        cmd.append("-v")

    if parallel:
        cmd.extend(["-n", "auto"])

    print(f"Running all tests with command: {' '.join(cmd)}")

    result = subprocess.run(cmd, capture_output=False)
    return result.returncode == 0


def run_tests_by_marker(marker, verbose=False, parallel=False):
    """Run tests by marker."""
    cmd = [sys.executable, "-m", "pytest", "tests/", "-m", marker]

    if verbose:
        cmd.append("-v")

    if parallel:
        cmd.extend(["-n", "auto"])

    print(f"Running tests with marker '{marker}' with command: {' '.join(cmd)}")

    result = subprocess.run(cmd, capture_output=False)
    return result.returncode == 0


def run_specific_test(test_path, verbose=False):
    """Run a specific test."""
    cmd = [sys.executable, "-m", "pytest", test_path]

    if verbose:
        cmd.append("-v")

    print(f"Running specific test with command: {' '.join(cmd)}")

    result = subprocess.run(cmd, capture_output=False)
    return result.returncode == 0


def main():
    """Main function."""
    parser = argparse.ArgumentParser(description="Run comprehensive WhaleTUI e2e tests")
    parser.add_argument(
        "--category", "-c",
        choices=["unit", "integration", "ui", "performance", "docker", "search", "error_handling"],
        help="Run tests for a specific category"
    )
    parser.add_argument(
        "--marker", "-m",
        help="Run tests with specific marker (e.g., 'slow', 'docker', 'ui')"
    )
    parser.add_argument(
        "--test", "-t",
        help="Run a specific test file or test function"
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Verbose output"
    )
    parser.add_argument(
        "--parallel", "-p",
        action="store_true",
        help="Run tests in parallel"
    )
    parser.add_argument(
        "--install-deps",
        action="store_true",
        help="Install dependencies before running tests"
    )
    parser.add_argument(
        "--setup-only",
        action="store_true",
        help="Only set up the environment without running tests"
    )
    parser.add_argument(
        "--exclude-slow",
        action="store_true",
        help="Exclude slow tests"
    )
    parser.add_argument(
        "--docker-only",
        action="store_true",
        help="Run only Docker tests"
    )
    parser.add_argument(
        "--ui-only",
        action="store_true",
        help="Run only UI tests"
    )
    parser.add_argument(
        "--performance-only",
        action="store_true",
        help="Run only performance tests"
    )

    args = parser.parse_args()

    # Set up environment
    setup_environment()

    if args.install_deps:
        if not install_dependencies():
            sys.exit(1)

    if args.setup_only:
        print("Environment setup complete")
        return

    # Determine which tests to run
    success = True

    if args.test:
        success = run_specific_test(args.test, args.verbose)
    elif args.category:
        success = run_test_category(args.category, args.verbose, args.parallel)
    elif args.marker:
        success = run_tests_by_marker(args.marker, args.verbose, args.parallel)
    elif args.exclude_slow:
        success = run_tests_by_marker("not slow", args.verbose, args.parallel)
    elif args.docker_only:
        success = run_tests_by_marker("docker", args.verbose, args.parallel)
    elif args.ui_only:
        success = run_tests_by_marker("ui", args.verbose, args.parallel)
    elif args.performance_only:
        success = run_tests_by_marker("performance", args.verbose, args.parallel)
    else:
        success = run_all_tests(args.verbose, args.parallel)

    if success:
        print("\n✅ All tests passed!")
    else:
        print("\n❌ Some tests failed!")
        sys.exit(1)


if __name__ == "__main__":
    main()
