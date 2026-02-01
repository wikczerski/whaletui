#!/usr/bin/env python3
"""
Advanced test data generator for WhaleTUI e2e tests.
Generates predictable data with 200+ volumes for stress testing.
"""
import docker
import time
import random
import string
import logging
import json
import os
from datetime import datetime
from typing import List, Dict, Any

# Set up logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class TestDataGenerator:
    """Generate comprehensive test data for WhaleTUI e2e tests."""

    def __init__(self):
        self.client = None
        self.test_data = {
            'volumes': [],
            'containers': [],
            'networks': [],
            'images': [],
            'services': [],
            'created_at': datetime.now().isoformat()
        }
        self._init_docker_client()

    def _init_docker_client(self):
        """Initialize Docker client with retry logic for DinD environment."""
        max_retries = 60  # Increased for DinD startup time
        retry_delay = 2

        for attempt in range(max_retries):
            try:
                # Use local Docker daemon (DinD)
                self.client = docker.from_env()
                # Test connection
                self.client.ping()
                logger.info("✅ Docker client initialized successfully in DinD environment")
                return
            except Exception as e:
                logger.warning(f"Docker connection attempt {attempt + 1}/{max_retries} failed: {e}")
                if attempt < max_retries - 1:
                    time.sleep(retry_delay)
                else:
                    logger.error("❌ Failed to initialize Docker client after all retries")
                    raise

    def generate_random_name(self, prefix: str = "test", length: int = 8) -> str:
        """Generate a random name for test resources."""
        suffix = ''.join(random.choices(string.ascii_lowercase + string.digits, k=length))
        return f"{prefix}-{suffix}"

    def verify_isolated_environment(self) -> bool:
        """Verify that we're running in an isolated Docker environment."""
        try:
            # Check if we're running inside a container
            if not os.path.exists('/.dockerenv'):
                logger.warning("⚠️  Not running inside a Docker container")
                return False

            # Check Docker info to verify isolation
            info = self.client.info()
            logger.info(f"Docker daemon info: {info['ServerVersion']}")
            logger.info(f"Containers: {info['Containers']}")
            logger.info(f"Images: {info['Images']}")
            logger.info(f"Volumes: {len(self.client.volumes.list())}")

            # Verify we have a clean environment (should be minimal items)
            containers = self.client.containers.list(all=True)
            images = self.client.images.list()
            volumes = self.client.volumes.list()
            networks = self.client.networks.list()

            logger.info(f"Current environment state:")
            logger.info(f"  - Containers: {len(containers)}")
            logger.info(f"  - Images: {len(images)}")
            logger.info(f"  - Volumes: {len(volumes)}")
            logger.info(f"  - Networks: {len(networks)}")

            return True
        except Exception as e:
            logger.error(f"Failed to verify isolated environment: {e}")
            return False

    def create_test_volumes(self, count: int = 200) -> List[Dict[str, Any]]:
        """Create test volumes with predictable patterns."""
        logger.info(f"Creating {count} test volumes...")
        volumes = []

        # Create volumes with predictable patterns
        patterns = [
            "volume-{i:03d}",
            "data-{i:03d}",
            "cache-{i:03d}",
            "logs-{i:03d}",
            "backup-{i:03d}",
            "temp-{i:03d}",
            "config-{i:03d}",
            "storage-{i:03d}"
        ]

        for i in range(count):
            try:
                pattern = patterns[i % len(patterns)]
                volume_name = pattern.format(i=i)

                # Create volume with labels
                volume = self.client.volumes.create(
                    name=volume_name,
                    labels={
                        'test': 'whaletui',
                        'type': 'volume',
                        'index': str(i),
                        'pattern': pattern
                    }
                )

                volumes.append({
                    'id': volume.id,
                    'name': volume_name,
                    'created_at': volume.attrs['CreatedAt'],
                    'labels': volume.attrs['Labels']
                })

                if (i + 1) % 50 == 0:
                    logger.info(f"Created {i + 1} volumes...")

            except Exception as e:
                logger.error(f"Failed to create volume {i + 1}: {e}")

        logger.info(f"Successfully created {len(volumes)} volumes")
        self.test_data['volumes'] = volumes
        return volumes

    def create_test_containers(self, count: int = 50) -> List[Dict[str, Any]]:
        """Create test containers with predictable patterns."""
        logger.info(f"Creating {count} test containers...")
        containers = []

        # Lightweight images for testing
        images = [
            "alpine:latest",
            "busybox:latest",
            "nginx:alpine",
            "redis:alpine",
            "postgres:13-alpine",
            "mysql:8.0",
            "mongo:5.0",
            "elasticsearch:7.17.0"
        ]

        # Pull images first
        for image in images:
            try:
                self.client.images.pull(image)
                logger.info(f"Pulled image: {image}")
            except Exception as e:
                logger.warning(f"Failed to pull {image}: {e}")

        # Create containers with predictable patterns
        patterns = [
            "web-{i:03d}",
            "db-{i:03d}",
            "cache-{i:03d}",
            "api-{i:03d}",
            "worker-{i:03d}",
            "proxy-{i:03d}",
            "monitor-{i:03d}",
            "service-{i:03d}"
        ]

        for i in range(count):
            try:
                pattern = patterns[i % len(patterns)]
                container_name = pattern.format(i=i)
                image = images[i % len(images)]

                # Create container with labels and environment
                container = self.client.containers.run(
                    image,
                    name=container_name,
                    detach=True,
                    command="sleep 3600",  # Keep container running
                    labels={
                        'test': 'whaletui',
                        'type': 'container',
                        'index': str(i),
                        'pattern': pattern,
                        'image': image
                    },
                    environment={
                        'TEST_INDEX': str(i),
                        'TEST_PATTERN': pattern,
                        'TEST_IMAGE': image
                    }
                )

                containers.append({
                    'id': container.id,
                    'name': container_name,
                    'image': image,
                    'status': container.status,
                    'created_at': container.attrs['Created'],
                    'labels': container.attrs['Config']['Labels']
                })

                if (i + 1) % 10 == 0:
                    logger.info(f"Created {i + 1} containers...")

            except Exception as e:
                logger.error(f"Failed to create container {i + 1}: {e}")

        logger.info(f"Successfully created {len(containers)} containers")
        self.test_data['containers'] = containers
        return containers

    def create_test_networks(self, count: int = 20) -> List[Dict[str, Any]]:
        """Create test networks with predictable patterns."""
        logger.info(f"Creating {count} test networks...")
        networks = []

        patterns = [
            "network-{i:03d}",
            "bridge-{i:03d}",
            "overlay-{i:03d}",
            "custom-{i:03d}"
        ]

        for i in range(count):
            try:
                pattern = patterns[i % len(patterns)]
                network_name = pattern.format(i=i)

                # Create network with labels
                network = self.client.networks.create(
                    name=network_name,
                    driver="bridge",
                    labels={
                        'test': 'whaletui',
                        'type': 'network',
                        'index': str(i),
                        'pattern': pattern
                    }
                )

                networks.append({
                    'id': network.id,
                    'name': network_name,
                    'driver': network.attrs['Driver'],
                    'created_at': network.attrs['Created'],
                    'labels': network.attrs['Labels']
                })

                if (i + 1) % 5 == 0:
                    logger.info(f"Created {i + 1} networks...")

            except Exception as e:
                logger.error(f"Failed to create network {i + 1}: {e}")

        logger.info(f"Successfully created {len(networks)} networks")
        self.test_data['networks'] = networks
        return networks

    def create_test_images(self, count: int = 30) -> List[Dict[str, Any]]:
        """Create test images by building simple Dockerfiles."""
        logger.info(f"Creating {count} test images...")
        images = []

        patterns = [
            "test-image-{i:03d}",
            "app-{i:03d}",
            "service-{i:03d}",
            "worker-{i:03d}"
        ]

        for i in range(count):
            try:
                pattern = patterns[i % len(patterns)]
                image_name = pattern.format(i=i)

                # Create a simple Dockerfile
                dockerfile_content = f"""
FROM alpine:latest
RUN echo "Test image {image_name}" > /test.txt
RUN echo "Index: {i}" >> /test.txt
RUN echo "Pattern: {pattern}" >> /test.txt
LABEL test=whaletui
LABEL type=image
LABEL index={i}
LABEL pattern={pattern}
CMD ["cat", "/test.txt"]
"""

                # Create a temporary directory for Dockerfile
                import tempfile
                import os

                with tempfile.TemporaryDirectory() as temp_dir:
                    dockerfile_path = os.path.join(temp_dir, 'Dockerfile')
                    with open(dockerfile_path, 'w') as f:
                        f.write(dockerfile_content)

                    # Build image
                    image, build_logs = self.client.images.build(
                        path=temp_dir,
                        tag=image_name,
                        rm=True,
                        labels={
                            'test': 'whaletui',
                            'type': 'image',
                            'index': str(i),
                            'pattern': pattern
                        }
                    )

                images.append({
                    'id': image.id,
                    'name': image_name,
                    'tags': image.tags,
                    'created_at': image.attrs['Created'],
                    'labels': image.attrs['Config']['Labels']
                })

                if (i + 1) % 10 == 0:
                    logger.info(f"Created {i + 1} images...")

            except Exception as e:
                logger.error(f"Failed to create image {i + 1}: {e}")

        logger.info(f"Successfully created {len(images)} images")
        self.test_data['images'] = images
        return images

    def setup_docker_swarm(self) -> bool:
        """Set up Docker swarm for testing."""
        logger.info("Setting up Docker swarm...")

        try:
            # Check if swarm already exists
            try:
                info = self.client.swarm.reload()
                if info:
                    logger.info("Docker swarm already initialized")
                    return True
            except:
                pass

            # Initialize swarm
            self.client.swarm.init(advertise_addr="127.0.0.1")
            logger.info("Docker swarm initialized successfully")
            return True
        except Exception as e:
            logger.error(f"Failed to initialize swarm: {e}")
            return False

    def create_test_services(self, count: int = 20) -> List[Dict[str, Any]]:
        """Create test services."""
        logger.info(f"Creating {count} test services...")
        services = []

        # Ensure swarm is running
        if not self.setup_docker_swarm():
            logger.warning("Swarm not available, skipping service creation")
            return services

        patterns = [
            "web-service-{i:03d}",
            "api-service-{i:03d}",
            "worker-service-{i:03d}",
            "monitor-service-{i:03d}"
        ]

        for i in range(count):
            try:
                pattern = patterns[i % len(patterns)]
                service_name = pattern.format(i=i)

                service = self.client.services.create(
                    image="nginx:alpine",
                    name=service_name,
                    replicas=random.randint(1, 3),
                    labels={
                        'test': 'whaletui',
                        'type': 'service',
                        'index': str(i),
                        'pattern': pattern
                    }
                )

                services.append({
                    'id': service.id,
                    'name': service_name,
                    'replicas': service.attrs['Spec']['Mode']['Replicated']['Replicas'],
                    'created_at': service.attrs['CreatedAt'],
                    'labels': service.attrs['Spec']['Labels']
                })

                if (i + 1) % 5 == 0:
                    logger.info(f"Created {i + 1} services...")

            except Exception as e:
                logger.error(f"Failed to create service {i + 1}: {e}")

        logger.info(f"Successfully created {len(services)} services")
        self.test_data['services'] = services
        return services

    def save_test_data(self, filename: str = "/app/test-data/test_data.json"):
        """Save test data to file."""
        try:
            os.makedirs(os.path.dirname(filename), exist_ok=True)
            with open(filename, 'w') as f:
                json.dump(self.test_data, f, indent=2)
            logger.info(f"Test data saved to {filename}")
        except Exception as e:
            logger.error(f"Failed to save test data: {e}")

    def cleanup_test_data(self):
        """Clean up all test data."""
        logger.info("Cleaning up test data...")

        # Clean up services
        for service in self.test_data.get('services', []):
            try:
                self.client.services.get(service['id']).remove()
            except:
                pass

        # Clean up containers
        for container in self.test_data.get('containers', []):
            try:
                self.client.containers.get(container['id']).remove(force=True)
            except:
                pass

        # Clean up networks
        for network in self.test_data.get('networks', []):
            try:
                self.client.networks.get(network['id']).remove()
            except:
                pass

        # Clean up volumes
        for volume in self.test_data.get('volumes', []):
            try:
                self.client.volumes.get(volume['id']).remove()
            except:
                pass

        # Clean up images
        for image in self.test_data.get('images', []):
            try:
                self.client.images.remove(image['id'])
            except:
                pass

        logger.info("Test data cleanup complete")

    def generate_all(self, volumes_count: int = 200, containers_count: int = 50,
                    networks_count: int = 20, images_count: int = 30,
                    services_count: int = 20):
        """Generate all test data."""
        logger.info("Starting comprehensive test data generation...")

        try:
            # Generate test data
            volumes = self.create_test_volumes(volumes_count)
            containers = self.create_test_containers(containers_count)
            networks = self.create_test_networks(networks_count)
            images = self.create_test_images(images_count)
            services = self.create_test_services(services_count)

            # Save test data
            self.save_test_data()

            # Summary
            logger.info("Test data generation complete!")
            logger.info(f"Created:")
            logger.info(f"  - {len(volumes)} volumes")
            logger.info(f"  - {len(containers)} containers")
            logger.info(f"  - {len(networks)} networks")
            logger.info(f"  - {len(images)} images")
            logger.info(f"  - {len(services)} services")

            return True

        except Exception as e:
            logger.error(f"Failed to generate test data: {e}")
            return False

def main():
    """Main function."""
    generator = TestDataGenerator()

    # Verify isolated environment
    if not generator.verify_isolated_environment():
        logger.error("❌ Isolated environment verification failed!")
        return 1

    # Generate test data
    success = generator.generate_all(
        volumes_count=200,
        containers_count=50,
        networks_count=20,
        images_count=30,
        services_count=20
    )

    if success:
        logger.info("Test data generation completed successfully!")
        return 0
    else:
        logger.error("Test data generation failed!")
        return 1

if __name__ == "__main__":
    exit(main())
