package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1645",
		Title:    "Style: `lsb_release` — prefer sourcing `/etc/os-release` (no dependency, no fork)",
		Severity: SeverityStyle,
		Description: "`lsb_release` is provided by the `lsb-release` / `redhat-lsb-core` " +
			"package, which is missing on most minimal / container images (Alpine does not " +
			"ship it at all). Scripts that depend on `lsb_release` fail the moment they hit " +
			"a stripped image. `/etc/os-release` is standardized by systemd and always " +
			"present on modern Linux — `source /etc/os-release; print -r -- $ID $VERSION_ID` " +
			"gives the same distribution info without the extra package, and without forking.",
		Check: checkZC1645,
	})
}

func checkZC1645(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "lsb_release" {
		return nil
	}

	return []Violation{{
		KataID: "ZC1645",
		Message: "`lsb_release` needs an optional package. Use `source /etc/os-release` and " +
			"read `$ID` / `$VERSION_ID` / `$PRETTY_NAME` instead — always present, no fork.",
		Line:   cmd.Token.Line,
		Column: cmd.Token.Column,
		Level:  SeverityStyle,
	}}
}
