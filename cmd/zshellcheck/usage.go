// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/config"
	"github.com/afadesigns/zshellcheck/pkg/version"
)

// printUsage renders the -h / --help screen. Flags are grouped into
// output / filter / fix / diagnostics columns rather than printed in
// flag-package's default alphabetical run, and a short example list
// follows. Colour is applied when the destination is a TTY and the
// NO_COLOR env var is unset; plain text everywhere else.
func printUsage(out io.Writer, fset *flag.FlagSet, showBanner bool) {
	c := newPalette(out)

	if showBanner {
		fmt.Fprint(out, config.Banner)
	}
	fmt.Fprintf(out, "%s %s — static analysis and auto-fix for Zsh.\n",
		c.bold("zshellcheck"), c.dim("v"+version.Version))
	fmt.Fprintln(out)

	fmt.Fprintf(out, "%s\n  %s %s\n",
		c.section("USAGE"),
		c.bold("zshellcheck"),
		c.dim("[flags] <path> [<path> ...]"))
	fmt.Fprintln(out)

	groups := []flagGroup{
		{
			title: "OUTPUT",
			names: []string{"format", "no-color", "no-banner", "verbose"},
			blurb: "Shape what lands on stdout / stderr.",
		},
		{
			title: "FILTER",
			names: []string{"severity"},
			blurb: "Drop violations below a threshold.",
		},
		{
			title: "AUTO-FIX",
			names: []string{"fix", "diff", "dry-run"},
			blurb: "Apply or preview deterministic rewrites.",
		},
		{
			title: "DIAGNOSTICS",
			names: []string{"cpuprofile", "version"},
			blurb: "Profile, inspect, or print metadata.",
		},
	}

	for _, g := range groups {
		fmt.Fprintln(out, c.section(g.title))
		fmt.Fprintf(out, "  %s\n", c.dim(g.blurb))
		for _, name := range g.names {
			f := fset.Lookup(name)
			if f == nil {
				continue
			}
			renderFlag(out, c, f)
		}
		fmt.Fprintln(out)
	}

	fmt.Fprintln(out, c.section("EXAMPLES"))
	examples := []struct {
		comment, command string
	}{
		{"Lint a single script", "zshellcheck path/to/script.zsh"},
		{"Lint a tree, suppress style-level findings", "zshellcheck -severity warning ./scripts"},
		{"Emit SARIF for GitHub Code Scanning", "zshellcheck -format sarif ./scripts > zshellcheck.sarif"},
		{"Preview every available auto-fix as a diff", "zshellcheck -diff path/to/script.zsh"},
		{"Apply auto-fixes in place", "zshellcheck -fix path/to/script.zsh"},
		{"CI-friendly run (no banner, errors only)", "zshellcheck -no-banner -severity error ./scripts"},
	}
	for _, ex := range examples {
		fmt.Fprintf(out, "  %s\n", c.dim("# "+ex.comment))
		fmt.Fprintf(out, "  %s\n\n", c.bold(ex.command))
	}

	fmt.Fprintf(out, "%s\n  Full guide: %s\n  Katas:      %s\n  Source:     %s\n",
		c.section("DOCUMENTATION"),
		c.link("https://github.com/afadesigns/zshellcheck/blob/main/docs/USER_GUIDE.md"),
		c.link("https://github.com/afadesigns/zshellcheck/blob/main/KATAS.md"),
		c.link("https://github.com/afadesigns/zshellcheck"))
	fmt.Fprintln(out)
}

type flagGroup struct {
	title string
	names []string
	blurb string
}

func renderFlag(out io.Writer, c palette, f *flag.Flag) {
	// Build the name + value-type column. Bool flags show no value placeholder.
	name := "-" + f.Name
	display := name
	if t := flagValueType(f); t != "" {
		display = fmt.Sprintf("%s %s", name, c.dim("<"+t+">"))
	}

	// First line: bold flag name, optional default in dim trailing.
	fmt.Fprintf(out, "  %s", c.flagName(display))
	if def := f.DefValue; def != "" && def != "false" && def != "0" {
		fmt.Fprintf(out, "  %s", c.dim(fmt.Sprintf("(default %q)", def)))
	}
	fmt.Fprintln(out)

	// Second line: usage indented under the name. Wraps long descriptions
	// at ~76 columns to stay readable in 80-col terminals.
	for _, line := range wrap(f.Usage, 76) {
		fmt.Fprintf(out, "      %s\n", line)
	}
}

// flagValueType returns a hint like "string", "duration", or "" for bool.
// Mirrors what flag.PrintDefaults derives but keeps it explicit.
func flagValueType(f *flag.Flag) string {
	getter, ok := f.Value.(flag.Getter)
	if !ok {
		return "string"
	}
	switch getter.Get().(type) {
	case bool:
		return ""
	case string:
		return "string"
	case int, int64:
		return "int"
	case float64:
		return "float"
	default:
		return "value"
	}
}

func wrap(s string, width int) []string {
	if len(s) <= width {
		return []string{s}
	}
	words := strings.Fields(s)
	var out []string
	var line strings.Builder
	for _, w := range words {
		if line.Len() > 0 && line.Len()+1+len(w) > width {
			out = append(out, line.String())
			line.Reset()
		}
		if line.Len() > 0 {
			line.WriteByte(' ')
		}
		line.WriteString(w)
	}
	if line.Len() > 0 {
		out = append(out, line.String())
	}
	return out
}

// palette is a tiny ANSI helper that respects NO_COLOR and TTY status.
type palette struct {
	enabled bool
}

func newPalette(out io.Writer) palette {
	if os.Getenv("NO_COLOR") != "" {
		return palette{}
	}
	f, ok := out.(*os.File)
	if !ok {
		return palette{}
	}
	stat, err := f.Stat()
	if err != nil {
		return palette{}
	}
	if stat.Mode()&os.ModeCharDevice == 0 {
		return palette{}
	}
	return palette{enabled: true}
}

func (p palette) wrap(code, s string) string {
	if !p.enabled {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func (p palette) bold(s string) string     { return p.wrap("1", s) }
func (p palette) dim(s string) string      { return p.wrap("38;5;243", s) }
func (p palette) section(s string) string  { return p.wrap("1;36", s) }
func (p palette) flagName(s string) string { return p.wrap("1;33", s) }
func (p palette) link(s string) string     { return p.wrap("4", s) }
