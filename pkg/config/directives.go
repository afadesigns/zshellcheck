package config

import (
	"bufio"
	"regexp"
	"strings"
)

// Directives captures per-file and per-line `# zshellcheck disable=…`
// annotations found in a source file. Populated by ParseDirectives and
// consumed alongside the config-level DisabledKatas list.
type Directives struct {
	// File contains kata IDs disabled for every line in the file. Any
	// `# zshellcheck disable=…` comment with no associated target line
	// (e.g. on its own line at the top of the file, before any code)
	// lands here.
	File []string
	// PerLine maps a 1-based line number to the kata IDs disabled for
	// just that line. Two forms produce entries:
	//  1. Trailing comment: `cmd  # zshellcheck disable=ZC1234`
	//     disables ZC1234 on the same line.
	//  2. Preceding directive: a line whose sole content is
	//     `# zshellcheck disable=ZC1234` disables the IDs on the *next*
	//     non-blank, non-comment line.
	PerLine map[int][]string
}

// HasAny returns true if the directive set disables any kata, anywhere.
func (d Directives) HasAny() bool {
	return len(d.File) > 0 || len(d.PerLine) > 0
}

// directiveRe matches the `zshellcheck disable=<ids>` directive inside
// any Zsh comment. The leading `#` (and optional whitespace) is handled
// by the caller — the regexp captures from the directive keyword on.
var directiveRe = regexp.MustCompile(`zshellcheck\s+disable\s*=\s*([A-Za-z0-9_,\s]+)`)

// ParseDirectives scans source text for `# zshellcheck disable=…`
// annotations and returns the file-wide and per-line sets.
func ParseDirectives(source string) Directives {
	d := Directives{PerLine: make(map[int][]string)}

	scanner := bufio.NewScanner(strings.NewReader(source))
	// Guard against single very-long lines overflowing the default buffer.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lineNo := 0
	pendingFrom := 0 // line number of a "preceding" directive waiting for a target
	var pendingIDs []string

	for scanner.Scan() {
		lineNo++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		hashIdx := strings.Index(raw, "#")
		if hashIdx < 0 {
			// No comment on this line. If a pending directive is waiting,
			// it applies to this line (assuming the line has content).
			if len(pendingIDs) > 0 && trimmed != "" {
				d.PerLine[lineNo] = append(d.PerLine[lineNo], pendingIDs...)
				pendingIDs = nil
				pendingFrom = 0
			}
			continue
		}

		comment := raw[hashIdx+1:]
		match := directiveRe.FindStringSubmatch(comment)
		if match == nil {
			// A regular comment on this line "consumes" a pending directive
			// only if the line has code before the comment.
			if len(pendingIDs) > 0 && strings.TrimSpace(raw[:hashIdx]) != "" {
				d.PerLine[lineNo] = append(d.PerLine[lineNo], pendingIDs...)
				pendingIDs = nil
				pendingFrom = 0
			}
			continue
		}

		ids := splitIDs(match[1])

		// Determine whether the directive is trailing (has code on the same
		// line) or standalone (comment-only line).
		before := strings.TrimSpace(raw[:hashIdx])
		if before != "" {
			// Trailing comment — applies to this line.
			d.PerLine[lineNo] = append(d.PerLine[lineNo], ids...)
			continue
		}

		// Comment-only directive. If it sits above another directive, the
		// earlier one also applied to the upcoming code line; merge them.
		pendingIDs = append(pendingIDs, ids...)
		pendingFrom = lineNo
	}

	// Any directive that never found a target line (empty file tail, or
	// the whole source is just directives) becomes file-wide.
	if len(pendingIDs) > 0 && pendingFrom > 0 {
		d.File = append(d.File, pendingIDs...)
	}

	return d
}

func splitIDs(list string) []string {
	raw := strings.FieldsFunc(list, func(r rune) bool {
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
	return out
}

// IsDisabledOn reports whether the given kata ID is silenced on the
// 1-based line number via this directive set.
func (d Directives) IsDisabledOn(kataID string, line int) bool {
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
