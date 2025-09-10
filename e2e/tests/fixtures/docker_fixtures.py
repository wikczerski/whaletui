"""
Docker fixtures for WhaleTUI e2e tests.
"""
import pytest
import subprocess
import time
import logging
from typing import List, Dict, Any


class DockerTestEnvironment:
    """Docker test environment manager."""

    def __init__(self):
        self.containers = []
        self.networks = []
        self.volumes = []
        self.services = []
        self.logger = logging.getLogger(__name__)

    def setup_containers(self, container_configs: List[Dict[str, Any]]) -> List[str]:
        """
        Set up test containers.

        Args:
            container_configs: List of container configuration dictionaries

        Returns:
            List of container IDs
        """
        container_ids = []

        for config in container_configs:
            try:
                # Build docker run command
                cmd = ["docker", "run", "-d"]

                # Add name
                if "name" in config:
                    cmd.extend(["--name", config["name"]])

                # Add ports
                if "ports" in config:
                    for port in config["ports"]:
                        cmd.extend(["-p", port])

                # Add environment variables
                if "environment" in config:
                    for key, value in config["environment"].items():
                        cmd.extend(["-e", f"{key}={value}"])

                # Add networks
                if "networks" in config:
                    for network in config["networks"]:
                        cmd.extend(["--network", network])

                # Add volumes
                if "volumes" in config:
                    for volume in config["volumes"]:
                        cmd.extend(["-v", volume])

                # Add image
                cmd.append(config["image"])

                # Run container
                result = subprocess.run(cmd, capture_output=True, text=True)
                if result.returncode == 0:
                    container_id = result.stdout.strip()
                    container_ids.append(container_id)
                    self.containers.append(container_id)
                    self.logger.info(f"Created container: {config.get('name', container_id)}")
                else:
                    self.logger.error(f"Failed to create container {config.get('name', 'unknown')}: {result.stderr}")

            except Exception as e:
                self.logger.error(f"Error creating container {config.get('name', 'unknown')}: {e}")

        return container_ids

    def setup_networks(self, network_names: List[str]) -> List[str]:
        """
        Set up test networks.

        Args:
            network_names: List of network names to create

        Returns:
            List of network IDs
        """
        network_ids = []

        for name in network_names:
            try:
                result = subprocess.run(
                    ["docker", "network", "create", name],
                    capture_output=True,
                    text=True
                )
                if result.returncode == 0:
                    network_id = result.stdout.strip()
                    network_ids.append(network_id)
                    self.networks.append(network_id)
                    self.logger.info(f"Created network: {name}")
                else:
                    self.logger.error(f"Failed to create network {name}: {result.stderr}")

            except Exception as e:
                self.logger.error(f"Error creating network {name}: {e}")

        return network_ids

    def setup_volumes(self, volume_names: List[str]) -> List[str]:
        """
        Set up test volumes.

        Args:
            volume_names: List of volume names to create

        Returns:
            List of volume names
        """
        created_volumes = []

        for name in volume_names:
            try:
                result = subprocess.run(
                    ["docker", "volume", "create", name],
                    capture_output=True,
                    text=True
                )
                if result.returncode == 0:
                    created_volumes.append(name)
                    self.volumes.append(name)
                    self.logger.info(f"Created volume: {name}")
                else:
                    self.logger.error(f"Failed to create volume {name}: {result.stderr}")

            except Exception as e:
                self.logger.error(f"Error creating volume {name}: {e}")

        return created_volumes

    def setup_swarm(self) -> bool:
        """
        Set up Docker swarm.

        Returns:
            True if successful, False otherwise
        """
        try:
            result = subprocess.run(
                ["docker", "swarm", "init", "--advertise-addr", "127.0.0.1"],
                capture_output=True,
                text=True
            )
            if result.returncode == 0:
                self.logger.info("Docker swarm initialized")
                return True
            else:
                self.logger.error(f"Failed to initialize swarm: {result.stderr}")
                return False

        except Exception as e:
            self.logger.error(f"Error initializing swarm: {e}")
            return False

    def setup_services(self, service_configs: List[Dict[str, Any]]) -> List[str]:
        """
        Set up test services.

        Args:
            service_configs: List of service configuration dictionaries

        Returns:
            List of service IDs
        """
        service_ids = []

        for config in service_configs:
            try:
                # Build docker service create command
                cmd = ["docker", "service", "create"]

                # Add name
                if "name" in config:
                    cmd.extend(["--name", config["name"]])

                # Add replicas
                if "replicas" in config:
                    cmd.extend(["--replicas", str(config["replicas"])])

                # Add networks
                if "networks" in config:
                    for network in config["networks"]:
                        cmd.extend(["--network", network])

                # Add image
                cmd.append(config["image"])

                # Run service
                result = subprocess.run(cmd, capture_output=True, text=True)
                if result.returncode == 0:
                    service_id = result.stdout.strip()
                    service_ids.append(service_id)
                    self.services.append(service_id)
                    self.logger.info(f"Created service: {config.get('name', service_id)}")
                else:
                    self.logger.error(f"Failed to create service {config.get('name', 'unknown')}: {result.stderr}")

            except Exception as e:
                self.logger.error(f"Error creating service {config.get('name', 'unknown')}: {e}")

        return service_ids

    def cleanup(self):
        """Clean up all test resources."""
        self.logger.info("Cleaning up Docker test environment...")

        # Remove services
        for service in self.services:
            try:
                subprocess.run(["docker", "service", "rm", service], capture_output=True)
                self.logger.info(f"Removed service: {service}")
            except Exception as e:
                self.logger.error(f"Error removing service {service}: {e}")

        # Leave swarm
        try:
            subprocess.run(["docker", "swarm", "leave", "--force"], capture_output=True)
            self.logger.info("Left Docker swarm")
        except Exception as e:
            self.logger.error(f"Error leaving swarm: {e}")

        # Remove containers
        for container in self.containers:
            try:
                subprocess.run(["docker", "rm", "-f", container], capture_output=True)
                self.logger.info(f"Removed container: {container}")
            except Exception as e:
                self.logger.error(f"Error removing container {container}: {e}")

        # Remove networks
        for network in self.networks:
            try:
                subprocess.run(["docker", "network", "rm", network], capture_output=True)
                self.logger.info(f"Removed network: {network}")
            except Exception as e:
                self.logger.error(f"Error removing network {network}: {e}")

        # Remove volumes
        for volume in self.volumes:
            try:
                subprocess.run(["docker", "volume", "rm", volume], capture_output=True)
                self.logger.info(f"Removed volume: {volume}")
            except Exception as e:
                self.logger.error(f"Error removing volume {volume}: {e}")

        self.logger.info("Docker test environment cleanup complete")


