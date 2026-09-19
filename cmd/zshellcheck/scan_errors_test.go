// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func runCaptured(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	oldArgs, oldFlags := os.Args, flag.CommandLine
	oldOut, oldErr := os.Stdout, os.Stderr
	defer func() {
		os.Args, flag.CommandLine = oldArgs, oldFlags
		os.Stdout, os.Stderr = oldOut, oldErr
	}()
	os.Args = append([]string{"zshellcheck"}, args...)
	resetFlags()
	os.Stdout = captureFile(t, "stdout")
	os.Stderr = captureFile(t, "stderr")
	code := run()
	return code, readTestFile(t, os.Stdout.Name()), readTestFile(t, os.Stderr.Name())
}

func captureFile(t *testing.T, name string) *os.File {
	t.Helper()
	f, err := os.Create(filepath.Join(t.TempDir(), name))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })
	return f
}

func readTestFile(t *testing.T, filename string) string {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func writeTestFile(t *testing.T, filename, content string) {
	t.Helper()
	if err := os.WriteFile(filename, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRunScanFailureModes(t *testing.T) {
	dir := t.TempDir()
	clean := filepath.Join(dir, "clean.zsh")
	bad := filepath.Join(dir, "bad.zsh")
	missing := filepath.Join(dir, "missing.zsh")
	writeTestFile(t, clean, "true\n")
	writeTestFile(t, bad, ") ]\n")
	for _, inputs := range [][]string{{missing}, {clean, missing}, {bad}} {
		for _, mode := range []string{"text", "json", "sarif", "statistics", "add-noka", "baseline-write"} {
			t.Run(mode+"/"+strings.Join(inputs, "+"), func(t *testing.T) {
				base := filepath.Join(t.TempDir(), "baseline.json")
				writeTestFile(t, base, "previous baseline\n")
				args := scanModeArgs(mode, base)
				code, out, errOut := runCaptured(t, append(args, inputs...)...)
				if code != 1 || errOut == "" {
					t.Errorf("failed scan: exit=%d stderr=%q", code, errOut)
				}
				if got := readTestFile(t, base); got != "previous baseline\n" {
					t.Errorf("failed scan replaced baseline: %q", got)
				}
				if (mode == "json" || mode == "sarif") && !json.Valid([]byte(out)) {
					t.Errorf("invalid structured output: %q", out)
				}
			})
		}
	}
}

func scanModeArgs(mode, baseline string) []string {
	args := []string{"-no-banner"}
	switch mode {
	case "text", "json", "sarif":
		return append(args, "-format", mode)
	case "baseline-write":
		return append(args, "-baseline-write", baseline)
	default:
		return append(args, "-"+mode)
	}
}

func TestRunStaleStructuredOutput(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "stale.zsh")
	writeTestFile(t, filename, "true # noka: ZC1523\n")
	for _, format := range []string{"json", "sarif"} {
		t.Run(format, func(t *testing.T) {
			code, out, errOut := runCaptured(t, "-detect-stale-noka", "-format", format, filename)
			if code != 1 || !json.Valid([]byte(out)) {
				t.Errorf("stale report: exit=%d stdout=%q", code, out)
			}
			if !strings.Contains(errOut, "stale `# noka: ZC1523`") {
				t.Errorf("missing stale diagnostic: %q", errOut)
			}
		})
	}
}

func TestRunUnreadableInput(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "private.zsh")
	writeTestFile(t, filename, "true\n")
	if err := os.Chmod(filename, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(filename, 0o600) })
	if _, err := os.ReadFile(filename); err == nil {
		t.Skip("filesystem permits reading mode-000 files")
	}
	code, _, errOut := runCaptured(t, "-no-banner", filename)
	if code != 1 || !strings.Contains(errOut, "Error reading file") {
		t.Errorf("unreadable input: exit=%d stderr=%q", code, errOut)
	}
}

func TestAddNokaWriteFailure(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "readonly.zsh")
	writeTestFile(t, filename, "echo hello\n")
	if err := os.Chmod(filename, 0o400); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(filename, 0o600) })
	if f, err := os.OpenFile(filename, os.O_WRONLY, 0); err == nil {
		_ = f.Close()
		t.Skip("filesystem permits writing mode-400 files")
	}
	code, _, errOut := runCaptured(t, "-no-banner", "-add-noka", filename)
	if code != 1 || !strings.Contains(errOut, "add-noka:") {
		t.Errorf("failed directive write: exit=%d stderr=%q", code, errOut)
	}
}

func TestDirectoryWalkFailure(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chmod(dir, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
	if _, err := os.ReadDir(dir); err == nil {
		t.Skip("filesystem permits reading mode-000 directories")
	}
	code, _, errOut := runCaptured(t, "-no-banner", "-statistics", dir)
	if code != 1 || !strings.Contains(errOut, "Error walking directory") {
		t.Errorf("failed directory walk: exit=%d stderr=%q", code, errOut)
	}
}

func TestFailedScanDoesNotCreateBaseline(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "baseline.json")
	code, _, _ := runCaptured(t, "-no-banner", "-baseline-write", base, filepath.Join(dir, "missing.zsh"))
	if _, err := os.Stat(base); code != 1 || !os.IsNotExist(err) {
		t.Errorf("failed scan created baseline: exit=%d stat=%v", code, err)
	}
}

func TestReportFailureState(t *testing.T) {
	opts := buildFixOpts(false, false, false, false)
	var errOut bytes.Buffer
	violations := []katas.Violation{{KataID: "ZC1001", Line: 1, Column: 1, Message: "test", Level: katas.SeverityWarning}}
	emitReport("test.zsh", failingWriter{}, &errOut, "text", config.DefaultConfig(), violations, []byte("test\n"), katas.Registry, opts)
	if !opts.hasFailed() || !strings.Contains(errOut.String(), "Error reporting") {
		t.Errorf("report failure lost: %q", errOut.String())
	}
}
