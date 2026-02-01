#!/usr/bin/env python3
"""
Test script to verify Docker-in-Docker (DinD) setup is working correctly.
This script tests the isolated Docker environment.
"""
import docker
import time
import logging
import os
import sys

# Set up logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def test_dind_setup():
    """Test the DinD setup."""
    logger.info("🧪 Testing Docker-in-Docker setup...")

    # Test 1: Check if we're running inside a container
    if not os.path.exists('/.dockerenv'):
        logger.error("❌ Not running inside a Docker container!")
        return False
    logger.info("✅ Running inside Docker container")

    # Test 2: Initialize Docker client
    try:
        client = docker.from_env()
        logger.info("✅ Docker client initialized")
    except Exception as e:
        logger.error(f"❌ Failed to initialize Docker client: {e}")
        return False

    # Test 3: Wait for Docker daemon to be ready
    max_retries = 30
    for attempt in range(max_retries):
        try:
            client.ping()
            logger.info("✅ Docker daemon is ready")
            break
        except Exception as e:
            if attempt < max_retries - 1:
                logger.info(f"⏳ Waiting for Docker daemon... (attempt {attempt + 1}/{max_retries})")
                time.sleep(2)
            else:
                logger.error(f"❌ Docker daemon not ready after {max_retries} attempts: {e}")
                return False

    # Test 4: Check Docker info
    try:
        info = client.info()
        logger.info(f"✅ Docker daemon info: {info['ServerVersion']}")
        logger.info(f"   Containers: {info['Containers']}")
        logger.info(f"   Images: {info['Images']}")
        logger.info(f"   Volumes: {len(client.volumes.list())}")
    except Exception as e:
        logger.error(f"❌ Failed to get Docker info: {e}")
        return False

    # Test 5: Create a test container
    try:
        logger.info("🧪 Creating test container...")
        container = client.containers.run(
            "alpine:latest",
            command="echo 'Hello from DinD!'",
            remove=True,
            detach=False
        )
        logger.info("✅ Test container created and ran successfully")
    except Exception as e:
        logger.error(f"❌ Failed to create test container: {e}")
        return False

    # Test 6: Create a test volume
    try:
        logger.info("🧪 Creating test volume...")
        volume = client.volumes.create(name="test-dind-volume")
        logger.info("✅ Test volume created successfully")

        # Clean up
        volume.remove()
        logger.info("✅ Test volume cleaned up")
    except Exception as e:
        logger.error(f"❌ Failed to create test volume: {e}")
        return False

    # Test 7: Create a test network
    try:
        logger.info("🧪 Creating test network...")
        network = client.networks.create(name="test-dind-network")
        logger.info("✅ Test network created successfully")

        # Clean up
        network.remove()
        logger.info("✅ Test network cleaned up")
    except Exception as e:
        logger.error(f"❌ Failed to create test network: {e}")
        return False

    logger.info("🎉 All DinD tests passed! The isolated Docker environment is working correctly.")
    return True

def main():
    """Main function."""
    success = test_dind_setup()
    if success:
        logger.info("✅ DinD setup verification completed successfully!")
        return 0
    else:
        logger.error("❌ DinD setup verification failed!")
        return 1

if __name__ == "__main__":
    sys.exit(main())
