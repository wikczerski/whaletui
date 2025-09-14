# Volume-Based E2E Testing

This directory contains volume-based e2e testing setup for WhaleTUI that allows for live code changes without rebuilding Docker images.

## Quick Start

### 1. Build the Base Image (One Time)
```bash
# Build the base Docker image
docker build -f e2e/docker/Dockerfile.test -t whaletui-e2e-test-isolated .
```

### 2. Run Tests with Volumes
```bash
# Run all tests
python e2e/docker/run_volume_tests.py

# Run specific test category
python e2e/docker/run_volume_tests.py --category ui

# Run single test file
python e2e/docker/run_volume_tests.py --test unit/test_basic.py

# Run with verbose output
python e2e/docker/run_volume_tests.py --verbose

# Generate test data
python e2e/docker/run_volume_tests.py --generate-data

# Run interactive shell
python e2e/docker/run_volume_tests.py --shell
```

### 3. Windows Users
```cmd
REM Use the batch file instead
e2e/docker/run_volume_tests.bat --verbose
```

## How It Works

The volume-based setup mounts your local source code into the Docker container, allowing you to:

- ✅ Make code changes and test immediately
- ✅ No need to rebuild Docker images
- ✅ Faster development cycle
- ✅ Live debugging capabilities

## File Structure

```
e2e/docker/
├── Dockerfile.test              # Base Docker image
├── docker-compose.volume.yml    # Volume-based compose file
├── run_volume_tests.py         # Python test runner
├── run_volume_tests.bat        # Windows batch file
├── test_data_generator.py      # Test data generator
└── README_VOLUME.md            # This file
```

## Volume Mounts

The setup mounts the following volumes:

- `../../:/app/whaletui:ro` - Your WhaleTUI source code (read-only)
- `./:/app/e2e:ro` - E2E test code (read-only)
- `whaletui-test-data:/app/test-data` - Test data (writable)
- `whaletui-screenshots:/app/screenshots` - Screenshots (writable)
- `whaletui-reports:/app/reports` - Test reports (writable)
- `whaletui-logs:/app/logs` - Logs (writable)

## Benefits Over Image Rebuilding

| Feature | Image Rebuilding | Volume Mounting |
|---------|------------------|-----------------|
| Code Changes | Rebuild required | Instant |
| Development Speed | Slow | Fast |
| Docker Layers | Many | Few |
| Storage Usage | High | Low |
| Debugging | Difficult | Easy |

## Commands

### Test Commands
```bash
# Run all tests
python e2e/docker/run_volume_tests.py

# Run specific categories
python e2e/docker/run_volume_tests.py --category ui
python e2e/docker/run_volume_tests.py --category unit
python e2e/docker/run_volume_tests.py --category performance

# Run single test
python e2e/docker/run_volume_tests.py --test unit/test_basic.py

# Verbose output
python e2e/docker/run_volume_tests.py --verbose
```

### Development Commands
```bash
# Interactive shell for debugging
python e2e/docker/run_volume_tests.py --shell

# Generate test data
python e2e/docker/run_volume_tests.py --generate-data

# Clean up Docker resources
python e2e/docker/run_volume_tests.py --cleanup
```

### Docker Compose Commands
```bash
# Start the test environment
docker-compose -f e2e/docker/docker-compose.volume.yml up

# Run tests manually
docker-compose -f e2e/docker/docker-compose.volume.yml run --rm whaletui-test python3 -m pytest /app/e2e/tests/unit/ -v

# Clean up
docker-compose -f e2e/docker/docker-compose.volume.yml down
```

## Troubleshooting

### Permission Issues
If you encounter permission issues with volumes:
```bash
# Fix ownership
docker-compose -f e2e/docker/docker-compose.volume.yml run --rm whaletui-test chown -R testuser:testuser /app/test-data /app/screenshots /app/reports /app/logs
```

### Clean Up
```bash
# Remove all test containers and volumes
docker-compose -f e2e/docker/docker-compose.volume.yml down -v

# Remove unused images
docker image prune -f
```

### Debug Mode
```bash
# Run with debug logging
PYTHONPATH=/app/e2e python e2e/docker/run_volume_tests.py --verbose --test unit/test_basic.py
```

## Integration with CI/CD

For CI/CD pipelines, you can still use the image-based approach:

```bash
# CI/CD: Build and run tests
docker build -f e2e/docker/Dockerfile.test -t whaletui-e2e-test-isolated .
docker run --rm --privileged -e DOCKER_TLS_CERTDIR= whaletui-e2e-test-isolated python3 -m pytest /app/e2e/tests/ -v
```

## Best Practices

1. **Use volumes for development** - Fast iteration and debugging
2. **Use image building for CI/CD** - Consistent, reproducible builds
3. **Clean up regularly** - Remove unused containers and volumes
4. **Monitor disk usage** - Volume mounts can accumulate data
5. **Use .gitignore** - Exclude test artifacts from version control
