# Threat model and assurance case

This document records the assurance argument for ZShellCheck.
It identifies the threat model, the trust boundaries, the secure-design principles in force, and how common implementation weaknesses are countered.
It is the canonical reference for the OpenSSF Best Practices `assurance_case` and `documentation_security` criteria.

## Assets

- The user's working tree of Zsh source files.
  These may be private and unreleased.
- The user's host filesystem.
  ZShellCheck has whatever filesystem rights the invoking process holds.
- The user's terminal output stream.
  Violation reports and SARIF documents are emitted on stdout.
- Release artefacts on `https://github.com/afadesigns/zshellcheck/releases`.
  Tarballs, container images, distro packages, signatures, attestations.

## Adversaries

1. **Malicious Zsh source under analysis.**
   The user runs `zshellcheck path/to/untrusted.zsh`.
   The file may be hostile.
2. **Network attacker on the install path.**
   The user runs `install.sh` or `install.ps1` from `raw.githubusercontent.com`.
3. **Compromised dependency.**
   A Go module pulled at build time has been replaced upstream.
4. **Compromised build runner.**
   GitHub-hosted runner or a transitive action has been tampered with.
5. **Stolen maintainer credential.**
   A signing key or PAT has leaked.

## Trust boundaries

| Boundary | Inside (trusted) | Outside (not trusted) |
| :--- | :--- | :--- |
| Filesystem read | Files the user passed on the CLI | Anything else under `$HOME` |
| Filesystem write | Same paths under `-fix`; report files under `-cpuprofile` | Everything else |
| Network | None at runtime | All network endpoints |
| Process exec | None | All external processes |
| Build pipeline | GitHub-hosted runners, signed cosign attestations, pinned action SHAs | Any unsigned artefact |
| Release distribution | cosign-verified archives + SHA-256 checksum file | Mirror copies, third-party rehosts |

## Threat model

| Threat | Mitigation |
| :--- | :--- |
| Hostile Zsh source triggers RCE via parser bug | Parser is recursive-descent over a typed token stream; no `eval`, no `exec`, no shell-out. Fuzz harness runs nightly against the lexer and parser. |
| Hostile input causes infinite loop or memory exhaustion | Parser has bounded recursion. Lexer is single-pass. Fuzz harness exercises pathological inputs. |
| `-fix` corrupts user files | Fixes are byte-exact, idempotent, and context-free per kata. Preview with `-diff` and `-dry-run` before write. Multi-pass cap prevents oscillation. |
| MITM on install script | `install.sh` and `install.ps1` are served over HTTPS, verify SHA-256 against `checksums.txt`, and verify cosign signatures when cosign is on PATH. |
| Compromised release artefact | Every archive ships with a cosign signature pinned to the GitHub Actions OIDC issuer. SLSA Level 3 build provenance is queryable from the attestations index. |
| Dependency hijack | `go.mod` and `go.sum` pin every direct and transitive module. Dependabot opens PRs for security advisories. OSV-Scanner runs on every PR. |
| Malicious GitHub Action | Every third-party action is pinned to a 40-char commit SHA, never a tag. `actionlint` validates workflow syntax; CodeQL scans every push. |
| Leaked credential | Repository has GitHub Secret Scanning and Push Protection enabled. The trace-hygiene rule blocks credential-shaped patterns at write time. |
| Stolen signing key | Releases use cosign keyless OIDC; there is no long-lived signing key to steal. Each signature is bound to a single workflow run. |

## Secure design principles applied

The eight Saltzer/Schroeder principles, ordered as in the original paper, with the project's instance.

1. **Economy of mechanism.**
   Single static binary (≈2 MB), no runtime dependencies, no plugin host, no network listener.
2. **Fail-safe defaults.**
   `-fix` is opt-in; the default is read-only analysis.
   Severity filter defaults to all-on so dangerous findings cannot be silenced by accident.
3. **Complete mediation.**
   The parser refuses every byte it does not have a grammar rule for.
   No bypass path between lexer and kata registry.
4. **Open design.**
   MIT-licensed, public source, every kata is a single Go file readable by any user.
5. **Separation of privilege.**
   The release pipeline requires both a GPG-signed tag and a passing cosign keyless flow.
6. **Least privilege.**
   The binary runs as the invoking user, never escalates.
   The `GITHUB_TOKEN` in workflows uses minimum-necessary scopes per job.
7. **Least common mechanism.**
   Kata state is per-violation; nothing global is mutated during a walk.
8. **Psychological acceptability.**
   The CLI prints actionable messages with file, line, column, and a suggested fix.
   Inline `# noka` directives let users silence one kata without editing config.

A ninth principle from later work — **input validation with allowlists** — applies at every parser branch: the parser only accepts known token patterns and rejects everything else.

## Common implementation weaknesses countered

Mapping against [CWE/SANS Top 25](https://cwe.mitre.org/top25/) and [OWASP Top 10](https://owasp.org/Top10/) for the categories that apply to a CLI static-analysis tool:

| Class | Counter |
| :--- | :--- |
| Injection (CWE-77, CWE-78, CWE-89) | Linter never invokes a shell or external process. No SQL surface. |
| Out-of-bounds read or write (CWE-125, CWE-787) | Go is memory-safe; no `unsafe` imports in project code. |
| Use after free (CWE-416) | Garbage-collected runtime; no manual lifetime management. |
| Path traversal (CWE-22) | The CLI accepts only paths the user passes; relative paths are resolved against the working directory; no symlink chase by default. |
| Deserialisation of untrusted data (CWE-502) | YAML config uses `yaml.Unmarshal` with a typed struct; no arbitrary type instantiation. |
| Insecure dependencies (CWE-1395) | `go.sum` pins every module; OSV-Scanner and Dependabot run continuously. |
| Insufficient logging (CWE-778) | Violations include file, line, column, kata ID, severity, and message; SARIF output preserves the same. |
| Improper certificate validation (CWE-295) | The release artefacts are validated by cosign; install scripts verify SHA-256 + cosign before execution. |
| Insecure defaults (CWE-1188) | `-fix` requires explicit opt-in; `-no-color` is auto-set on non-TTY; severity filter defaults to all. |
| Missing authentication / authorization (CWE-306, CWE-862) | Not applicable; the CLI has no authentication surface. |

## Assurance summary

Threat model documented above.
Trust boundaries enumerated.
Saltzer/Schroeder principles each have a concrete project counter-example.
Common implementation weaknesses each have a concrete counter or an explicit not-applicable.

This document is the canonical assurance case.
Updates to it are required whenever the trust boundaries change — for example, if the linter ever gains a network-fetching feature, the boundaries table must be revised in the same PR.

## Security review log

The project conducts a security review against this assurance case at least every twelve months and after any change that crosses a trust boundary.
Each review re-walks the threats table, the trust boundaries, and the CWE/SANS-Top-25 mapping; updates this document with any new mitigations; and is recorded below with date, scope, and reviewer.

| Date | Scope | Reviewer | Notes |
| :--- | :--- | :--- | :--- |
| 2026-04-26 | Full review against the OpenSSF Best Practices gold criteria. Threat model, trust boundaries, secure-design principles, and CWE/SANS Top 25 mapping re-walked end-to-end. CodeQL, gosec, govulncheck, OSV-Scanner, and golangci-lint clean on `main` at commit `aefcedc`. Fuzz-nightly green. | Andreas Fahl (@afadesigns) | Initial review. No outstanding issues. |

A formal external review may be commissioned once the project crosses adoption thresholds that justify it.
