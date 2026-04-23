# Security Policy

## Supported Versions

| Version | Supported |
| ------- | --------- |
| `v1.0.x` (latest minor) | Yes |
| `< v1.0.0` | No |

Only the latest `v1.0.x` release receives security fixes. Upgrade to the latest tag before reporting — fixes land as new patch releases, not as backports.

## Reporting a Vulnerability

We take security seriously. If you have discovered a vulnerability in ZShellCheck, we appreciate your help in disclosing it to us in a responsible manner.

### Process

1.  **Do NOT open a public GitHub issue.** This allows us to assess the risk and fix the issue before it can be exploited.
2.  **Email**: Please email the maintainer directly at `github@afadesign.co` (or contact **@afadesigns** via GitHub private (if available) or other social channels linked on the profile).
3.  **Details**: Please include as much information as possible:
    - The type of vulnerability.
    - Full steps to reproduce.
    - Any special configuration required.
    - Potential impact.

### Response

Acknowledgement target: **7 days**. ZShellCheck is maintained by a solo developer; critical issues will be triaged sooner, but please do not assume a same-day response.

## Vulnerability Categories

ZShellCheck is a static analysis tool. Security vulnerabilities generally fall into these categories:

1.  **Code Execution**: A malicious Zsh script causing ZShellCheck to execute arbitrary code on the machine running the linter.
2.  **DoS**: A malicious Zsh script causing ZShellCheck to hang or crash (Denial of Service).
3.  **False Negatives**: Failing to report a critical security flaw in a Zsh script (e.g., missed `eval` or injection). While this is technically a bug, we treat high-impact misses with high priority.
