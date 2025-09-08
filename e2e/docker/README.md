# Docker-Based E2E Testing for WhaleTUI

This directory contains Docker-based end-to-end testing for WhaleTUI, including Docker-in-Docker (DinD) setup and comprehensive test data generation with 200+ volumes for stress testing.

## 🐳 Overview

The Docker-based testing environment provides:
- **Isolated testing environment** with Docker-in-Docker
- **Predictable test data** with 200+ volumes, 50+ containers, 20+ networks, 30+ images, and 20+ services
- **Stress testing capabilities** for low-resource consumption Docker items
- **Comprehensive test coverage** across all WhaleTUI features
- **Automated test data generation** and cleanup

## 📁 Directory Structure

```
docker/
├── Dockerfile.test              # Test environment Docker image
├── docker-compose.test.yml      # Docker Compose configuration
├── run_tests_in_docker.sh       # Test runner script (inside container)
├── run_docker_tests.py          # Docker test runner (host)
├── test_data_generator.py       # Test data generator
└── README.md                    # This file
```

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Go 1.19+ for building WhaleTUI
- Python 3.8+ for test scripts

### Basic Usage

1. **Run all tests in Docker:**
   ```bash
   make e2e-docker-test
   ```

2. **Run tests with verbose output:**
   ```bash
   make e2e-docker-test-verbose
   ```

3. **Run specific test category:**
   ```bash
   make e2e-docker-test-unit
   make e2e-docker-test-ui
   make e2e-docker-test-performance
   ```

4. **Generate test data only (200+ volumes):**
   ```bash
   make e2e-docker-generate-data
   ```

5. **Run with Docker Compose:**
   ```bash
   make e2e-docker-test-compose
   ```

## 🧪 Test Categories

### Unit Tests (`--category unit`)
- Basic functionality tests
- Application startup/shutdown
- Command-line arguments
- **Execution time**: ~2-5 minutes

### UI Tests (`--category ui`)
- User interface navigation
- Keyboard interactions
- Screen rendering
- **Execution time**: ~5-10 minutes

### Performance Tests (`--category performance`)
- Startup time testing
- Memory usage monitoring
- CPU usage under load
- **Execution time**: ~10-15 minutes

### Integration Tests (`--category integration`)
- Docker operations
- Container management
- Network operations
- **Execution time**: ~5-10 minutes

### Search Tests (`--category search`)
- Search functionality
- Filtering operations
- Search performance
- **Execution time**: ~3-8 minutes

### Error Handling Tests (`--category error_handling`)
- Error scenarios
- Recovery testing
- Edge case handling
- **Execution time**: ~5-10 minutes

## 📊 Test Data Generation

The test environment generates comprehensive test data:

### Volumes (200+)
- **Patterns**: `volume-001`, `data-002`, `cache-003`, etc.
- **Labels**: Test metadata for identification
- **Purpose**: Stress testing volume operations

### Containers (50+)
- **Images**: Alpine, BusyBox, Nginx, Redis, PostgreSQL, MySQL, MongoDB, Elasticsearch
- **Patterns**: `web-001`, `db-002`, `cache-003`, etc.
- **Status**: Running containers for testing

### Networks (20+)
- **Types**: Bridge networks
- **Patterns**: `network-001`, `bridge-002`, etc.
- **Purpose**: Network management testing

### Images (30+)
- **Custom built**: Simple Alpine-based images
- **Patterns**: `test-image-001`, `app-002`, etc.
- **Labels**: Test metadata

### Services (20+)
- **Swarm services**: Nginx-based services
- **Replicas**: 1-3 replicas per service
- **Patterns**: `web-service-001`, `api-service-002`, etc.

## 🔧 Advanced Usage

### Direct Python Script Usage

```bash
# Run all tests
cd e2e/docker
python run_docker_tests.py

# Run specific category
python run_docker_tests.py --category ui --verbose

# Generate test data only
python run_docker_tests.py --generate-data-only

# Use Docker Compose
python run_docker_tests.py --compose

# Clean up resources
python run_docker_tests.py --cleanup

# Show test reports
python run_docker_tests.py --show-reports
```

### Docker Compose Usage

```bash
# Build and run tests
cd e2e/docker
docker-compose -f docker-compose.test.yml up --build

# Run in background
docker-compose -f docker-compose.test.yml up -d

# View logs
docker-compose -f docker-compose.test.yml logs -f

# Clean up
docker-compose -f docker-compose.test.yml down -v
```

### Manual Docker Usage

```bash
# Build test image
docker build -f e2e/docker/Dockerfile.test -t whaletui-e2e-test .

# Run tests
docker run --rm --privileged \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd):/app/whaletui \
  -v whaletui-test-data:/app/test-data \
  -v whaletui-screenshots:/app/screenshots \
  -v whaletui-reports:/app/reports \
  whaletui-e2e-test \
  /app/docker/run_tests_in_docker.sh
```

## 📈 Test Reports

