# Polyloft Roadmap

This document outlines planned features and enhancements for the Polyloft build system and ecosystem.

## Phase 1: Core Build System ✅ (Completed 100%)

- [x] TOML-based configuration (polyloft.toml)
- [x] `polyloft init` command for project initialization
- [x] `polyloft build` command to compile to executables
- [x] `polyloft install` command for dependency management
- [x] Support for Go and Polyloft library dependencies
- [x] Basic documentation and examples
- [x] `polyloft run` with no file mention will try to run the current directory as a project
- [x] register .pfx extension on linux/macOS desktop and mime types for file associations

## Phase 2: Dependency Server & Registry ✅ (Completed 100%)

### Package Registry Infrastructure
Leave the server/api on a folder called /server , later we will create a separate repo for it. and on the polyloft ENV variable we will set the URL to point to it.
adding to the cli the auth system, and like polyloft publish will require user to auth.
later polyloft install will also allow user to install online packages from the registry.
like `polyloft install vectores@arubiku` but if the lib is a go module install it with
`polyloft install github.com/arubiku/vectores`

- [x] **Centralized Package Registry**
  - [x] Complete /server directory structure with full implementation
  - [x] RESTful API endpoints (health, packages, search, auth)
  - [x] In-memory database for package storage
  - [x] Full HTTP server with CORS support
  - [x] Standalone server command (server/cmd/main.go)
  - [ ] Web interface for browsing packages (future enhancement)
  - [ ] Package statistics and metrics (future enhancement)

- [x] **Package Identification System**
  - [x] Packages identified as `{package}@{author}` (npm-style)
  - [x] Version management: `{package}@{author}@{version}`
  - [x] Latest version retrieval when version not specified
  - [x] Support for package@author syntax in installer

- [x] **Library Metadata (libs/ structure)**
  - [x] Package metadata in polyloft.toml
  - [x] Metadata includes:
    - Package name
    - Current version
    - Author/organization
    - Description
    - Entry point
    - License
  - [x] Package archiving with tar.gz compression
  - [x] Checksum calculation (SHA-256)

### Package Operations

- [x] **Upload and Publish**
  - [x] `polyloft publish` command with full implementation
  - [x] Complete package validation before publishing
  - [x] File packaging with tar.gz compression
  - [x] Checksum generation for package integrity
  - [x] Automated file collection (.pf files)
  - [x] Binary upload to server with Base64 encoding
  - [x] Server storage of package binaries

- [x] **Search and Discovery**
  - [x] `polyloft search <query>` command
  - [x] Search by package name, author, and description
  - [x] Package listing with pagination
  - [x] Direct package retrieval by name@author
  - [x] Package download endpoint (`/api/download/`)
  - [x] Automatic extraction of downloaded packages

- [x] **User Management**
  - [x] `polyloft register` command for account creation
  - [x] `polyloft login` command for authentication
  - [x] `polyloft logout` command to clear credentials
  - [x] Local credential storage (~/.polyloft/credentials.json)
  - [x] POLYLOFT_REGISTRY_URL environment variable support
  - [x] Full authentication system in /server/auth
  - [x] Server-side user account creation (registration)
  - [x] Server-side authentication with password hashing
  - [x] Bearer token authentication for API
  - [x] Token expiration (30-day default)
  - [x] Interactive CLI registration workflow

### Version Management

- [x] **Semantic Versioning Support**
  - [x] Version storage and retrieval
  - [x] Latest version detection
  - [x] Specific version retrieval
  - [ ] Automatic semver validation (future enhancement)
  - [ ] Breaking change detection (future enhancement)
  - [ ] Version compatibility checks (future enhancement)
  - [ ] Deprecation warnings (future enhancement)

- [ ] **Conflict Resolution** (Future)
  - [ ] Dependency graph analysis
  - [ ] Automatic conflict detection
  - [ ] Resolution strategies (latest-compatible, strict, etc.)
  - [ ] Manual override options

- [ ] **Lock Files** (Future)
  - [ ] `polyloft.lock` for reproducible builds
  - Pin exact versions of all dependencies
  - Integrity checksums
  - Platform-specific locks

## Phase 3: Advanced Build Options 🔮 (Future)

### Cross-Compilation

- [ ] Support for multiple target platforms
  - `polyloft build --target linux/amd64` --include-libraries --include-source
  - `polyloft build --target windows/amd64` --include-libraries --include-source
  - `polyloft build --target darwin/arm64` --include-libraries --include-source
- [ ] Better platform and selection of include/exclude options
- [ ] Platform-specific builds with conditional code
- [ ] Cross-compilation for embedded systems

