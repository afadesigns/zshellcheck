package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1542",
		Title:    "Error on `snap install --dangerous` — installs unsigned snap",
		Severity: SeverityError,
		Description: "`snap install --dangerous` tells snapd to install a snap that is not " +
			"assertion-verified. That bypass is named after the risk: any `.snap` file on disk " +
			"can register system services, confinement profiles, and hooks, running as whatever " +
			"user the snap declares. Use `--devmode` for developer work (still verified) or " +
			"ship the snap through the store / a private brand store for production rollouts.",
		Check: checkZC1542,
	})
}

func checkZC1542(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "snap" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--dangerous" {
			return []Violation{{
				KataID: "ZC1542",
				Message: "`snap install --dangerous` installs an assertion-unverified snap — " +
					"any .snap on disk can register system services. Use --devmode or the " +
					"store.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
