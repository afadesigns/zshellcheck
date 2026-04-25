<div align="center">

<img src="docs/assets/banner.png" alt="zshellcheck" width="100%" />

### The quiet linter for a quiet shell.

Static analysis and auto-fix for the setopts, hooks, and globs Bash never learned.

[![CI](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml/badge.svg)](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/afadesigns/zshellcheck?color=blue)](https://github.com/afadesigns/zshellcheck/releases/latest)
[![Marketplace](https://img.shields.io/badge/Marketplace-ZshellCheck%20v1-2ea44f?logo=githubactions&logoColor=white)](https://github.com/marketplace/actions/zshellcheck-v1)
[![Auto-fix](https://img.shields.io/badge/auto--fix-67%20katas-2ea44f)](KATAS.md)
[![Go Report](https://goreportcard.com/badge/github.com/afadesigns/zshellcheck)](https://goreportcard.com/report/github.com/afadesigns/zshellcheck)
[![codecov](https://codecov.io/gh/afadesigns/zshellcheck/graph/badge.svg)](https://codecov.io/gh/afadesigns/zshellcheck)
[![Scorecard](https://api.securityscorecards.dev/projects/github.com/afadesigns/zshellcheck/badge)](https://securityscorecards.dev/viewer/?uri=github.com/afadesigns/zshellcheck)
[![SLSA](https://img.shields.io/badge/SLSA-Level%203-brightgreen)](https://slsa.dev)

[**Install**](INSTALL.md) · [**User guide**](docs/USER_GUIDE.md) · [**Katas**](KATAS.md) · [**Integrations**](INTEGRATIONS.md) · [**Roadmap**](ROADMAP.md) · [**Changelog**](CHANGELOG.md)

</div>

---

## See it in action

<p align="center">
  <img src="docs/assets/demo.gif" alt="zshellcheck demo" width="100%" />
</p>

## Install

**The binary is the same wherever it lands.** Three ways to put it there:

```bash
# macOS, Linux, WSL
curl -fsSL https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.sh | bash
```

```powershell
# Windows
irm https://raw.githubusercontent.com/afadesigns/zshellcheck/main/install.ps1 | iex
```

```bash
# Anywhere Go is installed
go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
```

`--uninstall` reverses any of them.

**On every tag:** native `.deb`, `.rpm`, `.apk`, and multi-arch container → `ghcr.io/afadesigns/zshellcheck`.

Pinning, cosign verification, distro one-liners → [INSTALL.md](INSTALL.md).

## Run

```bash
# Lint
zshellcheck path/to/script.zsh

# Write SARIF for GitHub Code Scanning
zshellcheck -severity warning -format sarif ./scripts > zshellcheck.sarif

# Preview every auto-fix as a unified diff
zshellcheck -diff path/to/script.zsh

# Apply the fixes
zshellcheck -fix path/to/script.zsh
```

Exits *0* on a clean run, *1* when anything was flagged.
`zshellcheck -h` lists every flag, grouped by intent.

Silence inline with `# noka: ZC1234` — bare `# noka` silences every kata on the line.

Three forms (trailing, preceding, file-wide) → [USER_GUIDE](docs/USER_GUIDE.md#inline-noka-directives).

### CI/CD

```yaml
# .github/workflows/lint.yml
- uses: afadesigns/zshellcheck@v1.0.15
  with:
    args: -format sarif -severity warning ./scripts
```

```yaml
# .pre-commit-config.yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: v1.0.15
    hooks:
      - id: zshellcheck
```

## Integrations

Each release runs a parse + lint sweep across the script trees of every Zsh project listed in [INTEGRATIONS.md](INTEGRATIONS.md) — no panics, no crashes, deterministic output. Featured today:

[oh-my-zsh](https://github.com/ohmyzsh/ohmyzsh) ·
[prezto](https://github.com/sorin-ionescu/prezto) ·
[powerlevel10k](https://github.com/romkatv/powerlevel10k) ·
[zinit](https://github.com/zdharma-continuum/zinit) ·
[fzf](https://github.com/junegunn/fzf) ·
[zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting)

19 verified projects today, 300+ targeted before v2 — see the [full matrix](INTEGRATIONS.md).

## Documentation

| Doc | What's inside |
| :--- | :--- |
| [INSTALL.md](INSTALL.md) | Install + uninstall paths for macOS, Windows, Linux, Docker |
| [USER_GUIDE.md](docs/USER_GUIDE.md) | CLI reference, configuration, inline directives, FAQ |
| [KATAS.md](KATAS.md) | Every kata with description, severity, and auto-fix status |
| [INTEGRATIONS.md](INTEGRATIONS.md) | Verified Zsh projects — frameworks, plugins, prompts |
| [DEVELOPER.md](docs/DEVELOPER.md) | Architecture, AST reference, kata authoring, auto-fix catalog |
| [REFERENCE.md](docs/REFERENCE.md) | Governance, glossary, ShellCheck comparison |
| [ROADMAP.md](ROADMAP.md) | LSP, distribution channels, plugin system |
| [CHANGELOG.md](CHANGELOG.md) | Per-release history |

## Contributing

PRs welcome.
Start with [CONTRIBUTING.md](CONTRIBUTING.md).
Questions and ideas: [discussions](https://github.com/afadesigns/zshellcheck/discussions).
Bugs: [issues](https://github.com/afadesigns/zshellcheck/issues).

## License

Distributed under the MIT License.
See [LICENSE](LICENSE).

## Credits

Andreas Fahl (**@afadesigns**) — author and maintainer.
Inspired by [ShellCheck](https://www.shellcheck.net/); ZShellCheck is an independent Go implementation focused on Zsh-specific semantics.

<div align="center">
  <a href="https://github.com/afadesigns/zshellcheck/graphs/contributors">
    <img src="https://contrib.rocks/image?repo=afadesigns/zshellcheck" alt="Contributors" />
  </a>
</div>
