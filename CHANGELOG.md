# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project setup
- Docker container management interface
- Docker image management interface
- Docker volume management interface
- Docker network management interface
- Terminal UI with k9s-style interface
- Command mode for view switching
- Real-time container logs viewing
- Resource inspection capabilities
- Cross-platform build support

### Changed
- None

### Deprecated
- None

### Removed
- None

### Fixed
- None

### Security
- None

## [0.1.0a] - 2025-08-15

### Added
- Alpha release of D5r Docker CLI Dashboard
- Complete Docker management functionality
- Modern terminal UI with tview
- Comprehensive test coverage
- GitHub Actions CI/CD pipeline
- Cross-platform builds (Windows, Linux, macOS)
- Automated releases
- Code quality tools (golangci-lint, govulncheck)
- Comprehensive documentation

---

## Release Process

### Creating a New Release

1. **Update CHANGELOG.md**: Add a new version section with changes
2. **Commit changes**: `git add . && git commit -m "chore: prepare release v1.0.1"`
3. **Create tag**: `git tag v1.0.1`
4. **Push tag**: `git push origin v1.0.1`
5. **GitHub Actions** will automatically:
   - Run tests on all platforms
   - Build release binaries
   - Create GitHub release with assets

### Version Format

- **MAJOR.MINOR.PATCH**
- **MAJOR**: Breaking changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible

### Changelog Categories

- **Added**: New features
- **Changed**: Changes in existing functionality
- **Deprecated**: Soon-to-be removed features
- **Removed**: Removed features
- **Fixed**: Bug fixes
- **Security**: Security-related changes
