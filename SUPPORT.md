# Support

Thanks for using ZShellCheck. The fastest path to the right answer is below.

## Where to file what

| You have… | Open it as… |
| --- | --- |
| A **bug** (panic, crash, false positive, false negative, parser failure on valid Zsh) | [GitHub issue](https://github.com/afadesigns/zshellcheck/issues/new) — include the binary version (`zshellcheck -version`), a minimal repro, and the full stderr including the banner and any stack trace. |
| A **kata request** (Zsh anti-pattern that ZShellCheck does not yet catch) | [GitHub issue](https://github.com/afadesigns/zshellcheck/issues/new) — describe the pattern, why it bites, and a code sample that should trip the new rule. |
| A **question or design idea** | [GitHub Discussions](https://github.com/afadesigns/zshellcheck/discussions). Discussions are public and searchable; use issues only when there is something to fix. |
| A **security vulnerability** | Do **not** file a public issue. See [SECURITY.md](SECURITY.md) for the private disclosure flow. |
| A **documentation gap** (typo, broken link, stale fact, missing example) | A small PR is more useful than an issue here. The [Contributing guide](CONTRIBUTING.md#pull-request-workflow) shows the workflow; doc-only PRs skip most of the gates. |

## Before you open a bug

1. Run on the **latest tagged release** if you can — `go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest`.
2. Re-run with `-no-banner -severity error` to confirm the issue is not silenced by the noise filter.
3. Check the [issue tracker](https://github.com/afadesigns/zshellcheck/issues?q=is%3Aissue) — the same parser shape may already be filed.
4. If the bug is in detection on a specific corpus (oh-my-zsh, prezto, etc.), include the file path inside the corpus repo + the exact line.

## What you'll get back

- **Bugs**: triaged within a few days. Confirmed bugs get a milestone and a fix in the next patch release.
- **Kata requests**: triaged into the [ROADMAP](ROADMAP.md). Severity and fixability are decided per-pattern.
- **Discussions**: best-effort, usually within a week. The author is the only maintainer; please be patient.

## Sponsoring / commercial support

ZShellCheck is MIT-licensed and free to use commercially. There is no paid support tier today. If you depend on it heavily and want to support development, sponsor the author at [github.com/sponsors/afadesigns](https://github.com/sponsors/afadesigns) once that page is live; until then a star and a thoughtful issue or PR is the best way to help.
