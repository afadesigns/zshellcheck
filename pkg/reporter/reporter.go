// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/katas"
)

// Reporter defines the interface for reporting violations.
type Reporter interface {
	Report(violations []katas.Violation) error
}

// TextReporter is a simple reporter that writes plain text to an io.Writer.
type TextReporter struct {
	writer   io.Writer
	filename string
	lines    []string
	config   config.Config
	// fixable, when set, reports whether a kata ID ships an auto-fix; a
	// fixable finding is tagged with a trailing ` [*]` marker.
	fixable func(string) bool
}

// MarkFixable sets the predicate used to tag findings whose kata ships an
// auto-fix with a ` [*]` marker. Passing nil disables the marker.
func (r *TextReporter) MarkFixable(fn func(string) bool) {
	r.fixable = fn
}

// NewTextReporter creates a new TextReporter.
func NewTextReporter(writer io.Writer, filename, source string, config config.Config) *TextReporter {
	return &TextReporter{
		writer:   writer,
		filename: filename,
		lines:    strings.Split(source, "\n"),
		config:   config,
	}
}

func (r *TextReporter) getColor(code string) string {
	if r.config.NoColor {
		return ""
	}
	return code
}

// Report prints the violations to the writer.
// levelColor returns the ANSI colour for a severity, honouring NoColor.
func (r *TextReporter) levelColor(level katas.Severity) string {
	switch level {
	case katas.SeverityError:
		return r.getColor("\033[31m") // Red
	case katas.SeverityWarning:
		return r.getColor("\033[33m") // Yellow
	case katas.SeverityInfo:
		return r.getColor("\033[34m") // Blue
	case katas.SeverityStyle:
		return r.getColor("\033[36m") // Cyan
	}
	return ""
}

func (r *TextReporter) Report(violations []katas.Violation) error {
	for _, v := range violations {
		color := r.levelColor(v.Level)
		reset := r.getColor("\033[0m")
		bold := r.getColor("\033[1m")

		// Location: filename:line:col
		if _, err := fmt.Fprintf(r.writer, "%s:%d:%d: ", r.filename, v.Line, v.Column); err != nil {
			return err
		}

		// Severity, ID, Message, and a ` [*]` marker when the kata is
		// auto-fixable. Example: warning: [ZC1001] Some message [*]
		mark := ""
		if r.fixable != nil && r.fixable(v.KataID) {
			mark = " [*]"
		}
		if _, err := fmt.Fprintf(r.writer, "%s%s%s: [%s] %s%s\n", color, v.Level, reset, v.KataID, v.Message, mark); err != nil {
			return err
		}

		// Code snippet
		if v.Line > 0 && v.Line <= len(r.lines) {
			lineContent := r.lines[v.Line-1]
			// Replace tabs with spaces for correct alignment of caret (simple approach)
			// Or keep it simple for now.
			if _, err := fmt.Fprintf(r.writer, "  %s\n", lineContent); err != nil {
				return err
			}

			// Caret
			padding := v.Column - 1
			if padding < 0 {
				padding = 0
			}
			// Use a simple space padding. Note: this might be slightly off if tabs are present,
			// but it's a standard starting point.
			if _, err := fmt.Fprintf(r.writer, "  %s%s↑%s\n", strings.Repeat(" ", padding), bold, reset); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(r.writer); err != nil {
			return err
		}
	}
	return nil
}
