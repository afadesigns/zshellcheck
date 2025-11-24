# Security Policy

## Supported Versions

Only the latest major version of ZShellCheck is officially supported with security updates.

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |
| < 0.0.x | :x:                |

## Reporting a Vulnerability

We take security seriously. If you have discovered a vulnerability in ZShellCheck, we appreciate your help in disclosing it to us in a responsible manner.

### Process

1.  **Do NOT open a public GitHub issue.** This allows us to assess the risk and fix the issue before it can be exploited.
2.  **Email**: Please email the maintainer directly at `security@afadesigns.com` (or contact **@afadesigns** via GitHub private (if available) or other social channels linked on the profile).
3.  **Details**: Please include as much information as possible:
    - The type of vulnerability.
    - Full steps to reproduce.
    - Any special configuration required.
    - Potential impact.

### Response

We will acknowledge your report within **48 hours** and provide an estimated timeline for the fix.

## Vulnerability Categories

ZShellCheck is a static analysis tool. Security vulnerabilities generally fall into these categories:

1.  **Code Execution**: A malicious Zsh script causing ZShellCheck to execute arbitrary code on the machine running the linter.
2.  **DoS**: A malicious Zsh script causing ZShellCheck to hang or crash (Denial of Service).
3.  **False Negatives**: Failing to report a critical security flaw in a Zsh script (e.g., missed `eval` or injection). While this is technically a bug, we treat high-impact misses with high priority.
