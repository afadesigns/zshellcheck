<div align="center">

```
 mmmmmm  mmmm  #             ""#    ""#      mmm  #                    #
     #" #"   " # mm    mmm     #      #    m"   " # mm    mmm    mmm   #   m
   m#   "#mmm  #"  #  #"  #    #      #    #      #"  #  #"  #  #"  "  # m"
  m"        "# #   #  #""""    #      #    #      #   #  #""""  #      #"#
 ##mmmm "mmm#" #   #  "#mm"    "mm    "mm   "mmm" #   #  "#mm"  "#mm"  #  "m
```

**Native static analysis for Zsh.** 1000 Zsh-specific checks covering syntax, security, portability, and style — the counterpart to ShellCheck for code that uses Zsh-only features.

[![CI](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml/badge.svg)](https://github.com/afadesigns/zshellcheck/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/afadesigns/zshellcheck?color=blue)](https://github.com/afadesigns/zshellcheck/releases/latest)
[![Marketplace](https://img.shields.io/badge/Marketplace-ZshellCheck%20v1-2ea44f?logo=githubactions&logoColor=white)](https://github.com/marketplace/actions/zshellcheck-v1)
[![Go Report](https://goreportcard.com/badge/github.com/afadesigns/zshellcheck)](https://goreportcard.com/report/github.com/afadesigns/zshellcheck)
[![codecov](https://codecov.io/gh/afadesigns/zshellcheck/graph/badge.svg)](https://codecov.io/gh/afadesigns/zshellcheck)
[![Scorecard](https://api.securityscorecards.dev/projects/github.com/afadesigns/zshellcheck/badge)](https://securityscorecards.dev/viewer/?uri=github.com/afadesigns/zshellcheck)
[![SLSA](https://img.shields.io/badge/SLSA-Level%203-brightgreen)](https://slsa.dev)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

---

## At a glance

- **1000** Zsh-specific katas (checks) — `error` / `warning` / `info` / `style`
- **3** output formats — text (coloured), JSON, SARIF (GitHub Code Scanning)
- **Zero** runtime dependencies — single static Go binary
- **Signed** releases — cosign keyless + SBOM + SLSA Level 3 provenance
- **Cross-platform** — Linux / macOS / Windows × x86_64 / arm64 / i386

## Quick start

**Install**

```bash
# Automatic — downloads signed binary, or builds if Go is present
./install.sh

# Or via Go toolchain
go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
```

**Run**

```bash
zshellcheck path/to/script.zsh
zshellcheck -severity warning -format sarif ./scripts > zshellcheck.sarif
```

**GitHub Actions**

```yaml
- uses: afadesigns/zshellcheck@v1.0.13
  with:
    args: -format sarif -severity warning ./scripts
```

**Pre-commit**

```yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: v1.0.13
    hooks:
      - id: zshellcheck
```

## Example output

```text
scripts/backup.zsh:14:5: warning: [ZC1136] Avoid `rm -rf $path` without a guard — an empty `$path` deletes `/`.
  rm -rf $path
      ^

scripts/backup.zsh:22:1: style: [ZC1030] Prefer `print -r --` over `echo` for predictable output.
  echo "done"
  ^

Found 2 violations.
```

## Severity

Four levels: `error`, `warning`, `info`, `style`. Filter with `--severity <level>`. Full rubric with examples: [docs/USER_GUIDE.md#severity-levels](docs/USER_GUIDE.md#severity-levels).

## Why ZShellCheck, not ShellCheck?

Run **ShellCheck** for portable `sh` / `bash`. Run **ZShellCheck** for native Zsh — parameter-expansion flags (`${(U)x}`, `${(f)x}`), glob qualifiers (`*.zsh(.)`), `[[`, `(( ))`, `print -r --`, modifiers (`:t`, `:h`, `:r`), associative arrays, `setopt` options, hook functions. Full table: [docs/REFERENCE.md#comparison-vs-shellcheck](docs/REFERENCE.md#comparison-vs-shellcheck).

## Documentation

| Doc | What's inside |
| --- | --- |
| [USER_GUIDE.md](docs/USER_GUIDE.md) | CLI reference, configuration, inline directives, integrations, FAQ |
| [DEVELOPER.md](docs/DEVELOPER.md) | Architecture, AST reference, kata authoring, release process |
| [REFERENCE.md](docs/REFERENCE.md) | Governance, glossary, ShellCheck comparison table |
| [KATAS.md](KATAS.md) | Every kata with description and severity |
| [CHANGELOG.md](CHANGELOG.md) | Per-release history |
| [SECURITY.md](SECURITY.md) | Vulnerability disclosure |
| [CONTRIBUTING.md](CONTRIBUTING.md) | PR workflow, local checks, conventions |
| [ROADMAP.md](ROADMAP.md) | What's next (LSP, auto-fixer, plugins) |

## Contributing

PRs welcome. Start with [CONTRIBUTING.md](CONTRIBUTING.md). Issues and discussions live on GitHub: [issues](https://github.com/afadesigns/zshellcheck/issues), [discussions](https://github.com/afadesigns/zshellcheck/discussions).

## License

MIT. See [LICENSE](LICENSE).

## Credits

- Andreas Fahl (**@afadesigns**) — author and maintainer.
- Inspired by [ShellCheck](https://www.shellcheck.net/) — independent implementation focused on Zsh-specific semantics.

<div align="center">
  <a href="https://github.com/afadesigns/zshellcheck/graphs/contributors">
    <img src="https://contrib.rocks/image?repo=afadesigns/zshellcheck" alt="Contributors" />
  </a>
</div>
