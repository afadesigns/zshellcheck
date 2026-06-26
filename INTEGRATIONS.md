# Integrations

ZShellCheck is verified against the script trees of widely used Zsh integrations.
Every release runs a parse and lint sweep over the pinned corpus matrix — no panics, no crashes, deterministic output.
The heaviest Featured trees are swept manually rather than in the per-release gate, because their clones are too large to pin in CI.
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
| [zgen](https://github.com/tarjoilija/zgen) | 1 |
| [zsh-snap](https://github.com/marlonrichert/zsh-snap) | 3 |
| [zpm](https://github.com/zpm-zsh/zpm) | 3 |
| [zcomet](https://github.com/agkozak/zcomet) | 1 |

## Plugin and theme tooling

| Integration | Files |
| :--- | ---: |
| [fzf](https://github.com/junegunn/fzf) | 2 |
| [fzf-tab](https://github.com/Aloxaf/fzf-tab) | 5 |
| [fast-syntax-highlighting](https://github.com/zdharma-continuum/fast-syntax-highlighting) | 4 |
| [F-Sy-H](https://github.com/z-shell/F-Sy-H) | 4 |
| [fzf-zsh-plugin](https://github.com/unixorn/fzf-zsh-plugin) | 2 |
| [zsh-async](https://github.com/mafredri/zsh-async) | 4 |

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
| [atuin](https://github.com/atuinsh/atuin) | 2 |
| [navi](https://github.com/denisidoro/navi) | 1 |
| [autojump](https://github.com/wting/autojump) | 1 |
| [mcfly](https://github.com/cantino/mcfly) | 2 |
| [forgit](https://github.com/wfxr/forgit) | 1 |
| [git-open](https://github.com/paulirish/git-open) | 1 |
| [git-flow-completion](https://github.com/bobthecow/git-flow-completion) | 1 |
| [zsh-nvm](https://github.com/lukechilds/zsh-nvm) | 1 |
| [zsh-z](https://github.com/agkozak/zsh-z) | 1 |
| [zsh-histdb](https://github.com/larkery/zsh-histdb) | 3 |
| [git-extra-commands](https://github.com/unixorn/git-extra-commands) | 1 |
| [alias-tips](https://github.com/djui/alias-tips) | 1 |
| [zsh-abbr](https://github.com/olets/zsh-abbr) | 16 |
| [zsh-autoenv](https://github.com/Tarrasch/zsh-autoenv) | 4 |
| [wd](https://github.com/mfaerevaag/wd) | 1 |
| [zsh-autopair](https://github.com/hlissner/zsh-autopair) | 3 |
| [zaw](https://github.com/zsh-users/zaw) | 27 |
| [zsh-notify](https://github.com/marzocchi/zsh-notify) | 1 |
| [zsh-better-npm-completion](https://github.com/lukechilds/zsh-better-npm-completion) | 1 |
| [zsh-defer](https://github.com/romkatv/zsh-defer) | 1 |
| [zsh-fzf-history-search](https://github.com/joshskidmore/zsh-fzf-history-search) | 2 |
| [zsh-bd](https://github.com/Tarrasch/zsh-bd) | 1 |
| [emoji-cli](https://github.com/b4b4r07/emoji-cli) | 2 |
| [zsh-interactive-cd](https://github.com/changyuheng/zsh-interactive-cd) | 1 |
| [zsh-completion-generator](https://github.com/RobSis/zsh-completion-generator) | 1 |
| [zsh-shift-select](https://github.com/jirutka/zsh-shift-select) | 1 |
| [zsh-hist](https://github.com/marlonrichert/zsh-hist) | 1 |
| [zsh-edit](https://github.com/marlonrichert/zsh-edit) | 2 |
| [zsh-256color](https://github.com/chrissicool/zsh-256color) | 1 |
| [history-search-multi-word](https://github.com/zdharma-continuum/history-search-multi-word) | 3 |
| [zsh-dircolors-solarized](https://github.com/joel-porquet/zsh-dircolors-solarized) | 2 |
| [zsh-eza](https://github.com/z-shell/zsh-eza) | 3 |
| [zsh-vimode-visual](https://github.com/b4b4r07/zsh-vimode-visual) | 2 |
| [zcolors](https://github.com/marlonrichert/zcolors) | 1 |
| [zhooks](https://github.com/agkozak/zhooks) | 1 |
| [zsh-colored-man-pages](https://github.com/ael-code/zsh-colored-man-pages) | 1 |
| [zsh-navigation-tools](https://github.com/z-shell/zsh-navigation-tools) | 2 |

## Prompts

| Integration | Files |
| :--- | ---: |
| [powerlevel10k](https://github.com/romkatv/powerlevel10k) | 16 |
| [spaceship-prompt](https://github.com/spaceship-prompt/spaceship-prompt) | 119 |
| [starship](https://github.com/starship/starship) | 1 |
| [gitstatus](https://github.com/romkatv/gitstatus) | 2 |
| [typewritten](https://github.com/reobin/typewritten) | 5 |
| [agkozak-zsh-prompt](https://github.com/agkozak/agkozak-zsh-prompt) | 2 |
| [bullet-train.zsh](https://github.com/caiogondim/bullet-train.zsh) | 1 |
| [geometry](https://github.com/geometry-zsh/geometry) | 1 |
| [purer](https://github.com/dfurnes/purer) | 2 |
| [minimal](https://github.com/subnixr/minimal) | 1 |

## Roadmap — targeted next

The pinned sweep matrix admits only corpora that parse with zero errors and never crash.
A popular integration that still surfaces parser errors stays on this list until the parser drains them, then graduates to the matrix.

Pending parser support:

- [zsh](https://github.com/zsh-users/zsh) — the reference shell. Its `Completion/`
  tree leans on glob-pattern alternation (`(foo|bar)`), pending parser support.
- [antigen](https://github.com/zsh-users/antigen)
- [lean](https://github.com/miekg/lean)
- [powerlevel9k](https://github.com/Powerlevel9k/powerlevel9k)
- [pure](https://github.com/sindresorhus/pure)
- [you-should-use](https://github.com/MichaelAquilina/zsh-you-should-use)
- [zgenom](https://github.com/jandamm/zgenom)
- [zplug](https://github.com/zplug/zplug)
- [zsh-auto-notify](https://github.com/MichaelAquilina/zsh-auto-notify)
- [zsh-autoswitch-virtualenv](https://github.com/MichaelAquilina/zsh-autoswitch-virtualenv)
- [zsh-bench](https://github.com/romkatv/zsh-bench)
- [zsh-lazyload](https://github.com/qoomon/zsh-lazyload)
- [zsh-system-clipboard](https://github.com/kutsan/zsh-system-clipboard)

To propose another corpus, follow [Adding an integration](#adding-an-integration).

## How the sweep runs

Each release tag triggers a parse and lint pass over every integration in the pinned corpus matrix.
The heaviest Featured trees, whose clones are too large to pin in CI, are swept manually instead.
Each pass produces:

- `parse_errors` — total parser failures across the integration.
- `violations` — total kata hits across all severities.

When the sweep surfaces a bug, a GitHub issue is filed, a PR fixes it, and the integration stays in the matrix on every subsequent release.

## Adding an integration

To get a popular Zsh integration covered, open an issue tagged `integration` with the repo URL and a short note on what it covers.
The next sweep adds it and the changelog entry credits the request.
