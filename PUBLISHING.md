# Publishing Polyloft Releases

This guide is for **maintainers** who want to publish new Polyloft releases. For installation instructions, see [README.md](README.md).

## Overview

Publishing a new Polyloft release allows users to install it via:
```bash
go install github.com/ArubikU/polyloft/cmd/polyloft@latest
```

Go uses **git tags** as version identifiers, so publishing is primarily about creating and pushing properly formatted tags.

## Prerequisites

1. **Git Repository**: Code hosted on GitHub (`github.com/ArubikU/polyloft`)
2. **Go Module**: Valid `go.mod` file at repository root
3. **Semantic Versioning**: Understanding of semver (e.g., `v1.0.0`, `v0.1.0`)
4. **Public Repository**: Repository must be publicly accessible
5. **Write Access**: Permission to push tags and create releases

## Release Process

### Step 1: Prepare for Release

Before creating a release, ensure everything is ready:

```bash
# 1. Ensure all changes are committed
git status

# 2. Run all tests
go test ./...

# 3. Test the CLI locally
go run ./cmd/polyloft version
go run ./cmd/polyloft build -o test_build

# 4. Update CHANGELOG.md (if exists) with new changes
```

### Step 2: Create and Push Git Tag

Git tags are how Go identifies versions:

```bash
# Create a new tag (replace with your version)
git tag v0.1.0

# Push the tag to GitHub
git push origin v0.1.0
```

**Tag Naming Rules:**
- Must start with `v` (e.g., `v1.0.0`, not `1.0.0`)
- Follow semantic versioning: `vMAJOR.MINOR.PATCH`
- Use pre-release tags for testing: `v0.1.0-beta`, `v0.1.0-rc1`

### Step 3: Create GitHub Release (Recommended)

1. Go to: https://github.com/ArubikU/polyloft/releases/new
2. Select the tag you just created
3. Add release title (e.g., "Release v0.1.0")
4. Add release notes describing changes
5. Click "Publish release"

### Step 4: Verify Module is Discoverable

Wait a few minutes, then verify:

```bash
# Check if Go can find your module
go list -m github.com/ArubikU/polyloft@latest

# Check specific version
go list -m github.com/ArubikU/polyloft@v0.2.6c

# View all available versions
go list -m -versions github.com/ArubikU/polyloft
```

### Step 5: Test Installation

Test from a clean environment (different machine or directory):

```bash
# Install the new version
go install github.com/ArubikU/polyloft/cmd/polyloft@latest

# Verify it works
polyloft version

# Test basic functionality
mkdir test-project
cd test-project
polyloft init
polyloft build -o test
```
> On Windows the default build artifact uses the `.pfx` extension so it runs like a native executable.

### Step 6: Update Documentation

After successful release:

1. Verify pkg.go.dev has indexed the new version: https://pkg.go.dev/github.com/ArubikU/polyloft
2. Update README.md if installation instructions changed
3. Announce the release (if applicable)

## Quick Release Checklist

- [ ] All changes committed and pushed
- [ ] All tests passing: `go test ./...`
- [ ] CHANGELOG.md updated (if exists)
- [ ] Version bumped appropriately (major/minor/patch)
- [ ] Git tag created: `git tag vX.Y.Z`
- [ ] Tag pushed: `git push origin vX.Y.Z`
- [ ] GitHub release created with notes
- [ ] Module verified: `go list -m github.com/ArubikU/polyloft@vX.Y.Z`
- [ ] Installation tested: `go install github.com/ArubikU/polyloft/cmd/polyloft@vX.Y.Z`
- [ ] pkg.go.dev indexed (wait ~5 minutes)

## Semantic Versioning Guide

Follow [semver](https://semver.org/) for version numbers:

**Major version (v1.0.0 → v2.0.0)** - Breaking changes:
- Changes to CLI command structure
- Incompatible language syntax changes
- API changes that break existing code

**Minor version (v1.0.0 → v1.1.0)** - New features:
- New commands (e.g., `polyloft test`)
- New language features
- Backwards-compatible enhancements

**Patch version (v1.0.0 → v1.0.1)** - Bug fixes:
- Bug fixes
- Documentation updates
- Performance improvements
- Security patches

## Troubleshooting Publishing Issues

### Module Not Found After Tag Push

**Problem**: `go list -m github.com/ArubikU/polyloft@vX.Y.Z` fails

**Solutions**:
1. Wait 2-5 minutes for Go proxy to sync
2. Verify tag exists on GitHub: `git ls-remote --tags origin`
3. Check tag format starts with `v`: `v1.0.0` not `1.0.0`
4. Verify repository is public

### Users Report Old Version

**Problem**: Users run `@latest` but get old version

**Solutions**:
1. Check latest tag: `git describe --tags --abbrev=0`
2. Verify tag is pushed: `git ls-remote --tags origin`
3. Clear Go proxy cache (users): `go clean -modcache`
4. Check Go proxy: https://proxy.golang.org/github.com/!arubik!u/polyloft/@v/list

### pkg.go.dev Not Updating

**Problem**: Documentation on pkg.go.dev is outdated

**Solutions**:
1. Visit: https://pkg.go.dev/github.com/ArubikU/polyloft
2. Click "Request" to re-index
3. Or trigger manually: `GOPROXY=https://proxy.golang.org go list -m github.com/ArubikU/polyloft@vX.Y.Z`

## Post-Release Tasks

1. **Monitor**:
   - Watch for installation issues on GitHub
   - Check pkg.go.dev renders correctly
   - Monitor Go proxy for version availability

2. **Announce** (optional):
   - Create release announcement
   - Update social media
   - Notify users in relevant channels

3. **Maintain**:
   - Respond to issues
   - Plan next release
   - Keep dependencies updated: `go get -u ./...`

## Complete Release Example

Here's a complete workflow for releasing v0.2.0:

```bash
# 1. Verify everything is ready
git status                    # Should be clean
go test ./...                 # All tests passing

# 2. Commit any final changes
git add .
git commit -m "Prepare release v0.2.0"
git push

# 3. Create and push tag
git tag v0.2.0
git push origin v0.2.0

# 4. Create GitHub release at:
# https://github.com/ArubikU/polyloft/releases/new

# 5. Wait 2-5 minutes, then verify
go list -m github.com/ArubikU/polyloft@v0.2.0

# 6. Test installation in clean environment
cd /tmp
go install github.com/ArubikU/polyloft/cmd/polyloft@v0.2.0
polyloft version                # Should show v0.2.0

# 7. Done! Users can now install with @latest
```

## Quick Reference

**Create release:**
```bash
git tag vX.Y.Z && git push origin vX.Y.Z
```

**Verify release:**
```bash
go list -m github.com/ArubikU/polyloft@vX.Y.Z
```

**Test installation:**
```bash
go install github.com/ArubikU/polyloft/cmd/polyloft@vX.Y.Z
```

**View all versions:**
```bash
go list -m -versions github.com/ArubikU/polyloft
```

For user installation instructions, see [README.md](README.md).
