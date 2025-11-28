package reporter

import (
	"fmt"
	"io"
	"strings"

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
	verbose  bool
	noColor  bool
}

// NewTextReporter creates a new TextReporter.
func NewTextReporter(writer io.Writer, filename, source string, verbose bool, noColor bool) *TextReporter {
	return &TextReporter{
		writer:   writer,
		filename: filename,
		lines:    strings.Split(source, "\n"),
		verbose:  verbose,
		noColor:  noColor,
	}
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func (r *TextReporter) cReset() string  { return r.getColor(colorReset) }
func (r *TextReporter) cRed() string    { return r.getColor(colorRed) }
func (r *TextReporter) cYellow() string { return r.getColor(colorYellow) }
func (r *TextReporter) cCyan() string   { return r.getColor(colorCyan) }
func (r *TextReporter) cBold() string   { return r.getColor(colorBold) }

func (r *TextReporter) getColor(code string) string {
	if r.noColor {
		return ""
	}
	return code
}

// Report prints the violations to the writer.
func (r *TextReporter) Report(violations []katas.Violation) error {
	for _, v := range violations {
		kata, ok := katas.Registry.GetKata(v.KataID)
		if !ok {
			return fmt.Errorf("kata with ID %s not found", v.KataID)
		}

		// Format: file:line:col: [Level] [ID] Title: Message
		// Example: demo.zsh:10:5: [Warning] [ZC1001] Array Access: Invalid array access...

		levelColor := r.cYellow()
		if v.Level == katas.Error {
			levelColor = r.cRed()
		} else if v.Level == katas.Info {
			levelColor = r.cCyan()
		}

		fmt.Fprintf(r.writer, "%s%s:%d:%d:%s %s[%s]%s %s[%s]%s %s%s:%s %s\n",
			r.cBold(), r.filename, v.Line, v.Column, r.cReset(),
			levelColor, v.Level, r.cReset(),
			r.cRed(), v.KataID, r.cReset(),
			r.cCyan(), kata.Title, r.cReset(),
			v.Message,
		)

		// Print source line context with gutter
		if v.Line > 0 && v.Line <= len(r.lines) {
			line := r.lines[v.Line-1]
			lineNumStr := fmt.Sprintf("%d", v.Line)
			
			// Gutter padding based on line number length
			// For simple alignment, we can assume line numbers won't differ wildly in width for adjacent errors,
			// but let's just print it directly.
			
			fmt.Fprintf(r.writer, "  %s%s |%s %s\n", r.cCyan(), lineNumStr, r.cReset(), line)
			
			// Calculate pointer padding
			pad := ""
			for i := 0; i < v.Column-1 && i < len(line); i++ {
				if line[i] == '\t' {
					pad += "\t"
				} else {
					pad += " "
				}
			}
			
			gutterSpace := strings.Repeat(" ", len(lineNumStr))
			fmt.Fprintf(r.writer, "  %s %s|%s %s%s^%s\n", gutterSpace, r.cCyan(), r.cReset(), pad, r.cYellow(), r.cReset())
		}

		if r.verbose {
			fmt.Fprintf(r.writer, "  %sDescription:%s %s\n", r.cBold(), r.cReset(), kata.Description)
		}
		fmt.Fprintln(r.writer) // Add blank line between violations
	}
	return nil
}
