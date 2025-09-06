---
id: coding-standards
title: Coding Standards
sidebar_label: Coding Standards
description: Development guidelines and coding standards for WhaleTUI contributors
---

# Coding Standards

This document outlines the coding standards and best practices that all WhaleTUI contributors must follow to maintain code quality and consistency.

## Table of Contents

- [General Principles](#general-principles)
- [Go Language Standards](#go-language-standards)
- [Code Structure](#code-structure)
- [Naming Conventions](#naming-conventions)
- [Documentation](#documentation)
- [Testing](#testing)
- [Error Handling](#error-handling)
- [Performance](#performance)
- [Security](#security)

## General Principles

### DRY (Don't Repeat Yourself)
- Avoid code duplication
- Extract common functionality into reusable functions
- Use interfaces and abstractions appropriately

### KISS (Keep It Simple, Stupid)
- Prefer simple solutions over complex ones
- Reduce complexity wherever possible
- Write code that's easy to understand and maintain

### Single Responsibility
- Each function should do one thing well
- Classes should have a single, clear purpose
- Keep functions small and focused

## Go Language Standards

### Code Formatting
- Use `gofumpt` for automatic formatting (via `make fmt`)
- Follow Go's standard formatting conventions
- Use 4-space indentation (tabs converted to spaces)

### Error Handling
- Always check and handle errors explicitly
- Use meaningful error messages
- Wrap errors with context when appropriate
- Use `errors.Is` and `errors.As` for error comparison

```go
// Good
if err != nil {
    return fmt.Errorf("failed to start container %s: %w", containerID, err)
}

// Avoid
if err != nil {
    return err
}
```

### Logging
- Use the centralized logger from `internal/logger.go`
- Never use `fmt.Print*` methods
- Use appropriate log levels (DEBUG, INFO, WARN, ERROR)
- Include relevant context in log messages

```go
// Good
logger.Info("Container started successfully", "container_id", containerID, "image", imageName)

// Avoid
fmt.Printf("Container %s started\n", containerID)
```

## Code Structure

### File Organization
- Group related concepts vertically
  - Struct -> Constructor -> Public methods -> private methods
- Keep dependent functions close together
- Place functions in logical downward flow
- Use whitespace to group related code

### Function Design
- Keep functions small (ideally under 30 lines, golangci-lint limits to 40)
- Use descriptive names that explain intent
- Minimize arguments (prefer structs for multiple parameters)
- Avoid side effects
- No flag arguments - split into separate methods

```go
// Good
func (cm *ContainerManager) StartContainer(containerID string) error {
    // Implementation
}

func (cm *ContainerManager) StartContainerWithOptions(containerID string, opts ContainerOptions) error {
    // Implementation
}

// Avoid
func (cm *ContainerManager) StartContainer(containerID string, withLogs, withMetrics bool) error {
    // Implementation
}
```

### Variable Declaration
- Declare variables close to their usage
- Use meaningful variable names
- Use constants for magic numbers

## Naming Conventions

### Packages
- Use lowercase, single-word names
- Avoid underscores or mixed caps
- Use descriptive names that indicate purpose

### Functions and Methods
- Use camelCase
- Start with a verb for actions
- Be descriptive and specific

```go
// Good
func StartContainer()
func GetContainerStatus()
func ValidateImageName()

// Avoid
func Container()
func Status()
func Validate()
```

### Variables and Constants
- Use camelCase for variables
- Use PascalCase for exported constants
- Use descriptive names that explain purpose

```go
// Good
const MaxRetryAttempts = 3
const DefaultTimeout = 30 * time.Second

var containerStatus string
var maxMemoryUsage int64

// Avoid
const MAX_RETRY = 3
const TIMEOUT = 30

var status string
var mem int64
```

## Documentation

### Code Comments
- Explain intent, not implementation
- Use comments to clarify complex logic
- Avoid obvious or redundant comments
- No closing brace comments

```go
// Good
// Retry operation up to MaxRetryAttempts times with exponential backoff
for attempt := 0; attempt < MaxRetryAttempts; attempt++ {
    // Implementation
}

// Avoid
// Loop through attempts
for attempt := 0; attempt < MaxRetryAttempts; attempt++ {
    // Implementation
} // End of loop
```

### Function Documentation
- Document all exported functions
- Use Go's standard comment format
- Include examples for complex functions

```go
// StartContainer starts a Docker container with the given ID.
// Returns an error if the container cannot be started or if the
// container ID is invalid.
//
// Example:
//
//	err := manager.StartContainer("abc123")
//	if err != nil {
//	    log.Printf("Failed to start container: %v", err)
//	}
func (cm *ContainerManager) StartContainer(containerID string) error {
    // Implementation
}
```

## Testing

### Test Structure
- One assertion per test
- Tests should be readable, fast, independent, and repeatable
- Use descriptive test names
- Group related tests together

```go
func TestContainerManager_StartContainer(t *testing.T) {
    t.Run("successful start", func(t *testing.T) {
        // Test successful case
    })

    t.Run("invalid container ID", func(t *testing.T) {
        // Test error case
    })

    t.Run("container already running", func(t *testing.T) {
        // Test edge case
    })
}
```

### Mocking
- Use mockery to generate mocks
- Never create custom mocks manually
- Use existing mocks when available
- Mock external dependencies, not internal logic

## Error Handling

### Error Types
- Use custom error types for specific error conditions
- Implement `Error()` method for custom errors
- Use `errors.Is` and `errors.As` for error handling

```go
type ContainerNotFoundError struct {
    ContainerID string
}

func (e ContainerNotFoundError) Error() string {
    return fmt.Sprintf("container %s not found", e.ContainerID)
}

// Usage
if errors.Is(err, &ContainerNotFoundError{}) {
    // Handle container not found
}
```

### Error Context
- Add context to errors when appropriate
- Use `fmt.Errorf` with `%w` verb for wrapping
- Include relevant parameters in error messages

## Performance

### Memory Management
- Avoid unnecessary allocations
- Use object pools for frequently allocated objects

### Concurrency
- Use goroutines appropriately
- Use channels for communication between goroutines
- Separate multi-threading code

## Security

### Input Validation
- Validate all user inputs
- Sanitize data before processing

## Code Review Checklist

Before submitting code for review, ensure:

- [ ] All pre-commit checks pass
- [ ] Code follows formatting standards
- [ ] All tests pass
- [ ] Error handling is appropriate
- [ ] Logging is used instead of fmt.Print*
- [ ] Functions are small and focused
- [ ] Variable names are descriptive
- [ ] Comments explain intent, not implementation
- [ ] No code duplication
- [ ] Security considerations addressed
- [ ] Performance impact considered

## Tools and Automation

### Required Tools
- `pre-commit` - Automation of running all code quality checks
- `golangci-lint` - Code quality checks
- `mockery` - Mock generation
- `make` - Build automation

### Pre-commit Workflow
- Run `make test-all` before committing
- This includes formatting, linting, and testing
- Fix any issues before pushing code

## Getting Help

If you have questions about coding standards:

1. Check this document first
2. Review existing code for examples
3. Open an issue for clarification

Remember: **Good code is readable, maintainable, and follows established patterns. When in doubt, prioritize clarity over cleverness.**
