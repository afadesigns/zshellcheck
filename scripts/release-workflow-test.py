#!/usr/bin/env python3
"""Exercise release input validation and publication failure handling offline."""

import base64
import hashlib
import os
from pathlib import Path
import subprocess
import tempfile
import unittest


ROOT = Path(__file__).resolve().parents[1]


def step_script(workflow, name):
    """Read an indented run block without adding a YAML dependency."""
    lines = (ROOT / ".github/workflows" / workflow).read_text().splitlines()
    start = lines.index("      - name: " + name)
    start = lines.index("        run: |", start) + 1
    result = []
    for line in lines[start:]:
        if line and not line.startswith("          "):
            break
        result.append(line[10:])
    if not result:
        raise ValueError("Missing step script: " + name)
    return "\n".join(result) + "\n"


VALIDATE = step_script("release-provenance.yml", "Validate subject records")
PUBLISH = step_script("release-build.yml", "Verify draft assets and publish")
NAMES = ["checksums.txt", "release.tar.gz", "release.zip", "release.deb",
         "release.rpm", "release.apk", "release.tar.gz.sbom.json"]


class ReleaseWorkflowTests(unittest.TestCase):
    def setUp(self):
        self.scratch = tempfile.TemporaryDirectory()
        self.addCleanup(self.scratch.cleanup)
        self.root = Path(self.scratch.name)
        self.env = dict(os.environ, RUNNER_TEMP=str(self.root))

    def run_step(self, script):
        return subprocess.run(
            ["bash", "-euo", "pipefail", "-c", script], env=self.env,
            cwd=self.root, capture_output=True, text=True, check=False,
        )

    def validate(self, records):
        self.env["SUBJECTS"] = base64.b64encode(records.encode()).decode()
        return self.run_step(VALIDATE)

    def test_subject_validation(self):
        valid = "".join("a" * 64 + "  " + name + "\n" for name in NAMES)
        cases = {
            "valid": (valid, True),
            "empty": ("", False),
            "duplicate": (valid + valid.splitlines()[0] + "\n", False),
            "malformed": (valid + "ignored-by-upstream\n", False),
            "path": (valid + "a" * 64 + "  ../escape\n", False),
            "shell": (valid + "a" * 64 + "  $(touch-owned)\n", False),
            "sha512": (valid.replace("a" * 64, "a" * 128), False),
            "missing-class": (valid.replace("a" * 64 + "  release.rpm\n", ""), False),
            "too-many": (valid + "".join("a" * 64 + f"  extra-{n}\n" for n in range(1024)), False),
        }
        for name, (records, succeeds) in cases.items():
            with self.subTest(name=name):
                result = self.validate(records)
                self.assertEqual(result.returncode == 0, succeeds, result.stderr)
                if succeeds:
                    self.assertEqual((self.root / "subjects.txt").read_text(), records)
        self.env["SUBJECTS"] = "invalid!base64"
        self.assertNotEqual(self.run_step(VALIDATE).returncode, 0)

    def prepare_publish(self, mode):
        assets = self.root / "assets"
        assets.mkdir()
        records = []
        for name in NAMES:
            (assets / name).write_text(name)
            records.append(hashlib.sha256(name.encode()).hexdigest() + "  " + name)
            for suffix in (".pem", ".sig"):
                (assets / (name + suffix)).write_text("fixture")
        (assets / "zshellcheck-provenance.json").write_text("fixture")
        self.env.update(
            SUBJECTS=base64.b64encode(("\n".join(records) + "\n").encode()).decode(),
            GITHUB_REPOSITORY="example/project", GITHUB_SHA="a" * 40,
            GITHUB_REF="refs/tags/v1.8.0", GITHUB_REF_NAME="v1.8.0",
            VERIFY_FAIL=str(mode == "verification-failure"),
            IS_DRAFT=str(mode != "published").lower(),
        )
        if mode == "tampered":
            (assets / "release.zip").write_text("tampered")
        if mode == "missing-signature":
            (assets / "release.zip.sig").unlink()
        if mode == "extra-archive":
            (assets / "unattested.tar.gz").write_text("unattested")
        if mode == "extra-file":
            (assets / "unattested.txt").write_text("unattested")
        bin_dir = self.root / "bin"
        bin_dir.mkdir()
        fake_gh = bin_dir / "gh"
        fake_gh.write_text('''#!/usr/bin/env python3
import json, os, pathlib, shutil, sys
root = pathlib.Path(os.environ['RUNNER_TEMP'])
args = sys.argv[1:]
with (root / 'calls.jsonl').open('a') as log:
    log.write(json.dumps(args) + '\\n')
if args[:2] == ['release', 'view']:
    print(os.environ['IS_DRAFT'])
elif args[:2] == ['release', 'download']:
    shutil.copytree(root / 'assets', root / 'release', dirs_exist_ok=True)
elif args[:2] == ['attestation', 'verify']:
    flags = dict(zip(args[3::2], args[4::2]))
    assert flags['--repo'] == os.environ['GITHUB_REPOSITORY']
    assert flags['--signer-digest'] == os.environ['GITHUB_SHA']
    assert flags['--source-digest'] == os.environ['GITHUB_SHA']
    assert flags['--source-ref'] == os.environ['GITHUB_REF']
    assert flags['--bundle'] == 'zshellcheck-provenance.json'
    expected = ('https://github.com/' + os.environ['GITHUB_REPOSITORY'] +
                '/.github/workflows/release-provenance.yml@' + os.environ['GITHUB_REF'])
    assert flags['--cert-identity'] == expected
    assert args[-1] == '--deny-self-hosted-runners'
    if os.environ['VERIFY_FAIL'] == 'True':
        sys.exit(1)
elif args[:2] == ['release', 'edit']:
    assert '--draft=false' in args
    (root / 'published').touch()
else:
    sys.exit('Unexpected command')
''')
        fake_gh.chmod(0o755)
        self.env["PATH"] = str(bin_dir) + os.pathsep + self.env["PATH"]

    def test_publication_requires_every_check(self):
        for mode in ("valid", "tampered", "verification-failure", "missing-signature",
                     "extra-archive", "extra-file", "published"):
            with self.subTest(mode=mode), tempfile.TemporaryDirectory(dir=self.root) as case:
                previous_root = self.root
                self.root = Path(case)
                self.env["RUNNER_TEMP"] = case
                self.prepare_publish(mode)
                result = self.run_step(PUBLISH)
                self.assertEqual(result.returncode == 0, mode == "valid", result.stderr)
                self.assertEqual((self.root / "published").exists(), mode == "valid")
                self.root = previous_root


if __name__ == "__main__":
    unittest.main()
