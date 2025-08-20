# Contributing to whaletui

Thank you for your interest in contributing to whaletui! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Code Style](#code-style)
- [Pull Request Process](#pull-request-process)
- [Release Process](#release-process)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- **Go 1.25.0+** - [Download Go](https://golang.org/dl/)
- **Git** - [Download Git](https://git-scm.com/)
- **Docker Desktop** - [Download Docker](https://docker.com/products/docker-desktop/)
- **Windows 10/11** (for testing)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/whaletui.git
   cd whaletui
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/wikczerski/whaletui.git
   ```

## Development Setup

### 1. Install Dependencies

```bash
# Install Go dependencies
go mod download
go mod tidy

# Install development tools
go install golang.org/x/vuln/cmd/govulncheck@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 2. Verify Setup

```bash
# Run tests
go test ./...

# Run linter
golangci-lint run

# Build the application
go build -o whaletui.exe .
```



## Making Changes

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

- Write clear, readable code
- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation as needed

### 3. Commit Your Changes

```bash
git add .
git commit -m "feat: add new feature description"
```

**Commit Message Format:**
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Writing Tests

- Write tests for all new functionality
- Use descriptive test names
- Test both success and failure cases
- Aim for good test coverage

### Test Files

- Test files should be named `*_test.go`
- Place tests in the same package as the code being tested
- Use the `testing` package and follow Go testing conventions

## Code Style

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for code formatting
- Use `goimports` for import organization
- Follow Go naming conventions

### Code Quality

- Run `golangci-lint` before submitting PRs
- Fix all linting issues
- Ensure code passes `go vet`
- Write clear, documented code

### Project Structure

- Keep the clean architecture pattern
- Separate concerns between layers
- Use interfaces for dependency injection
- Follow existing patterns in the codebase

## Pull Request Process

### 1. Prepare Your PR

- Ensure your branch is up to date with upstream
- Run all tests and linting
- Update documentation if needed
- Write a clear PR description

### 2. PR Description Template

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] Added tests for new functionality
- [ ] All existing tests pass

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Code is commented where necessary
- [ ] Documentation updated
```

### 3. Submit and Review

- Submit your PR
- Respond to review comments promptly
- Make requested changes
- Ensure CI checks pass

## Release Process

### Creating a Release

1. **Update Version**: Update version in relevant files
2. **Create Tag**: Create and push a version tag
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. **Automated Release**: GitHub Actions will automatically:
   - Run tests on all platforms
   - Build binaries for Windows, Linux, and macOS
   - Create a GitHub release with assets

### Versioning

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR.MINOR.PATCH**
- **MAJOR**: Breaking changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible

## Getting Help

- **Issues**: Use GitHub Issues for bug reports and feature requests
- **Discussions**: Use GitHub Discussions for questions and ideas
- **Code Review**: Ask questions in PR reviews

## Thank You

Thank you for contributing to whaletui! Your contributions help make this project better for everyone.
