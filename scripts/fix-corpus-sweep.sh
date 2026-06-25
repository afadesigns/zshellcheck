#!/usr/bin/env bash
# fix-corpus-sweep.sh — Auto-fix safety gate. For every file in the
# pinned corpora, apply `zshellcheck -fix` (safe tier) and
# `zshellcheck -fix -unsafe-fixes` (every tier), then assert the fixer
# never corrupts the source and always converges:
#
#   1. No corruption — the rewritten file must not raise the parser-error
#      count above the original file's own count. A fix that turns
#      parseable source into unparseable source is a destructive bug.
#   2. Idempotent — a second `-fix` pass over the already-fixed file must
#      produce no further change. A fix that keeps rewriting (for example
#      stacking a glob qualifier every pass) never converges.
#   3. Panic-free — the linter must never crash (exit >= 2 or a Go stack
#      trace) on any invocation.
#
# This is the canonical autofix-safety check, run on real-world code at
# scale. Unlike the parser and violation sweeps it has no baseline: fix
# safety is an absolute invariant, so any violation fails the gate.
#
# Manifest: .github/parser-corpus-manifest.txt (shared with the parser sweep)
#
# Usage:
#   scripts/fix-corpus-sweep.sh
#   ZSHELLCHECK_BIN=/path/to/binary scripts/fix-corpus-sweep.sh
#
# Exit status:
#   0  every file round-trips cleanly
#   1  one or more files were corrupted or did not converge
#   2  manifest / clone / IO error, or a linter panic (never tolerated)

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MANIFEST="${REPO_ROOT}/.github/parser-corpus-manifest.txt"
WORK_DIR="${PARSER_SWEEP_WORK:-${REPO_ROOT}/testdata/external-corpora}"
BIN="${ZSHELLCHECK_BIN:-${REPO_ROOT}/zshellcheck}"

if [[ ! -f "${MANIFEST}" ]]; then
    echo "::error::manifest not found: ${MANIFEST}" >&2
    exit 2
fi

if [[ ! -x "${BIN}" ]]; then
    echo ">>> building zshellcheck into ${BIN}"
    (cd "${REPO_ROOT}" && go build -o "${BIN}" ./cmd/zshellcheck)
fi

mkdir -p "${WORK_DIR}"
SCRATCH="$(mktemp -d)"
trap 'rm -rf "${SCRATCH}"' EXIT

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

# Count parser errors a file raises, suppressing the banner.
parser_errors() {
    "${BIN}" -no-banner "$1" 2>/dev/null | grep -c "^Parser Error" || true
}

total_files=0
corruptions=()
non_idempotent=()
PANICS=()

while IFS=$'\t' read -r name sha url glob_list; do
    [[ -z "${name}" || "${name}" == \#* ]] && continue
    fetch_corpus "${name}" "${sha}" "${url}" || exit 2

    read -r -a globs <<< "${glob_list}"
    while IFS= read -r f; do
        [[ -z "${f}" ]] && continue
        rel="${name}/${f#${WORK_DIR}/${name}/}"
        total_files=$((total_files+1))
        orig_errors=$(parser_errors "${f}")

        for mode in safe unsafe; do
            if [[ "${mode}" == unsafe ]]; then
                fix_args=(-fix -unsafe-fixes -no-banner)
            else
                fix_args=(-fix -no-banner)
            fi

            cp "${f}" "${SCRATCH}/work"
            set +e
            out=$("${BIN}" "${fix_args[@]}" "${SCRATCH}/work" 2>&1)
            rc=$?
            set -e
            if (( rc >= 2 )) || grep -qiE 'panic:|^goroutine [0-9]+ ' <<< "${out}"; then
                PANICS+=("${rel} [${mode}] (exit ${rc})")
                continue
            fi

            fixed_errors=$(parser_errors "${SCRATCH}/work")
            if (( fixed_errors > orig_errors )); then
                corruptions+=("${rel} [${mode}]: parser errors ${orig_errors} -> ${fixed_errors}")
            fi

            cp "${SCRATCH}/work" "${SCRATCH}/work2"
            "${BIN}" "${fix_args[@]}" "${SCRATCH}/work2" > /dev/null 2>&1 || true
            if ! cmp -s "${SCRATCH}/work" "${SCRATCH}/work2"; then
                non_idempotent+=("${rel} [${mode}]: second -fix pass changed the file")
            fi
        done
    done < <(list_files "${WORK_DIR}/${name}" "${globs[@]}")
done < "${MANIFEST}"

echo "fix-corpus sweep: files=${total_files} corruptions=${#corruptions[@]} non_idempotent=${#non_idempotent[@]} panics=${#PANICS[@]}"

# A panic is always fatal.
if (( ${#PANICS[@]} > 0 )); then
    echo "::error::fix-corpus gate: ${#PANICS[@]} file(s) panicked the linter under -fix — fixes must be panic-free"
    for p in "${PANICS[@]}"; do
        echo "  ! ${p}"
    done
    exit 2
fi

fail=0
if (( ${#corruptions[@]} > 0 )); then
    echo "::error::fix-corpus gate: ${#corruptions[@]} file(s) corrupted by -fix (parseable source rewritten into unparseable source)"
    for c in "${corruptions[@]}"; do
        echo "  + ${c}"
    done
    fail=1
fi

if (( ${#non_idempotent[@]} > 0 )); then
    echo "::error::fix-corpus gate: ${#non_idempotent[@]} file(s) did not converge (a second -fix pass kept rewriting)"
    for n in "${non_idempotent[@]}"; do
        echo "  + ${n}"
    done
    fail=1
fi

exit "${fail}"
