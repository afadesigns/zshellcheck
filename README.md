<div align="center">

<img src="docs/assets/banner.png" alt="zshellcheck" width="100%" />

### The quiet linter for a quiet shell.

Static analysis and auto-fix for the setopts, hooks, and globs Bash never learned.

[![CI](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml/badge.svg)](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/afadesigns/zshellcheck?color=blue)](https://github.com/afadesigns/zshellcheck/releases/latest)
[![Marketplace](https://img.shields.io/badge/Marketplace-ZshellCheck%20v1-2ea44f?logo=githubactions&logoColor=white)](https://github.com/marketplace/actions/zshellcheck-v1)
[![Auto-fix](https://img.shields.io/badge/auto--fix-137%20katas-2ea44f)](KATAS.md)
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
- uses: afadesigns/zshellcheck@v1.0.16
  with:
    args: -format sarif -severity warning ./scripts
```

```yaml
# .pre-commit-config.yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: v1.0.16
    hooks:
      - id: zshellcheck
```

## Integrations

| Category | Integrations |
| :--- | :--- |
| Frameworks | [oh-my-zsh](https://github.com/ohmyzsh/ohmyzsh) `497`, [prezto](https://github.com/sorin-ionescu/prezto) `41`, [zephyr](https://github.com/mattmc3/zephyr) `21`, [zsh-utils](https://github.com/belak/zsh-utils) `5`, [zimfw](https://github.com/zimfw/zimfw) `1` |
| Plugin managers | [antidote](https://github.com/mattmc3/antidote) `24`, [zinit](https://github.com/zdharma-continuum/zinit) `9` |
| Plugin tooling | [fzf-tab](https://github.com/Aloxaf/fzf-tab) `5`, [fast-syntax-highlighting](https://github.com/zdharma-continuum/fast-syntax-highlighting) `4`, [fzf](https://github.com/junegunn/fzf) `2` |
| Plugins | [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting) `301`, [zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions) `13`, [zsh-autocomplete](https://github.com/marlonrichert/zsh-autocomplete) `3`, [zsh-history-substring-search](https://github.com/zsh-users/zsh-history-substring-search) `2`, [zsh-vi-mode](https://github.com/jeffreytse/zsh-vi-mode) `2`, [zsh-completions](https://github.com/zsh-users/zsh-completions) `1` |
| Prompts | [spaceship-prompt](https://github.com/spaceship-prompt/spaceship-prompt) `119`, [powerlevel10k](https://github.com/romkatv/powerlevel10k) `16`, [starship](https://github.com/starship/starship) `1` |

## Documentation

**Use it**
- [INSTALL.md](INSTALL.md) — Install + uninstall paths for macOS, Windows, Linux, Docker.
- [USER_GUIDE.md](docs/USER_GUIDE.md) — CLI reference, configuration, inline directives, FAQ.
- [KATAS.md](KATAS.md) — Every kata with description, severity, and auto-fix status.
- [INTEGRATIONS.md](INTEGRATIONS.md) — Verified Zsh projects — frameworks, plugins, prompts.

**Develop with it**
- [DEVELOPER.md](docs/DEVELOPER.md) — Architecture, AST reference, kata authoring, auto-fix catalog.
- [REFERENCE.md](docs/REFERENCE.md) — Governance, glossary, ShellCheck comparison.
- [ROADMAP.md](ROADMAP.md) — LSP, distribution channels, plugin system.
- [CHANGELOG.md](CHANGELOG.md) — Per-release history.

**Contribute**
- [CONTRIBUTING.md](CONTRIBUTING.md) — Workflow, signing requirements, kata standards.
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) — Community expectations.
- [SECURITY.md](SECURITY.md) — Vulnerability disclosure.
- [SUPPORT.md](SUPPORT.md) — Bug, kata, discussion routing.

## Contributing

Contributions of all kinds are welcome — read [CONTRIBUTING.md](CONTRIBUTING.md) to get started.

- Have a question or idea? Join [discussions](https://github.com/afadesigns/zshellcheck/discussions).
- Found a bug? Report it in [issues](https://github.com/afadesigns/zshellcheck/issues).
- Want to write a kata? See the kata-authoring guide in [CONTRIBUTING.md](CONTRIBUTING.md).

## License

ZShellCheck is licensed under the [MIT License](LICENSE).

## Credits

Authored and maintained by Andreas Fahl ([@afadesigns](https://github.com/afadesigns)). Inspired by [ShellCheck](https://www.shellcheck.net/).

<div align="center">

[![Website](https://img.shields.io/badge/Website-afadesign.co-262626?style=flat-square&logo=googlechrome&logoColor=white&labelColor=262626)](https://afadesign.co)
[![GitHub](https://img.shields.io/badge/GitHub-afadesigns-262626?style=flat-square&logo=github&logoColor=white&labelColor=262626)](https://github.com/afadesigns)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-andreasfahl-262626?style=flat-square&logo=linkedin&logoColor=white&labelColor=262626)](https://linkedin.com/in/andreasfahl)
[![Instagram](https://img.shields.io/badge/Instagram-afadesign.official-262626?style=flat-square&logo=instagram&logoColor=white&labelColor=262626)](https://instagram.com/afadesign.official)
[![Facebook](https://img.shields.io/badge/Facebook-andreas.fahl.5-262626?style=flat-square&logo=facebook&logoColor=white&labelColor=262626)](https://facebook.com/andreas.fahl.5)

</div>
