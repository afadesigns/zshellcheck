package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

var zc1810RecursiveFlags = map[string]bool{
	"-r":          true,
	"--recursive": true,
	"-m":          true,
	"--mirror":    true,
}

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1810",
		Title:    "Warn on `wget -r` / `--mirror` without `--level=N` — unbounded recursive download",
		Severity: SeverityWarning,
		Description: "`wget -r` and `wget --mirror` (short `-m`) follow links to arbitrary depth. " +
			"Without `--level=N` or `-l N` the crawl keeps going until `wget` hits the " +
			"remote server's limits, fills the local disk, or climbs into a parent directory " +
			"the author did not intend to mirror (add `--no-parent` to block that too). " +
			"Pin a depth (`--level=3`), restrict siblings (`--no-parent`, `--accept=` / " +
			"`--reject=`), and cap the byte budget (`--quota=1G`) before running a recursive " +
			"wget in automation.",
		Check: checkZC1810,
	})
}

func checkZC1810(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}
	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "wget" {
		return nil
	}

	hasRecursive := false
	for _, arg := range cmd.Arguments {
		if zc1810RecursiveFlags[arg.String()] {
			hasRecursive = true
			break
		}
	}
	if !hasRecursive {
		return nil
	}
	if zc1810HasLevel(cmd) {
		return nil
	}
	return zc1810Hit(cmd)
}

func zc1810HasLevel(cmd *ast.SimpleCommand) bool {
	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch {
		case v == "-l" || v == "--level":
			return true
		case strings.HasPrefix(v, "--level="):
			return true
		case strings.HasPrefix(v, "-l") && len(v) > 2:
			return true
		}
	}
	return false
}

func zc1810Hit(cmd *ast.SimpleCommand) []Violation {
	line, col := FlagArgPosition(cmd, zc1810RecursiveFlags)
	return []Violation{{
		KataID: "ZC1810",
		Message: "`wget -r` / `--mirror` without `--level=N` follows links to " +
			"arbitrary depth — the crawl can exhaust disk and climb into parent " +
			"paths. Pin `--level=3`, add `--no-parent`, and cap with `--quota=1G`.",
		Line:   line,
		Column: col,
		Level:  SeverityWarning,
	}}
}
