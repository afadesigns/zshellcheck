// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package config

import (
	"bufio"
	"regexp"
	"strings"
)

// Directives captures per-file and per-line `# noka` annotations found in a
// source file. Populated by ParseDirectives and consumed alongside the
// config-level DisabledKatas list.
//
// Three forms are recognised, all spelt with the `noka` keyword:
//
//  1. Trailing:    `cmd  # noka`              silences every kata on that line.
//     `cmd  # noka: ZC1234`      silences ZC1234 on that line.
//     `cmd  # noka: ZC1234, ZC1075` silences both on that line.
//  2. Preceding:   a comment-only line whose sole content is `# noka` (or
//     `# noka: …`) silences the next non-blank, non-comment
//     code line.
//  3. File-wide:   a `# noka …` directive at file tail with no following
//     code line silences for the whole file.
//
// `# noka` (no IDs) means "silence everything in scope". Listed IDs limit
// the suppression to those katas only.
type Directives struct {
	// File contains kata IDs disabled file-wide. Populated by trailing
	// directives at file end with no code after them.
	File []string
	// FileAll is true when a file-wide directive used the bare `# noka`
	// (no IDs) form, suppressing every kata for every line.
	FileAll bool
	// PerLine maps a 1-based line number to the kata IDs disabled for
	// just that line.
	PerLine map[int][]string
	// PerLineAll marks lines that carried a bare `# noka` (no IDs)
	// directive, suppressing every kata on the line.
	PerLineAll map[int]bool
}

// HasAny returns true if the directive set disables any kata, anywhere.
func (d Directives) HasAny() bool {
	return d.FileAll || len(d.File) > 0 || len(d.PerLine) > 0 || len(d.PerLineAll) > 0
}

// directiveRe matches the `noka` directive inside any Zsh comment. The
// keyword stands alone or is followed by `:` and a comma/space-separated
// list of kata IDs. The leading `#` is handled by the caller — this regex
// captures from the keyword on.
//
// Group 1 captures the optional ID list. When empty/missing, the directive
// silences every kata in scope.
var directiveRe = regexp.MustCompile(`\bnoka\b\s*(?::\s*([A-Za-z0-9_,\s]+))?`)

// ParseDirectives scans source text for `# noka` annotations and returns
// the file-wide and per-line sets.
type directiveScan struct {
	d           Directives
	pendingIDs  []string
	pendingAll  bool
	pendingFrom int
}

func ParseDirectives(source string) Directives {
	scan := &directiveScan{
		d: Directives{
			PerLine:    make(map[int][]string),
			PerLineAll: make(map[int]bool),
		},
	}
	scanner := bufio.NewScanner(strings.NewReader(source))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		scan.consumeLine(lineNo, scanner.Text())
	}
	scan.absorbTail()
	return scan.d
}

func (s *directiveScan) consumeLine(lineNo int, raw string) {
	hashIdx := strings.Index(raw, "#")
	if hashIdx < 0 {
		s.consumeCodeLine(lineNo, strings.TrimSpace(raw) != "")
		return
	}
	match := directiveRe.FindStringSubmatch(raw[hashIdx+1:])
	if match == nil {
		s.consumeCodeLine(lineNo, strings.TrimSpace(raw[:hashIdx]) != "")
		return
	}
	ids, all := parseDirectiveIDs(match[1])
	if strings.TrimSpace(raw[:hashIdx]) != "" {
		s.recordTrailing(lineNo, ids, all)
		return
	}
	s.recordPending(lineNo, ids, all)
}

func (s *directiveScan) consumeCodeLine(lineNo int, hasContent bool) {
	if !hasContent || (!s.pendingAll && len(s.pendingIDs) == 0) {
		return
	}
	s.flushPending(lineNo)
}

func (s *directiveScan) flushPending(targetLine int) {
	if s.pendingAll {
		s.d.PerLineAll[targetLine] = true
	}
	if len(s.pendingIDs) > 0 {
		s.d.PerLine[targetLine] = append(s.d.PerLine[targetLine], s.pendingIDs...)
	}
	s.pendingIDs = nil
	s.pendingAll = false
	s.pendingFrom = 0
}

func (s *directiveScan) recordTrailing(lineNo int, ids []string, all bool) {
	if all {
		s.d.PerLineAll[lineNo] = true
	}
	if len(ids) > 0 {
		s.d.PerLine[lineNo] = append(s.d.PerLine[lineNo], ids...)
	}
}

func (s *directiveScan) recordPending(lineNo int, ids []string, all bool) {
	if all {
		s.pendingAll = true
	}
	if len(ids) > 0 {
		s.pendingIDs = append(s.pendingIDs, ids...)
	}
	s.pendingFrom = lineNo
}

func (s *directiveScan) absorbTail() {
	if s.pendingFrom == 0 {
		return
	}
	if s.pendingAll {
		s.d.FileAll = true
	}
	if len(s.pendingIDs) > 0 {
		s.d.File = append(s.d.File, s.pendingIDs...)
	}
}

// parseDirectiveIDs extracts the kata-ID list from the directive's optional
// `: ID, ID …` tail. When the tail is absent the directive silences every
// kata in scope (`all` = true, `ids` = nil).
func parseDirectiveIDs(tail string) (ids []string, all bool) {
	tail = strings.TrimSpace(tail)
	if tail == "" {
		return nil, true
	}
	raw := strings.FieldsFunc(tail, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t'
	})
	out := make([]string, 0, len(raw))
	for _, s := range raw {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	if len(out) == 0 {
		return nil, true
	}
	return out, false
}

// IsDisabledOn reports whether the given kata ID is silenced on the
// 1-based line number via this directive set.
func (d Directives) IsDisabledOn(kataID string, line int) bool {
	if d.FileAll || d.PerLineAll[line] {
		return true
	}
	for _, id := range d.File {
		if id == kataID {
			return true
		}
	}
	for _, id := range d.PerLine[line] {
		if id == kataID {
			return true
		}
	}
	return false
}
