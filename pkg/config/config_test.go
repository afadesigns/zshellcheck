package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.ErrorColor != ColorRed {
		t.Errorf("expected ErrorColor=%q, got %q", ColorRed, cfg.ErrorColor)
	}
	if cfg.WarningColor != ColorYellow {
		t.Errorf("expected WarningColor=%q, got %q", ColorYellow, cfg.WarningColor)
	}
	if cfg.InfoColor != ColorCyan {
		t.Errorf("expected InfoColor=%q, got %q", ColorCyan, cfg.InfoColor)
	}
	if cfg.IDColor != ColorRed {
		t.Errorf("expected IDColor=%q, got %q", ColorRed, cfg.IDColor)
	}
	if cfg.TitleColor != ColorCyan {
		t.Errorf("expected TitleColor=%q, got %q", ColorCyan, cfg.TitleColor)
	}
	if cfg.MessageColor != ColorReset {
		t.Errorf("expected MessageColor=%q, got %q", ColorReset, cfg.MessageColor)
	}
	if cfg.LineColor != ColorCyan {
		t.Errorf("expected LineColor=%q, got %q", ColorCyan, cfg.LineColor)
	}
	if cfg.ColumnColor != ColorYellow {
		t.Errorf("expected ColumnColor=%q, got %q", ColorYellow, cfg.ColumnColor)
	}
	if cfg.NoColor {
		t.Error("expected NoColor=false")
	}
	if cfg.Verbose {
		t.Error("expected Verbose=false")
	}
	if len(cfg.DisabledKatas) != 0 {
		t.Errorf("expected empty DisabledKatas, got %v", cfg.DisabledKatas)
	}
}

func TestMergeConfig_OverridesAllFields(t *testing.T) {
	base := DefaultConfig()
	override := Config{
		DisabledKatas: []string{"ZC1001"},
		ErrorColor:    "custom-error",
		WarningColor:  "custom-warning",
		InfoColor:     "custom-info",
		IDColor:       "custom-id",
		TitleColor:    "custom-title",
		MessageColor:  "custom-message",
		LineColor:     "custom-line",
		ColumnColor:   "custom-column",
		NoColor:       true,
		Verbose:       true,
	}

	merged := MergeConfig(base, override)

	if len(merged.DisabledKatas) != 1 || merged.DisabledKatas[0] != "ZC1001" {
		t.Errorf("expected DisabledKatas=[ZC1001], got %v", merged.DisabledKatas)
	}
	if merged.ErrorColor != "custom-error" {
		t.Errorf("expected ErrorColor=custom-error, got %s", merged.ErrorColor)
	}
	if merged.WarningColor != "custom-warning" {
		t.Errorf("expected WarningColor=custom-warning, got %s", merged.WarningColor)
	}
	if merged.InfoColor != "custom-info" {
		t.Errorf("expected InfoColor=custom-info, got %s", merged.InfoColor)
	}
	if merged.IDColor != "custom-id" {
		t.Errorf("expected IDColor=custom-id, got %s", merged.IDColor)
	}
	if merged.TitleColor != "custom-title" {
		t.Errorf("expected TitleColor=custom-title, got %s", merged.TitleColor)
	}
	if merged.MessageColor != "custom-message" {
		t.Errorf("expected MessageColor=custom-message, got %s", merged.MessageColor)
	}
	if merged.LineColor != "custom-line" {
		t.Errorf("expected LineColor=custom-line, got %s", merged.LineColor)
	}
	if merged.ColumnColor != "custom-column" {
		t.Errorf("expected ColumnColor=custom-column, got %s", merged.ColumnColor)
	}
	if !merged.NoColor {
		t.Error("expected NoColor=true")
	}
	if !merged.Verbose {
		t.Error("expected Verbose=true")
	}
}

func TestMergeConfig_EmptyOverridePreservesBase(t *testing.T) {
	base := DefaultConfig()
	override := Config{} // all zero-values

	merged := MergeConfig(base, override)

	// String fields with zero-value ("") should NOT override base
	if merged.ErrorColor != base.ErrorColor {
		t.Errorf("expected ErrorColor preserved as %q, got %q", base.ErrorColor, merged.ErrorColor)
	}
	if merged.WarningColor != base.WarningColor {
		t.Errorf("expected WarningColor preserved as %q, got %q", base.WarningColor, merged.WarningColor)
	}
	// Bool fields are always overridden (direct assignment)
	if merged.NoColor != false {
		t.Error("expected NoColor=false from zero-value override")
	}
}

func TestNewConfigFromYAML_ValidFile(t *testing.T) {
	content := []byte("disabled_katas:\n  - ZC1001\n  - ZC1002\nno_color: true\nverbose: true\n")
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := NewConfigFromYAML(path)
	if err != nil {
		t.Fatalf("NewConfigFromYAML() error: %v", err)
	}

	if len(cfg.DisabledKatas) != 2 {
		t.Fatalf("expected 2 disabled katas, got %d", len(cfg.DisabledKatas))
	}
	if cfg.DisabledKatas[0] != "ZC1001" || cfg.DisabledKatas[1] != "ZC1002" {
		t.Errorf("unexpected disabled katas: %v", cfg.DisabledKatas)
	}
	if !cfg.NoColor {
		t.Error("expected NoColor=true")
	}
	if !cfg.Verbose {
		t.Error("expected Verbose=true")
	}
	// Defaults should still be set for non-overridden fields
	if cfg.ErrorColor != ColorRed {
		t.Errorf("expected ErrorColor default preserved, got %q", cfg.ErrorColor)
	}
}

func TestNewConfigFromYAML_FileNotFound(t *testing.T) {
	_, err := NewConfigFromYAML("/nonexistent/path/config.yml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestNewConfigFromYAML_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yml")
	if err := os.WriteFile(path, []byte(":::invalid:::yaml\n\t\t[[["), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := NewConfigFromYAML(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestNewConfigFromYAML_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.yml")
	if err := os.WriteFile(path, []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := NewConfigFromYAML(path)
	if err != nil {
		t.Fatalf("NewConfigFromYAML() error: %v", err)
	}

	// Should return defaults
	if cfg.ErrorColor != ColorRed {
		t.Errorf("expected default ErrorColor, got %q", cfg.ErrorColor)
	}
}
