#!/usr/bin/env python3
"""Consolidate per-kata test files into hundred-bucket files.

Each `pkg/katas/katatests/zcXXXX_test.go` carries one TestZCXXXX(t *testing.T)
plus optional helpers. Layout per file:

    // SPDX-License-Identifier: MIT
    // Copyright the ZShellCheck contributors.
    package katas

    import (
            ...
    )

    func TestZCXXXX(t *testing.T) {
            ...
    }

This script merges every file into a bucket file `zcNNNNs_test.go`
keyed by `floor(id / 100) * 100`. Each bucket file gets a unified
header + the union of imports across its contributing files. The
function bodies are concatenated unchanged.

Run from repo root:
    uv run python scripts/consolidate-katatests.py
or:
    python3 scripts/consolidate-katatests.py
"""

from __future__ import annotations

import re
from pathlib import Path

DIR = Path("pkg/katas/katatests")
PER_FILE_RE = re.compile(r"zc(\d{4})_test\.go")
IMPORT_BLOCK_RE = re.compile(r"^import \(\n([\s\S]*?)\n\)\n", re.MULTILINE)


def parse_imports(block: str) -> list[str]:
    out: list[str] = []
    for line in block.splitlines():
        line = line.strip()
        if not line:
            continue
        out.append(line)
    return out


def main() -> None:
    files = sorted(DIR.glob("zc*_test.go"))
    if not files:
        raise SystemExit(f"no kata test files in {DIR}")

    buckets: dict[int, dict] = {}
    for path in files:
        match = PER_FILE_RE.fullmatch(path.name)
        if not match:
            continue
        kid = int(match.group(1))
        bucket = (kid // 100) * 100
        text = path.read_text()

        # Strip the SPDX/copyright/package preamble; keep the import
        # list and the body (function declarations).
        body_match = re.search(
            r"^package katas\n+import \(\n([\s\S]*?)\n\)\n+([\s\S]*)$",
            text,
            re.MULTILINE,
        )
        if not body_match:
            raise SystemExit(f"unexpected layout: {path}")
        imports_text = body_match.group(1)
        body = body_match.group(2).rstrip() + "\n"

        b = buckets.setdefault(
            bucket,
            {
                "imports": set(),
                "bodies": [],
            },
        )
        b["imports"].update(parse_imports(imports_text))
        b["bodies"].append(body)

    # Emit one bucket file per group.
    for bucket, payload in buckets.items():
        bucket_path = DIR / f"zc{bucket}s_test.go"
        imports_sorted = sorted(payload["imports"])
        # Group stdlib imports above third-party (gofmt convention):
        stdlib = [i for i in imports_sorted if not i.startswith('"github.com/')]
        third = [i for i in imports_sorted if i.startswith('"github.com/')]
        import_lines = ["\t" + line for line in stdlib]
        if stdlib and third:
            import_lines.append("")
        import_lines.extend("\t" + line for line in third)

        header = (
            "// SPDX-License-Identifier: MIT\n"
            "// Copyright the ZShellCheck contributors.\n"
            "package katas\n"
            "\n"
            "import (\n"
            + "\n".join(import_lines)
            + "\n"
            ")\n"
            "\n"
        )
        bucket_path.write_text(header + "\n".join(payload["bodies"]).rstrip() + "\n")

    # Drop the originals.
    for path in files:
        path.unlink()

    print(f"consolidated {len(files)} files into {len(buckets)} buckets")


if __name__ == "__main__":
    main()
