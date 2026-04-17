package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1453",
		Title:    "Avoid `sudo pip` / `sudo npm` / `sudo gem` — language package managers as root",
		Severity: SeverityWarning,
		Description: "Running a language package manager as root installs third-party code with " +
			"full privileges, may overwrite distro-managed libs, and can execute arbitrary " +
			"install-time hooks as root. Use `--user`, a virtualenv/venv, or a version manager " +
			"(nvm, pyenv, rbenv) instead.",
		Check: checkZC1453,
	})
}

func checkZC1453(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "sudo" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		switch v {
		case "pip", "pip3", "npm", "yarn", "pnpm", "gem", "cpan", "luarocks":
			return []Violation{{
				KataID: "ZC1453",
				Message: "`sudo " + v + "` runs a language package manager as root. Prefer " +
					"`--user`, a virtualenv/venv, or a version manager (nvm/pyenv/rbenv).",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}

	return nil
}
