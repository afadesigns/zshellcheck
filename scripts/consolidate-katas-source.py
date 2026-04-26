#!/usr/bin/env python3
"""Consolidate per-kata source files into hundred-bucket files.

Each `pkg/katas/zcXXXX.go` carries one or more `init()` blocks
registering the kata, plus its `checkZCXXXX` / `fixZCXXXX` /
helper functions. The file layout is:

    // SPDX-License-Identifier: MIT
    // Copyright the ZShellCheck contributors.
    package katas

    import (
            ...
    )

    func init() { ... }
    func checkZCXXXX(...) { ... }
    func fixZCXXXX(...) { ... }
    ...

This script merges every kata source file into a bucket file
`zcNNNNs.go` keyed by `floor(id / 100) * 100`. The bucket file
gets a unified header + the union of imports across its members.
Function bodies, init blocks, and package-level helpers are
concatenated unchanged.

GitHub's web view truncates directory listings at 1,000 entries.
With 1,000 kata source files in `pkg/katas/`, anything else in
that directory was hidden. Bucketing drops the kata source count
to 11 files and keeps the directory navigable.

Run from repo root:
    python3 scripts/consolidate-katas-source.py
"""

from __future__ import annotations

import re
from pathlib import Path

DIR = Path("pkg/katas")
PER_FILE_RE = re.compile(r"zc(\d{4})\.go")


def parse_imports(block: str) -> list[str]:
    out: list[str] = []
    for line in block.splitlines():
        line = line.strip()
        if not line:
            continue
        out.append(line)
    return out


def main() -> None:
    files = sorted(p for p in DIR.glob("zc*.go") if PER_FILE_RE.fullmatch(p.name))
    if not files:
        raise SystemExit(f"no kata source files in {DIR}")

    buckets: dict[int, dict] = {}
    for path in files:
        match = PER_FILE_RE.fullmatch(path.name)
        if not match:
            continue
        kid = int(match.group(1))
        bucket = (kid // 100) * 100
        text = path.read_text()
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

    for bucket, payload in buckets.items():
        bucket_path = DIR / f"zc{bucket}s.go"
        imports_sorted = sorted(payload["imports"])
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

    for path in files:
        path.unlink()

    print(f"consolidated {len(files)} files into {len(buckets)} buckets")


if __name__ == "__main__":
    main()
