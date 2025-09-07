---
id: intro
title: Introduction
sidebar_label: Introduction
---

# Welcome to WhaleTUI

WhaleTUI is a powerful, terminal-based Docker management tool that simplifies container operations through an intuitive user interface. Built with Go and featuring a modern TUI (Terminal User Interface), WhaleTUI provides comprehensive Docker management capabilities right from your terminal.

## What is WhaleTUI?

WhaleTUI is designed to make Docker container management more accessible and efficient. Whether you're a developer, DevOps engineer, or system administrator, WhaleTUI provides the tools you need to manage your Docker environment effectively.

## Key Features

- **Container Management**: Start, stop, restart, delete, and manage containers with ease
- **Image Management**: Browse, inspect, and manage Docker images
- **Network Management**: View and manage Docker networks
- **Volume Management**: Handle Docker volumes and data persistence
- **Swarm Support**: Manage Docker Swarm services and nodes
- **SSH Integration**: Connect to remote Docker hosts securely
- **Real-time Logs**: Monitor container logs in real-time
- **Interactive Shell**: Built-in shell for advanced operations
- **Theme Support**: Customizable color schemes and UI appearance
- **Column Configuration**: Customize table columns with responsive widths, alignment, and visibility controls

## Why Choose WhaleTUI?

### Simplicity
WhaleTUI provides an intuitive interface that makes Docker operations straightforward, even for complex tasks.

### Performance
Built with Go, WhaleTUI is fast and efficient, handling large numbers of containers and images without performance degradation.

### Flexibility
Support for both local and remote Docker hosts, with SSH integration for secure remote management.

### Extensibility
Modular architecture allows for easy extension and customization of functionality.

## Getting Started

Ready to get started with WhaleTUI? Check out our [Installation Guide](installation.md) to set up WhaleTUI on your system, or jump into the [Quick Start Guide](quick-start.md) for immediate hands-on experience.

## Architecture Overview

WhaleTUI is built with a modular architecture that separates concerns and promotes maintainability:

- **Core Application**: Main application logic and coordination
- **Domain Services**: Business logic for Docker operations (containers, images, networks, volumes, swarm)
- **UI Layer**: Terminal user interface components with keyboard-driven navigation
- **Docker Client**: Integration with Docker Engine API
- **SSH Client**: Remote host connectivity for secure management
- **Theme System**: Customizable appearance and color schemes

## Contributing

WhaleTUI is an open-source project, and we welcome contributions! Whether you're fixing bugs, adding features, or improving documentation, your help is appreciated. See our [Contributing Guide](https://github.com/wikczerski/whaletui/blob/master/CONTRIBUTING.md) for more information.

## License

WhaleTUI is licensed under the MIT License. See the [LICENSE](https://github.com/wikczerski/whaletui/blob/main/LICENSE) file for details.
