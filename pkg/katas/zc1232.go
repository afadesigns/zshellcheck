package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1232",
		Title:    "Avoid bare `pip install` — use `--user` or virtualenv",
		Severity: SeverityWarning,
		Description: "Bare `pip install` may modify system Python packages. " +
			"Use `pip install --user`, `pipx`, or a virtualenv to isolate dependencies.",
		Check: checkZC1232,
	})
}

func checkZC1232(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}

	if ident.Value != "pip" && ident.Value != "pip3" {
		return nil
	}

	hasInstall := false
	hasSafe := false

	for _, arg := range cmd.Arguments {
		val := arg.String()
		if val == "install" {
			hasInstall = true
		}
		if val == "--user" || val == "-t" || val == "--target" || val == "--prefix" {
			hasSafe = true
		}
	}

	if hasInstall && !hasSafe {
		return []Violation{{
			KataID: "ZC1232",
			Message: "Use `pip install --user` or a virtualenv instead of bare `pip install`. " +
				"System-wide pip installs can break OS package managers.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityWarning,
		}}
	}

	return nil
}
