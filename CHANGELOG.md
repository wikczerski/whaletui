# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.1] - 2025-09-14

### Added
- **Date-Time-Specific Log Files**: Implemented log file naming with millisecond precision
- **Enhanced Log Organization**: Log files now use format `whaletui-YYYY-MM-DD_HH-MM-SS.mmm.log`
- **Improved Debugging**: Better log file identification for debugging and troubleshooting

### Changed
- **Log File Naming**: Default log files now include timestamp with milliseconds
- **Log Path Generation**: Updated `generateDateSpecificLogPath()` function for date-time format

### Fixed
- **Test Coverage**: Updated `TestDefaultLogPath` to validate new date-time format
- **Backward Compatibility**: Custom log file paths continue to work as before

### Technical Improvements
- **Log Management**: Enhanced log organization for multiple application instances
- **Debugging Support**: Improved debugging capabilities with precise timestamp identification
- **Test Validation**: Added regex pattern validation for date-time log file format

## [0.5.0] - 2025-09-10

### Added
- **Comprehensive Search Functionality**: Real-time filtering and search across all resource types
- **Search Handler**: New SearchHandler with UI-based filtering and search persistence
- **Enhanced SSH Authentication**: Password support and custom SSH key paths for remote connections
- **Search Navigation**: `/` key for quick search access across all views
- **Search State Persistence**: Search state maintained and restored with `/` key

### Changed
- **Search Implementation**: Moved from service-level to UI-level filtering for better performance
- **SSH Client**: Enhanced SSH client with improved authentication methods
- **Navigation Controls**: Updated all views with search functionality and consistent behavior
- **Base View Architecture**: Enhanced base view with search capabilities

### Fixed
- **SSH Authentication**: Improved SSH key authentication test compatibility with CI environments
- **Search Functionality**: Enhanced search for swarm nodes and services
- **UI Consistency**: Consistent search behavior across all resource types

### Technical Improvements
- **Search Performance**: Local UI filtering for faster search results
- **SSH Security**: Enhanced SSH authentication with multiple methods
- **Test Coverage**: Improved test coverage for SSH authentication scenarios
- **Code Organization**: Better separation of search functionality

## [0.4.0] - 2025-09-07

### Added
- **SSH Tunneling as Primary Connection Method**: Implemented SSH tunneling with local TCP port creation on remote machines
- **Comprehensive Column Configuration System**: Added percentage-based column widths with min/max constraints
- **Per-View Column Configurations**: Support for containers, images, volumes, networks, Swarm Nodes, and Swarm Services
- **Column Visibility and Alignment Controls**: Custom column visibility, alignment (left, right, center), and display names

### Changed
- **SSH Connection Architecture**: Removed socat/netcat fallback methods, SSH tunneling is now the only connection method
- **Column Width System**: Replaced fixed character limits with flexible percentage-based width system
- **Navigation Method**: Updated documentation to reflect command mode (`:`) navigation instead of single-letter shortcuts

### Fixed
- **SSH Tunnel Connection Management**: Proper connection management and cleanup for SSH tunnels
- **UI Template Parameter Count**: Fixed template parameter count issues in UI components
- **View Name Case Sensitivity**: Fixed case sensitivity issues in TableFormatter
- **Documentation Accuracy**: Updated all documentation to reflect actual implementation

### Technical Improvements
- **Code Cleanup**: Removed deprecated socat/netcat files and fallback methods
- **Test Coverage**: Updated tests to match new SSH client signature and port range logic
- **Architecture Simplification**: Streamlined connection architecture with single SSH tunneling method

## [0.3.0] - 2025-09-03

### Added
- **Configurable Table Columns**: New character limit configuration for table columns to improve readability
- **Enhanced Configuration System**: Updated theme configuration files with improved structure
- **Swarm Domain Architecture**: Moved swarm views to domains structure for better organization

### Changed
- **Architecture Refactoring**: Major code structure improvements and cleanup
- **Table Column Handling**: Refactored `createTableColumn` into smaller, focused functions
- **Action Keys**: Updated action keys and configuration paths for better consistency

### Fixed
- **Code Organization**: Better separation of concerns and reduced complexity
- **Configuration Paths**: Fixed action keys and config paths for proper functionality

## [0.2.0] - 2025-08-29

### Added
- **Complete Documentation Site**: New Docusaurus-based documentation with comprehensive guides
  - Concept documentation for containers, images, networks, nodes, swarm, and volumes
  - Development guides including coding standards and setup instructions
  - Installation and quick-start guides for new users
  - Professional styling with custom CSS and components
- **Enhanced Navigation**: Backspace and escape key support for navigating back from subviews
- **Improved User Experience**: Better modal management and header handling

