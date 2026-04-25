# Integrations

ZShellCheck is verified against the script corpora of these popular Zsh ecosystems. Every release runs a parse + lint sweep over each — no panics, no crashes, deterministic output. The list grows each release; the goal is 300+ corpora before v2.

## Summary

|   |   |
| ---: | :--- |
| **19** | corpora verified today |
| **0** | panics on the current sweep |
| **300+** | targeted before v2 — see [ROADMAP.md](../ROADMAP.md) |

## Featured

The corpora we test most heavily and link from the docs.

| Project | Category | Files |
| :--- | :--- | ---: |
| [oh-my-zsh](https://github.com/ohmyzsh/ohmyzsh) | Framework | 497 |
| [prezto](https://github.com/sorin-ionescu/prezto) | Framework | 41 |
| [powerlevel10k](https://github.com/romkatv/powerlevel10k) | Prompt | 16 |
| [zinit](https://github.com/zdharma-continuum/zinit) | Plugin manager | 9 |
| [fzf](https://github.com/junegunn/fzf) | Tooling | 2 |
| [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting) | Plugin | 301 |

## Frameworks

| Project | Files |
| :--- | ---: |
| [oh-my-zsh](https://github.com/ohmyzsh/ohmyzsh) | 497 |
| [prezto](https://github.com/sorin-ionescu/prezto) | 41 |
| [zimfw](https://github.com/zimfw/zimfw) | 1 |
| [zephyr](https://github.com/mattmc3/zephyr) | 21 |
| [zsh-utils](https://github.com/belak/zsh-utils) | 5 |

## Plugin / theme managers

| Project | Files |
| :--- | ---: |
| [antidote](https://github.com/mattmc3/antidote) | 24 |
| [zinit](https://github.com/zdharma-continuum/zinit) | 9 |

## Plugin / theme tooling

| Project | Files |
| :--- | ---: |
| [fzf](https://github.com/junegunn/fzf) | 2 |
| [fzf-tab](https://github.com/Aloxaf/fzf-tab) | 5 |
| [fast-syntax-highlighting](https://github.com/zdharma-continuum/fast-syntax-highlighting) | 4 |

## Plugins

| Project | Files |
| :--- | ---: |
| [zsh-autosuggestions](https://github.com/zsh-users/zsh-autosuggestions) | 13 |
| [zsh-syntax-highlighting](https://github.com/zsh-users/zsh-syntax-highlighting) | 301 |
| [zsh-history-substring-search](https://github.com/zsh-users/zsh-history-substring-search) | 2 |
| [zsh-vi-mode](https://github.com/jeffreytse/zsh-vi-mode) | 2 |
| [zsh-autocomplete](https://github.com/marlonrichert/zsh-autocomplete) | 3 |
| [zsh-completions](https://github.com/zsh-users/zsh-completions) | 1 |

## Prompts

| Project | Files |
| :--- | ---: |
| [powerlevel10k](https://github.com/romkatv/powerlevel10k) | 16 |
| [spaceship-prompt](https://github.com/spaceship-prompt/spaceship-prompt) | 119 |
| [starship](https://github.com/starship/starship) | 1 |

## Roadmap — targeted next

- [zsh-users/zsh](https://github.com/zsh-users/zsh) — `Functions/` and `Completion/` directories full of canonical Zsh.
- [romkatv/zsh-bench](https://github.com/romkatv/zsh-bench)
- [romkatv/gitstatus](https://github.com/romkatv/gitstatus)
- [sorin-ionescu/prezto-contrib](https://github.com/sorin-ionescu/prezto-contrib)
- [ohmyzsh-incubator](https://github.com/ohmyzsh-incubator)
- [Freed-Wu/zsh-help](https://github.com/Freed-Wu/zsh-help)

## How the sweep runs

Each release tag triggers a parse + lint pass over every corpus listed in the **Featured** + per-category tables above. Sandbox lives at `~/.config/zsh/zshellcheck/integration/` (local-only) and never touches the user's actual `~/.config/zsh` modules. A run produces:

- `parse_errors` — total parser failures across the corpus.
- `violations` — total kata hits (all severities).
- A per-corpus markdown report under `~/.config/zsh/zshellcheck/integration/reports/`.

A bug surfaced by the sweep gets a GitHub issue, a PR fixes it, and the corpus stays in the matrix on every subsequent release.

## Adding an integration

If you maintain (or rely on) a popular Zsh project not listed above and want it covered by every release sweep, open an issue tagged `integration` with the repo URL and a short note on what it covers. We add it to the next sweep and credit the request in the changelog entry.
