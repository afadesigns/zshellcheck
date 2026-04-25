# User Guide

This guide covers the configuration, usage, and troubleshooting of ZShellCheck.

ZShellCheck currently implements **1000 Katas** (checks) covering syntax errors, security issues, performance improvements, and Zsh best practices.

## Table of Contents

- [CLI Reference](#cli-reference)
- [Severity Levels](#severity-levels)
- [Configuration](#configuration)
- [Integrations](#integrations)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)
- [Support](#support)

---

## CLI Reference

```
zshellcheck [flags] <path> [<path> ...]
```

Paths may be files or directories. Directories are walked recursively; `.go`, `.md`, `.json`, `.yml`, `.yaml`, `.txt`, and hidden directories (anything starting with `.`) are skipped.

| Flag | Default | Purpose |
| --- | --- | --- |
| `-format <text\|json\|sarif>` | `text` | Output format. `sarif` is for GitHub Security / Code Scanning ingestion. |
| `-severity <level[,level...]>` | (all) | Comma-separated filter. Accepts `error`, `warning`, `info`, `style`. |
| `-verbose` | off | Emit full kata descriptions in text output. |
| `-no-color` | off | Disable ANSI colours. Also implied when stdout is not a TTY. |
| `-cpuprofile <path>` | — | Write a Go `pprof` CPU profile to `<path>` for benchmarking. |
| `-fix` | off | Apply auto-fixes in place for katas that declare one. Safe, deterministic rewrites only. |
| `-diff` | off | Preview the fixes as a unified diff instead of writing them. Implies dry-run. |
| `-dry-run` | off | With `-fix`, report what would change without modifying files. |
| `-version` | — | Print the version and exit. |
| `-h` / `--help` | — | Print usage and exit. |

### Exit Codes

| Code | Meaning |
| --- | --- |
| `0` | No violations. |
| `1` | One or more violations found, or a parse error, or a usage error. |

### Examples

```bash
# Lint a single script, text output
zshellcheck ./install.sh

# Lint a tree, silence style-level findings
zshellcheck -severity error,warning,info ./scripts

# Emit SARIF for CI upload
zshellcheck -format sarif -severity warning ./scripts > zshellcheck.sarif

# Preview auto-fixes as a diff
zshellcheck -diff ./scripts

# Apply auto-fixes in place
zshellcheck -fix ./scripts
```

### Auto-fixes

Katas with a deterministic, reversible rewrite ship a `Fix` implementation. Run `zshellcheck -fix <path>` to apply them in place, or `zshellcheck -diff <path>` to preview the result. The fixer only rewrites the exact span the kata points at — arguments, quoting, and surrounding whitespace are preserved byte-for-byte.

Silenced violations (via `.zshellcheckrc` or inline `# zshellcheck disable=…` directives) keep their fixes silenced too.

---

## Severity Levels

Every kata declares a severity. Canonical rubric:

| Level | Go constant | When to use | Example kata |
| --- | --- | --- | --- |
| `error` | `SeverityError` | Code is broken or will crash under Zsh. Output is wrong. | `ZC2000` — `kubectl taint nodes …:NoExecute` |
| `warning` | `SeverityWarning` | Dangerous behaviour; data loss, security risk, or silent subtle bug. | `ZC1136` — `rm -rf $var` without guard |
| `info` | `SeverityInfo` | Works, but brittle or non-portable. Heads-up, not a must-fix. | `ZC1075` — implicit word-splitting reliance |
| `style` | `SeverityStyle` | Convention / idiomatic Zsh. Cosmetic. | `ZC1030` — `echo` vs `print -r --` |

### Filtering by Severity

```bash
# Show only errors
zshellcheck --severity error my_script.zsh

# Show errors and warnings
zshellcheck --severity warning my_script.zsh

# Show everything (default)
zshellcheck --severity style my_script.zsh
```

### Output Formats

- **Text (default)**: Human-readable with ANSI colors and source context. Use `--no-color` to disable colors.
- **JSON**: `zshellcheck -format json file.zsh`
- **SARIF**: `zshellcheck -format sarif file.zsh` (GitHub Security integration)

---

## Configuration

ZShellCheck looks for a file named `.zshellcheckrc` in the current working directory. The file uses **YAML** syntax.

Global settings can be placed in `~/.config/zshellcheck/config.yml` or `${XDG_CONFIG_HOME}/zshellcheck/config.yml`.

### Disabling Katas

To suppress specific checks (Katas), use the `disabled_katas` list:

```yaml
# .zshellcheckrc
disabled_katas:
  - ZC1005 # We prefer 'which' over 'whence'
  - ZC1042 # Ignore specific rule
```

Refer to `KATAS.md` for the list of IDs.

### Inline Disable Directives

Silence katas directly inside a script without touching `.zshellcheckrc`. Comments are recognised in three forms:

```zsh
# Trailing — silence only this line:
rm -rf /tmp/noise  # zshellcheck disable=ZC1136,ZC1075

# Preceding — silence only the next non-blank code line:
# zshellcheck disable=ZC1030
echo "ok"

# File-tail — a directive with no code after it disables the IDs file-wide:
# zshellcheck disable=ZC1092
```

Multiple IDs may be separated by commas or whitespace. IDs disabled inline are merged with any `disabled_katas` from `.zshellcheckrc`.

---

## Integrations

ZShellCheck can be integrated into editors and workflows.

### VS Code (Run on Save)

Install the **Run on Save** extension and add to `settings.json`:

```json
"emeraldwalk.runonsave": {
    "commands": [
        {
            "match": "\\.zsh$",
            "cmd": "zshellcheck ${file}"
        }
    ]
}
```

### Neovim (nvim-lint)

`null-ls` is archived; use [mfussenegger/nvim-lint](https://github.com/mfussenegger/nvim-lint) instead. It knows how to parse ZShellCheck's JSON output natively via a small parser:

```lua
require("lint").linters.zshellcheck = {
    cmd = "zshellcheck",
    stdin = false,
    args = { "-format", "json" },
    stream = "stdout",
    ignore_exitcode = true,
    parser = require("lint.parser").from_errorformat(
        -- fallback to regex on stderr if JSON isn't available
        "%f:%l:%c: %t%m",
        { source = "zshellcheck" }
    ),
}

require("lint").linters_by_ft.zsh = { "zshellcheck" }

vim.api.nvim_create_autocmd({ "BufWritePost", "BufReadPost" }, {
    pattern = { "*.zsh", ".zshrc", ".zshenv" },
    callback = function() require("lint").try_lint() end,
})
```

### Neovim (conform.nvim + LSP — future)

An official LSP is on the roadmap but not yet shipped. Track [ROADMAP.md](../ROADMAP.md) for status.

### Pre-commit Hook

Add to `.pre-commit-config.yaml`:

```yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: v1.0.13
    hooks:
    -   id: zshellcheck
```

---

## Troubleshooting

### Common Issues

1.  **"command not found"**: 
    - Ensure `zshellcheck` is in your `$PATH`. 
    - If you installed as a user, add `$HOME/.local/bin`.
    - If you installed as root, it should be in `/usr/local/bin`.
    - **Fix:** Run `./install.sh` again; it will offer to automatically fix your `$PATH`.
2.  **Parser Errors**: Use `zsh -n` to verify syntax first. Open an issue if valid code fails.
3.  **False Positives**: Disable the Kata via `.zshellcheckrc`.

---

## FAQ

### Why does ZShellCheck error on `${var:-default}`?

The parser doesn't yet handle Zsh / POSIX parameter-expansion modifiers (`:-`, `:=`, `:+`, `:?`, `##`, `%%`, `/pat/rep`, `:offset:length`). Tracked in [#129](https://github.com/afadesigns/zshellcheck/issues/129). Until that lands, wrap the expansion in a guard block or refactor to a temporary variable.

### Should I use ZShellCheck or ShellCheck?

Both. Run ShellCheck for anything targeting `sh` / `bash` portability. Run ZShellCheck for anything using Zsh-only features: parameter-expansion flags (`${(U)x}`, `${(f)x}`), glob qualifiers (`*.zsh(.)`), `[[`, `(( ))`, `print -r --`, modifiers (`:t`, `:h`, `:r`), associative arrays, `setopt` flags, hook functions. See [REFERENCE.md#comparison-vs-shellcheck](REFERENCE.md#comparison-vs-shellcheck).

### How do I exempt one line without editing the whole file?

Add a trailing comment: `some-command # zshellcheck disable=ZC1234`. See [Inline Disable Directives](#inline-disable-directives) above.

### Is there an auto-fixer (`-fix`)?

Yes. Run `zshellcheck -fix path/to/script.zsh` to apply every available rewrite, or `zshellcheck -diff path/to/script.zsh` to preview the unified diff without writing. 67 katas ship deterministic rewrites today; coverage grows each release. The set of fixable katas is listed in [`KATAS.md`](../KATAS.md) — each entry carries an explicit `Auto-fix: yes/no` line.

`-fix` runs multi-pass (up to five) so nested rewrites resolve in a single invocation. Pair `-fix` with `-dry-run` to report what would change without writing.

### The SARIF output is empty after a parse error. Why?

When the parser rejects a file, ZShellCheck exits before katas run — so there is nothing to emit. Fix the syntax first (`zsh -n file.zsh` is a fast sanity check) or open an issue if valid Zsh is being rejected.

### Where does ZShellCheck look for config?

In order, merged with project-local winning:

1. `$XDG_CONFIG_HOME/zshellcheck/config.yml` (or `.yaml`)
2. `~/.config/zshellcheck/config.yml` (or `.yaml`)
3. `~/.zshellcheckrc`
4. `./.zshellcheckrc`

---

## Support

- **Discussions**: https://github.com/afadesigns/zshellcheck/discussions — questions and ideas.
- **Issues**: https://github.com/afadesigns/zshellcheck/issues — bugs, feature requests.
- **Security**: report vulnerabilities privately per [SECURITY.md](../SECURITY.md).
