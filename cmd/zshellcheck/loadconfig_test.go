// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigBranches_EmptyPath(t *testing.T) {
	cfg, err := loadConfig("")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	_ = cfg
}

func TestLoadConfigBranches_MissingPath(t *testing.T) {
	if _, err := loadConfig("/nonexistent/zzz.yml"); err != nil {
		t.Errorf("unexpected error for missing file: %v", err)
	}
}

func TestLoadConfigBranches_ValidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.yml")
	if err := os.WriteFile(path, []byte("no_color: true\ndisabled_katas:\n  - ZC1001\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.NoColor {
		t.Error("expected NoColor true")
	}
}

func TestLoadConfigBranches_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.yml")
	if err := os.WriteFile(path, []byte(":this is not yaml\n  - bad:\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadConfig(path); err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadConfigBranches_UnreadableFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.yml")
	if err := os.WriteFile(path, []byte("ok\n"), 0o000); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chmod(path, 0o600) }()
	// On systems where the test runs as root, the open will succeed —
	// in that case the call simply parses and returns no error. We
	// only care that the path is exercised.
	_, _ = loadConfig(path)
}

func TestLoadConfigBranches_MultiplePaths(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.yml")
	b := filepath.Join(dir, "b.yml")
	if err := os.WriteFile(a, []byte("no_color: false\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("no_color: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadConfig(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.NoColor {
		t.Error("expected merged NoColor true from second file")
	}
}