Test reports are generated in Docker volumes:

- **`whaletui-reports`**: HTML test reports
- **`whaletui-screenshots`**: Test screenshots
- **`whaletui-test-data`**: Generated test data JSON

### Accessing Reports

```bash
# List available reports
make e2e-docker-reports

# Copy reports to local filesystem
docker run --rm -v whaletui-reports:/reports -v $(pwd):/local alpine cp -r /reports/* /local/

# View specific report
docker run --rm -v whaletui-reports:/reports -p 8080:80 nginx:alpine
# Then visit http://localhost:8080/combined_report.html
```

## 🧹 Cleanup

### Automatic Cleanup
Tests automatically clean up generated data after completion.

### Manual Cleanup
```bash
# Clean up Docker resources
make e2e-docker-clean

# Or use Python script
python run_docker_tests.py --cleanup

# Clean up specific resources
docker system prune -f
docker volume prune -f
docker network prune -f
```

## 🔍 Debugging

### View Test Logs
```bash
# View container logs
docker logs whaletui-e2e-test

# Follow logs in real-time
docker logs -f whaletui-e2e-test
```

### Interactive Debugging
```bash
# Run interactive container
docker run --rm -it --privileged \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v $(pwd):/app/whaletui \
  whaletui-e2e-test bash

# Inside container, run tests manually
python3 /app/docker/test_data_generator.py
python3 -m pytest /app/e2e/tests/unit/ -v
```

### Test Data Inspection
```bash
# Inspect generated test data
docker run --rm -v whaletui-test-data:/data alpine cat /data/test_data.json

# List all volumes
docker volume ls | grep whaletui

# Inspect specific volume
docker volume inspect whaletui-test-data
```

## ⚡ Performance Considerations

### Resource Requirements
- **RAM**: Minimum 4GB, recommended 8GB+
- **CPU**: 2+ cores recommended
- **Disk**: 2GB+ for test data and images
- **Docker**: 2GB+ available space

### Optimization Tips
1. **Use SSD storage** for better I/O performance
2. **Increase Docker memory limit** if available
3. **Run tests in parallel** when possible
4. **Clean up regularly** to free disk space

### Test Data Scaling
You can modify the test data generation in `test_data_generator.py`:

```python
# Generate more volumes for stress testing
generator.generate_all(
    volumes_count=500,      # Increase from 200
    containers_count=100,   # Increase from 50
    networks_count=50,      # Increase from 20
    images_count=60,        # Increase from 30
    services_count=40       # Increase from 20
)
```

## 🚨 Troubleshooting

### Common Issues

1. **Docker daemon not running**
   ```bash
   # Start Docker daemon
   sudo systemctl start docker
   # Or on Windows/Mac, start Docker Desktop
   ```

2. **Permission denied errors**
   ```bash
   # Add user to docker group
   sudo usermod -aG docker $USER
   # Log out and back in
   ```

3. **Out of disk space**
   ```bash
   # Clean up Docker resources
   docker system prune -a -f
   docker volume prune -f
   ```

4. **Tests failing due to resource constraints**
   - Increase Docker memory limit
   - Reduce test data volume
   - Run tests sequentially instead of parallel

5. **Docker-in-Docker issues**
   ```bash
   # Ensure privileged mode is enabled
   docker run --privileged ...
   ```

### Getting Help

1. **Check logs**: Always check container logs first
2. **Verify Docker**: Ensure Docker is running and accessible
3. **Check resources**: Ensure sufficient RAM and disk space
4. **Test data**: Verify test data generation completed successfully

## 🔄 CI/CD Integration

### GitHub Actions Example
```yaml
name: E2E Tests
on: [push, pull_request]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      - name: Run E2E tests
        run: make e2e-docker-test-compose
      - name: Upload test reports
        uses: actions/upload-artifact@v3
        with:
          name: test-reports
          path: whaletui-reports/
```

### Jenkins Pipeline Example
```groovy
pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh 'make build'
            }
        }
        stage('E2E Tests') {
            steps {
                sh 'make e2e-docker-test-compose'
            }
        }
        stage('Publish Reports') {
            steps {
                publishHTML([
                    allowMissing: false,
                    alwaysLinkToLastBuild: true,
                    keepAll: true,
                    reportDir: 'whaletui-reports',
                    reportFiles: 'combined_report.html',
                    reportName: 'E2E Test Report'
                ])
            }
        }
    }
}
```

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [pytest Documentation](https://docs.pytest.org/)
- [WhaleTUI Project Documentation](../../README.md)

## 🤝 Contributing

When adding new tests:

1. **Follow naming conventions**: Use descriptive test names
2. **Add appropriate markers**: Mark tests with `@pytest.mark.docker`
3. **Update test data**: Add new patterns to test data generator
4. **Document changes**: Update this README if needed
5. **Test locally**: Run tests in Docker before submitting

## 📝 License

This testing framework is part of the WhaleTUI project and follows the same license terms.
