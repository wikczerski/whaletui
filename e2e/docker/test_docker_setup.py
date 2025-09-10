#!/usr/bin/env python3
"""
Simple test script to verify Docker-based testing setup.
"""
import subprocess
import sys
import time
from pathlib import Path

def test_docker_binary():
    """Test if the Docker binary exists and is executable."""
    # Check multiple possible locations
    possible_paths = [
        Path("/app/whaletui/whaletui"),
        Path("/app/whaletui/whaletui.exe"),
        Path("/app/whaletui/whaletui"),
    ]

    binary_path = None
    for path in possible_paths:
        if path.exists():
            binary_path = path
            break

    if binary_path:
        print(f"Testing Docker binary at: {binary_path}")
        print(f"Binary exists: {binary_path.exists()}")
        print(f"Binary is file: {binary_path.is_file()}")
        print(f"Binary is executable: {binary_path.stat().st_mode & 0o111 != 0}")
        print("✅ Docker binary found and ready!")
        return True
    else:
        print("❌ Docker binary not found in any expected location!")
        print("Searched paths:")
        for path in possible_paths:
            print(f"  - {path}: {path.exists()}")
        return False

def test_docker_connection():
    """Test Docker daemon connection."""
    try:
        result = subprocess.run(
            ["docker", "version", "--format", "{{.Server.Version}}"],
            capture_output=True,
            text=True,
            timeout=5
        )
        if result.returncode == 0:
            print(f"✅ Docker daemon connected! Version: {result.stdout.strip()}")
            return True
        else:
            print(f"❌ Docker daemon connection failed: {result.stderr}")
            return False
    except Exception as e:
        print(f"❌ Docker daemon connection error: {e}")
        return False

def test_whaletui_basic():
    """Test basic WhaleTUI functionality."""
    # Check multiple possible locations
    possible_paths = [
        Path("/app/whaletui/whaletui"),
        Path("/app/whaletui/whaletui.exe"),
    ]

    binary_path = None
    for path in possible_paths:
        if path.exists():
            binary_path = path
            break

    if not binary_path:
        print("❌ Binary not found, skipping test")
        return False

    try:
        # Test help command
        result = subprocess.run(
            [str(binary_path), "--help"],
            capture_output=True,
            text=True,
            timeout=10
        )

        if result.returncode == 0:
            print("✅ WhaleTUI help command works!")
            print(f"Help output: {result.stdout[:100]}...")
            return True
        else:
            print(f"❌ WhaleTUI help command failed: {result.stderr}")
            return False
    except Exception as e:
        print(f"❌ WhaleTUI test error: {e}")
        return False

def test_docker_resources():
    """Test Docker resources (volumes, containers, etc.)."""
    try:
        # List volumes
        result = subprocess.run(
            ["docker", "volume", "ls", "--format", "{{.Name}}"],
            capture_output=True,
            text=True,
            timeout=10
        )

        if result.returncode == 0:
            volumes = result.stdout.strip().split('\n')
            test_volumes = [v for v in volumes if 'test' in v.lower() or 'whaletui' in v.lower()]
            print(f"✅ Found {len(test_volumes)} test volumes")
            if test_volumes:
                print(f"Test volumes: {test_volumes[:5]}...")
            return True
        else:
            print(f"❌ Failed to list volumes: {result.stderr}")
            return False
    except Exception as e:
        print(f"❌ Docker resources test error: {e}")
        return False

def main():
    """Run all tests."""
    print("🐳 Testing Docker-based WhaleTUI e2e setup...")
    print("=" * 50)

    tests = [
        ("Docker Binary", test_docker_binary),
        ("Docker Connection", test_docker_connection),
        ("WhaleTUI Basic", test_whaletui_basic),
        ("Docker Resources", test_docker_resources),
    ]

    results = []
    for test_name, test_func in tests:
        print(f"\n🧪 Running {test_name} test...")
        try:
            result = test_func()
            results.append((test_name, result))
        except Exception as e:
            print(f"❌ {test_name} test failed with exception: {e}")
            results.append((test_name, False))

    print("\n" + "=" * 50)
    print("📊 Test Results Summary:")
    print("=" * 50)

    passed = 0
    total = len(results)

    for test_name, result in results:
        status = "✅ PASS" if result else "❌ FAIL"
        print(f"{test_name:20} {status}")
        if result:
            passed += 1

    print("=" * 50)
    print(f"Total: {passed}/{total} tests passed")

    if passed == total:
        print("🎉 All tests passed! Docker-based testing is ready!")
        return 0
    else:
        print("⚠️  Some tests failed. Check the output above.")
        return 1

if __name__ == "__main__":
    sys.exit(main())
