# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0-alpha] - 2025-01-27

### Added
- **SSH Client for Remote Docker**: New SSH client functionality for secure remote Docker connections
- **Comprehensive Testing Infrastructure**: Extensive test coverage for all major components
- **Mock Service Implementations**: Mock interfaces for improved testing capabilities
- **Enhanced UI Testing**: Test infrastructure for UI components and interactions
- **Security Considerations**: Proper security documentation for SSH host key handling
- **Named Return Values**: Improved function readability with named return values
- **Helper Function Library**: Reusable helper functions for common operations

### Changed
- **Docker Client Architecture**: Major refactoring to follow single responsibility principle
- **Function Structure**: Long functions broken down into focused, single-responsibility functions
- **Code Organization**: Better separation of concerns and improved maintainability
- **Type System**: Updated to use modern, non-deprecated Docker API types
- **Error Handling**: Standardized error handling patterns across all components
- **Parameter Optimization**: Resolved unused parameters and huge parameter issues
- **Code Formatting**: Fixed all gofmt issues and improved code consistency

### Fixed
- **All golangci-lint Issues**: Complete resolution of all code quality warnings
- **Deprecated Types**: Fixed usage of deprecated Docker API types
- **Code Quality**: Resolved parameter optimization and naming issues
- **Linting Errors**: All static analysis issues resolved

### Technical Improvements
- **Function Refactoring**: Extracted common logic into reusable helper functions
- **Type Safety**: Updated from `types.Port` to `container.Port`, `types.IDResponse` to `container.ExecCreateResponse`
- **Code Maintainability**: Significantly improved code structure and readability
- **Architecture**: Better separation of concerns and single responsibility principle
- **Future Development**: Solid foundation for continued development and feature additions

### Files Changed
- **Modified**: 50+ files with improvements and refactoring
- **Added**: 20+ new test files and mock implementations
- **New**: SSH client functionality and enhanced testing infrastructure
- **Removed**: Obsolete test utilities and deprecated code

## [0.2.0a] - 2025-08-16

### Added
- **Theme Configuration System**: YAML and JSON theme support with configurable colors for all UI elements
- **Advanced Container Shell**: Interactive shell with command history, tab completion, and multi-line support
- **Container Exec Commands**: Dynamic command execution with pipes, redirects, and shell operators
- **Dedicated Logs View**: Modular logs component for better separation of concerns
- **Enhanced Action Handling**: View-specific actions with dynamic header updates
- **Interactive Command Protection**: Blocks TUI-freezing commands with helpful alternatives
- **Terminal Cleanup**: Proper cleanup with ANSI escape codes and state restoration
- **Image and Network Removal**: Full Docker API integration for removing images and networks
- **Network Metadata Enhancement**: Complete network information including creation time, internal status, and labels
- **Windows Docker Host Detection**: Automatic detection of Docker Desktop pipe paths for Windows users
- **Enhanced Theme Management**: Theme system with multiple preset themes (default, dark, custom)
- **Improved Action Handler Patterns**: Common action handling patterns for consistent resource operations
- **Enhanced Constants Management**: Color and UI constant definitions with theme integration

### Changed
- **UI Architecture**: Refactored logs system into modular LogsView component
- **Header Management**: Dynamic action updates based on current view and mode
- **Shell Command Execution**: Simplified synchronous execution for better reliability
- **Key Binding Management**: Improved focus management and shortcut handling
- **Shutdown Flow**: Proper shutdown signals and graceful cleanup
- **Theme System**: Migrated from hardcoded colors to configurable theme system
- **Action Handling**: Standardized action patterns across all resource types
- **Constants Organization**: Reorganized UI constants with theme-aware alternatives

### Fixed
- Header actions not updating when switching views
- Shell view display and navigation issues
- Command parsing and execution problems
- UI focus management in modals and input fields
- Key binding conflicts between shell and application shortcuts
- Application exit and terminal state cleanup
- Windows Docker connection issues with automatic host detection
- Theme consistency across different UI components

### Technical Improvements
- Dynamic view registry access using reflection
- Enhanced error handling in theme configuration
- Robust command parsing for shell operators
- Tab completion system
- Multi-layer cleanup architecture
- Enhanced signal handling for graceful shutdown
- **Docker Client Enhancement**: Complete field mapping for networks with proper creation time support
- **Service Layer Integration**: Full implementation of image and network removal operations
- **API Completeness**: Leveraging all available Docker API fields instead of partial implementations
- **Cross-Platform Support**: Enhanced Windows compatibility with Docker Desktop detection
- **Theme Configuration**: Multiple theme formats (YAML/JSON) with fallback mechanisms
- **Action Handler Patterns**: Reusable action handling for consistent user experience

## [0.1.0a] - 2025-08-15

### Added
- Initial alpha release of D5r Docker management tool
- Basic container, image, volume, and network management
- TUI interface with tview
- Docker client integration
- Basic operations for Docker resources
- Logs viewing functionality
- Resource inspection capabilities

### Features
- Container lifecycle management (start, stop, restart, delete)
- Resource listing
- Real-time status updates
- Error handling and user feedback
- Cross-platform support (Windows, Linux, macOS)
