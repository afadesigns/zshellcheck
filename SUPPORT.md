# Support

Thanks for using ZShellCheck.
The fastest path to the right answer is below.

## Where to file what

| You have | Open it as |
| --- | --- |
| A **bug** — panic, crash, false positive, false negative, parser failure on valid Zsh | A [GitHub issue](https://github.com/afadesigns/zshellcheck/issues/new). Include the binary version (`zshellcheck -version`), a minimal reproduction, and the full stderr including the banner and any stack trace. |
| A **kata request** — a Zsh anti-pattern ZShellCheck does not yet catch | A [GitHub issue](https://github.com/afadesigns/zshellcheck/issues/new). Describe the pattern, why it bites, and a code sample that should trip the new rule. |
| A **question or design idea** | A [GitHub Discussion](https://github.com/afadesigns/zshellcheck/discussions). Discussions are public and searchable; use issues only when there is something to fix. |
| A **security vulnerability** | Do not file a public issue. See [SECURITY.md](SECURITY.md) for the private disclosure flow. |
| A **documentation gap** — typo, broken link, stale fact, missing example | A small PR is more useful than an issue. The [contributing guide](CONTRIBUTING.md#pull-request-workflow) shows the workflow; doc-only PRs skip most gates. |

## Before opening a bug

1. Run the latest tagged release.
   `go install github.com/afadesigns/zshellcheck/cmd/zshellcheck@latest` is the fastest path.
2. Re-run with `-no-banner -severity error` to confirm the issue is not silenced by the noise filter.
3. Check the [issue tracker](https://github.com/afadesigns/zshellcheck/issues?q=is%3Aissue); the same parser shape may already be filed.
4. When the bug is in detection on a specific corpus — oh-my-zsh, prezto, and similar — include the file path inside the corpus repo and the exact line.

## What to expect back

- **Bugs.**
  Triaged within a few days.
  Confirmed bugs get a milestone and a fix in the next patch release.
- **Kata requests.**
  Triaged into the [roadmap](ROADMAP.md).
  Severity and fixability are decided per pattern.
- **Discussions.**
  Best effort, usually within a week.
  The author is the only maintainer; please be patient.

## Sponsoring and commercial support

ZShellCheck is MIT-licensed and free to use commercially.
There is no paid support tier.
To support development, sponsor the author at [github.com/sponsors/afadesigns](https://github.com/sponsors/afadesigns) once the page is live.
Until then, a star and a thoughtful issue or PR is the best way to help.