@pytest.fixture(scope="session")
def docker_environment():
    """Docker test environment fixture."""
    env = DockerTestEnvironment()
    yield env
    env.cleanup()


@pytest.fixture(scope="session")
def docker_containers(docker_environment):
    """Docker containers fixture."""
    from tests.utils.test_helpers import TestHelpers

    container_configs = TestHelpers.get_docker_test_containers()
    container_ids = docker_environment.setup_containers(container_configs)

    # Wait for containers to be ready
    time.sleep(5)

    yield container_ids


@pytest.fixture(scope="session")
def docker_networks(docker_environment):
    """Docker networks fixture."""
    from tests.utils.test_helpers import TestHelpers

    network_names = TestHelpers.get_docker_test_networks()
    network_ids = docker_environment.setup_networks(network_names)

    yield network_ids


@pytest.fixture(scope="session")
def docker_volumes(docker_environment):
    """Docker volumes fixture."""
    from tests.utils.test_helpers import TestHelpers

    volume_names = TestHelpers.get_docker_test_volumes()
    volume_names = docker_environment.setup_volumes(volume_names)

    yield volume_names


@pytest.fixture(scope="session")
def docker_swarm(docker_environment):
    """Docker swarm fixture."""
    success = docker_environment.setup_swarm()
    yield success


@pytest.fixture(scope="session")
def docker_services(docker_environment, docker_swarm):
    """Docker services fixture."""
    if not docker_swarm:
        pytest.skip("Docker swarm not available")

    service_configs = [
        {
            "name": "test-nginx-service",
            "image": "nginx:alpine",
            "replicas": 2,
            "networks": ["test-network-1"]
        },
        {
            "name": "test-redis-service",
            "image": "redis:alpine",
            "replicas": 1,
            "networks": ["test-network-2"]
        }
    ]

    service_ids = docker_environment.setup_services(service_configs)

    # Wait for services to be ready
    time.sleep(10)

    yield service_ids
