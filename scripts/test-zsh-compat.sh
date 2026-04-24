#!/usr/bin/env bash
#
# test-zsh-compat.sh — Run zshellcheck against a curated corpus of
# well-known Zsh projects and report parser errors + violation counts.
#
# The corpora live under testdata/external-corpora/<name>/ which is
# gitignored (see .gitignore). This script is local-only: it never
# commits or pushes anything from those clones.
#
# Usage:
#   scripts/test-zsh-compat.sh                  # full matrix
#   scripts/test-zsh-compat.sh omz              # oh-my-zsh only
#   scripts/test-zsh-compat.sh p10k prezto      # two specific corpora
#
# Exit status: 0 when every requested corpus parsed without errors,
# 1 when any file in any corpus produced a parser error.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CORPORA_DIR="${REPO_ROOT}/testdata/external-corpora"
BIN="${ZSHELLCHECK_BIN:-${REPO_ROOT}/zshellcheck}"

# Corpus catalogue: name|upstream|file-glob
# File globs are relative to the corpus root and are shell-expanded
# via `find` below.
declare -A CORPUS_REPO=(
    [omz]="https://github.com/ohmyzsh/ohmyzsh.git"
    [p10k]="https://github.com/romkatv/powerlevel10k.git"
    [autosuggestions]="https://github.com/zsh-users/zsh-autosuggestions.git"
    [syntax-highlighting]="https://github.com/zsh-users/zsh-syntax-highlighting.git"
    [completions]="https://github.com/zsh-users/zsh-completions.git"
    [prezto]="https://github.com/sorin-ionescu/prezto.git"
    [spaceship]="https://github.com/spaceship-prompt/spaceship-prompt.git"
)

declare -A CORPUS_GLOB=(
    [omz]='*.zsh *.zsh-theme *.plugin.zsh'
    [p10k]='*.zsh *.zsh-theme'
    [autosuggestions]='*.zsh'
    [syntax-highlighting]='*.zsh'
    [completions]='*.zsh'
    [prezto]='*.zsh *.zsh-theme'
    [spaceship]='*.zsh *.zsh-theme'
)

# Pick the corpora to run. When no args are passed we iterate the
# whole catalogue in insertion-stable order.
if [[ $# -eq 0 ]]; then
    CORPORA=(omz p10k autosuggestions syntax-highlighting completions prezto spaceship)
else
    CORPORA=("$@")
fi

# Build the binary locally when one wasn't supplied via env.
if [[ ! -x "${BIN}" ]]; then
    echo ">>> building zshellcheck into ${BIN}"
    (cd "${REPO_ROOT}" && go build -o "${BIN}" ./cmd/zshellcheck)
fi

mkdir -p "${CORPORA_DIR}"

ensure_corpus() {
    local name="$1"
    local repo="${CORPUS_REPO[$name]:-}"
    if [[ -z "${repo}" ]]; then
        echo "!! unknown corpus: ${name}" >&2
        return 1
    fi
    local dir="${CORPORA_DIR}/${name}"
    if [[ -d "${dir}/.git" ]]; then
        (cd "${dir}" && git fetch --quiet --depth=1 origin HEAD && git reset --quiet --hard FETCH_HEAD) \
            || echo "!! refresh failed for ${name} (keeping existing clone)" >&2
    else
        echo ">>> cloning ${name} from ${repo}"
        git clone --quiet --depth=1 "${repo}" "${dir}"
    fi
}

scan_corpus() {
    local name="$1"
    local dir="${CORPORA_DIR}/${name}"
    local glob="${CORPUS_GLOB[$name]}"
    local parse_errors=0
    local files=0

    echo
    echo "=== ${name} ==="
    echo "corpus: ${dir}"

    # Build a file list that matches the configured glob. xargs -0 keeps
    # paths with spaces intact.
    local find_args=()
    for pat in ${glob}; do
        find_args+=(-o -name "${pat}")
    done
    # Drop the leading '-o' so the expression starts with a name clause.
    find_args=("${find_args[@]:1}")

    while IFS= read -r -d '' f; do
        files=$((files + 1))
        # Only parser errors go to stderr with the "Parser Error" prefix.
        # Count them per file without failing the pipeline on kata
        # violations (those are informational for the compat matrix).
        if "${BIN}" -format json -- "${f}" >/dev/null 2>"${REPO_ROOT}/.compat-err"; then
            :
        fi
        if grep -q 'Parser Error' "${REPO_ROOT}/.compat-err"; then
            parse_errors=$((parse_errors + 1))
            echo "  PARSE-ERROR: ${f#"${dir}/"}"
            head -n1 "${REPO_ROOT}/.compat-err" | sed 's/^/    /'
        fi
    done < <(find "${dir}" -type f \( "${find_args[@]}" \) -not -path '*/.git/*' -print0)

    rm -f "${REPO_ROOT}/.compat-err"
    echo "files scanned: ${files}"
    echo "parse errors : ${parse_errors}"
    if [[ ${parse_errors} -gt 0 ]]; then
        return 1
    fi
    return 0
}

EXIT=0
for name in "${CORPORA[@]}"; do
    ensure_corpus "${name}" || { EXIT=1; continue; }
    scan_corpus "${name}" || EXIT=1
done

exit "${EXIT}"
