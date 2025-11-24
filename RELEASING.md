# Releasing Guide

This document outlines the process for releasing a new version of ZShellCheck.

## Prerequisites

- **Goreleaser**: Installed locally or available in CI.
- **Git**: Clean working tree.
- **GitHub Token**: With `repo` permissions (for local release) or configured in GitHub Actions secrets.

## Versioning

ZShellCheck follows [Semantic Versioning](https://semver.org/).

## Release Process

We use [GoReleaser](https://goreleaser.com/) to automate the build and release process. The release is triggered by pushing a semver tag.

1.  **Prepare the Changelog**:
    - Update `CHANGELOG.md` with the new version and release date.
    - Move "Unreleased" changes to the new version section.

2.  **Commit changes**:
    ```bash
    git add CHANGELOG.md
    git commit -m "chore: Prepare release vX.Y.Z"
    git push origin main
    ```

3.  **Tag the Release**:
    Create a lightweight or annotated tag.
    ```bash
    git tag vX.Y.Z
    git push origin vX.Y.Z
    ```

4.  **CI/CD Pipeline**:
    The GitHub Actions workflow (`.github/workflows/release.yml` or similar, depending on setup) will detect the tag and run `goreleaser`.
    
    This will:
    - Cross-compile binaries for Linux, Windows, and macOS (AMD64/ARM64).
    - Create archives (`tar.gz`/`zip`).
    - Generate checksums.
    - Create a GitHub Release.
    - Upload artifacts to the release.

## Manual Release (Testing)

To test the release process locally without publishing:

```bash
goreleaser release --snapshot --clean
```

Binaries and archives will be in the `dist/` directory.

## Post-Release

- Verify the release on the [GitHub Releases](https://github.com/afadesigns/zshellcheck/releases) page.
- Update any external package managers (e.g., Homebrew tap) if not automated.