### Include/Exclude Options
  - [ ] **Selective Inclusion**
    - `--include-libraries` - Include all dependent libraries in the build
      Include all the libraries used by the project in the final build output. 
      This makes the executable self-contained and portable, not requiring external libraries at runtime.
    - `--include-source` - Include source files for debugging 
      When the source is not included, only compiled bytecode is packaged
      so it depends 100% on the installed libraries at runtime. 
      and the user to have polyloft installed with the same libraries.
      This makes the file bigger, and platform dependent.
    - `--exclude-tests` - Exclude test files from the build
      Future option to exclude test files from the build output.
    - `--exclude-docs` - Exclude documentation files
      Future option to exclude documentation files from the build output.

  This means users can choose to create fully self-contained executables 
  or just "zips" of the hycode to run using the polyloft installation on their system.
### Build Optimization

- [ ] **Optimization Flags**
  - `-O` levels for different optimization strategies
  - Dead code elimination
  - Inlining optimizations
  - Size vs. speed tradeoffs

- [ ] **Debug Information**
  - `--strip` flag to remove debug symbols
  - Source maps for debugging
  - Profiling hooks
  - Debug vs. release builds

- [ ] **Custom Build Scripts**
  - Pre-build and post-build hooks
  - Custom compilation steps
  - Integration with build tools (make, ninja, etc.)
  - Plugin system for extensibility

### Packaging & Distribution

- [ ] Generate installers (`.deb`, `.rpm`, `.msi`)
- [ ] Docker image generation
- [ ] Static binary generation
- [ ] Library packaging (`.a`, `.so`, `.dll`)

## Phase 4: Package Management Features 🔮 (Future)

### Enhanced Dependency Management

- [ ] **Private Package Repositories**
  - Self-hosted registry support
  - Private package hosting
  - Access control and permissions
  - Mirror support for offline builds

- [ ] **Workspace Support**
  - Multi-package projects
  - Shared dependencies across packages
  - Monorepo support
  - Workspace-aware builds

- [ ] **Dependency Auditing**
  - Security vulnerability scanning
  - License compliance checking
  - Dependency update notifications
  - Automated dependency updates (like Dependabot)

### Developer Experience

- [ ] **Interactive Dependency Manager**
  - `polyloft add <package>` - Interactive package selection
  - `polyloft update` - Interactive update wizard
  - `polyloft audit fix` - Automated security fixes

- [ ] **Package Templates**
  - `polyloft init --template <name>` - Start from templates
  - Community-maintained templates
  - Custom template repositories

## Phase 5: Ecosystem Integration 🔮 (Future)

### IDE & Editor Support

- [ ] Language Server Protocol (LSP) integration
- [ ] VS Code extension
- [ ] IntelliJ/IDEA plugin
- [ ] Vim/Emacs packages
- [ ] Auto-completion for dependencies

### CI/CD Integration

- [ ] GitHub Actions integration
- [ ] GitLab CI templates
- [ ] Jenkins plugins
- [ ] Automated testing for packages
- [ ] Coverage reporting

### Documentation

- [ ] Auto-generated API documentation
- [ ] Package documentation hosting
- [ ] Interactive examples
- [ ] Tutorial system

## Publishing and Distribution

For users to install Polyloft globally, the project must be published with proper Git tags:

1. **Create Git Tags**: `git tag v0.1.0 && git push origin v0.1.0`
2. **Go Module System**: Automatically makes the module available via `go install`
3. **Users Install**: `go install github.com/ArubikU/polyloft/cmd/polyloft@latest`

See [PUBLISHING.md](PUBLISHING.md) for complete publishing guide for developers.

## Implementation Notes

### Current Limitations

1. **Module Availability**: The build system currently requires the polyloft module to be available locally (development mode) or via Go module proxy (production mode). 

2. **Installation**: Users need both:
   - Go installed on their system (version 1.22.0+)
   - polyloft CLI (install via `go install github.com/ArubikU/polyloft/cmd/polyloft@latest`)

### Development vs. Production Modes

**Development Mode** (current):
- Uses local polyloft module via replace directive
- Searches for go.mod in parent directories
- Suitable for development and testing

**Production Mode** (planned):
- Polyloft module published to Go module proxy
- Users install via `go install`
- No local source code required
- Works with just CLI binary

### Migration Path

To make the build system work with just the CLI binary:

1. **Short-term**: Document installation via `go install`
2. **Medium-term**: Publish polyloft module to pkg.go.dev
3. **Long-term**: Create fully self-contained CLI that bundles runtime

## Contributing

This roadmap is a living document. Contributions and feedback are welcome! 

To suggest features or changes:
1. Open an issue for discussion
2. Submit a PR to update this roadmap
3. Join community discussions

## Timeline

- **Q4 2024**: Phase 1 completion ✅
- **Q1 2025**: Phase 2 planning and infrastructure
- **Q2-Q3 2025**: Phase 2 implementation
- **Q4 2025**: Phase 3-4 planning
- **2026+**: Phases 3-5 implementation

Note: Timeline is tentative and subject to change based on community needs and contributions.
