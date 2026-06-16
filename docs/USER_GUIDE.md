# User guide

This guide covers configuration, usage, and troubleshooting for ZShellCheck.

ZShellCheck implements a comprehensive set of katas — checks that cover syntax errors, security issues, performance pitfalls, and Zsh idioms.
The full list lives in [KATAS.md](../KATAS.md).

## Contents

- [CLI reference](#cli-reference)
- [Baseline ratchet](#baseline-ratchet)
- [Severity levels](#severity-levels)
- [Configuration](#configuration)
- [Inline `noka` directives](#inline-noka-directives)
- [Integrations](#integrations)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)
- [Support](#support)

---

## CLI reference

```
zshellcheck [flags] <path> [<path> ...]
```

Paths may be files or directories.
Directories are walked recursively.
Files with `.go`, `.md`, `.json`, `.yml`, `.yaml`, or `.txt` extensions are skipped, as are hidden directories.

| Flag | Default | Purpose |
| --- | --- | --- |
| `-format <text\|json\|sarif>` | `text` | Output format. `sarif` is for GitHub Code Scanning ingestion. |
| `-statistics` | off | Print a per-kata count of findings, sorted by frequency, instead of individual reports. |
| `-baseline <path>` | — | Suppress findings recorded in the baseline file; report only findings new since it. |
| `-baseline-write <path>` | — | Write a baseline snapshot of the current findings and exit 0. |
| `-severity <level[,level...]>` | (all) | Comma-separated filter. Accepts `error`, `warning`, `info`, `style`. |
| `-rule-severity <ZC####:level[,...]>` | — | Re-grade specific katas to a chosen severity, for example `ZC1037:error`. |
| `-add-noka` | off | Append a `# noka: ZC####` directive to every line with a finding, write the files, and exit. |
| `-detect-stale-noka` | off | Report `# noka` directives that suppress no actual finding; exit non-zero if any. |
| `-verbose` | off | Emit full kata descriptions in text output. |
| `-no-color` | off | Disable ANSI colours. Implied when stdout is not a TTY. |
| `-cpuprofile <path>` | — | Write a Go pprof CPU profile to `<path>` for benchmarking. |
| `-fix` | off | Apply auto-fixes in place. Safe (value-preserving) fixes only, unless `-unsafe-fixes` is set. |
| `-unsafe-fixes` | off | Also apply fixes that may change runtime behavior — command and flag swaps, scope changes, glob qualifiers. |
| `-diff` | off | Preview the fixes as a unified diff instead of writing them. Implies dry-run. |
| `-dry-run` | off | With `-fix`, report what would change without modifying files. |
| `-list-rules` | — | Print every kata (ID, severity, title) and exit. |
| `-explain <ZC####>` | — | Print one kata's full description and exit. Case-insensitive. |
| `-version` | — | Print the version and exit. |
| `-h`, `--help` | — | Print usage and exit. |

### Exit codes

| Code | Meaning |
| ---: | --- |
| `0` | No violations. |
| `1` | One or more violations, a parse error, or a usage error. |

### Examples

```bash
# Lint a single script with text output
zshellcheck ./install.sh

# Lint a tree, suppress style-level findings
zshellcheck -severity error,warning,info ./scripts

# Emit SARIF for CI upload
zshellcheck -format sarif -severity warning ./scripts > zshellcheck.sarif

# Preview auto-fixes as a diff
zshellcheck -diff ./scripts

# Apply auto-fixes in place
zshellcheck -fix ./scripts

# List every kata, or explain one by ID
zshellcheck -list-rules
zshellcheck -explain ZC1001
```

### Auto-fixes

Katas with a deterministic, reversible rewrite ship a `Fix` implementation.
A finding from such a kata is tagged with a trailing `[*]` in the text report, and the run ends with a `[*] N fixable with the -fix option.` summary, so you can see at a glance how much the fixer resolves before touching anything.

Fixes come in two tiers, following the model Ruff popularised.
A **safe** fix is value-preserving: applying it cannot change runtime behavior or drop a comment, for example `` `cmd` `` to `$(cmd)` or `$arr[i]` to `${arr[i]}`.
An **unsafe** fix may change behavior — a command or flag swap (`which` to `whence`, `netstat` to `ss`), a scope change (adding `local`), or a glob qualifier.
`-fix` applies only safe fixes by default; pass `-unsafe-fixes` to apply the rest.
The report counts each tier separately so you can apply the mechanical fixes confidently and review the behavior-changing ones by hand.

Run `zshellcheck -fix <path>` to apply the safe rewrites in place, `zshellcheck -fix -unsafe-fixes <path>` to apply every fix, or `zshellcheck -diff <path>` to preview the unified diff.
The fixer rewrites only the exact span the kata points at — arguments, quoting, and surrounding whitespace are preserved byte-for-byte.

Silenced violations (via `.zshellcheckrc` or inline `# noka` directives) keep their fixes silenced too.

The fixer runs multi-pass with a default cap of five iterations.
Nested rewrites — for example `` result=`which git` `` collapsing to `result=$(whence git)` — converge in a single invocation.

Combine flags freely:

| Combination | Effect |
| --- | --- |
| `-fix` | Apply safe rewrites to disk. |
| `-fix -unsafe-fixes` | Apply every rewrite, including behavior-changing ones. |
| `-diff` | Print a unified diff. Source unchanged. |
| `-fix -dry-run` | Report which files would change without writing. |
| `-fix -severity warning` | Apply safe rewrites; suppress style-level findings from the human-facing report. |
| `-no-banner -fix` | Apply rewrites without the startup banner — useful in CI. |

[KATAS.md](../KATAS.md) lists every kata with an explicit `Auto-fix: yes/no` line, and the summary table reports the current count.

---

## Baseline ratchet

Adopting ZShellCheck on a large existing codebase need not mean fixing every finding at once.
Snapshot the current findings, then let CI fail only on findings introduced after the snapshot.

```bash
# Record the current findings once and commit the file.
zshellcheck -baseline-write .zshellcheck-baseline ./scripts

# In CI: only findings new since the baseline cause a non-zero exit.
zshellcheck -baseline .zshellcheck-baseline ./scripts
```

A baseline entry identifies a finding by its kata, file, and the trimmed source line — not the line number — so inserting or removing unrelated lines elsewhere in a file does not resurrect a suppressed finding.
Re-run `-baseline-write` to refresh the snapshot after you fix some findings.

## Severity levels

Every kata declares a severity.
The canonical rubric:

| Level | Go constant | When to use | Example |
| --- | --- | --- | --- |
| `error` | `SeverityError` | Code is broken or crashes under Zsh; output is wrong. | `ZC2000` — `kubectl taint nodes …:NoExecute` |
| `warning` | `SeverityWarning` | Dangerous behaviour: data loss, security risk, or silent subtle bug. | `ZC1136` — `rm -rf $var` without guard |
| `info` | `SeverityInfo` | Works, but brittle or non-portable. Heads-up, not a must-fix. | `ZC1005` — `which` vs `whence` |
| `style` | `SeverityStyle` | Convention or idiomatic Zsh. Cosmetic. | `ZC1030` — `echo` vs `print -r --` |

### Filter by severity

```bash
# Errors only
zshellcheck -severity error my_script.zsh

# Errors and warnings
zshellcheck -severity warning my_script.zsh

# Everything
zshellcheck -severity style my_script.zsh
```

### Output formats

- **Text** (default).
  Human-readable, ANSI-coloured, with source context.
  `-no-color` disables colour.
- **JSON.**
  `zshellcheck -format json file.zsh` for tooling and editor integrations.
- **SARIF.**
  `zshellcheck -format sarif file.zsh` for GitHub Code Scanning.

---

## Configuration

ZShellCheck reads `.zshellcheckrc` from the working directory.
The file is YAML.
Global settings live at `~/.config/zshellcheck/config.yml` or `${XDG_CONFIG_HOME}/zshellcheck/config.yml`.

### Disabling katas

Use the `disabled_katas` list to suppress specific checks:

```yaml
# .zshellcheckrc
disabled_katas:
  - ZC1005  # Prefer 'which' over 'whence' in this codebase
  - ZC1042  # Internal exception
```

Refer to [KATAS.md](../KATAS.md) for the full kata list.

---

## Inline `noka` directives

Silence katas inside a script with a `# noka` comment.
No `.zshellcheckrc` edit required.
The bare keyword silences every kata in scope.
The colon-prefixed form narrows to a list:

```zsh
# Trailing — silence specific katas on this line
rm -rf /tmp/noise  # noka: ZC1136, ZC1075

# Trailing — silence every kata on this line
rm -rf /tmp/noise  # noka

# Preceding — applies to the next non-blank code line
# noka: ZC1030
echo "ok"

# File-tail — a directive with no code after it goes file-wide
# noka: ZC1092
```

Multiple IDs may be separated by commas or whitespace.
Inline IDs are merged with `disabled_katas` from `.zshellcheckrc`.

To silence an existing codebase in bulk, `-add-noka` appends a `# noka` directive to every line that carries a finding, then exits — review the diff before committing.
Over time directives go stale as the code around them changes; `-detect-stale-noka` reports any `# noka` that no longer suppresses a finding so you can remove it.

---

## Integrations

ZShellCheck plugs into editors and CI workflows.

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

`null-ls` is archived; use [mfussenegger/nvim-lint](https://github.com/mfussenegger/nvim-lint) instead.
It parses ZShellCheck's JSON output natively:

```lua
require("lint").linters.zshellcheck = {
    cmd = "zshellcheck",
    stdin = false,
    args = { "-format", "json" },
    stream = "stdout",
    ignore_exitcode = true,
    parser = require("lint.parser").from_errorformat(
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

### LSP

An official LSP is on the [roadmap](../ROADMAP.md) but has not shipped.

### pre-commit hook

```yaml
# .pre-commit-config.yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: latest
    hooks:
      - id: zshellcheck
```

Pin `rev` to an exact release tag for reproducible CI.

---

## Troubleshooting

**`command not found`.**
Ensure `zshellcheck` is on `$PATH`.
A user install lives at `$HOME/.local/bin`; a root install lives at `/usr/local/bin`.
Re-running `install.sh` offers to repair `$PATH`.

**Parser errors.**
Run `zsh -n file.zsh` to verify the syntax independently.
Open an issue when valid Zsh code is rejected.

**False positives.**
Silence the kata inline with `# noka: ZCxxxx`, or add it to `disabled_katas` in `.zshellcheckrc`.

---

## FAQ

### Why does ZShellCheck error on `${var:-default}`?

The parser does not yet handle Zsh and POSIX parameter-expansion modifiers (`:-`, `:=`, `:+`, `:?`, `##`, `%%`, `/pat/rep`, `:offset:length`).
Tracked in [#129](https://github.com/afadesigns/zshellcheck/issues/129).
Until the parser lands the modifier set, wrap the expansion in a guard block or refactor to a temporary variable.

### Should I use ZShellCheck or ShellCheck?

Both.
ShellCheck targets `sh` and `bash` portability.
ZShellCheck targets Zsh-specific features: parameter-expansion flags (`${(U)x}`, `${(f)x}`), glob qualifiers (`*.zsh(.)`), `[[`, `(( ))`, `print -r --`, modifiers (`:t`, `:h`, `:r`), associative arrays, `setopt` flags, and hook functions.
See [REFERENCE.md → comparison vs ShellCheck](REFERENCE.md#comparison-vs-shellcheck).

### How do I exempt one line without editing the whole file?

Add a trailing `# noka` comment: `some-command  # noka: ZC1234`.
Bare `# noka` silences every kata on the line.
See [Inline `noka` directives](#inline-noka-directives) above.

### Is there an auto-fixer?

Yes.
Run `zshellcheck -fix path/to/script.zsh` to apply every available rewrite.
Use `-diff` to preview the unified diff without writing.
The set of fixable katas is listed in [KATAS.md](../KATAS.md) — every entry carries an explicit `Auto-fix: yes/no` line, and the summary table reports the count for the current release.

`-fix` runs multi-pass (up to five iterations) so nested rewrites resolve in a single invocation.
Pair `-fix` with `-dry-run` to report what would change without writing.

### The SARIF output is empty after a parse error. Why?

When the parser rejects a file, ZShellCheck exits before katas run; there is nothing to emit.
Fix the syntax (`zsh -n file.zsh` is a fast sanity check), or open an issue when valid Zsh is being rejected.

### Where does ZShellCheck look for config?

In order, with project-local winning:

1. `$XDG_CONFIG_HOME/zshellcheck/config.yml` (or `.yaml`)
2. `~/.config/zshellcheck/config.yml` (or `.yaml`)
3. `~/.zshellcheckrc`
4. `./.zshellcheckrc`

---

## Support

- [Discussions](https://github.com/afadesigns/zshellcheck/discussions) — questions and ideas.
- [Issues](https://github.com/afadesigns/zshellcheck/issues) — bugs and feature requests.
- Vulnerabilities — disclose privately per [SECURITY.md](../SECURITY.md).
