package katas

import (
	"strings"

	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1507",
		Title:    "Warn on `rsync -l` / default symlink handling — follows escaping symlinks",
		Severity: SeverityWarning,
		Description: "By default rsync copies symlinks as-is but does not prevent one from " +
			"pointing outside the source tree. When the destination is rooted elsewhere (or " +
			"the receiver creates a file at the symlink's resolved path) this becomes a path " +
			"traversal primitive. Use `--safe-links` to skip symlinks pointing outside the " +
			"transfer set, or `--copy-unsafe-links` to materialise them as regular files.",
		Check: checkZC1507,
	})
}

func checkZC1507(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rsync" {
		return nil
	}

	// Trigger only when rsync is actually asked to handle symlinks:
	// -l (preserve symlinks) or -a (archive, which includes -l).
	var hasSymlinkMode, hasSafeHandling bool
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "-a" || v == "-l" || v == "--archive" || v == "--links" ||
			strings.HasPrefix(v, "-a") && strings.ContainsAny(v[1:], "lavx") {
			hasSymlinkMode = true
		}
		if v == "--safe-links" || v == "--copy-unsafe-links" || v == "--no-links" ||
			v == "--munge-links" {
			hasSafeHandling = true
		}
	}
	if !hasSymlinkMode || hasSafeHandling {
		return nil
	}

	return []Violation{{
		KataID: "ZC1507",
		Message: "`rsync` preserving symlinks without `--safe-links` follows ones pointing " +
			"outside the source tree — path traversal vector. Add `--safe-links` or " +
			"`--copy-unsafe-links`.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityWarning,
	}}
}
