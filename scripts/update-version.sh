#!/usr/bin/env bash
# Derive pkg/version/version.go from the current kata count.
#
# Base formula: MAJOR.MINOR.PATCH from (count / 1000, (count % 1000) / 100, count % 100).
# Hotfix offset: scripts/HOTFIX holds a monotonic counter that is added to
# PATCH so bugfix releases can ship between kata additions without waiting
# for the next kata. Increment the HOTFIX file by 1 for each patch release.
set -euo pipefail

KATA_COUNT=$(ls -1 pkg/katas/zc*.go 2>/dev/null | wc -l)

HOTFIX_FILE="scripts/HOTFIX"
if [[ -f "$HOTFIX_FILE" ]]; then
    HOTFIX=$(<"$HOTFIX_FILE")
else
    HOTFIX=0
fi

MAJOR=$((KATA_COUNT / 1000))
MINOR=$(( (KATA_COUNT % 1000) / 100 ))
PATCH=$(( (KATA_COUNT % 100) + HOTFIX ))
VERSION="${MAJOR}.${MINOR}.${PATCH}"

cat > pkg/version/version.go << EOF
package version

// Version is the current version of ZShellCheck.
// It is derived from the kata count plus the scripts/HOTFIX offset.
// ${KATA_COUNT} Katas + ${HOTFIX} hotfix = ${VERSION}
const Version = "${VERSION}"
EOF

echo "Updated version to ${VERSION} (${KATA_COUNT} katas, hotfix offset ${HOTFIX})"
