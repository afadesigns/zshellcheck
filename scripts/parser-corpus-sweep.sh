#!/usr/bin/env bash
# parser-corpus-sweep.sh — Run zshellcheck against pinned upstream
# corpora and compare per-file parser-error counts to a baseline.
#
# Manifest: .github/parser-corpus-manifest.txt
# Baseline: .github/parser-error-baseline.txt
#
# Usage:
#   scripts/parser-corpus-sweep.sh                   # gate (compare to baseline)
#   scripts/parser-corpus-sweep.sh --update-baseline # rewrite baseline from current run
#   ZSHELLCHECK_BIN=/path/to/binary scripts/parser-corpus-sweep.sh
#
# Exit status:
#   0  every file matches its baseline (or improves below it)
#   1  one or more files regressed
#   2  manifest / clone / IO error

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MANIFEST="${REPO_ROOT}/.github/parser-corpus-manifest.txt"
BASELINE="${REPO_ROOT}/.github/parser-error-baseline.txt"
WORK_DIR="${PARSER_SWEEP_WORK:-${REPO_ROOT}/testdata/external-corpora}"
BIN="${ZSHELLCHECK_BIN:-${REPO_ROOT}/zshellcheck}"

UPDATE_BASELINE=0
if [[ "${1:-}" == "--update-baseline" ]]; then
    UPDATE_BASELINE=1
fi

if [[ ! -f "${MANIFEST}" ]]; then
    echo "::error::manifest not found: ${MANIFEST}" >&2
    exit 2
fi

if [[ ! -x "${BIN}" ]]; then
    echo ">>> building zshellcheck into ${BIN}"
    (cd "${REPO_ROOT}" && go build -o "${BIN}" ./cmd/zshellcheck)
fi

mkdir -p "${WORK_DIR}"

# Fetch a corpus at a pinned SHA. Idempotent.
fetch_corpus() {
    local name="$1" sha="$2" url="$3"
    local dir="${WORK_DIR}/${name}"
    if [[ -d "${dir}/.git" ]]; then
        local cur
        cur=$(cd "${dir}" && git rev-parse HEAD 2>/dev/null || echo "")
        if [[ "${cur}" == "${sha}"* ]]; then
            return 0
        fi
        (cd "${dir}" && git fetch --quiet --depth=1 origin "${sha}" 2>/dev/null \
            && git reset --quiet --hard "${sha}") || {
            rm -rf "${dir}"
        }
    fi
    if [[ ! -d "${dir}/.git" ]]; then
        git clone --quiet "${url}" "${dir}"
        (cd "${dir}" && git reset --quiet --hard "${sha}") || {
            echo "::error::failed to checkout ${sha} in ${name}" >&2
            return 2
        }
    fi
}

# Build the file list for one corpus given its glob list.
list_files() {
    local dir="$1"
    shift
    local globs=("$@")
    local find_args=()
    local first=1
    for g in "${globs[@]}"; do
        if (( first )); then
            find_args+=(-name "${g}")
            first=0
        else
            find_args+=(-o -name "${g}")
        fi
    done
    find "${dir}" -type f \( "${find_args[@]}" \) -print 2>/dev/null | sort
}

# Run zshellcheck on one file, return parser-error count.
parser_errors_for() {
    local f="$1"
    "${BIN}" "$f" 2>&1 | grep -c "^Parser Error" || true
}

declare -A CURRENT
declare -A BASELINE_MAP

if [[ -f "${BASELINE}" ]]; then
    while IFS=$'\t' read -r path count; do
        [[ -z "${path}" || "${path}" == \#* ]] && continue
        BASELINE_MAP["${path}"]="${count}"
    done < "${BASELINE}"
fi

total_files=0
total_errors=0
regressions=()
improvements=()

while IFS=$'\t' read -r name sha url glob_list; do
    [[ -z "${name}" || "${name}" == \#* ]] && continue
    fetch_corpus "${name}" "${sha}" "${url}" || exit 2

    # shellcheck disable=SC2206
    globs=( ${glob_list} )
    while IFS= read -r f; do
        [[ -z "${f}" ]] && continue
        rel="${name}/${f#${WORK_DIR}/${name}/}"
        n=$(parser_errors_for "${f}")
        CURRENT["${rel}"]="${n}"
        total_files=$((total_files+1))
        total_errors=$((total_errors+n))
        prior="${BASELINE_MAP[${rel}]:-0}"
        if (( n > prior )); then
            regressions+=("${rel}: ${prior} -> ${n}")
        elif (( n < prior )); then
            improvements+=("${rel}: ${prior} -> ${n}")
        fi
    done < <(list_files "${WORK_DIR}/${name}" "${globs[@]}")
done < "${MANIFEST}"

echo "parser-corpus sweep: files=${total_files} parser_errors=${total_errors}"

if (( UPDATE_BASELINE )); then
    {
        echo "# Parser-error baseline. Regenerate with scripts/parser-corpus-sweep.sh --update-baseline"
        echo "# Format: <relpath>\t<count>"
        for k in $(printf '%s\n' "${!CURRENT[@]}" | sort); do
            v="${CURRENT[$k]}"
            (( v == 0 )) && continue
            printf '%s\t%s\n' "${k}" "${v}"
        done
    } > "${BASELINE}"
    echo "wrote baseline: ${BASELINE}"
    exit 0
fi

if (( ${#regressions[@]} > 0 )); then
    echo "::error::parser-corpus gate: ${#regressions[@]} file(s) regressed"
    for r in "${regressions[@]}"; do
        echo "  + ${r}"
    done
    exit 1
fi

if (( ${#improvements[@]} > 0 )); then
    echo "::warning::parser-corpus gate: ${#improvements[@]} file(s) improved below baseline. Run scripts/parser-corpus-sweep.sh --update-baseline to lock the gain in."
    for i in "${improvements[@]}"; do
        echo "  - ${i}"
    done
fi

exit 0
