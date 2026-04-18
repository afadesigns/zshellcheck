package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1625",
		Title:    "Error on `rm --no-preserve-root` — disables GNU rm safeguard against `rm -rf /`",
		Severity: SeverityError,
		Description: "GNU `rm` refuses to remove `/` by default — the `--preserve-root` " +
			"safeguard added in coreutils 8.4. `--no-preserve-root` explicitly disables that " +
			"check so `rm -rf /` actually recurses and wipes the filesystem. Scripts that pass " +
			"the flag are asking `rm` to go ahead if the argument happens to evaluate to `/`. " +
			"Remove the flag; if a specific path genuinely needs deletion, list it explicitly " +
			"and leave the safeguard in place.",
		Check: checkZC1625,
	})
}

func checkZC1625(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "rm" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		if arg.String() == "--no-preserve-root" {
			return []Violation{{
				KataID: "ZC1625",
				Message: "`rm --no-preserve-root` disables the GNU safeguard against `rm -rf " +
					"/`. Remove the flag; if a specific path needs deletion, list it " +
					"explicitly.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
