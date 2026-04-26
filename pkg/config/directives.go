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
func ParseDirectives(source string) Directives {
	d := Directives{
		PerLine:    make(map[int][]string),
		PerLineAll: make(map[int]bool),
	}

	scanner := bufio.NewScanner(strings.NewReader(source))
	// Guard against single very-long lines overflowing the default buffer.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lineNo := 0
	pendingFrom := 0 // line number of a "preceding" directive waiting for a target
	var pendingIDs []string
	pendingAll := false

	flushPending := func(targetLine int) {
		if pendingAll {
			d.PerLineAll[targetLine] = true
		}
		if len(pendingIDs) > 0 {
			d.PerLine[targetLine] = append(d.PerLine[targetLine], pendingIDs...)
		}
		pendingIDs = nil
		pendingAll = false
		pendingFrom = 0
	}

	for scanner.Scan() {
		lineNo++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		hashIdx := strings.Index(raw, "#")
		if hashIdx < 0 {
			// No comment on this line. If a pending directive is waiting,
			// it applies to this line (assuming the line has content).
			if (pendingAll || len(pendingIDs) > 0) && trimmed != "" {
				flushPending(lineNo)
			}
			continue
		}

		comment := raw[hashIdx+1:]
		match := directiveRe.FindStringSubmatch(comment)
		if match == nil {
			// A regular comment on this line "consumes" a pending directive
			// only if the line has code before the comment.
			if (pendingAll || len(pendingIDs) > 0) && strings.TrimSpace(raw[:hashIdx]) != "" {
				flushPending(lineNo)
			}
			continue
		}

		ids, all := parseDirectiveIDs(match[1])

		// Determine whether the directive is trailing (has code on the same
		// line) or standalone (comment-only line).
		before := strings.TrimSpace(raw[:hashIdx])
		if before != "" {
			// Trailing comment — applies to this line.
			if all {
				d.PerLineAll[lineNo] = true
			}
			if len(ids) > 0 {
				d.PerLine[lineNo] = append(d.PerLine[lineNo], ids...)
			}
			continue
		}

		// Comment-only directive. If it sits above another directive, the
		// earlier one also applied to the upcoming code line; merge them.
		if all {
			pendingAll = true
		}
		if len(ids) > 0 {
			pendingIDs = append(pendingIDs, ids...)
		}
		pendingFrom = lineNo
	}

	// Any directive that never found a target line (empty file tail, or
	// the whole source is just directives) becomes file-wide.
	if pendingFrom > 0 {
		if pendingAll {
			d.FileAll = true
		}
		if len(pendingIDs) > 0 {
			d.File = append(d.File, pendingIDs...)
		}
	}

	return d
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
