#!/usr/bin/env bash
# Exercise the corpus gate without downloading or changing an external corpus.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRATCH="$(mktemp -d)"
trap 'rm -rf "${SCRATCH}"' EXIT
mkdir -p "${SCRATCH}/repo/scripts" "${SCRATCH}/repo/.github" \
    "${SCRATCH}/corpora/fixture/.git" "${SCRATCH}/bin"
cp "${REPO_ROOT}/scripts/fix-corpus-sweep.sh" "${SCRATCH}/repo/scripts/"
printf 'fixture\taaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\tunused\t*.zsh\n' \
    > "${SCRATCH}/repo/.github/parser-corpus-manifest.txt"
printf 'print ok\n' > "${SCRATCH}/corpora/fixture/sample.zsh"

cat > "${SCRATCH}/bin/git" <<'EOF'
#!/usr/bin/env bash
[[ "$*" == 'rev-parse HEAD' ]] || exit 2
printf '%s\n' aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
EOF

cat > "${SCRATCH}/bin/zshellcheck" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
file=${!#}
case "${SWEEP_CASE}:$1:${file##*/}" in
    panic-original:-no-banner:sample.zsh)
        echo 'panic: original scan' >&2; exit 1 ;;
    crash-fixed:-no-banner:work)
        exit 2 ;;
    crash-first:-fix:work)
        exit 3 ;;
    write-failure:-fix:work)
        echo 'fix: write failed for fixture: permission denied' >&2; exit 1 ;;
    crash-second:-fix:work2)
        exit 2 ;;
    trace-second:-fix:work2)
        echo 'goroutine 1 [running]:'; exit 1 ;;
esac
if [[ "${SWEEP_CASE}" == corruption ]]; then
    if [[ "$1" == -fix ]] && ! grep -q '# broken' "${file}"; then
        printf '# broken\n' >> "${file}"
    elif [[ "$1" == -no-banner ]] && grep -q '# broken' "${file}"; then
        echo 'Parser Error at line 2: fixture failure' >&2
        exit 1
    fi
fi
if [[ "${SWEEP_CASE}" == findings ]]; then
    echo 'fixture finding'
    exit 1
fi
EOF
chmod +x "${SCRATCH}/bin/git" "${SCRATCH}/bin/zshellcheck"

for spec in clean:0 findings:0 corruption:1 panic-original:2 crash-fixed:2 \
    crash-first:2 write-failure:2 crash-second:2 trace-second:2; do
    rc=0
    PATH="${SCRATCH}/bin:${PATH}" SWEEP_CASE="${spec%:*}" \
        ZSHELLCHECK_BIN="${SCRATCH}/bin/zshellcheck" \
        PARSER_SWEEP_WORK="${SCRATCH}/corpora" \
        bash "${SCRATCH}/repo/scripts/fix-corpus-sweep.sh" > "${SCRATCH}/result" 2>&1 || rc=$?
    if [[ "${rc}" != "${spec#*:}" ]]; then
        printf 'FAIL %s: expected exit %s, got %s\n' "${spec%:*}" "${spec#*:}" "${rc}" >&2
        cat "${SCRATCH}/result" >&2
        exit 1
    fi
    if [[ "${spec%:*}" == corruption ]]; then
        grep -q 'parser errors 0 -> 1' "${SCRATCH}/result"
    fi
    printf 'PASS %s\n' "${spec%:*}"
done
