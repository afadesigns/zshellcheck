package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1483",
		Title:    "Warn on `pip install --break-system-packages` — bypasses PEP 668 externally-managed guard",
		Severity: SeverityWarning,
		Description: "`--break-system-packages` tells pip to ignore the distro's PEP 668 marker " +
			"and install into `/usr/lib/python*`, overwriting files the package manager owns. " +
			"The next `apt`/`dnf` upgrade clobbers or gets clobbered by the pip-installed " +
			"version, and you now have two sources of truth for Python dependencies. Install " +
			"into a virtualenv (`python -m venv`), use `pipx` for application scripts, or use " +
			"`uv` / `poetry` for project dependencies.",
		Check: checkZC1483,
	})
}

func checkZC1483(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "pip" && ident.Value != "pip3" && ident.Value != "pipx" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--break-system-packages" {
			return []Violation{{
				KataID: "ZC1483",
				Message: "`--break-system-packages` installs into distro-managed paths and " +
					"collides with apt/dnf. Use a venv, pipx, or uv/poetry instead.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
