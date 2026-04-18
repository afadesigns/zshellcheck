package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1730",
		Title:    "Warn on `brew install --HEAD <pkg>` — pulls upstream HEAD, no version stability",
		Severity: SeverityWarning,
		Description: "`brew install --HEAD <pkg>` (also `reinstall --HEAD`, `upgrade --HEAD`) " +
			"builds the formula from the upstream source repository's HEAD branch. The " +
			"build is unrepeatable — every run pulls a different commit — and any " +
			"compromised upstream commit lands directly on the install host. Pin to a " +
			"stable release of the formula, or if HEAD is genuinely required, vendor the " +
			"build into a private tap that fixes a specific revision.",
		Check: checkZC1730,
	})
}

func checkZC1730(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok || ident.Value != "brew" {
		return nil
	}
	if len(cmd.Arguments) < 2 {
		return nil
	}

	switch cmd.Arguments[0].String() {
	case "install", "reinstall", "upgrade":
	default:
		return nil
	}

	for _, arg := range cmd.Arguments[1:] {
		v := arg.String()
		if v == "--HEAD" {
			return []Violation{{
				KataID: "ZC1730",
				Message: "`brew " + cmd.Arguments[0].String() + " --HEAD` builds from " +
					"upstream HEAD — every run pulls a different commit. Pin to a " +
					"stable formula release or vendor a private tap with a fixed " +
					"revision.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityWarning,
			}}
		}
	}
	return nil
}
