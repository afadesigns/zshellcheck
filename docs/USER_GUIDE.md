# User Guide

This guide covers the configuration, usage, and troubleshooting of ZShellCheck.

## Table of Contents

- [Configuration](#configuration)
- [Integrations](#integrations)
- [Troubleshooting](#troubleshooting)
- [Support](#support)

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
    rev: v0.0.93
    hooks:
    -   id: zshellcheck
```

---

## Troubleshooting

### Common Issues

1.  **"command not found"**: Ensure `$GOPATH/bin` is in your `$PATH`.
2.  **Parser Errors**: Use `zsh -n` to verify syntax first. Open an issue if valid code fails.
3.  **False Positives**: Disable the Kata via `.zshellcheckrc`.

---

## Support

-   **Discussions**: For Q&A and ideas.
-   **Issues**: For bugs and feature requests.
-   **Security**: Report vulnerabilities privately. See `SECURITY.md`.
