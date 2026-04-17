package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1452",
		Title:    "Avoid `npm install -g` — global installs need root, break under multiple Node versions",
		Severity: SeverityStyle,
		Description: "`npm install -g` places packages in a system-wide prefix (typically " +
			"`/usr/local`). That requires sudo, conflicts with Node version managers (nvm, " +
			"asdf, volta), and is rarely what you want in a project. Prefer project-local " +
			"installs (`npm i`), or `pnpm dlx`/`npx` for one-off tools.",
		Check: checkZC1452,
	})
}

func checkZC1452(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "npm" && ident.Value != "yarn" && ident.Value != "pnpm" {
		return nil
	}

	hasInstall := false
	hasGlobal := false
	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "install" || v == "i" || v == "add" {
			hasInstall = true
		}
		if v == "-g" || v == "--global" {
			hasGlobal = true
		}
	}
	if hasInstall && hasGlobal {
		return []Violation{{
			KataID: "ZC1452",
			Message: "`" + ident.Value + " install -g` installs system-wide. Prefer project-local " +
				"install or `npx`/`pnpm dlx` for one-off tools.",
			Line:   cmd.Token.Line,
			Column: cmd.Token.Column,
			Level:  SeverityStyle,
		}}
	}

	return nil
}
