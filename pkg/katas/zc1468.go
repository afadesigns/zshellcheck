package katas

import (
	"github.com/afadesigns/zshellcheck/pkg/ast"
)

func init() {
	RegisterKata(ast.SimpleCommandNode, Kata{
		ID:       "ZC1468",
		Title:    "Error on apt `--allow-unauthenticated` / `--force-yes` — installs unsigned packages",
		Severity: SeverityError,
		Description: "`--allow-unauthenticated` and the deprecated `--force-yes` disable APT's " +
			"package-signature verification, turning any MITM or typo-squat into arbitrary " +
			"code execution as root. Always sign internal packages and leave verification on.",
		Check: checkZC1468,
	})
}

func checkZC1468(node ast.Node) []Violation {
	cmd, ok := node.(*ast.SimpleCommand)
	if !ok {
		return nil
	}

	ident, ok := cmd.Name.(*ast.Identifier)
	if !ok {
		return nil
	}
	if ident.Value != "apt" && ident.Value != "apt-get" && ident.Value != "aptitude" {
		return nil
	}

	for _, arg := range cmd.Arguments {
		v := arg.String()
		if v == "--allow-unauthenticated" ||
			v == "--force-yes" ||
			v == "--allow-downgrades" ||
			v == "--allow-remove-essential" ||
			v == "--allow-change-held-packages" {
			return []Violation{{
				KataID: "ZC1468",
				Message: "APT installing unsigned or override-policy packages (" + v + ") — " +
					"disables signature verification, MITM-to-root trivial.",
				Line:   cmd.Token.Line,
				Column: cmd.Token.Column,
				Level:  SeverityError,
			}}
		}
	}
	return nil
}
