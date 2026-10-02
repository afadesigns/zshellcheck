// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// parseRuleSeverity reads a `-rule-severity` value of the form
// `ZC1037:error,ZC1075:style` into a map from kata ID to severity. An
// empty input yields a nil map; a malformed entry or unknown level is an
// error. A repeated ID keeps its highest severity and yields a warning.
func parseRuleSeverity(spec string) (map[string]katas.Severity, []string, error) {
	if strings.TrimSpace(spec) == "" {
		return nil, nil, nil
	}
	out := map[string]katas.Severity{}
	var warnings []string
	for _, pair := range strings.Split(spec, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		id, level, ok := strings.Cut(pair, ":")
		id = strings.TrimSpace(id)
		level = strings.TrimSpace(level)
		if !ok || id == "" {
			return nil, nil, fmt.Errorf("rule-severity: expected ZC####:level, got %q", pair)
		}
		sev, valid := katas.ParseSeverity(level)
		if !valid {
			return nil, nil, fmt.Errorf("rule-severity: %q is not error, warning, info, style, or disabled", level)
		}
		if w := katas.SetSeverity(out, id, sev); w != "" {
			warnings = append(warnings, "rule-severity: "+w)
		}
	}
	return out, warnings, nil
}

// reportStaleNoka writes a line for every per-line `# noka` directive that
// suppresses no actual finding, and tallies them. raw is the file's
// findings before suppression. A bare `# noka` (no IDs) is never stale: it
// is an intentional blanket silence.
func reportStaleNoka(out io.Writer, filename string, raw []katas.Violation, directives config.Directives, counter *int) {
	fires := map[string]bool{}
	for _, v := range raw {
		fires[fmt.Sprintf("%d\x00%s", v.Line, v.KataID)] = true
	}
	type stale struct {
		line int
		id   string
	}
	var stales []stale
	for line, ids := range directives.PerLine {
		for _, id := range ids {
			if !fires[fmt.Sprintf("%d\x00%s", line, id)] {
				stales = append(stales, stale{line, id})
			}
		}
	}
	sort.Slice(stales, func(i, j int) bool {
		if stales[i].line != stales[j].line {
			return stales[i].line < stales[j].line
		}
		return stales[i].id < stales[j].id
	})
	for _, s := range stales {
		fmt.Fprintf(out, "%s:%d: stale `# noka: %s` — no such finding on this line.\n", filename, s.line, s.id)
		if counter != nil {
			*counter++
		}
	}
}

// addNokaDirectives appends a `# noka: ZC####` directive to every line that
// carries a finding and does not already have a `# noka`, then writes the
// file back. It is the bulk "silence what exists" companion to the baseline
// ratchet.
func addNokaDirectives(filename string, data []byte, violations []katas.Violation) error {
	byLine := map[int]map[string]bool{}
	for _, v := range violations {
		if byLine[v.Line] == nil {
			byLine[v.Line] = map[string]bool{}
		}
		byLine[v.Line][v.KataID] = true
	}
	if len(byLine) == 0 {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	changed := false
	for ln, ids := range byLine {
		idx := ln - 1
		if idx < 0 || idx >= len(lines) || strings.Contains(lines[idx], "# noka") {
			continue
		}
		sorted := make([]string, 0, len(ids))
		for id := range ids {
			sorted = append(sorted, id)
		}
		sort.Strings(sorted)
		lines[idx] += "  # noka: " + strings.Join(sorted, ", ")
		changed = true
	}
	if !changed {
		return nil
	}
	mode := os.FileMode(0o600)
	if info, err := os.Stat(filename); err == nil {
		mode = info.Mode().Perm()
	}
	return os.WriteFile(filename, []byte(strings.Join(lines, "\n")), mode)
}
