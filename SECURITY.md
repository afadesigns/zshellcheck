# Security policy

## Supported versions

| Version | Supported |
| --- | --- |
| `v1.0.x` (latest minor) | Yes |
| `< v1.0.0` | No |

Only the latest `v1.0.x` release receives security fixes.
Upgrade to the latest tag before reporting; fixes land as new patch releases, not as backports.

## Reporting a vulnerability

If you discover a vulnerability in ZShellCheck, disclose it responsibly using the process below.

### Process

1. **Do not open a public GitHub issue.**
   Public disclosure before a fix lets the issue be exploited.
2. **Use one of the two private channels.**
   - GitHub Private Vulnerability Reporting (preferred): submit at [Security → Advisories → Report a vulnerability](https://github.com/afadesigns/zshellcheck/security/advisories/new).
     The form is encrypted in transit and visible only to maintainers.
   - Email the maintainer at `github@afadesign.co`.
     GitHub private contact is also available via [@afadesigns](https://github.com/afadesigns).
3. **Include as much detail as possible.**
   - The type of vulnerability.
   - Full reproduction steps.
   - Any special configuration required.
   - Potential impact.

### Response

Acknowledgement target: 7 days.
ZShellCheck is maintained by a solo developer; critical issues are triaged sooner, but a same-day response is not guaranteed.

## Vulnerability categories

ZShellCheck is a static-analysis tool.
Vulnerabilities fall into three categories:

1. **Code execution.**
   A malicious Zsh script causing ZShellCheck to execute arbitrary code on the host running the linter.
2. **Denial of service.**
   A malicious Zsh script causing ZShellCheck to hang or crash.
3. **False negatives.**
   Failure to report a critical security flaw in a Zsh script — for example a missed `eval` or injection.
   This is a bug class; high-impact misses are treated with high priority.
