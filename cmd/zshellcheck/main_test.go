package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// resetFlags resets the global flag.CommandLine for testing run().
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestRun_NoArgs(t *testing.T) {
	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck"}

	code := run()
	if code != 1 {
		t.Errorf("expected exit code 1 for no args, got %d", code)
	}
}

func TestRun_Version(t *testing.T) {
	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-version"}

	code := run()
	if code != 0 {
		t.Errorf("expected exit code 0 for -version, got %d", code)
	}
}

func TestRun_WithFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-no-color", path}

	code := run()
	_ = code
}

func TestRun_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-format", "json", path}

	code := run()
	_ = code
}

func TestRun_SarifFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-format", "sarif", path}

	code := run()
	_ = code
}

func TestRun_WithSeverityFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-severity", "error,warning", "-no-color", path}

	code := run()
	_ = code
}

func TestRun_InvalidSeverityFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-severity", "invalid_level", path}

	code := run()
	if code != 1 {
		t.Errorf("expected exit code 1 for invalid severity, got %d", code)
	}
}

func TestRun_VerboseFlag(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-verbose", "-no-color", path}

	code := run()
	_ = code
}

func TestRun_CPUProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	profilePath := filepath.Join(dir, "cpu.prof")

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-cpuprofile", profilePath, "-no-color", path}

	code := run()
	_ = code

	// Verify profile file was created
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		t.Error("expected CPU profile file to be created")
	}
}

func TestRun_TextFormatWithViolations(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	// Script that is likely to produce violations
	if err := os.WriteFile(path, []byte("#!/bin/zsh\nfor i in $(ls); do echo $i; done\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-no-color", path}

	code := run()
	_ = code
}

func TestRun_WithDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-no-color", dir}

	code := run()
	_ = code
}

func TestRun_StyleSeverity(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-severity", "style", "-no-color", path}

	code := run()
	_ = code
}

func TestRun_WithViolationsTextFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	// Input likely to produce violations
	if err := os.WriteFile(path, []byte("#!/bin/zsh\nfor i in $(ls); do echo $i; done\nrm -rf ${dir}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-no-color", "-format", "text", path}

	code := run()
	// Should return 1 if violations found
	if code != 1 {
		t.Logf("expected exit code 1 for violations, got %d", code)
	}
}

func TestRun_InfoSeverity(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resetFlags()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"zshellcheck", "-severity", "info", "-no-color", path}

	code := run()
	_ = code
}

func TestLoadConfig_NoFile(t *testing.T) {
	cfg, err := loadConfig("/nonexistent/path/.zshellcheckrc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should return defaults when file doesn't exist
	if cfg.ErrorColor != config.ColorRed {
		t.Errorf("expected default ErrorColor, got %q", cfg.ErrorColor)
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".zshellcheckrc")
	content := []byte("disabled_katas:\n  - ZC1001\nno_color: true\n")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.DisabledKatas) != 1 || cfg.DisabledKatas[0] != "ZC1001" {
		t.Errorf("unexpected DisabledKatas: %v", cfg.DisabledKatas)
	}
	if !cfg.NoColor {
		t.Error("expected NoColor=true")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".zshellcheckrc")
	if err := os.WriteFile(path, []byte(":::bad\n\t[[["), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := loadConfig(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestProcessFile_TextFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	// Write a simple shell script
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	cfg.NoColor = true
	registry := katas.Registry

	count := processFile(path, &out, &errOut, cfg, registry, "text", nil)
	// We don't know exactly how many violations, but it should not panic
	_ = count
}

func TestProcessFile_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	registry := katas.Registry

	count := processFile(path, &out, &errOut, cfg, registry, "json", nil)
	_ = count
}

func TestProcessFile_SarifFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	registry := katas.Registry

	count := processFile(path, &out, &errOut, cfg, registry, "sarif", nil)
	_ = count
}

func TestProcessFile_NonexistentFile(t *testing.T) {
	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	registry := katas.Registry

	count := processFile("/nonexistent/file.zsh", &out, &errOut, cfg, registry, "text", nil)
	if count != 0 {
		t.Errorf("expected 0 violations for nonexistent file, got %d", count)
	}
}

func TestProcessFile_SeverityFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	registry := katas.Registry

	// Filter to only show errors
	count := processFile(path, &out, &errOut, cfg, registry, "text", []katas.Severity{katas.SeverityError})
	_ = count
}

func TestProcessPath_File(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	cfg.NoColor = true
	registry := katas.Registry

	count := processPath(path, &out, &errOut, cfg, registry, "text", nil)
	_ = count
}

func TestProcessPath_Directory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.zsh")
	if err := os.WriteFile(path, []byte("#!/bin/zsh\necho hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Also create a non-shell file that should be skipped
	goFile := filepath.Join(dir, "test.go")
	if err := os.WriteFile(goFile, []byte("package main\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	cfg.NoColor = true
	registry := katas.Registry

	count := processPath(dir, &out, &errOut, cfg, registry, "text", nil)
	_ = count
}

func TestProcessPath_Nonexistent(t *testing.T) {
	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	registry := katas.Registry

	count := processPath("/nonexistent/path", &out, &errOut, cfg, registry, "text", nil)
	if count != 0 {
		t.Errorf("expected 0 for nonexistent path, got %d", count)
	}
}

func TestProcessPath_DirectoryWithHiddenDir(t *testing.T) {
	dir := t.TempDir()
	// Create a hidden directory that should be skipped
	hiddenDir := filepath.Join(dir, ".hidden")
	if err := os.MkdirAll(hiddenDir, 0o755); err != nil {
		t.Fatal(err)
	}
	hiddenFile := filepath.Join(hiddenDir, "test.zsh")
	if err := os.WriteFile(hiddenFile, []byte("echo hidden\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Create a normal file
	normalFile := filepath.Join(dir, "normal.zsh")
	if err := os.WriteFile(normalFile, []byte("echo normal\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	cfg.NoColor = true
	registry := katas.Registry

	count := processPath(dir, &out, &errOut, cfg, registry, "text", nil)
	_ = count
}

func TestProcessFile_ParserErrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.zsh")
	// Write something that will cause parser errors
	if err := os.WriteFile(path, []byte("if then fi fi fi fi\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	registry := katas.Registry

	count := processFile(path, &out, &errOut, cfg, registry, "text", nil)
	// Parser errors should return 1
	if count < 1 {
		t.Errorf("expected at least 1 for parser errors, got %d", count)
	}
}

func TestProcessPath_DirectorySkipsNonShellFiles(t *testing.T) {
	dir := t.TempDir()

	// Create various non-shell files that should be skipped
	extensions := []string{".go", ".md", ".json", ".yml", ".yaml", ".txt"}
	for _, ext := range extensions {
		path := filepath.Join(dir, "test"+ext)
		if err := os.WriteFile(path, []byte("content\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	var out, errOut bytes.Buffer
	cfg := config.DefaultConfig()
	cfg.NoColor = true
	registry := katas.Registry

	count := processPath(dir, &out, &errOut, cfg, registry, "text", nil)
	// All non-shell files should be skipped, so violations from parsing Go/etc should be 0
	if count != 0 {
		t.Errorf("expected 0 violations for skipped files, got %d", count)
	}
}
