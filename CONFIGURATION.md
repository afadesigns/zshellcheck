# Configuration Guide

ZShellCheck is designed to work out-of-the-box with sensible defaults, but it can be customized to fit your project's specific needs using a configuration file.

## The `.zshellcheckrc` File

ZShellCheck looks for a file named `.zshellcheckrc` in the current working directory where you run the command. The configuration file uses **YAML** syntax.

### disabling Katas (Checks)

The primary configuration option is `disabled_katas`, which allows you to suppress specific checks. This is useful if:
- You disagree with a specific rule.
- You have a legacy codebase where fixing a specific issue is not a priority.
- You have a specific use case that requires violating a rule.

**Syntax:**

```yaml
disabled_katas:
  - ZC1001 # Description (optional comment)
  - ZC1002
```

**Example:**

```yaml
# .zshellcheckrc
disabled_katas:
  - ZC1005 # We prefer 'which' over 'whence' in this project
  - ZC1042 # We specifically need to iterate over arguments this way
```

To find the ID of a Kata (`ZCXXXX`), refer to the [KATAS.md](KATAS.md) documentation or the CLI output.

### Future Configuration Options

We plan to add more configuration options in the future, including:
- **`include` / `exclude` paths**: To control which files are analyzed.
- **Severity Levels**: To treat some checks as warnings instead of errors.
- **Custom Globals**: To define global variables that should be ignored by "unused variable" checks.

Check the [Changelog](CHANGELOG.md) for updates.
