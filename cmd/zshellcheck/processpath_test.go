// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

func TestProcessPathSingleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	processPath(path, &out, &errOut, config.DefaultConfig(), katas.Registry, "text", nil, fixOptions{})
}

func TestProcessPathDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.zsh"), []byte("#!/bin/zsh\necho a\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.zsh"), []byte("#!/bin/zsh\necho b\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "skip.go"), []byte("package main\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "skip.md"), []byte("hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, ".hidden"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".hidden", "h.zsh"), []byte("echo hidden\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	processPath(dir, &out, &errOut, config.DefaultConfig(), katas.Registry, "text", nil, fixOptions{})
}

func TestProcessPathMissing(t *testing.T) {
	var out, errOut bytes.Buffer
	got := processPath("/no/such/dir/zzz", &out, &errOut, config.DefaultConfig(), katas.Registry, "text", nil, fixOptions{})
	if got != 0 {
		t.Errorf("expected 0 on stat error, got %d", got)
	}
}
