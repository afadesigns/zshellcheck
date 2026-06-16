// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"os"
	"sort"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// baselineState drives the `-baseline` / `-baseline-write` ratchet. In
// write mode it accumulates a fingerprint for every finding so a snapshot
// can be saved; otherwise known holds the snapshot's fingerprints and any
// finding already in it is suppressed, leaving only new findings.
type baselineState struct {
	write   bool
	known   map[string]bool
	collect []string
}

// loadBaseline reads a snapshot file into a filtering baselineState.
func loadBaseline(path string) (*baselineState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	known := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" {
			known[line] = true
		}
	}
	return &baselineState{known: known}, nil
}

// baselineFingerprint identifies a finding by kata, file, and the trimmed
// source line — not the line number — so inserting or removing unrelated
// lines elsewhere in the file does not invalidate the snapshot.
func baselineFingerprint(file string, lines []string, v katas.Violation) string {
	content := ""
	if v.Line >= 1 && v.Line <= len(lines) {
		content = strings.TrimSpace(lines[v.Line-1])
	}
	return v.KataID + "\t" + file + "\t" + content
}

// applyBaseline records or filters this file's findings against the
// baseline. In write mode it returns the findings unchanged after
// collecting them; otherwise it returns only the findings absent from the
// snapshot.
func (b *baselineState) applyBaseline(filename string, data []byte, violations []katas.Violation) []katas.Violation {
	lines := strings.Split(string(data), "\n")
	if b.write {
		for _, v := range violations {
			b.collect = append(b.collect, baselineFingerprint(filename, lines, v))
		}
		return violations
	}
	kept := violations[:0]
	for _, v := range violations {
		if !b.known[baselineFingerprint(filename, lines, v)] {
			kept = append(kept, v)
		}
	}
	return kept
}

// writeBaseline saves the collected fingerprints as a sorted, de-duplicated
// snapshot.
func (b *baselineState) writeBaseline(path string) error {
	seen := map[string]bool{}
	uniq := make([]string, 0, len(b.collect))
	for _, fp := range b.collect {
		if !seen[fp] {
			seen[fp] = true
			uniq = append(uniq, fp)
		}
	}
	sort.Strings(uniq)
	out := strings.Join(uniq, "\n")
	if out != "" {
		out += "\n"
	}
	return os.WriteFile(path, []byte(out), 0o600)
}
