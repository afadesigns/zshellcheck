# Project Governance

This document outlines the governance model for the ZShellCheck project.

## Overview

ZShellCheck is an open-source project maintained by the community, with leadership provided by the project founder. We aim to be inclusive, transparent, and meritocratic.

## Roles

### Founder / Lead Maintainer
The project founder (**@afadesigns**) acts as the "Benevolent Dictator for Life" (BDFL) but aims to operate by consensus.
*   Has final say on architectural decisions and roadmap.
*   Manages access to the repository, package managers, and domain names.

### Maintainers
Maintainers are active contributors who have been granted write access to the repository.
*   Can merge Pull Requests.
*   Can triage issues and manage labels.
*   Participate in strategic decisions.

### Contributors
Anyone who submits a Pull Request, opens an issue, or improves documentation.
*   Reviewers: Contributors who consistently review PRs may be invited to become Maintainers.

## Decision Making

### Consensus
We strive for consensus on all technical decisions. Major changes (breaking API changes, new architectural components) should be discussed in a GitHub Issue or Discussion before implementation.

### Deadlocks
If consensus cannot be reached, the Lead Maintainer has the casting vote to resolve the deadlock and move the project forward.

## Contribution Process

1.  **Fork & Clone**: All contributions come via Pull Requests.
2.  **Review**: Every PR requires at least one review from a Maintainer (or the Lead Maintainer).
3.  **CI/CD**: All automated checks must pass.
4.  **Merge**: Squash merges are preferred to keep the history clean.

## Code of Conduct

All participants must adhere to the [Code of Conduct](CODE_OF_CONDUCT.md). Instances of abusive behavior will be dealt with swiftly by the project leadership.

## Changes to Governance

This governance model can be amended via a Pull Request, subject to the approval of the Lead Maintainer.
