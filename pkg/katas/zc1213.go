package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1213",
		Title:    "Use `apt-get -y` in scripts for non-interactive installs",
		Severity: SeverityWarning,
		Description: "`apt-get install` without `-y` prompts for confirmation which hangs scripts. " +
			"Use `-y` or set `DEBIAN_FRONTEND=noninteractive` for unattended installs.",
		Check: checkZC1213,
	})
}

func checkZC1213(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "apt-get" {
		return nil
	}

	hasInstall := false
	hasYes := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "install" || val == "upgrade" || val == "dist-upgrade" {
			hasInstall = true
		}
		if val == "-y" || val == "--yes" || val == "-qq" {
			hasYes = true
		}
	}

	if hasInstall && !hasYes {
		return []Violation{{
			KataID: "ZC1213",
			Message: "Use `apt-get -y` in scripts. Without `-y`, apt-get prompts for confirmation " +
				"which hangs in non-interactive execution.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