### Changed
- **UI Architecture**: Major refactoring of the monolithic `ui.go` (1,468 lines) into 8 focused modules
  - `ui_api.go` - API interface definitions
  - `ui_navigation.go` - Navigation logic
  - `ui_keybindings.go` - Keyboard shortcuts
  - `ui_modals.go` - Modal management
  - `ui_views.go` - View handling
  - `ui_utilities.go` - Helper functions
- **Header Manager**: Simplified architecture with better maintainability and readability
- **Code Quality**: Applied DRY principles, smaller focused functions, better separation of concerns

### Fixed
- **Code Maintainability**: Reduced function complexity
- **Performance**: Optimized header manager operations
- **Architecture**: Better separation of concerns and dependency injection

## [0.1.1] - 2025-08-25

### Fixed
- **Windows Performance**: Updated tcell version improving windows performance

## [0.1.0] - 2025-08-25

### Added
- **Swarm Domain Support**:  Swarm node and service basic implementation added
- **Enhanced UI System**: Improved UI interfaces, managers, and view implementations
- **Docker SSH Implementation**: New dockerssh package for remote Docker connections
- **Error Handling**: Comprehensive error handling and user interaction improvements
- **Mock System**: Complete mockery-generated mock system for all services
- **Build System**: Makefile and build scripts for cross-platform development
- **CI/CD Pipeline**: GitHub Actions workflow for continuous integration
- **Pre-commit Hooks**: Automated code quality checks and formatting

### Changed
- **Codebase Restructure**: Major refactoring to domains structure for better separation of concerns
- **Logger System**: Migrated from custom logger to Go's standard slog package
- **Service Architecture**: Consolidated shared functionality and improved service interfaces
- **UI Components**: Enhanced UI builders, handlers, and managers for better maintainability
- **Docker Client**: Refactored Docker client implementation with improved SSH support
- **Configuration Management**: Enhanced configuration handling with better validation

### Fixed
- **Windows Compatibility**: Logger path validation and Docker host detection for Windows
- **Type Naming**: Resolved stuttering type names throughout the codebase
- **Code Quality**: Resolved all golangci-lint issues and improved code consistency

### Technical Improvements
- **Architecture**: Better separation of concerns and single responsibility principle
- **Code Maintainability**: Significantly improved code structure and readability
- **Mock Generation**: Standardized on mockery-generated mocks for all services
- **Error Handling**: Standardized error handling patterns across all components
- **Future Development**: Solid foundation for continued development and feature additions

### Breaking Changes
- **Logger Interface**: Changed from custom logger to Go's standard slog package
- **Service Structure**: Services moved to domains structure with updated interfaces
- **Mock System**: All mocks now generated by mockery instead of custom implementations

### Migration
- **Logger Usage**: Update logger calls to use slog instead of custom logger
- **Service Imports**: Update import paths for services moved to domains structure
- **Mock Usage**: Use mockery-generated mocks instead of custom mock implementations

## [1.0.0-alpha] - 2025-08-19

### Added
- **Subcommand Architecture**: New `connect` subcommand for remote Docker host connections
- **Required Parameter Enforcement**: Automatic validation of required `--host` and `--user` flags for SSH connections
- **Enhanced SSH Client**: Username parameter support in SSH client for secure remote connections
- **Improved Error Handling**: Clean error messages without automatic help display on connection failures

### Changed
- **CLI Restructure**: Migrated from flag-based remote connection to intuitive subcommand approach
- **Command Structure**:
  - `whaletui` - Local Docker instance (default)
- `whaletui connect --host <host> --user <user>` - Remote Docker host via SSH
- `whaletui theme` - Theme configuration management
- **SSH Connection Flow**: Updated SSH client to accept username from configuration
- **Help Suppression**: Disabled automatic help display on errors for cleaner user experience

### Fixed
- **Help Message Display**: Prevented help messages from appearing when SSH connections fail

### Technical Improvements
- **Configuration Management**: Added `RemoteUser` field to config structure for SSH connections
- **Test Coverage**: Updated test suite to reflect new subcommand architecture

### Breaking Changes
- **CLI Interface**: Remote connections now require `whaletui connect` subcommand instead of `--host` and `--user` flags
- **Required Flags**: Both `--host` and `--user` are now required when using the connect command

### Migration Guide
- **Old Usage**: `whaletui --host 192.168.1.100 --user admin`
- **New Usage**: `whaletui connect --host 192.168.1.100 --user admin`
- **Local Usage**: `whaletui` (unchanged)

## [0.3.0-alpha] - 2025-08-19

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
- Initial alpha release of whaletui Docker management tool
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
