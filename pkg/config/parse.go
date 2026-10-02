// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package config

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// fileConfig is the on-disk YAML shape: the scalar settings plus the two
// ways to set kata severities, `kata_severity` (severity name to ID list)
// and the legacy `disabled_katas` list.
type fileConfig struct {
	Config        `yaml:",inline"`
	KataSeverity  map[string][]string `yaml:"kata_severity"`
	DisabledKatas []string            `yaml:"disabled_katas"`
}

// Parse reads a ZShellCheck configuration from YAML.
func Parse(data []byte) (Config, error) {
	cfg, _, err := ParseWithWarnings(data)
	return cfg, err
}

// ParseWithWarnings is Parse plus the non-fatal problems found: unknown
// `kata_severity` keys (ignored) and kata IDs listed more than once (the
// highest severity wins). Unknown top-level keys are ignored silently for
// forward compatibility.
func ParseWithWarnings(data []byte) (Config, []string, error) {
	var raw fileConfig
	if err := yaml.NewDecoder(strings.NewReader(string(data))).Decode(&raw); err != nil && !errors.Is(err, io.EOF) {
		return Config{}, nil, fmt.Errorf("config: %w", err)
	}
	cfg := raw.Config
	cfg.ModifiedSeverities = make(map[string]katas.Severity)

	var warnings []string
	add := func(source string, severity katas.Severity, ids []string) {
		for _, id := range ids {
			if id == "" {
				continue
			}
			if w := katas.SetSeverity(cfg.ModifiedSeverities, id, severity); w != "" {
				warnings = append(warnings, fmt.Sprintf("config: %s: %s", source, w))
			}
		}
	}
	add("disabled_katas", katas.SeverityDisabled, raw.DisabledKatas)
	keys := make([]string, 0, len(raw.KataSeverity))
	for key := range raw.KataSeverity {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		severity, ok := katas.ParseSeverity(key)
		if !ok {
			warnings = append(warnings, fmt.Sprintf("config: kata_severity: unknown key %q ignored", key))
			continue
		}
		add("kata_severity."+key, severity, raw.KataSeverity[key])
	}
	return cfg, warnings, nil
}
