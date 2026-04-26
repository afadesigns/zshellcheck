// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCpuprofileError(t *testing.T) {
	resetFlags()
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"zshellcheck", "-cpuprofile", "/no/such/dir/cpu.out", "x.zsh"}
	got := run()
	if got != 1 {
		t.Errorf("expected exit 1 for unwritable cpuprofile path, got %d", got)
	}
}

func TestRunCpuprofileSuccess(t *testing.T) {
	dir := t.TempDir()
	prof := filepath.Join(dir, "cpu.out")
	src := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(src, []byte("echo hi\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"zshellcheck", "-cpuprofile", prof, "-no-banner", src}
	_ = run()
}

func TestRunInvalidSeverity(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(src, []byte("echo hi\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"zshellcheck", "-severity", "frobozz", "-no-banner", src}
	got := run()
	if got != 1 {
		t.Errorf("expected exit 1 for invalid severity, got %d", got)
	}
}

func TestRunValidSeverity(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(src, []byte("echo hi\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"zshellcheck", "-severity", "error,warning", "-no-banner", src}
	_ = run()
}

func TestRunDiffMode(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(src, []byte("result=`which git`\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"zshellcheck", "-diff", "-no-banner", src}
	_ = run()
}

func TestRunFixMode(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(src, []byte("result=`which git`\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"zshellcheck", "-fix", "-no-banner", src}
	_ = run()
}

func TestRunFixModeMultipleFiles(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.zsh")
	b := filepath.Join(dir, "b.zsh")
	if err := os.WriteFile(a, []byte("result=`which git`\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("echo $arr[1]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"zshellcheck", "-fix", "-no-banner", a, b}
	_ = run()
}

func TestRunVerboseAndNoColor(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(src, []byte("echo hi\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	resetFlags()
	old := os.Args
	defer func() { os.Args = old }()
	os.Args = []string{"zshellcheck", "-verbose", "-no-color", "-no-banner", src}
	_ = run()
}
