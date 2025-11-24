# Editor Integrations

ZShellCheck can be integrated into your favorite text editor to provide real-time feedback while you write Zsh scripts.

## VS Code

Currently, there is no official VS Code extension for ZShellCheck. However, you can use the **Run on Save** extension or a generic linter extension.

### Using `Run on Save`

1.  Install the [Run on Save](https://marketplace.visualstudio.com/items?itemName=emeraldwalk.RunOnSave) extension.
2.  Add the following to your `.vscode/settings.json`:

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
*Note: This will run the command but the output will be in the output panel, not inline.*

### Planned: Official Extension

We are planning to build an official VS Code extension that implements the Language Server Protocol (LSP) for inline diagnostics.

## Neovim / Vim

### Using `null-ls.nvim` (Neovim)

If you use `null-ls`, you can configure it to use `zshellcheck` as a custom source.

```lua
local null_ls = require("null-ls")
local helpers = require("null-ls.helpers")

local zshellcheck = {
    name = "zshellcheck",
    method = null_ls.methods.DIAGNOSTICS,
    filetypes = { "zsh" },
    generator = helpers.generator_factory({
        command = "zshellcheck",
        args = { "-format", "json", "$FILENAME" },
        format = "json",
        to_stdin = false,
        check_exit_code = function(code)
            return code <= 1
        end,
        on_output = function(params)
            local diags = {}
            for _, violation in ipairs(params.output) do
                table.insert(diags, {
                    row = violation.line,
                    col = violation.column,
                    message = violation.message .. " (" .. violation.kata_id .. ")",
                    severity = vim.diagnostic.severity.WARN,
                    source = "zshellcheck",
                })
            end
            return diags
        end,
    }),
}

null_ls.register(zshellcheck)
```

### Using `ALE` (Vim/Neovim)

You can define a custom linter for ALE:

```vim
call ale#linter#define('zsh', {
\   'name': 'zshellcheck',
\   'executable': 'zshellcheck',
\   'command': 'zshellcheck -format json %t',
\   'callback': 'ale#handlers#unix#HandleAsJSON',
\})
```
*(Note: You might need to adjust the callback to parse the specific JSON structure of ZShellCheck).*

## Pre-commit Hook

The easiest way to integrate ZShellCheck into your workflow is via `pre-commit`. See the [README](README.md#pre-commit-hook-recommended) for setup instructions.

```
