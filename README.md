<div align="center">

<img src="docs/assets/banner.png" alt="ZShellCheck" width="100%" />

### The quiet linter for a quiet shell.

Static analysis and auto-fix for the setopts, hooks, and globs Bash never learned.

[![CI](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml/badge.svg)](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml)
[![Release](https://img.shields.io/badge/release-v1.6.0-blue)](https://github.com/afadesigns/zshellcheck/releases/latest)
[![Marketplace](https://img.shields.io/badge/Marketplace-ZshellCheck%20v1-2ea44f?logo=githubactions&logoColor=white)](https://github.com/marketplace/actions/zshellcheck-v1)
[![Auto-fix](https://img.shields.io/badge/auto--fix-132%20katas-2ea44f)](KATAS.md)
[![Dependencies](https://img.shields.io/badge/dependencies-0-2ea44f)](go.mod)
[![Go Report](https://goreportcard.com/badge/github.com/afadesigns/zshellcheck)](https://goreportcard.com/report/github.com/afadesigns/zshellcheck)
[![codecov](https://codecov.io/gh/afadesigns/zshellcheck/graph/badge.svg)](https://codecov.io/gh/afadesigns/zshellcheck)
[![Scorecard](https://api.securityscorecards.dev/projects/github.com/afadesigns/zshellcheck/badge)](https://securityscorecards.dev/viewer/?uri=github.com/afadesigns/zshellcheck)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/12657/badge?v=gold)](https://www.bestpractices.dev/projects/12657)
[![SLSA](https://img.shields.io/badge/SLSA-Level%203-brightgreen)](https://slsa.dev)

[**Install**](INSTALL.md) · [**User guide**](docs/USER_GUIDE.md) · [**Katas**](KATAS.md) · [**Integrations**](INTEGRATIONS.md) · [**Roadmap**](ROADMAP.md) · [**Changelog**](CHANGELOG.md)

</div>

---

## See it in action

<p align="center">
  <img src="docs/assets/demo.gif" alt="ZShellCheck demo" width="100%" />
</p>

## Install

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

Native `.deb`, `.rpm`, `.apk`, and a multi-arch container at `ghcr.io/afadesigns/zshellcheck` ship on every release tag.

Pinning, cosign verification, and distro one-liners are in [INSTALL.md](INSTALL.md).

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

Exits `0` on a clean run, `1` when anything was flagged.
`zshellcheck -h` lists every flag, grouped by intent.

Silence inline with `# noka: ZC1234`.
Bare `# noka` silences every kata on the line.
Trailing, preceding, and file-wide forms are documented in [USER_GUIDE.md](docs/USER_GUIDE.md#inline-noka-directives).

### CI/CD

The published action checks out your repository, installs a signed release binary, runs it, and fails the job on any finding.
Add the SARIF upload to surface results in the repository Security tab:

```yaml
# .github/workflows/lint.yml
name: zshellcheck
on: [push, pull_request]
permissions:
  contents: read
  security-events: write
jobs:
  zshellcheck:
    runs-on: ubuntu-latest
    steps:
      - uses: afadesigns/zshellcheck@latest
        with:
          args: -format sarif -severity warning ./scripts > zshellcheck.sarif
      - uses: github/codeql-action/upload-sarif@v3
        if: always()
        with:
          sarif_file: zshellcheck.sarif
```

Run it as a pre-commit hook instead:

```yaml
# .pre-commit-config.yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: latest
    hooks:
      - id: zshellcheck
```

Pin `@latest` and `rev: latest` to a tag from [Releases](https://github.com/afadesigns/zshellcheck/releases/latest) for reproducible CI.

## Integrations

ZShellCheck is verified against widely used Zsh frameworks, plugin managers, plugins, and prompts on every release.
Each runs a parse-and-findings sweep: zero parser errors, zero crashes, and kata findings locked to a reviewed baseline.
The full catalog with file counts lives in [INTEGRATIONS.md](INTEGRATIONS.md).

| Category | Examples |
| :--- | :--- |
| Frameworks | [oh-my-zsh](https://github.com/ohmyzsh/ohmyzsh), [prezto](https://github.com/sorin-ionescu/prezto), [prezto-contrib](https://github.com/belak/prezto-contrib), [zephyr](https://github.com/mattmc3/zephyr), [zimfw](https://github.com/zimfw/zimfw) |
| Plugin managers | [antidote](https://github.com/mattmc3/antidote), [zinit](https://github.com/zdharma-continuum/zinit) |
| Plugins | [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting), [zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions), [zsh-autocomplete](https://github.com/marlonrichert/zsh-autocomplete), [atuin](https://github.com/atuinsh/atuin), [zsh-help](https://github.com/Freed-Wu/zsh-help) |
| Prompts | [powerlevel10k](https://github.com/romkatv/powerlevel10k), [spaceship-prompt](https://github.com/spaceship-prompt/spaceship-prompt), [starship](https://github.com/starship/starship), [gitstatus](https://github.com/romkatv/gitstatus) |
| Tooling | [fzf](https://github.com/junegunn/fzf), [fzf-tab](https://github.com/Aloxaf/fzf-tab), [fast-syntax-highlighting](https://github.com/zdharma-continuum/fast-syntax-highlighting) |

## Quality

Every release replays the linter over the pinned integration corpora and gates on two snapshots:

- Parser errors and crashes stay at zero.
- Kata findings match a reviewed baseline; a new finding on known-good code fails the build as a candidate false positive.

Semantic-preserving rewrites — added blank lines, comments, or variable renames — must not change which katas fire.
See the [local checks](CONTRIBUTING.md#local-checks) for the commands.

## Documentation

**Use it**

- [INSTALL.md](INSTALL.md) — install and uninstall paths for macOS, Windows, Linux, and Docker.
- [USER_GUIDE.md](docs/USER_GUIDE.md) — CLI reference, configuration, inline directives, FAQ.
- [KATAS.md](KATAS.md) — every kata with description, severity, and auto-fix status.
- [INTEGRATIONS.md](INTEGRATIONS.md) — verified Zsh frameworks, plugins, and prompts.

**Develop with it**

- [DEVELOPER.md](docs/DEVELOPER.md) — architecture, AST reference, kata authoring, auto-fix catalog.
- [REFERENCE.md](docs/REFERENCE.md) — governance, glossary, ShellCheck comparison.
- [ROADMAP.md](ROADMAP.md) — LSP, distribution channels, plugin system.
- [CHANGELOG.md](CHANGELOG.md) — per-release history.

**Contribute**

- [CONTRIBUTING.md](CONTRIBUTING.md) — workflow, signing requirements, kata standards.
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) — community expectations.
- [SECURITY.md](SECURITY.md) — vulnerability disclosure.
- [SUPPORT.md](SUPPORT.md) — bug, kata, and discussion routing.

## Contributing

Contributions of all kinds are welcome.
Start with [CONTRIBUTING.md](CONTRIBUTING.md).

- A question or idea? Open a [discussion](https://github.com/afadesigns/zshellcheck/discussions).
- A bug? File an [issue](https://github.com/afadesigns/zshellcheck/issues).
- A new kata? See the kata-authoring guide in [CONTRIBUTING.md](CONTRIBUTING.md).

## License

ZShellCheck is licensed under the [MIT License](LICENSE).

## Credits

Authored and maintained by Andreas Fahl ([@afadesigns](https://github.com/afadesigns)).
Inspired by [ShellCheck](https://www.shellcheck.net/).

<div align="center">

[![Website](https://img.shields.io/badge/Website-afadesign.co-262626?style=flat-square&logo=googlechrome&logoColor=white&labelColor=262626)](https://afadesign.co)
[![GitHub](https://img.shields.io/badge/GitHub-afadesigns-262626?style=flat-square&logo=github&logoColor=white&labelColor=262626)](https://github.com/afadesigns)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-andreasfahl-262626?style=flat-square&logo=linkedin&logoColor=white&labelColor=262626)](https://linkedin.com/in/andreasfahl)
[![Instagram](https://img.shields.io/badge/Instagram-afadesign.official-262626?style=flat-square&logo=instagram&logoColor=white&labelColor=262626)](https://instagram.com/afadesign.official)
[![Facebook](https://img.shields.io/badge/Facebook-andreas.fahl.5-262626?style=flat-square&logo=facebook&logoColor=white&labelColor=262626)](https://facebook.com/andreas.fahl.5)

</div>
