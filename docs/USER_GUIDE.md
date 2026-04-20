# User Guide

This guide covers the configuration, usage, and troubleshooting of ZShellCheck.

ZShellCheck currently implements **1000 Katas** (checks) covering syntax errors, security issues, performance improvements, and Zsh best practices.

## Table of Contents

- [Severity Levels](#severity-levels)
- [Configuration](#configuration)
- [Integrations](#integrations)
- [Troubleshooting](#troubleshooting)
- [Support](#support)

---

## Severity Levels

Every check is assigned a severity level. Use the `--severity` flag to filter output by minimum severity.

| Level | Description |
| :--- | :--- |
| **error** | Bugs or dangerous constructs that will likely cause incorrect behavior |
| **warning** | Risky patterns that may cause subtle issues or security concerns |
| **info** | Suggestions for improved practices and platform compatibility |
| **style** | Cosmetic or idiomatic improvements for cleaner Zsh code |

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

### Neovim (null-ls)

```lua
local null_ls = require("null-ls")
local zshellcheck = {
    name = "zshellcheck",
    method = null_ls.methods.DIAGNOSTICS,
    filetypes = { "zsh" },
    generator = null_ls.generator({
        command = "zshellcheck",
        args = { "-format", "json", "$FILENAME" },
        format = "json",
        check_exit_code = function(c) return c <= 1 end,
        on_output = function(params) 
             -- parsing logic here
        end,
    }),
}
null_ls.register(zshellcheck)
```

### Pre-commit Hook

Add to `.pre-commit-config.yaml`:

```yaml
-   repo: https://github.com/afadesigns/zshellcheck
    rev: v1.0.0
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

## Support

-   **Discussions**: For Q&A and ideas.
-   **Issues**: For bugs and feature requests.
-   **Security**: Report vulnerabilities privately. See `SECURITY.md`.
