# Integrations

ZShellCheck is verified against the script trees of widely used Zsh integrations.
Every release runs a parse and lint sweep over each — no panics, no crashes, deterministic output.
The list grows with every release.

## Featured

The integrations the project tests most heavily and links from the docs.

| Integration | Category | Files |
| :--- | :--- | ---: |
| [oh-my-zsh](https://github.com/ohmyzsh/ohmyzsh) | Framework | 497 |
| [prezto](https://github.com/sorin-ionescu/prezto) | Framework | 41 |
| [powerlevel10k](https://github.com/romkatv/powerlevel10k) | Prompt | 16 |
| [zinit](https://github.com/zdharma-continuum/zinit) | Plugin manager | 9 |
| [fzf](https://github.com/junegunn/fzf) | Tooling | 2 |
| [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting) | Plugin | 301 |

## Frameworks

| Integration | Files |
| :--- | ---: |
| [oh-my-zsh](https://github.com/ohmyzsh/ohmyzsh) | 497 |
| [prezto](https://github.com/sorin-ionescu/prezto) | 41 |
| [prezto-contrib](https://github.com/belak/prezto-contrib) | 20 |
| [zimfw](https://github.com/zimfw/zimfw) | 1 |
| [zephyr](https://github.com/mattmc3/zephyr) | 21 |
| [zsh-utils](https://github.com/belak/zsh-utils) | 5 |

## Plugin and theme managers

| Integration | Files |
| :--- | ---: |
| [antidote](https://github.com/mattmc3/antidote) | 24 |
| [zinit](https://github.com/zdharma-continuum/zinit) | 9 |

## Plugin and theme tooling

| Integration | Files |
| :--- | ---: |
| [fzf](https://github.com/junegunn/fzf) | 2 |
| [fzf-tab](https://github.com/Aloxaf/fzf-tab) | 5 |
| [fast-syntax-highlighting](https://github.com/zdharma-continuum/fast-syntax-highlighting) | 4 |

## Plugins

| Integration | Files |
| :--- | ---: |
| [zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions) | 13 |
| [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting) | 301 |
| [zsh-history-substring-search](https://github.com/zsh-users/zsh-history-substring-search) | 2 |
| [zsh-vi-mode](https://github.com/jeffreytse/zsh-vi-mode) | 2 |
| [zsh-autocomplete](https://github.com/marlonrichert/zsh-autocomplete) | 3 |
| [zsh-completions](https://github.com/zsh-users/zsh-completions) | 1 |
| [zsh-help](https://github.com/Freed-Wu/zsh-help) | 1 |

## Prompts

| Integration | Files |
| :--- | ---: |
| [powerlevel10k](https://github.com/romkatv/powerlevel10k) | 16 |
| [spaceship-prompt](https://github.com/spaceship-prompt/spaceship-prompt) | 119 |
| [starship](https://github.com/starship/starship) | 1 |
| [gitstatus](https://github.com/romkatv/gitstatus) | 2 |

## Roadmap — targeted next

[gitstatus](https://github.com/romkatv/gitstatus),
[prezto-contrib](https://github.com/belak/prezto-contrib), and
[zsh-help](https://github.com/Freed-Wu/zsh-help) joined the pinned sweep matrix.

Targeted next:

- [zsh](https://github.com/zsh-users/zsh) — the reference shell. Its `Completion/`
  tree leans on glob-pattern alternation (`(foo|bar)`), pending parser support.
- [zsh-bench](https://github.com/romkatv/zsh-bench)

To propose another corpus, follow [Adding an integration](#adding-an-integration).

## How the sweep runs

Each release tag triggers a parse and lint pass over every integration listed in the **Featured** and per-category tables.
Each pass produces:

- `parse_errors` — total parser failures across the integration.
- `violations` — total kata hits across all severities.

When the sweep surfaces a bug, a GitHub issue is filed, a PR fixes it, and the integration stays in the matrix on every subsequent release.

## Adding an integration

To get a popular Zsh integration covered, open an issue tagged `integration` with the repo URL and a short note on what it covers.
The next sweep adds it and the changelog entry credits the request.
