# Install ZShellCheck

ZShellCheck ships as a single static binary (≈ 2 MB).
No runtime dependencies, no admin rights, clean uninstall on every supported platform.

Skip to your platform: [macOS](#macos) · [Windows](#windows) · [Linux](#linux) · [Cross-platform](#cross-platform).

Once installed, verify with `zshellcheck -version`.
The full CLI reference lives in [docs/USER_GUIDE.md](docs/USER_GUIDE.md).

> **Roadmap.** `brew install zshellcheck` (homebrew-core), AUR `zshellcheck-bin`, and an `.app`-free Cask are scheduled — see [ROADMAP.md → Distribution channels](ROADMAP.md#version-1x---beyond-the-milestone).

---

## macOS

The signed installer below covers Apple Silicon and Intel.

### Recommended — automated installer

```bash
curl -fsSL https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.sh | bash
```

What it does:

- Resolves the latest GitHub release tag.
- Verifies the SHA-256 checksum against `checksums.txt`.
- Verifies the cosign signature **if `cosign` is on `$PATH`** (otherwise the SHA check is the floor).
- Drops the binary into `~/.local/bin/zshellcheck` (no `sudo` needed).
- Adds a man page to `~/.local/share/man/man1` and shell completions to `~/.local/share/zsh/site-functions` and `~/.local/share/bash-completion/completions`.
- Updates your shell `fpath` for completions when you accept the prompt.

When run as root, the same script installs to `/usr/local/bin` for system-wide use.

### Pin a specific version

```bash
curl -fsSL https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.sh | bash -s -- --version v1.0.14
```

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

`uname -m` returns `arm64` on Apple Silicon and `x86_64` on Intel — both are pre-built and signed.
The archive includes the binary, `LICENSE`, `README.md`, `CHANGELOG.md`, the man page, and shell completions.

---

## Windows

The signed installer below covers x64, ARM64, and x86 (i386) without admin rights.

### Recommended — automated installer

```powershell
irm https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.ps1 | iex
```

What it does:

- Resolves the latest GitHub release tag.
- Verifies the SHA-256 checksum against `checksums.txt`.
- Verifies the cosign signature **if `cosign.exe` is on `$PATH`** (otherwise the SHA check is the floor).
- Extracts the binary to `%LOCALAPPDATA%\Programs\zshellcheck\bin\zshellcheck.exe`.
- Adds that directory to your **user** PATH (no admin, no system-wide changes).
- Confirms the binary runs by calling `zshellcheck -version`.

### Pin a specific version

```powershell
$ErrorActionPreference = 'Stop'
$installer = "$env:TEMP\zshellcheck-install.ps1"
irm https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.ps1 -OutFile $installer
& $installer -Version v1.0.14 -Yes
```

### Uninstall

```powershell
$ErrorActionPreference = 'Stop'
$installer = "$env:TEMP\zshellcheck-install.ps1"
irm https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.ps1 -OutFile $installer
& $installer -Uninstall
```

Removes `%LOCALAPPDATA%\Programs\zshellcheck` and the user PATH entry.
Nothing else touched.

### Manual archive

1. Download the appropriate `zshellcheck_Windows_<arch>.zip` from [Releases](https://github.com/afadesigns/zshellcheck/releases/latest).
2. Verify SHA-256 against `checksums.txt` (`Get-FileHash -Algorithm SHA256`).
3. Extract `zshellcheck.exe` to a directory on your `$PATH`.

### Inside WSL

Treat as Linux — pick the matching distribution path below.
The `install.sh` one-liner works inside WSL exactly as on native Linux.

---

## Linux

ZShellCheck ships native packages for the major package managers, plus a tarball fallback.
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
sudo dnf remove zshellcheck   # Fedora/RHEL/Rocky
sudo zypper remove zshellcheck # openSUSE
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

The `zshellcheck-bin` AUR package is in flight (see [ROADMAP.md → Distribution channels](ROADMAP.md#version-1x---beyond-the-milestone)).
Until it lands, use the automated installer below.

### Universal — automated installer

Works on every distribution with `bash`, `curl`, and `tar`.
Skips the package manager.

```bash
curl -fsSL https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.sh | bash
```

Installs to `~/.local/bin/zshellcheck` (or `/usr/local/bin/zshellcheck` when run as root).
Same SHA-256 + cosign verification as macOS.

Pin a version: `bash -s -- --version v1.0.14`.
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

Multi-arch image (`linux/amd64`, `linux/arm64`), `FROM scratch`, ≈ 2 MB, signed with cosign.
Use `--read-only --user 65532:65532` for hardened CI.

### Go toolchain

```bash
go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
```

Builds the latest tagged release into `$GOBIN`.
This is the only channel that requires a Go toolchain — every other channel is dependency-free.

### GitHub Actions

```yaml
- uses: afadesigns/zshellcheck@v1.0.14
  with:
    args: -format sarif -severity warning ./scripts
```

The action installs the matching release on the runner.
Pin to a tag, never `@main`, for reproducibility.

### pre-commit

```yaml
# .pre-commit-config.yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: v1.0.14
    hooks:
      - id: zshellcheck
```

`pre-commit install` registers the hook; subsequent `git commit` invocations run ZShellCheck against staged Zsh files.

---

## Verifying a release manually

Every Releases archive ships with three sibling files: `<archive>.pem` (cosign certificate), `<archive>.sig` (cosign signature), and `checksums.txt` (SHA-256 of every artifact).

```bash
# Single-archive verification
cosign verify-blob \
  --certificate zshellcheck_Linux_x86_64.tar.gz.pem \
  --signature   zshellcheck_Linux_x86_64.tar.gz.sig \
  --certificate-identity-regexp 'https://github.com/afadesigns/zshellcheck/.*' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  zshellcheck_Linux_x86_64.tar.gz
```

The SHA-256 sum file `checksums.txt` is itself signed (`checksums.txt.pem`, `checksums.txt.sig`).
SLSA Level 3 build provenance lives at `https://github.com/afadesigns/zshellcheck/attestations`.

---

## Troubleshooting

- **`zshellcheck: command not found` after install.** Open a new shell — `~/.local/bin` and `%LOCALAPPDATA%\Programs\zshellcheck\bin` are added to PATH per shell.
  Verify the directory is on `$PATH` with `echo $PATH | tr ':' '\n' | grep zshellcheck`.
- **`cosign not on PATH — skipping signature verification`.** Optional but recommended.
  Install cosign from <https://github.com/sigstore/cosign>; the installer falls back to SHA-256 only.
- **Apple Silicon vs Intel.** The installer detects `uname -m`.
  Force a binary with `--arch arm64` or `--arch x86_64` if your shell reports a translated value (Rosetta).
- **WSL 1 vs WSL 2.** Both run the Linux tarball.
  WSL 1 lacks some `nfpms` post-install scriptlets — use the universal installer or the manual tarball.

Open a [GitHub issue](https://github.com/afadesigns/zshellcheck/issues/new) if your platform is missing or an install path fails.
We aim to keep every channel reproducible from a clean machine in under 60 seconds.
