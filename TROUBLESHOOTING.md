# Troubleshooting Guide

This guide addresses common issues you might encounter while using ZShellCheck.

## Common Issues

### 1. "zshellcheck: command not found"

**Symptoms:**
- You run `zshellcheck` in the terminal and get an error saying the command is not found.

**Causes:**
- The Go binary directory (`$GOPATH/bin` or `$HOME/go/bin`) is not in your system's `$PATH`.
- The installation failed.

**Solutions:**
- **Add Go bin to PATH:**
  Add this to your `~/.zshrc` or `~/.bashrc`:
  ```bash
  export PATH=$PATH:$(go env GOPATH)/bin
  ```
- **Reinstall:**
  ```bash
  go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest
  ```

### 2. "Parser Error: unexpected token"

**Symptoms:**
- ZShellCheck fails to parse a file that runs correctly in Zsh.
- Errors like `Parser Error in script.zsh: expected next token to be ...`.

**Causes:**
- ZShellCheck's parser might not support a specific, complex, or very new Zsh syntax feature yet.
- The script might actually have a syntax error that Zsh tolerates (or hasn't hit yet).

**Solutions:**
- **Validate with Zsh:** Run `zsh -n script.zsh` to check if Zsh considers it valid.
- **Report a Bug:** If valid in Zsh but fails in ZShellCheck, please open an issue with a minimal reproduction snippet.

### 3. False Positives (Incorrect Flags)

**Symptoms:**
- ZShellCheck reports a violation (e.g., "Use ${} for array access") but your code is actually correct or intentional.

**Solutions:**
- **Disable the Kata:**
  Create or edit `.zshellcheckrc` in your project root:
  ```yaml
  disabled_katas:
    - ZC1001
  ```
- **Inline Disabling (Not Supported Yet):**
  *Note: ZShellCheck does not currently support disabling checks via comments (e.g., `# zshellcheck disable=...`). Please use the configuration file.*

### 4. Integration Issues (VS Code / Vim)

**Symptoms:**
- Linter doesn't show up in the editor.

**Solutions:**
- Ensure `zshellcheck` is in your `$PATH` and accessible by the editor.
- Check the [Integrations Guide](INTEGRATIONS.md) for correct configuration.

## Getting Help

If you're still stuck:
1.  Search [Existing Issues](https://github.com/afadesigns/zshellcheck/issues).
2.  Start a [Discussion](https://github.com/afadesigns/zshellcheck/discussions).
