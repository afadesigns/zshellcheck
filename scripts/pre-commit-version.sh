#!/usr/bin/env bash
# Pre-commit hook to auto-update version based on kata count
bash scripts/update-version.sh
git add pkg/version/version.go
