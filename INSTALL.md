# Install ZShellCheck

ZShellCheck ships as a single static binary, around 2 MB, with no runtime dependencies.
Installation requires no admin rights and provides a clean uninstall on every supported platform.

Skip to your platform: [macOS](#macos) · [Windows](#windows) · [Linux](#linux) · [Cross-platform](#cross-platform).

Verify the install with `zshellcheck -version`.
The full CLI reference lives in [docs/USER_GUIDE.md](docs/USER_GUIDE.md).

> **Roadmap.** `brew install zshellcheck` (homebrew-core), AUR `zshellcheck-bin`, and an `.app`-free Cask are scheduled.
> See [ROADMAP.md → Distribution channels](ROADMAP.md#version-1x--beyond-the-milestone).

---

## macOS

The installer covers Apple Silicon and Intel.

### Recommended: download installer

```bash
curl -fsSL https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.sh | bash
```

What the installer does:

- Resolves the latest GitHub release tag.
- Verifies a matching SHA-256 entry when checksum tools and `checksums.txt` are available.
  Missing verification support produces a warning and installation continues.
- Leaves signature and provenance verification to the [manual procedure](#verifying-a-release-manually).
- Drops the binary into `~/.local/bin/zshellcheck` without `sudo`.
- Installs a man page to `~/.local/share/man/man1` and shell completions to `~/.local/share/zsh/site-functions` and `~/.local/share/bash-completion/completions`.
- Updates the shell `fpath` for completions when the prompt is accepted.

Run as root and the same script installs to `/usr/local/bin` for system-wide use.

### Pin a specific version

```bash
curl -fsSL https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.sh | bash -s -- --version vX.Y.Z
```

Replace `vX.Y.Z` with a tag from [Releases](https://github.com/afadesigns/zshellcheck/releases/latest).

### Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.sh | bash -s -- --uninstall
```

Removes the binary, man page, and completions exactly where the installer wrote them.
No leftover dotfiles.

### Manual archive

```bash
curl -fsSLo zshellcheck.tar.gz \
  "https://github.com/afadesigns/zshellcheck/releases/latest/download/zshellcheck_Darwin_$(uname -m | sed 's/x86_64/x86_64/;s/arm64/arm64/').tar.gz"
tar -xzf zshellcheck.tar.gz
mv zshellcheck ~/.local/bin/
```

`uname -m` returns `arm64` on Apple Silicon and `x86_64` on Intel.
Both are pre-built and signed.
The archive includes the binary, `LICENSE`, `README.md`, `CHANGELOG.md`, the man page, and shell completions.

### Zsh completions from the binary

You can install the completion function without downloading another file:

```zsh
mkdir -p ~/.local/share/zsh/site-functions
zshellcheck --completions > ~/.local/share/zsh/site-functions/_zshellcheck
```

Add the directory to `fpath` before calling `compinit` in your `.zshrc`:

```zsh
fpath=(~/.local/share/zsh/site-functions $fpath)
autoload -Uz compinit
compinit
```

If your configuration already calls `compinit`, place the `fpath` line before that call.
Restart Zsh to load the completion function.
The output is an autoload file for `fpath`; do not source it directly.

---

## Windows

The signed installer covers x64, ARM64, and x86 (i386) without admin rights.

### Recommended: automated installer

```powershell
irm https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.ps1 | iex
```

What the installer does:

- Resolves the latest GitHub release tag.
- Verifies the SHA-256 checksum against `checksums.txt`.
- Verifies the cosign signature when `cosign.exe` is on `$PATH`; the SHA check is the floor when it is not.
- Extracts the binary to `%LOCALAPPDATA%\Programs\zshellcheck\bin\zshellcheck.exe`.
- Adds that directory to the user PATH without admin rights or system-wide changes.
- Confirms the binary runs by calling `zshellcheck -version`.

### Pin a specific version

```powershell
$ErrorActionPreference = 'Stop'
$installer = "$env:TEMP\zshellcheck-install.ps1"
irm https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.ps1 -OutFile $installer
& $installer -Version vX.Y.Z -Yes
```

Replace `vX.Y.Z` with a tag from [Releases](https://github.com/afadesigns/zshellcheck/releases/latest).

### Uninstall

```powershell
$ErrorActionPreference = 'Stop'
$installer = "$env:TEMP\zshellcheck-install.ps1"
irm https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.ps1 -OutFile $installer
& $installer -Uninstall
```

Removes `%LOCALAPPDATA%\Programs\zshellcheck` and the user PATH entry.
Nothing else is touched.

### Manual archive

1. Download the appropriate `zshellcheck_Windows_<arch>.zip` from [Releases](https://github.com/afadesigns/zshellcheck/releases/latest).
2. Verify SHA-256 against `checksums.txt` with `Get-FileHash -Algorithm SHA256`.
3. Extract `zshellcheck.exe` into a directory on `$PATH`.

### Inside WSL

Treat WSL as Linux and pick the matching distribution path below.
The `install.sh` one-liner works inside WSL exactly as on native Linux.

---

## Linux

ZShellCheck ships native packages for every major package manager, plus a tarball fallback.
Each release includes signed `.deb`, `.rpm`, and `.apk` artifacts alongside the tarballs.

### Debian, Ubuntu, derivatives — `.deb`

```bash
arch=$(dpkg --print-architecture)
ver=$(curl -fsSL https://api.github.com/repos/afadesigns/zshellcheck/releases/latest | sed -n 's/.*"tag_name": *"v\([^"]*\)".*/\1/p')
curl -fsSLo zshellcheck.deb \
  "https://github.com/afadesigns/zshellcheck/releases/download/v${ver}/zshellcheck_${ver}_linux_${arch}.deb"
sudo dpkg -i zshellcheck.deb
```

Uninstall:

```bash
sudo apt-get remove zshellcheck
```

### Fedora, RHEL, Rocky, openSUSE — `.rpm`

```bash
arch=$(uname -m)
ver=$(curl -fsSL https://api.github.com/repos/afadesigns/zshellcheck/releases/latest | sed -n 's/.*"tag_name": *"v\([^"]*\)".*/\1/p')
curl -fsSLo zshellcheck.rpm \
  "https://github.com/afadesigns/zshellcheck/releases/download/v${ver}/zshellcheck_${ver}_linux_${arch}.rpm"
sudo rpm -i zshellcheck.rpm
```

Uninstall:

```bash
sudo dnf remove zshellcheck     # Fedora/RHEL/Rocky
sudo zypper remove zshellcheck  # openSUSE
```

### Alpine — `.apk`

```bash
arch=$(uname -m)
ver=$(curl -fsSL https://api.github.com/repos/afadesigns/zshellcheck/releases/latest | sed -n 's/.*"tag_name": *"v\([^"]*\)".*/\1/p')
curl -fsSLo zshellcheck.apk \
  "https://github.com/afadesigns/zshellcheck/releases/download/v${ver}/zshellcheck_${ver}_linux_${arch}.apk"
sudo apk add --allow-untrusted zshellcheck.apk
```

Uninstall:

```bash
sudo apk del zshellcheck
```

### Arch, Manjaro, EndeavourOS — AUR

The `zshellcheck-bin` AUR package is in flight.
See [ROADMAP.md → Distribution channels](ROADMAP.md#version-1x--beyond-the-milestone).
Until it lands, use the universal installer below.

### Universal: download installer

Works on every distribution with `bash`, `curl`, and `tar`.
Skips the package manager.

```bash
curl -fsSL https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.sh | bash
```

Installs to `~/.local/bin/zshellcheck`, or `/usr/local/bin/zshellcheck` when run as root.
Uses the same optional SHA-256 verification as the macOS installer.
Verify signatures and provenance with the manual procedure below.

Pin a version: `bash -s -- --version vX.Y.Z`.
Uninstall: `bash -s -- --uninstall`.

### Manual archive

```bash
arch=$(uname -m)
curl -fsSLo zshellcheck.tar.gz \
  "https://github.com/afadesigns/zshellcheck/releases/latest/download/zshellcheck_Linux_${arch}.tar.gz"
tar -xzf zshellcheck.tar.gz
sudo install -m 0755 zshellcheck /usr/local/bin/
```

---

## Cross-platform

These channels run identically on every supported OS.

### Docker

```bash
docker run --rm -v "$PWD:/work" -w /work ghcr.io/afadesigns/zshellcheck:latest path/to/script.zsh
```

Multi-arch image (`linux/amd64`, `linux/arm64`), `FROM scratch`, around 2 MB, signed with cosign.
Use `--read-only --user 65532:65532` for hardened CI.

### Go toolchain

```bash
go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
```

Builds the latest tagged release into `$GOBIN`.
This channel requires a Go toolchain.
Every other channel is dependency-free.

### GitHub Actions

```yaml
- uses: afadesigns/zshellcheck@latest
  with:
    args: -format sarif -severity warning ./scripts
```

Pin to a tag for reproducibility.
Replace `latest` with the tag from [Releases](https://github.com/afadesigns/zshellcheck/releases/latest); never use `@main`.

### pre-commit

```yaml
# .pre-commit-config.yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: latest
    hooks:
      - id: zshellcheck
```

`pre-commit install` registers the hook.
Subsequent `git commit` invocations run ZShellCheck against staged Zsh files.

---

## Verifying a release manually

Every Releases archive ships with three sibling files: `<archive>.pem` (cosign certificate), `<archive>.sig` (cosign signature), and `checksums.txt` (SHA-256 of every artifact).

```bash
release_tag=vX.Y.Z
cosign verify-blob \
  --certificate zshellcheck_Linux_x86_64.tar.gz.pem \
  --signature   zshellcheck_Linux_x86_64.tar.gz.sig \
  --certificate-identity "https://github.com/afadesigns/zshellcheck/.github/workflows/release-build.yml@refs/tags/$release_tag" \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  zshellcheck_Linux_x86_64.tar.gz
```

Replace `vX.Y.Z` with the intended release tag.
The SHA-256 sum file `checksums.txt` is itself signed (`checksums.txt.pem`, `checksums.txt.sig`).

Release archives, packages, archive SBOMs, and checksums have SLSA build provenance from an isolated signing workflow.
Download `zshellcheck-provenance.json` alongside the archive and verify its signer and source:

```bash
release_tag=vX.Y.Z
release_commit=REVIEWED_COMMIT_SHA
gh attestation verify zshellcheck_Linux_x86_64.tar.gz \
  --bundle zshellcheck-provenance.json \
  --repo afadesigns/zshellcheck \
  --cert-identity "https://github.com/afadesigns/zshellcheck/.github/workflows/release-provenance.yml@refs/tags/$release_tag" \
  --signer-digest "$release_commit" \
  --source-ref "refs/tags/$release_tag" --source-digest "$release_commit" \
  --deny-self-hosted-runners
```

Replace both placeholders with the intended tag and its independently reviewed commit SHA.
The signer workflow must be reviewed at that commit; trusting a checksum alone does not establish the builder's identity.
The release workflow checks every attested file before publishing the GitHub release.
Container images retain separate cosign signatures and are outside this file-provenance guarantee.
See the [attestations index](https://github.com/afadesigns/zshellcheck/attestations) for published statements.

---

## Troubleshooting

**`zshellcheck: command not found` after install.**
Open a new shell.
`~/.local/bin` and `%LOCALAPPDATA%\Programs\zshellcheck\bin` are added to PATH per shell.
Verify the directory is on `$PATH` with `echo $PATH | tr ':' '\n' | grep zshellcheck`.

**`cosign not on PATH — skipping signature verification`.**
The Windows installer verifies signatures when cosign is available.
The POSIX installer requires the manual verification procedure above.

**Apple Silicon vs Intel.**
The installer detects `uname -m`.
Force a binary with `--arch arm64` or `--arch x86_64` if the shell reports a translated value (Rosetta).

**WSL 1 vs WSL 2.**
Both run the Linux tarball.
WSL 1 lacks some `nfpms` post-install scriptlets — use the universal installer or the manual tarball.

Open a [GitHub issue](https://github.com/afadesigns/zshellcheck/issues/new) when a platform is missing or an install path fails.
The goal is reproducibility from a clean machine in under 60 seconds on every channel.
