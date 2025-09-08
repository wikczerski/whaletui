#!/usr/bin/env python3
"""
Test script to verify the isolated Docker environment works correctly.
This script builds the DinD image and tests the isolated environment.
"""
import subprocess
import sys
import logging
import time

# Set up logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def run_command(cmd, timeout=300):
    """Run a command and return success status."""
    try:
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout)
        if result.returncode == 0:
            logger.info(f"✅ Command succeeded: {' '.join(cmd)}")
            if result.stdout:
                logger.info(f"Output: {result.stdout}")
            return True
        else:
            logger.error(f"❌ Command failed: {' '.join(cmd)}")
            logger.error(f"Error: {result.stderr}")
            return False
    except subprocess.TimeoutExpired:
        logger.error(f"❌ Command timed out: {' '.join(cmd)}")
        return False
    except Exception as e:
        logger.error(f"❌ Command error: {e}")
        return False

def test_dind_build():
    """Test building the DinD image."""
    logger.info("🧪 Testing DinD image build...")

    cmd = [
        "docker", "build",
        "-f", "e2e/docker/Dockerfile.test",
        "-t", "whaletui-e2e-test-isolated",
        ".",
        "--no-cache"
    ]

    return run_command(cmd, timeout=600)

def test_dind_run():
    """Test running the DinD container."""
    logger.info("🧪 Testing DinD container run...")

    cmd = [
        "docker", "run", "--rm",
        "--privileged",
        "-e", "DOCKER_TLS_CERTDIR=",
        "whaletui-e2e-test-isolated",
        "python3", "/app/docker/test_dind_setup.py"
    ]

    return run_command(cmd, timeout=300)

def test_data_generation():
    """Test generating test data in isolated environment."""
    logger.info("🧪 Testing test data generation in isolated environment...")

    cmd = [
        "docker", "run", "--rm",
        "--privileged",
        "-e", "DOCKER_TLS_CERTDIR=",
        "whaletui-e2e-test-isolated",
        "python3", "/app/docker/test_data_generator.py"
    ]

    return run_command(cmd, timeout=600)

def cleanup():
    """Clean up Docker resources."""
    logger.info("🧹 Cleaning up Docker resources...")

    # Remove test containers
    run_command(["docker", "container", "prune", "-f"])

    # Remove test images
    run_command(["docker", "image", "rm", "whaletui-e2e-test-isolated"])

def main():
    """Main function."""
    logger.info("🚀 Starting isolated Docker environment test...")

    # Test 1: Build DinD image
    if not test_dind_build():
        logger.error("❌ DinD image build failed!")
        return 1

    # Test 2: Run DinD container
    if not test_dind_run():
        logger.error("❌ DinD container run failed!")
        cleanup()
        return 1

    # Test 3: Generate test data
    if not test_data_generation():
        logger.error("❌ Test data generation failed!")
        cleanup()
        return 1

    logger.info("🎉 All tests passed! The isolated Docker environment is working correctly.")

    # Cleanup
    cleanup()
    return 0

if __name__ == "__main__":
    sys.exit(main())
