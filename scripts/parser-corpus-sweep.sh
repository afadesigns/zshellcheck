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
PANICS=()

while IFS=$'\t' read -r name sha url glob_list; do
    [[ -z "${name}" || "${name}" == \#* ]] && continue
    fetch_corpus "${name}" "${sha}" "${url}" || exit 2

    # Split the glob list on spaces without pathname expansion. A
    # bare `globs=( ${glob_list} )` would let a pattern such as `_*`
    # match files in the repo root (e.g. `_typos.toml`) before it
    # ever reaches find; read -a word-splits on IFS but never globs.
    read -r -a globs <<< "${glob_list}"
    while IFS= read -r f; do
        [[ -z "${f}" ]] && continue
        rel="${name}/${f#${WORK_DIR}/${name}/}"
        # Run the binary once in the main shell (not a command-
        # substitution subshell) so a detected panic can be recorded in
        # the PANICS global. A panic exits >= 2 with a Go stack trace; it
        # kills the whole run, so a corpus file that triggers one is a
        # release-blocking regression. Parser errors are baselined and
        # tolerated; panics never are.
        set +e
        out=$("${BIN}" "${f}" 2>&1)
        rc=$?
        set -e
        if (( rc >= 2 )) || grep -qiE 'panic:|^goroutine [0-9]+ ' <<< "${out}"; then
            PANICS+=("${rel} (exit ${rc})")
        fi
        n=$(grep -c "^Parser Error" <<< "${out}" || true)
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

# A panic is always fatal — even in --update-baseline mode, so a crash
# can never be silently baked into the baseline.
if (( ${#PANICS[@]} > 0 )); then
    echo "::error::parser-corpus gate: ${#PANICS[@]} file(s) panicked the linter — integrations must be panic-free"
    for p in "${PANICS[@]}"; do
        echo "  ! ${p}"
    done
    exit 2
fi

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
