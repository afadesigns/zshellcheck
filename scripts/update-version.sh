#!/usr/bin/env bash
# Automatically update version.go based on kata count
set -euo pipefail

KATA_COUNT=$(ls -1 pkg/katas/zc*.go 2>/dev/null | wc -l)
MAJOR=$((KATA_COUNT / 1000))
MINOR=$(( (KATA_COUNT % 1000) / 100 ))
PATCH=$((KATA_COUNT % 100))
VERSION="${MAJOR}.${MINOR}.${PATCH}"

cat > pkg/version/version.go << EOF
package version

// Version is the current version of ZShellCheck.
// It is calculated based on the number of implemented Katas.
// ${KATA_COUNT} Katas = ${VERSION}
const Version = "${VERSION}"
EOF

echo "Updated version to ${VERSION} (${KATA_COUNT} katas)"
